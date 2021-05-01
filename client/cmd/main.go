// Entrypoint for PSOStats Client
package main

import (
	"log"
	"os"

	"github.com/phelix-/psostats/v2/client/internal/client"
)

func main() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	version := "0.6.2"

	log.SetOutput(file)
	log.Printf("Starting Up version %v", version)

	c, err := client.New(version)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}
	c.Run()
}
