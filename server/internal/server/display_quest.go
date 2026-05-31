package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"net/url"
	"time"
)

func (s *Server) QuestPage(c *fiber.Ctx) error {
	quest, err := url.PathUnescape(c.Params("quest"))
	if err != nil {
		c.Status(500)
		return err
	}
	player := c.Query("player")

	questRecords, err := db.GetQuestRecordsForQuest(quest, s.dynamoClient)
	if err != nil {
		return err
	}

	recordsByCategory := make(map[string]model.FormattedGame)
	for _, game := range questRecords {
		fg := getFormattedGame(game)
		fg.Record = true
		recordsByCategory[game.Category] = fg
	}

	var playerGames []model.FormattedGame
	playerNotFound := false

	if player != "" {
		rawGames, err := db.GetAllGamesForPlayerQuest(player, quest, s.dynamoClient)
		if err != nil {
			return err
		}
		if len(rawGames) == 0 {
			playerNotFound = true
		}

		pbIds := make(map[string]string)
		bestTimes := make(map[string]time.Duration)
		for _, g := range rawGames {
			if best, ok := bestTimes[g.Category]; !ok || g.Time < best {
				bestTimes[g.Category] = g.Time
				pbIds[g.Category] = g.Id
			}
		}

		for _, game := range rawGames {
			fg := getFormattedGame(game)
			if pbIds[game.Category] == game.Id {
				fg.Pb = true
			}
			if rec, ok := recordsByCategory[game.Category]; ok && rec.Id == game.Id {
				fg.Record = true
				fg.Pb = false
			}
			playerGames = append(playerGames, fg)
		}
	}

	pageModel := struct {
		QuestName      string
		Records        map[string]model.FormattedGame
		Player         string
		PlayerGames    []model.FormattedGame
		PlayerNotFound bool
	}{
		QuestName:      quest,
		Records:        recordsByCategory,
		Player:         player,
		PlayerGames:    playerGames,
		PlayerNotFound: playerNotFound,
	}

	err = s.questTemplate.ExecuteTemplate(c.Response().BodyWriter(), "quest", pageModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}
