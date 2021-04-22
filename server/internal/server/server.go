package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
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
	s.app.Get("/api/game/:gameId", s.GetGame)
	s.app.Get("recent", s.GetRecentGames)
	s.app.Get("/records", s.RecordsPage)
	s.app.Get("/players/:player", s.PlayerPage)
	if certLocation, found := os.LookupEnv("CERT"); found {
		keyLocation := os.Getenv("KEY")
		if err := s.app.ListenTLS(":443", certLocation, keyLocation); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := s.app.Listen(":80"); err != nil {
			log.Fatal(err)
		}
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

func (s *Server) RecordsPage(c *fiber.Ctx) error {
	t, err := template.ParseFiles("./server/internal/templates/records.gohtml")
	if err != nil {
		c.Status(500)
		return err
	}
	games, err := db.GetQuestRecords(s.dynamoClient)
	sort.Slice(games, func(i, j int) bool { return games[i].Quest < games[j].Quest })

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

func (s *Server) PlayerPage(c *fiber.Ctx) error {
	player := c.Params("player")
	t, err := template.ParseFiles("./server/internal/templates/player.gohtml")
	if err != nil {
		c.Status(500)
		return err
	}

	player, err = url.PathUnescape(player)
	if err != nil {
		c.Status(500)
		return err
	}
	//player = string(playerBytes)
	pbs, err := db.GetPlayerPbs(player, s.dynamoClient)

	if err != nil {
		log.Print("get pbs")
		c.Status(500)
		return err
	}
	sort.Slice(pbs, func(i, j int) bool { return pbs[i].Quest < pbs[j].Quest })
	for i, game := range pbs {
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
		pbs[i] = game
	}

	recent, err := db.GetPlayerRecentGames(player, s.dynamoClient)

	if err != nil {
		log.Print("get recent")
		c.Status(500)
		return err
	}
	for i, game := range recent {
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
		recent[i] = game
	}

	model := struct {
		PlayerPbs []model.Game
		RecentGames []model.Game
	}{
		PlayerPbs: pbs,
		RecentGames: recent,
	}
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	err = t.ExecuteTemplate(c.Response().BodyWriter(), "index", model)
	return err
}

func (s *Server) PostGame(c *fiber.Ctx) error {
	user, pass, err := s.getUserFromBasicAuth(c.Request().Header.Peek("Authorization"))
	if err != nil {
		c.Status(401)
		return nil
	}
	userObject, err := db.GetUser(s.dynamoClient, user)
	if err != nil {
		c.Status(401)
		return nil
	}
	if passwordsMatch := doPasswordsMatch(userObject.Password, pass); !passwordsMatch {
		c.Status(401)
		return nil
	}
	var questRun model.QuestRun
	if err := c.BodyParser(&questRun); err != nil {
		log.Printf("body parser")
		c.Status(400)
		return err
	}
	questDuration, err := time.ParseDuration(questRun.QuestDuration)
	if err != nil {
		c.Status(400)
		return err
	}
	questRun.GuildCard = user
	gameId, err := db.WriteGame(&questRun, s.dynamoClient)
	if err != nil {
		log.Printf("write game %v", err)
		c.Status(500)
		return err
	}
	questRun.Id = gameId
	if isLeaderboardCandidate(questRun) {
		numPlayers := len(questRun.AllPlayers)
		topRun, err := db.GetQuestRecord(questRun.QuestName, numPlayers, questRun.PbCategory, s.dynamoClient)
		if err != nil {
			log.Printf("failed to get top quest runs for gameId:%v - %v", gameId, err)
		} else if topRun == nil || topRun.Time > questDuration {
			log.Printf("new record for %v %vp pb:%v - %v",
				questRun.QuestName, numPlayers, questRun.PbCategory, gameId)
			if err = db.WriteQuestRecord(&questRun, s.dynamoClient); err != nil {
				log.Printf("failed to update leaderboard for game %v - %v", gameId, err)
			}
		}
		if err = db.WriteGameByQuest(&questRun, s.dynamoClient); err != nil {
			log.Printf("failed to update games by quest for game %v - %v", gameId, err)
		}
		playerPb, err := db.GetPlayerPB(questRun.QuestName, user, numPlayers, questRun.PbCategory, s.dynamoClient)
		if err != nil {
			log.Printf("failed to get player pb for gameId:%v - %v", gameId, err)
		} else if playerPb == nil || playerPb.Time > questDuration {
			log.Printf("new pb for %v %v %vp pb:%v - %v",
				user, questRun.QuestName, numPlayers, questRun.PbCategory, gameId)
			if err = db.WritePlayerPb(&questRun, s.dynamoClient); err != nil {
				log.Printf("failed to update pb for game %v - %v", gameId, err)
			}
		}
	}
	if err = db.WriteGameByPlayer(&questRun, s.dynamoClient); err != nil {
		log.Printf("failed to update games by player for game %v - %v", gameId, err)
	}

	c.Response().AppendBodyString(gameId)
	log.Printf("got quest: %v %v, %v, %v, %v",
		gameId, questRun.QuestName, questRun.PlayerName, questRun.Server, questRun.GuildCard)
	return nil
}

func (s *Server) GetGame(c *fiber.Ctx) error {
	gameId := c.Params("gameId")
	game, err := db.GetGame(gameId, s.dynamoClient)
	if err != nil {
		return err
	}

	if game == nil {
		c.Status(404)
		return nil
	} else {
		jsonBytes, err := json.Marshal(game)
		if err != nil {
			return err
		}
		c.Response().AppendBody(jsonBytes)
		c.Response().Header.Set("Content-Type", "application/json")
		return nil
	}
}

func isLeaderboardCandidate(questRun model.QuestRun) bool {
	allowedDifficulty := questRun.Difficulty == "Ultimate" || strings.HasPrefix(questRun.QuestName, "Stage")
	return allowedDifficulty && questRun.QuestComplete && !questRun.IllegalShifta
}

func (s *Server) getUserFromBasicAuth(headerBytes []byte) (string, string, error) {
	headerString := string(headerBytes)
	if len(headerString) > 0 && strings.HasPrefix(headerString, "Basic ") {
		authBase64 := strings.TrimPrefix(headerString, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(authBase64)
		if err != nil {
			return "", "", err
		}
		auth := string(decoded)
		authSplit := strings.SplitN(auth, ":", 2)

		return authSplit[0], authSplit[1], nil
	} else {
		return "", "", errors.New("missing basic auth header")
	}
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
