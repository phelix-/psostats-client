package server

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
)

type Server struct {
	app          *fiber.App
	dynamoClient *dynamodb.DynamoDB
}

func New(dynamo *dynamodb.DynamoDB) *Server {
	f := fiber.New(fiber.Config{
		// modify config
	})
	return &Server{
		app:          f,
		dynamoClient: dynamo,
	}
}

func (s *Server) Run() {
	s.app.Static("/", "./static", fiber.Static{
		// modify config
	})
	s.app.Post("quests", s.PostGame)
	if err := s.app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) PostGame(c *fiber.Ctx) error {
	var questRun model.QuestRun
	if err := c.BodyParser(&questRun); err != nil {
		log.Printf("body parser")
		return err
	}
	log.Printf("got quest: %v", questRun)
	return nil
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
