package main

import (
	"log"

	"wallet_service/config"
	"wallet_service/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	server, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer server.DB.Close()

	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := server.Router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
