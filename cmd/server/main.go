package main

import (
	"log"
	"redis-go/internal/commands"
	"redis-go/internal/router"
	"redis-go/internal/store"
	"redis-go/pkg/connection"
)

func main() {
	s := store.New()
	r := router.New()
	commands.Register(r, s)
	server := connection.TCPServer{
		Handler: r.Handle,
	}
	log.Println("Server listening on :8080")
	server.Setup(":8080")
}
