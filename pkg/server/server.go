package server

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/phelix-/psostats/v2/pkg/model"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
	DynamoClient *dynamodb.DynamoDB
}

func (s *Server) Run() {
	ui := http.FileServer(http.Dir("./static"))
	http.Handle("/", ui)
	http.HandleFunc("/game", s.PostGame)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) PostGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
		r.Response.StatusCode = 400
	} else {
		questRun := model.QuestRun{}

		allBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Read error %v", err)
		}

		log.Printf("%v", len(allBytes))
		err = json.Unmarshal(allBytes, &questRun)
		if err != nil {
			log.Printf("got error %v", err)
		} else {
			log.Printf("Got quest: %v", questRun)
		}
	}
}

//func (c *Client) runHttp() {
//	fileServer := http.FileServer(http.Dir("./static"))
//	http.Handle("/", fileServer)
//	http.HandleFunc("/game/info", func(w http.ResponseWriter, r *http.Request) {
//		bytes, err := json.Marshal(c.GetGameInfo())
//		if err != nil {
//			r.Response.StatusCode = 500
//			fmt.Fprintf(w, "Error!")
//			return
//		}
//		w.Header().Add("Content-Type", "application/json")
//		w.Write(bytes)
//	})
//	addr := fmt.Sprintf(":%v", c.config.GetUiPort())
//	log.Printf("Hosting local ui at localhost%v", addr)
//}