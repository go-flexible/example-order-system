package domain

import (
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Dependencies interface {
	Redis() redis.Cmdable
	DB() *pgxpool.Pool
}

// Publisher wraps the requisite methods used from nats.
type NatsPublisher interface {
	Publish(subject string, v interface{}) error
}
