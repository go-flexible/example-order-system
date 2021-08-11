// Package domain_test contains all the required bootstrapping for testing the
// domain package as a consumer / user of it's public api.
package domain_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/postgres"
	"github.com/go-flexible/example-order-system/svc/config"
	"github.com/go-flexible/example-order-system/svc/migrations"
	"github.com/go-redis/redis"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats.go"
)

var embedSource = migration.EmbedMigrationSource{
	EmbedFS: migrations.FS,
	Dir:     ".",
}

var (
	db              *pgxpool.Pool
	migrationDriver migration.Driver
)

func TestMain(m *testing.M) {
	var spec config.Specification
	envconfig.MustProcess("", &spec)

	// migrate
	var err error
	migrationDriver, err = postgres.New(spec.DSN)
	if err != nil {
		log.Fatal(err)
	}
	count, err := migration.Migrate(migrationDriver, embedSource, migration.Up, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d migrations applied\n", count)

	// connect
	db, err = pgxpool.Connect(context.Background(), spec.DSN)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func dropAllTables(t *testing.T) {
	t.Helper()
	exec := func(_ pgconn.CommandTag, err error) {
		if err != nil {
			t.Log(err)
		}
	}
	exec(db.Exec(context.Background(), "delete from payments"))
	exec(db.Exec(context.Background(), "delete from line_items"))
	exec(db.Exec(context.Background(), "delete from order_totals"))
	exec(db.Exec(context.Background(), "delete from order_metadata"))
	exec(db.Exec(context.Background(), "delete from orders"))

	exec(db.Exec(context.Background(), "drop table if exists payments"))
	exec(db.Exec(context.Background(), "drop table if exists line_items"))
	exec(db.Exec(context.Background(), "drop table if exists order_totals"))
	exec(db.Exec(context.Background(), "drop table if exists order_metadata"))
	exec(db.Exec(context.Background(), "drop table if exists orders"))
	exec(db.Exec(context.Background(), "drop table if exists schema_migration"))
}

// provide mocked dependencies for the tests.
type dependencies struct{}

func (d dependencies) DB() *pgxpool.Pool    { return db }
func (d dependencies) Redis() redis.Cmdable { panic("not implemented, consider mocking") }

// mock the Publisher.
type natsPulisherMock struct{ *nats.EncodedConn }

func (n *natsPulisherMock) Publish(subject string, v interface{}) error { return nil }
