package main

import (
	"log"

	"github.com/zhashkevych/package-ordering/internal/api"
)

func main() {
	// Default packs; can be changed at runtime via PUT /packs
	// Not best for production, but for the sake of the exercise, it's fine. Should be loaded from a database or a config file.
	defaultPacks := []int{250, 500, 1000, 2000, 5000}

	server := api.NewServer(defaultPacks)
	if err := server.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
