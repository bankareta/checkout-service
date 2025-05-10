package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"base-golang/internal/config"
	"base-golang/internal/environtment"

	"sync"
	"syscall"
)

func init() {

	environtment.LoadConfig()
}

func main() {
	configEnv := environtment.Configs
	log := config.NewLogger(configEnv)
	_, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	// Setup channel to listen for interrupt signals for graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Running Worker")

	// Add the number of consumers to the WaitGroup
	wg.Add(1)

	// Wait for a shutdown signal
	go func() {
		<-shutdownChan
		log.Println("Shutting down gracefully...")
		cancel() // cancel context to stop consuming messages
	}()

	// Wait for all consumers to finish
	fmt.Println("waiting go routine to finish")
	wg.Wait()
}
