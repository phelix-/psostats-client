package server

import (
	"encoding/json"
	"fmt"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"log"
	"text/template"
	"time"

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
	s.app.Static("/main.css", "./static/main.css", fiber.Static{
		// modify config
	})
	s.app.Get("/", s.Index)
	s.app.Get("/game/:gameId", s.GamePage)
	s.app.Post("/api/game", s.PostGame)
	s.app.Get("recent", s.GetRecentGames)
	if err := s.app.Listen(":80"); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Index(c *fiber.Ctx) error {
	t, err := template.ParseFiles("./server/internal/templates/index.gohtml")
	if err != nil {
		c.Status(500)
		return err
	}
	games, err := db.GetRecentGames(s.dynamoClient)
	if err != nil {
		log.Print("get recent games")
		c.Status(500)
		return err
	}
	for i, game := range games {
		minutes := game.Time / time.Minute
		seconds := (game.Time % time.Minute) / time.Second
		game.FormattedTime = fmt.Sprintf("%01d:%02d", minutes, seconds)
		shortCategory := game.Category
		numPlayers := string(shortCategory[0])
		pbRun := string(shortCategory[1])
		pbText := ""
		if pbRun == "p" {
			pbText = " PB"
		}
		game.Category = numPlayers + " Player" + pbText
		game.FormattedDate = game.Timestamp.In(time.Local).Format("15:04 01/02/2006")
		games[i] = game
	}
	model := struct {
		Games []model.Game
	}{
		Games: games,
	}
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	err = t.ExecuteTemplate(c.Response().BodyWriter(), "index", model)
	return err
}
func (s *Server) GamePage(c *fiber.Ctx) error {
	gameId := c.Params("gameId")
	game, err := db.GetGame(gameId, s.dynamoClient)
	if err != nil {
		return err
	}

	if game == nil {
		t, err := template.ParseFiles("./server/internal/templates/gameNotFound.gohtml")
		if err != nil {
			return err
		}
		err = t.ExecuteTemplate(c.Response().BodyWriter(), "gameNotFound", game)
	} else {
		t, err := template.ParseFiles("./server/internal/templates/game.gohtml")
		if err != nil {
			return err
		}
		err = t.ExecuteTemplate(c.Response().BodyWriter(), "game", game)
	}
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func (s *Server) PostGame(c *fiber.Ctx) error {
	var questRun model.QuestRun
	if err := c.BodyParser(&questRun); err != nil {
		log.Printf("body parser")
		c.Status(400)
		return err
	}
	gameId, err := db.WriteGame(&questRun, s.dynamoClient)
	if err != nil {
		log.Printf("write game")
		c.Status(500)
		return err
	}
	c.Response().AppendBodyString(gameId)
	log.Printf("got quest: %v %v, %v, %v, %v",
		gameId, questRun.QuestName, questRun.PlayerName, questRun.Server, questRun.GuildCard)
	return nil
}

func (s *Server) GetRecentGames(c *fiber.Ctx) error {
	games, err := db.GetRecentGames(s.dynamoClient)
	if err != nil {
		log.Print("get recent games")
		c.Status(500)
		return err
	}
	bytes, err := json.Marshal(games)
	if err != nil {
		return err
	}
	_, err = c.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// func (c *Client) runHttp() {
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
// }
