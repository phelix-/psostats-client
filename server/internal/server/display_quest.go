package server

import (
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/server/internal/db"
)

func (s *Server) QuestPage(c *fiber.Ctx) error {
	quest, err := url.PathUnescape(c.Params("quest"))
	if err != nil {
		c.Status(500)
		return err
	}
	player := strings.TrimSpace(c.Query("player"))

	questRecords, err := db.GetQuestRecordsForQuest(quest, s.dynamoClient)
	if err != nil {
		return err
	}

	recordsByCategory := make(map[string]model.FormattedGame)
	for _, game := range questRecords {
		if len(game.Category) < 2 {
			continue
		}
		fg := getFormattedGame(game)
		fg.Record = true
		recordsByCategory[game.Category] = fg
	}
	sortedRecords := make([]model.FormattedGame, 0, len(recordsByCategory))
	for _, fg := range recordsByCategory {
		sortedRecords = append(sortedRecords, fg)
	}
	sort.Slice(sortedRecords, func(i, j int) bool {
		if sortedRecords[i].NumPlayers != sortedRecords[j].NumPlayers {
			return sortedRecords[i].NumPlayers < sortedRecords[j].NumPlayers
		}
		return sortedRecords[i].PbRun && !sortedRecords[j].PbRun
	})

	playerGamesByCount := make(map[int][]model.FormattedGame)
	var playerCounts []int
	playerNotFound := false
	resultsCapped := false

	if player != "" {
		const maxResults = 500
		rawGames, err := db.GetAllGamesForPlayerQuest(player, quest, maxResults, s.dynamoClient)
		if err != nil {
			return err
		}
		if len(rawGames) == 0 {
			runes := []rune(player)
			titleCased := strings.ToUpper(string(runes[:1])) + strings.ToLower(string(runes[1:]))
			if titleCased != player {
				retried, retryErr := db.GetAllGamesForPlayerQuest(titleCased, quest, maxResults, s.dynamoClient)
				if retryErr != nil {
					return retryErr
				}
				if len(retried) > 0 {
					rawGames = retried
					player = titleCased
				}
			}
		}
		if len(rawGames) == 0 {
			playerNotFound = true
		}
		resultsCapped = len(rawGames) == maxResults

		pbIds := make(map[string]string)
		bestTimes := make(map[string]time.Duration)
		for _, g := range rawGames {
			if g.Category == "" {
				continue
			}
			if best, ok := bestTimes[g.Category]; !ok || g.Time < best {
				bestTimes[g.Category] = g.Time
				pbIds[g.Category] = g.Id
			}
		}

		for _, game := range rawGames {
			if len(game.Category) < 2 {
				continue
			}
			fg := getFormattedGame(game)
			if pbIds[game.Category] == game.Id {
				fg.Pb = true
			}
			if rec, ok := recordsByCategory[game.Category]; ok && rec.Id == game.Id {
				fg.Record = true
			}
			playerGamesByCount[fg.NumPlayers] = append(playerGamesByCount[fg.NumPlayers], fg)
		}

		for np, games := range playerGamesByCount {
			sort.Slice(games, func(i, j int) bool {
				return games[i].Duration < games[j].Duration
			})
			playerGamesByCount[np] = games
		}

		for np := 1; np <= 4; np++ {
			if len(playerGamesByCount[np]) > 0 {
				playerCounts = append(playerCounts, np)
			}
		}
	}

	defaultPlayer := ""
	for _, rec := range sortedRecords {
		if rec.NumPlayers == 1 && rec.PbRun {
			for _, p := range rec.Players {
				if p.Name != "" {
					defaultPlayer = p.Name
					break
				}
			}
			break
		}
	}

	pageModel := struct {
		QuestName          string
		Records            []model.FormattedGame
		Player             string
		DefaultPlayer      string
		PlayerGamesByCount map[int][]model.FormattedGame
		PlayerCounts       []int
		PlayerNotFound     bool
		ResultsCapped      bool
	}{
		QuestName:          quest,
		Records:            sortedRecords,
		Player:             player,
		DefaultPlayer:      defaultPlayer,
		PlayerGamesByCount: playerGamesByCount,
		PlayerCounts:       playerCounts,
		PlayerNotFound:     playerNotFound,
		ResultsCapped:      resultsCapped,
	}

	err = s.questTemplate.ExecuteTemplate(c.Response().BodyWriter(), "quest", pageModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}
