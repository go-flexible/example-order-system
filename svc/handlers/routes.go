package handlers

import (
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/nats.go"
)

// Dependencies allows injecting the requisite dependencies.
type Dependencies interface {
	Redis() redis.Cmdable
	DB() *pgxpool.Pool
	Nats() *nats.EncodedConn
}

// RegisterRotues allows adding handlers to the provider router.
func RegisterRoutes(router *mux.Router, deps Dependencies) {
	router.Handle("/order", createOrder(deps)).Methods(http.MethodPost)
	router.Handle("/order/{id}", getOrderByID(deps)).Methods(http.MethodGet)
}
