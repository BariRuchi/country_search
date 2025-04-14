package main

import (
	"CountrySearch/countryhandler"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	mux := http.NewServeMux()
	handler := countryhandler.New()
	mux.HandleFunc("/api/countries/search", handler.CountryHandler)

	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	// Channel to listen for interrupt or terminate signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, os.Kill)

	// Run server in goroutine
	go func() {
		log.Println("Server started on :8000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	// Block until signal is received
	<-stopChan
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully")
}
