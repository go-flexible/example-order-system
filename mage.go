//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func Run() error {
	if err := sh.RunV("docker-compose", "up", "-d"); err != nil {
		return err
	}
	return sh.RunV("go", "run", "svc/main.go")
}

func Test() error {
	if err := sh.RunV("docker-compose", "up", "-d"); err != nil {
		return err
	}
	return sh.RunWithV(map[string]string{
		"DB_DSN": "postgres://root:@localhost:26258/defaultdb?sslmode=disable",
	}, "go", "test", "-v", "./...", "--cover")
}
