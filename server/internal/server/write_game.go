package server

import (
	"encoding/base64"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"github.com/valyala/fasthttp"
	"log"
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

	var matchingGame *model.QuestRun = nil
	for _, recentGame := range s.recentGames {
		if gamesMatch(recentGame, questRun) {
			log.Printf("matched game[%v]", recentGame.Id)
			matchingGame = &recentGame
			break
		}
	}

	if matchingGame != nil {
		questRun.Id = matchingGame.Id
		err := db.AttachGameToId(questRun, matchingGame.Id, s.dynamoClient)
		if err != nil {
			log.Printf("%v", err)
		}
	} else {
		gameId, err := db.WriteGameById(&questRun, s.dynamoClient)
		if err != nil {
			log.Printf("write game %v", err)
			c.Status(500)
			return err
		}
		questRun.Id = gameId
		s.recentGames[s.recentGamesCount%s.recentGamesSize] = questRun
		s.recentGamesCount++
	}

	if IsLeaderboardCandidate(questRun) {
		numPlayers := len(questRun.AllPlayers)
		if matchingGame == nil {
			topRun, err := db.GetQuestRecord(questRun.QuestName, numPlayers, questRun.PbCategory, s.dynamoClient)
			if err != nil {
				log.Printf("failed to get top quest runs for gameId:%v - %v", questRun.Id, err)
			} else if topRun == nil || topRun.Time > questDuration {
				log.Printf("new record for %v %vp pb:%v - %v",
					questRun.QuestName, numPlayers, questRun.PbCategory, questRun.Id)
				if err = db.WriteGameByQuestRecord(&questRun, s.dynamoClient); err != nil {
					log.Printf("failed to update leaderboard for game %v - %v", questRun.Id, err)
				}
			}
			if err = db.WriteGameByQuest(&questRun, s.dynamoClient); err != nil {
				log.Printf("failed to update games by quest for game %v - %v", questRun.Id, err)
			}
		}
		playerPb, err := db.GetPlayerPB(questRun.QuestName, user, numPlayers, questRun.PbCategory, s.dynamoClient)
		if err != nil {
			log.Printf("failed to get player pb for gameId:%v - %v", questRun.Id, err)
		} else if playerPb == nil || playerPb.Time > questDuration {
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

	c.Response().AppendBodyString(questRun.Id)
	log.Printf("got quest: %v %v, %v, %v, %v",
		questRun.Id, questRun.QuestName, questRun.PlayerName, questRun.Server, questRun.UserName)
	return nil
}
