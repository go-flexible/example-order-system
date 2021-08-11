package main

import (
	"context"
	"embed"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
	"github.com/go-flexible/example-order-system/svc/config"
	"github.com/go-flexible/example-order-system/svc/handlers"
	"github.com/go-flexible/example-order-system/svc/workers"
	"github.com/go-flexible/flex"
	"github.com/go-flexible/flexhttp"
	"github.com/go-flexible/flexmetrics"
	"github.com/go-flexible/flexready"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
)

//go:embed migrations
// TODO: find a good way to make migrations work as an embed.
var migrations embed.FS

var embedSourced = migration.EmbedMigrationSource{
	EmbedFS: migrations,
	Dir:     "migrations",
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var spec config.Specification
	envconfig.MustProcess("", &spec)

	services := backingServices{
		redis:     connectRedis(),
		cockroach: connectCockroach(spec.DSN),
		nats:      connectNats(),
	}

	migrateDatabase(spec.DSN)

	router := mux.NewRouter()
	handlers.RegisterRoutes(router, services)

	// setup our http server with custom timeouts.
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}

	// configure the readiness server.
	readysrvConfig := &flexready.Config{
		Server: &http.Server{Addr: "localhost:9999"},
	}
	// prepare the readiness checks.
	readysrv := flexready.New(readysrvConfig, flexready.Checks{
		"redis":       func() error { return redisCheck(services.redis) },
		"nats":        func() error { return natsCheck(services.nats) },
		"cockroachdb": func() error { return crdbCheck(services.cockroach) },
	})

	// start all of our workers.
	flex.MustStart(
		context.Background(),
		readysrv,
		flexmetrics.New(nil),
		&workers.NATSWorker{Nats: services.nats},
		flexhttp.NewHTTPServer(srv),
	)
}

// backingServices holds the services that are used by the application.
type backingServices struct {
	redis     redis.Cmdable
	cockroach *pgxpool.Pool
	nats      *nats.EncodedConn
}

// satisfy the Dependencies interface.
func (b backingServices) Redis() redis.Cmdable    { return b.redis }
func (b backingServices) DB() *pgxpool.Pool       { return b.cockroach }
func (b backingServices) Nats() *nats.EncodedConn { return b.nats }

func connectRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
}

func connectCockroach(dsn string) *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}

func connectNats() *nats.EncodedConn {
	natsClient, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	encodedConn, err := nats.NewEncodedConn(natsClient, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

	return encodedConn
}

func migrateDatabase(dsn string) {
	driver, err := postgres.New(dsn)
	if err != nil {
		log.Fatal(err)
	}
	count, err := migration.Migrate(driver, embedSourced, migration.Up, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d migrations applied\n", count)
}

// check readiness of redis.
func redisCheck(redisClient redis.Cmdable) error {
	if redisClient.Ping().Val() == "PONG" {
		return nil
	}
	return errors.New("redis ping did not respond with pong")
}

// check readiness of nats.
func natsCheck(natsClient *nats.EncodedConn) error {
	if !natsClient.Conn.IsConnected() {
		return errors.New("nats client is not connected")
	}
	return nil
}

// check readiness of cockroachdb.
func crdbCheck(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return pool.Ping(ctx)
}
