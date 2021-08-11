package config

type Specification struct {
	DSN string `envconfig:"DB_DSN" required:"true" default:"postgres://root:@localhost:26257/defaultdb?sslmode=disable"`
}
