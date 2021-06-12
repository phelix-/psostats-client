package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"strings"
	"time"
)

func getUserFromBasicAuth(headerBytes []byte) (string, string, error) {
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

func (s *Server) verifyAuth(header *fasthttp.RequestHeader) (bool, string) {
	user, pass, err := getUserFromBasicAuth(header.Peek("Authorization"))
	if err != nil || len(user) < 1 {
		return false, ""
	}
	userObject, err := s.userDb.GetUser(user)
	if err != nil {
		return false, ""
	}
	passwordsMatch := DoPasswordsMatch(userObject.Password, pass)
	return passwordsMatch, user
}

func (s *Server) PostGame(c *fiber.Ctx) error {
	authorized, user := s.verifyAuth(&c.Request().Header)
	if !authorized {
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
	questRun.UserName = user
	questRun.SubmittedTime = time.Now()

	matchingGame := s.findMatchingGame(questRun)

	if matchingGame == nil {
		s.recentGamesLock.Lock()
		// Check again inside the lock
		matchingGame = s.findMatchingGame(questRun)
		if matchingGame == nil {
			gameId, err := db.WriteGameById(&questRun, s.dynamoClient)
			if err != nil {
				log.Printf("write game %v", err)
				c.Status(500)
				s.recentGamesLock.Unlock()
				return err
			}
			questRun.Id = gameId
			s.recentGames[s.recentGamesCount%s.recentGamesSize] = questRun
			s.recentGamesCount++
		}
		s.recentGamesLock.Unlock()
	}
	if matchingGame != nil {
		questRun.Id = matchingGame.Id
		err := db.AttachGameToId(questRun, matchingGame.Id, s.dynamoClient)
		if err != nil {
			log.Printf("%v", err)
		}
	}

	record := false
	pb := false
	if IsLeaderboardCandidate(questRun) {
		s.recordsLock.Lock()
		numPlayers := len(questRun.AllPlayers)
		topRun, err := db.GetQuestRecord(questRun.QuestName, numPlayers, questRun.PbCategory, s.dynamoClient)
		if err != nil {
			log.Printf("failed to get top quest runs for gameId:%v - %v", questRun.Id, err)
		} else if matchingGame != nil {
			if topRun == nil {
				log.Printf("Matching game but no topRun, almost definitely a bug")
			} else {
				record = matchingGame.Id == topRun.Id
				if err := db.AddPovToRecord(questRun, s.dynamoClient); err != nil {
					log.Printf("failed to add pov to record")
				}
			}
		} else if topRun == nil || topRun.Time > questDuration {
			record = true
			s.QuestRecordWebhook(questRun, topRun)
			log.Printf("new record for %v %vp pb:%v - %v",
				questRun.QuestName, numPlayers, questRun.PbCategory, questRun.Id)
			if err = db.WriteGameByQuestRecord(&questRun, s.dynamoClient); err != nil {
				log.Printf("failed to update leaderboard for game %v - %v", questRun.Id, err)
			}
		}
		s.recordsLock.Unlock()
		if err = db.WriteGameByQuest(&questRun, s.dynamoClient); err != nil {
			log.Printf("failed to update games by quest for game %v - %v", questRun.Id, err)
		}

		playerPb, err := db.GetPlayerPB(questRun.QuestName, user, numPlayers, questRun.PbCategory, s.dynamoClient)
		if err != nil {
			log.Printf("failed to get player pb for gameId:%v - %v", questRun.Id, err)
		} else if playerPb == nil || playerPb.Time > questDuration {
			pb = true
			log.Printf("new pb for %v %v %vp pb:%v - %v",
				user, questRun.QuestName, numPlayers, questRun.PbCategory, questRun.Id)
			if err = db.WritePlayerPb(&questRun, s.dynamoClient); err != nil {
				log.Printf("failed to update pb for game %v - %v", questRun.Id, err)
			}
		}
	}
	if err = db.WriteGameByPlayer(&questRun, s.dynamoClient); err != nil {
		log.Printf("failed to update games by player for game %v - %v", questRun.Id, err)
	}

	jsonBytes, err := json.Marshal(model.PostGameResponse{
		Pb:     pb,
		Record: record,
		Id:     questRun.Id,
	})
	if err != nil {
		return err
	}
	c.Response().AppendBody(jsonBytes)
	c.Response().Header.Set("Content-Type", "application/json")

	log.Printf("got quest: %v %v, %v, %v, %v",
		questRun.Id, questRun.QuestName, questRun.PlayerName, questRun.Server, questRun.UserName)
	return nil
}

func (s *Server) findMatchingGame(questRun model.QuestRun) *model.QuestRun {
	var matchingGame *model.QuestRun = nil
	for _, recentGame := range s.recentGames {
		if GamesMatch(recentGame, questRun) {
			log.Printf("matched game[%v]", recentGame.Id)
			matchingGame = &recentGame
			break
		}
	}
	return matchingGame
}

func (s *Server) QuestRecordWebhook(questRun model.QuestRun, previousRecord *model.Game) {
	if len(s.webhookUrl) > 0 {
		duration, err := time.ParseDuration(questRun.QuestDuration)
		formattedDuration := formatDuration(duration)
		playersString := ""
		for _, player := range questRun.AllPlayers {
			playersString = fmt.Sprintf("%v%v - %v\n", playersString, player.Class, player.Name)
		}
		previousRecordText := ""
		if previousRecord != nil {
			difference := previousRecord.Time - duration
			previousRecordText = "\nbeating the previous record by " + formatDuration(difference)
		}
		jsonBytes, err := json.Marshal(Webhook{Embeds: []Embed{
			{
				Title: "New Record: " + questRun.QuestName,
				Description: fmt.Sprintf("%v https://psostats.com/game/%v%v",
					formattedDuration, questRun.Id, previousRecordText),
				Fields: []Field{
					{Name: "Players", Value: playersString, Inline: true},
				},
			},
		}})
		if err != nil {
			log.Printf("Failed to marshal data %v", err)
		}

		urls := strings.Split(s.webhookUrl, ",")
		for _, url := range urls {
			buf := bytes.NewBuffer(jsonBytes)
			_, err := http.Post(url, "application/json", buf)
			if err != nil {
				log.Printf("Failed to perform webhook %v", err)
			}
		}
	}
}

type Webhook struct {
	Embeds []Embed `json:"embeds"`
}
type Embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Fields      []Field `json:"fields"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
