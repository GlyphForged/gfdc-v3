package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gfdc-v3/modules/web"
)

func main() {
	// Parse listen address from env, if none, assume dev (good old localhost:8080)
	listenAddr := os.Getenv("ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}
	server, err := web.NewServer()
	if err != nil {
		log.Fatalf("failed to build server: %v", err)
	}

	msg := fmt.Sprintf("Serving on port %v", listenAddr)
	log.Println(msg)
	if err := http.ListenAndServe(listenAddr, server.Routes()); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
