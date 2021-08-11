package domain

import (
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/nats.go"
)

type Dependencies interface {
	Redis() redis.Cmdable
	DB() *pgxpool.Pool
	Nats() *nats.EncodedConn
}
