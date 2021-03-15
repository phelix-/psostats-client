package main

import (
	"log"
	"os"

	"github.com/phelix-/psostats/v2/pkg/client"
)

func Hello() string {
	return "Hello, world."
}

func main() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.Print("Starting Up")

	c, err := client.New("0.0.1")
	if err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	c.Run()
}
