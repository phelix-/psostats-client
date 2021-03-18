package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
	http.HandleFunc("/game/raw", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(client.GetFrames())
		if err != nil {
			r.Response.StatusCode = 500
			fmt.Fprintf(w, "Error!")
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
	})
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		minFrame, hasMinFrame := queryParams["min"]
		frames := client.GetFrames()
		gameInfo := client.GetGameInfo()

		start := 0
		log.Printf("%v %v", hasMinFrame, minFrame)
		if hasMinFrame && len(minFrame) > 0 {
			start, _ = strconv.Atoi(minFrame[0])
		}
		hp := make([]int, 0)
		tp := make([]int, 0)
		mesetaCharged := make([]int, 0)
		i := start
		for ; true; i++ {
			if value, exists := frames[i]; exists {
				hp = append(hp, int(value.HP))
				tp = append(tp, int(value.TP))
				mesetaCharged = append(mesetaCharged, int(value.MesetaCharged))
			} else {
				break
			}
		}
		ret := make(map[string]interface{})
		ret["QuestName"] = gameInfo.QuestName
		ret["QuestStarted"] = gameInfo.QuestStarted
		ret["QuestComplete"] = gameInfo.QuestComplete
		ret["QuestStartTime"] = gameInfo.QuestStartTime
		ret["QuestEndTime"] = gameInfo.QuestEndTime
		ret["QuestDuration"] = gameInfo.QuestDuration.String()

		ret["hp"] = hp
		ret["tp"] = tp
		ret["mesetaCharged"] = mesetaCharged
		ret["lastFrame"] = start + len(hp)
		bytes, err := json.Marshal(ret)
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
