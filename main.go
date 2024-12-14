package main

import (
	"log"

	"github.com/walteranderson/tromba/internal/config"
	"github.com/walteranderson/tromba/internal/project"
)

func main() {
	conf := config.Load()
	_, err := project.Build(conf)
	if err != nil {
		log.Fatal("ERROR: project build", err)
	}
}
