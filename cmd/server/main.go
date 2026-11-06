package main

import (
	"log"

	"github.com/rafialfauzan/layrd/internal/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatalf("Application terminated with error: %v", err)
	}
}
