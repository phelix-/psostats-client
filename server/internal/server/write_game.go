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
	"regexp"
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
	if err != nil || userObject == nil {
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
	// just making sure it parses
	_, err := time.ParseDuration(questRun.QuestDuration)
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
		otherPbCategory, _ := db.GetQuestRecord(questRun.QuestName, numPlayers, !questRun.PbCategory, s.dynamoClient)
		if err != nil {
			log.Printf("failed to get top quest runs for gameId:%v - %v", questRun.Id, err)
		} else if matchingGame != nil {
			if topRun == nil {
				log.Printf("Matching game but no topRun, almost definitely a bug")
			} else if matchingGame.Id == topRun.Id {
				record = true
				if err := db.AddPovToRecord(db.QuestRecordsTable, questRun, s.dynamoClient); err != nil {
					log.Printf("failed to add pov to record")
				}
			}
		} else if isNewRecord(questRun, topRun, otherPbCategory) {
			record = true
			s.QuestRecordWebhook(questRun, topRun)
			log.Printf("new record for %v %vp pb:%v - %v",
				questRun.QuestName, numPlayers, questRun.PbCategory, questRun.Id)
			if err = db.WriteGameByQuestRecord(&questRun, s.dynamoClient); err != nil {
				log.Printf("failed to update leaderboard for game %v - %v", questRun.Id, err)
			}
		}
		s.updateAnniv2020Record(questRun, matchingGame)
		s.recordsLock.Unlock()

		playerPb, err := db.GetPlayerPB(questRun.QuestName, user, numPlayers, questRun.PbCategory, s.dynamoClient)
		if err != nil {
			log.Printf("failed to get player pb for gameId:%v - %v", questRun.Id, err)
		} else if isBetterRun(questRun, playerPb) {
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

func (s *Server) updateAnniv2020Record(questRun model.QuestRun, matchingGame *model.QuestRun) {
	_,anniversaryQuest := s.anniversaryQuests[questRun.QuestName]
	if questRun.PbCategory || !anniversaryQuest || questRun.Server != "ephinea" {
		return
	}
	numPlayers := len(questRun.AllPlayers)
	topRun, err := db.GetAnniv2021Record(questRun.QuestName, numPlayers, questRun.PbCategory, s.dynamoClient)
	if err != nil {
		log.Printf("failed to get top quest runs for gameId:%v - %v", questRun.Id, err)
	} else if matchingGame != nil {
		if topRun == nil {
			log.Printf("Matching game but no topRun, almost definitely a bug")
		} else if matchingGame.Id == topRun.Id {
			if err := db.AddPovToRecord(db.Anniv2021RecordsTable, questRun, s.dynamoClient); err != nil {
				log.Printf("failed to add pov to anniv record")
			}
		}
	} else if isBetterRun(questRun, topRun) {
		log.Printf("new record for %v %vp pb:%v - %v",
			questRun.QuestName, numPlayers, questRun.PbCategory, questRun.Id)
		if err = db.WriteAnniv2021Record(&questRun, s.dynamoClient); err != nil {
			log.Printf("failed to update anniv leaderboard for game %v - %v", questRun.Id, err)
		}
	}
}

func isNewRecord(
	currentRun model.QuestRun,
	previousRecord *model.Game,
	otherPbCategory *model.Game,
) bool {
	if currentRun.PbCategory {
		return isBetterRun(currentRun, previousRecord) && isBetterRun(currentRun, otherPbCategory)
	} else {
		return isBetterRun(currentRun, previousRecord)
	}
}

func isBetterRun(
	currentRun model.QuestRun,
	other *model.Game,
) bool {
	if other == nil {
		return true
	}
	questDuration, _ := time.ParseDuration(currentRun.QuestDuration)
	if isRankedByScore(currentRun) {
		if int(currentRun.Points) > other.Points {
			return true
		}
		if int(currentRun.Points) == other.Points && questDuration < other.Time {
			// Same score but faster
			return true
		}
		return false
	} else {
		return questDuration < other.Time
	}
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

func IsLeaderboardCandidate(questRun model.QuestRun) bool {
	cmodeRegex := regexp.MustCompile("[12]c\\d")
	allowedDifficulty := questRun.Difficulty == "Ultimate" || cmodeRegex.MatchString(questRun.QuestName)
	return allowedDifficulty && questRun.QuestComplete && !questRun.IllegalShifta
}

func (s *Server) QuestRecordWebhook(questRun model.QuestRun, previousRecord *model.Game) {
	if len(s.webhookUrl) > 0 {
		duration, err := time.ParseDuration(questRun.QuestDuration)
		formattedScore := ""
		if isRankedByScore(questRun) {
			formattedScore = fmt.Sprintf("%d points in ", questRun.Points)
		}
		formattedDuration := formatDuration(duration)
		playersString := ""
		for _, player := range questRun.AllPlayers {
			playersString = fmt.Sprintf("%v%v - %v\n", playersString, player.Class, player.Name)
		}
		previousRecordText := ""
		if previousRecord != nil {
			timeDifference := previousRecord.Time - duration
			if isRankedByScore(questRun) {
				previousRecordText = fmt.Sprintf("\nbeating the previous record by %v points", int(questRun.Points) - previousRecord.Points)
				if timeDifference >= 0 {
					previousRecordText = fmt.Sprintf("%v (%v faster)", previousRecordText, formatDuration(timeDifference))
				} else {
					previousRecordText = fmt.Sprintf("%v (%v slower)", previousRecordText, formatDuration(-timeDifference))
				}
			} else {
				previousRecordText = "\nbeating the previous record by " + formatDuration(timeDifference)
			}
		}
		jsonBytes, err := json.Marshal(Webhook{Embeds: []Embed{
			{
				Title: "New Record: " + questRun.QuestName,
				Description: fmt.Sprintf("%v%v https://psostats.com/game/%v%v",
					formattedScore, formattedDuration, questRun.Id, previousRecordText),
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

func isRankedByScore(questRun model.QuestRun) bool {
	return questRun.QuestName == "Endless: Episode 1"
}