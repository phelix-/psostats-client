package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/pkg/psoclasses"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (s *Server) PlayerV2Page(c *fiber.Ctx) error {
	player := c.Params("player")
	player, err := url.PathUnescape(player)
	if err != nil {
		c.Status(500)
		return err
	}
	recentGames, err := db.GetPlayerRecentGames(player, s.dynamoClient, 15)
	if err != nil {
		return err
	}
	games, err := db.GetQuestRecords(db.QuestRecordsTable, s.dynamoClient)
	if err != nil {
		return err
	}
	sortedRecords := sortGames(games)
	playerPbs, err := db.GetPlayerPbs(player, s.dynamoClient)
	if err != nil {
		return err
	}
	sortedPbs := sortGames(playerPbs)
	maePbs := make(map[string]string)

	for _, pb := range playerPbs {
		if sortedRecords[pb.Episode][pb.Quest][pb.Category].Id == pb.Id {
			pbsForEp := sortedPbs[pb.Episode]
			pbsForQuest := pbsForEp[pb.Quest]
			pbForCategory := pbsForQuest[pb.Category]
			pbForCategory.Record = true
			pbsForQuest[pb.Category] = pbForCategory
			pbsForEp[pb.Quest] = pbsForQuest
			sortedPbs[pb.Episode] = pbsForEp
		}
	}
	validTotalDuration := true
	totalDuration := time.Duration(0)
	for name, _ := range s.anniversaryQuests {
		var bestGame = model.FormattedGame{}
		for _, quest := range sortedPbs[1][name] {
			if bestGame.Duration == 0 {
				bestGame = quest
			} else if quest.Duration < bestGame.Duration {
				bestGame = quest
			}
		}
		for _, quest := range sortedPbs[2][name] {
			if bestGame.Duration == 0 {
				bestGame = quest
			} else if quest.Duration < bestGame.Duration {
				bestGame = quest
			}
		}
		for _, quest := range sortedPbs[4][name] {
			if bestGame.Duration == 0 {
				bestGame = quest
			} else if quest.Duration < bestGame.Duration {
				bestGame = quest
			}
		}
		if bestGame.Duration > 0 {
			maePbs[name] = bestGame.Time
			totalDuration += bestGame.Duration
		} else {
			maePbs[name] = "No Game"
			validTotalDuration = false
		}
	}
	maeTotal := "Incomplete Lap"
	if validTotalDuration {
		maeTotal = formatDuration(totalDuration)
	}
	classUsage, err := db.GetPlayerClassCounts(player, s.dynamoClient)
	if err != nil {
		return err
	}
	for _, class := range psoclasses.GetAll() {
		if _, exists := classUsage[class.Name]; !exists {
			classUsage[class.Name] = 0
		}
	}
	questsPlayed, err := db.GetPlayerQuestCounts(player, s.dynamoClient)
	if err != nil {
		return err
	}
	gamesByEpisode := map[int]int{1: 0, 2: 0, 4: 0}
	questCounts := make([]QuestAndCount, 0)
	totalGames := 0
	for quest, count := range questsPlayed {
		split := strings.SplitN(quest, "_", 2)
		if len(split) == 2 {
			episode, err := strconv.Atoi(split[0])
			questName := split[1]
			if err == nil {
				totalGames += count
				gamesByEpisode[episode] += count
				questCounts = append(questCounts, QuestAndCount{questName, count})
			}
		}
	}
	favoriteQuest := QuestAndCount{"None", 0}
	if len(questCounts) > 0 {
		sort.Slice(questCounts, func(i, j int) bool { return questCounts[i].Count > questCounts[j].Count })
		favoriteQuest = questCounts[0]
	}

	infoModel := struct {
		PlayerName     string
		Classes        map[string]int
		TotalGames     int
		GamesByEpisode map[int]int
		FavoriteQuest  QuestAndCount
		RecentGames    []model.FormattedGame
		PbGames        map[int]map[string]map[string]model.FormattedGame
		MaePbs         map[string]string
		MaeTotal       string
	}{
		PlayerName:     player,
		Classes:        classUsage,
		TotalGames:     totalGames,
		GamesByEpisode: gamesByEpisode,
		FavoriteQuest:  favoriteQuest,
		RecentGames:    make([]model.FormattedGame, 0),
		PbGames:        sortedPbs,
		MaePbs:         maePbs,
		MaeTotal:       maeTotal,
	}
	for _, game := range recentGames {
		formattedGame := getFormattedGame(game)
		pbForQuestAndCategory := sortedPbs[game.Episode][game.Quest][game.Category]
		if pbForQuestAndCategory.Id == game.Id {
			formattedGame.Pb = true
			if pbForQuestAndCategory.Record {
				formattedGame.Record = true
			}
		}
		infoModel.RecentGames = append(infoModel.RecentGames, formattedGame)
	}
	err = s.playerTemplate.ExecuteTemplate(c.Response().BodyWriter(), "player", infoModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}
