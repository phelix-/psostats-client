package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	c, err := client.New("0.1.0")
	if err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	go runHttp(c)
	c.Run()
}

func runHttp(client *client.Client) {

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/game/info", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(client.GetGameInfo())
		if err != nil {
			r.Response.StatusCode = 500
			fmt.Fprintf(w, "Error!")
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
