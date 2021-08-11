// +build mage

package main

import "github.com/magefile/mage/sh"

func Run() error {
	if err := sh.Run("docker-compose", "up", "-d"); err != nil {
		return err
	}
	return sh.Run("go", "run", "svc/main.go")
}
