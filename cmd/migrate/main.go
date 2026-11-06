package main

import (
	"log"

	"github.com/rafialfauzan/layrd/internal/bootstrap"
)

func main() {
	if err := bootstrap.Migrate(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations completed successfully.")
}
