package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/pkg/psoclasses"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (s *Server) Anniv2022RecordsPage(c *fiber.Ctx) error {
	overallCounters, questCounters := s.getCounters()
	records, err := db.GetQuestRecords(db.Anniv2021RecordsTable, s.dynamoClient)
	if err != nil {
		log.Printf("get recent games %v", err)
		c.Status(500)
		return err
	}
	recordHistory, err := db.GetQuestRecords(db.AnnivRecordHistory, s.dynamoClient)
	sortedRecordHistory := s.sortRecordHistory(recordHistory)
	sortedRecs := sortAnnivGames(records)
	recordModel := struct {
		QuestNames      []string
		QuestShortNames []string
		TopLaps         []AnniversaryTimes
		OverallCounter  QuestCounters
		QuestCounters   map[string]QuestCounters
		Records         map[string]map[string]model.FormattedGame
		Classes         []psoclasses.PsoClass
		SectionIds      []string
		RecordHistory   map[string][]RecordHistoryPoint
	}{
		QuestNames: s.anniversaryNamesInOrder,
		QuestShortNames: []string{"Forest",
			"Caves",
			"Mines",
			"Ruins",
			"Temple",
			"Space",
			"CCA",
			"Seabed",
			"Tower",
			"Crater",
			"Desert"},
		TopLaps:        s.getTopLaps(),
		OverallCounter: overallCounters,
		QuestCounters:  questCounters,
		Records:        sortedRecs,
		Classes:        psoclasses.GetAll(),
		SectionIds: []string{
			"Viridia",
			"Greenill",
			"Skyly",
			"Bluefull",
			"Purplenum",
			"Pinkal",
			"Redria",
			"Oran",
			"Yellowboze",
			"Whitill",
		},
		RecordHistory: sortedRecordHistory,
	}

	err = s.anniversary2022Template.ExecuteTemplate(c.Response().BodyWriter(), "index", recordModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func sortAnnivGames(games []model.Game) map[string]map[string]model.FormattedGame {
	recordModel := make(map[string]map[string]model.FormattedGame)
	for _, game := range games {
		formattedGame := getFormattedGame(game)

		gamesForQuest := recordModel[game.Quest]
		if gamesForQuest == nil {
			gamesForQuest = make(map[string]model.FormattedGame)
		}
		gamesForQuest[game.Category] = formattedGame
		recordModel[game.Quest] = gamesForQuest
	}
	return recordModel
}

func (s *Server) sortRecordHistory(games []model.Game) map[string][]RecordHistoryPoint {
	sort.Slice(games, func(i, j int) bool {
		return games[i].Timestamp.Before(games[j].Timestamp)
	})
	gamesByQuest := make(map[string][]RecordHistoryPoint)
	for _, quest := range s.anniversaryNamesInOrder {
		gamesByQuest[quest] = make([]RecordHistoryPoint, 0)
	}
	for i, game := range games {
		gamesForQuest := gamesByQuest[game.Quest]
		nextGame := RecordHistoryPoint{Time: game.Timestamp}
		if len(gamesForQuest) > 0 {
			lastGame := gamesForQuest[len(gamesForQuest)-1]
			nextGame.P1 = lastGame.P1
			nextGame.P2 = lastGame.P2
			nextGame.P3 = lastGame.P3
			nextGame.P4 = lastGame.P4
		}
		switch game.Category {
		case "1n":
			nextGame.P1 = &games[i]
		case "2n":
			nextGame.P2 = &games[i]
		case "3n":
			nextGame.P3 = &games[i]
		case "4n":
			nextGame.P4 = &games[i]
		}
		gamesForQuest = append(gamesForQuest, nextGame)
		gamesByQuest[game.Quest] = gamesForQuest
	}
	for quest, gamesForQuest := range gamesByQuest {
		nextGame := RecordHistoryPoint{Time: time.Now()}
		if len(gamesForQuest) > 0 {
			lastGame := gamesForQuest[len(gamesForQuest)-1]
			nextGame.P1 = lastGame.P1
			nextGame.P2 = lastGame.P2
			nextGame.P3 = lastGame.P3
			nextGame.P4 = lastGame.P4
		}
		gamesForQuest = append(gamesForQuest, nextGame)
		gamesByQuest[quest] = gamesForQuest
	}
	return gamesByQuest
}

func (s *Server) getCounters() (QuestCounters, map[string]QuestCounters) {
	questCounters := make(map[string]QuestCounters)

	overallCounter := QuestCounters{
		ClassMesetaCharged: make(map[string]int64),
		ClassUse:           make(map[string]int64),
		SID:                make(map[string]int64),
		Players:            make(map[int]int64),
		Shifta:             make(map[int]int64),
		RunsByDay:          make(map[time.Time]int64),
	}
	for _, questName := range s.anniversaryNamesInOrder {
		questCounters[questName] = QuestCounters{
			ClassMesetaCharged: make(map[string]int64),
			ClassUse:           make(map[string]int64),
			SID:                make(map[string]int64),
			Players:            make(map[int]int64),
			Shifta:             make(map[int]int64),
			RunsByDay:          make(map[time.Time]int64),
		}
	}

	if counters, err := db.GetAnniversaryCounters(s.dynamoClient); err == nil {
		for _, counter := range counters {
			if counterForQuest, found := questCounters[counter.Key]; found {
				if counter.Counter == "MesetaCharged" {
					counterForQuest.MesetaCharged = counter.Count
					overallCounter.MesetaCharged += counter.Count
				} else if strings.HasPrefix(counter.Counter, "MesetaCharged.") {
					class := strings.TrimPrefix(counter.Counter, "MesetaCharged.")
					counterForQuest.ClassMesetaCharged[class] = counter.Count
					overallCounter.ClassMesetaCharged[class] = overallCounter.ClassMesetaCharged[class] + counter.Count
				} else if counter.Counter == "Runs" {
					counterForQuest.Runs = counter.Count
					overallCounter.Runs += counter.Count
				} else if counter.Counter == "Moving" {
					counterForQuest.Moving = counter.Count
					overallCounter.Moving += counter.Count
				} else if counter.Counter == "Standing" {
					counterForQuest.Standing = counter.Count
					overallCounter.Standing += counter.Count
				} else if counter.Counter == "Attacking" {
					counterForQuest.Attacking = counter.Count
					overallCounter.Attacking += counter.Count
				} else if counter.Counter == "Casting" {
					counterForQuest.Casting = counter.Count
					overallCounter.Casting += counter.Count
				} else if counter.Counter == "Deaths" {
					counterForQuest.Deaths = counter.Count
					overallCounter.Deaths += counter.Count
				} else if strings.HasPrefix(counter.Counter, "ClassUse.") {
					class := strings.TrimPrefix(counter.Counter, "ClassUse.")
					counterForQuest.ClassUse[class] = counter.Count
					overallCounter.ClassUse[class] = overallCounter.ClassUse[class] + counter.Count
				} else if strings.HasPrefix(counter.Counter, "SID.") {
					sid := strings.TrimPrefix(counter.Counter, "SID.")
					sidNumber, _ := strconv.Atoi(sid)
					counterForQuest.SID[getSectionId(sidNumber)] = counter.Count
					overallCounter.SID[getSectionId(sidNumber)] = overallCounter.SID[getSectionId(sidNumber)] + counter.Count
				} else if strings.HasPrefix(counter.Counter, "Players.") {
					players := strings.TrimPrefix(counter.Counter, "Players.")
					playerCount, _ := strconv.Atoi(players)
					counterForQuest.Players[playerCount] = counter.Count
					overallCounter.Players[playerCount] = overallCounter.Players[playerCount] + counter.Count
				} else if strings.HasPrefix(counter.Counter, "Shifta.") {
					shifta := strings.TrimPrefix(counter.Counter, "Shifta.")
					shiftaLevel, _ := strconv.Atoi(shifta)
					counterForQuest.Shifta[shiftaLevel] = counter.Count
					overallCounter.Shifta[shiftaLevel] = overallCounter.Shifta[shiftaLevel] + counter.Count
				} else if strings.HasPrefix(counter.Counter, "Day.") {
					day := strings.TrimPrefix(counter.Counter, "Day.")
					date, _ := time.Parse("060102", day)
					counterForQuest.RunsByDay[date] = counter.Count
					overallCounter.RunsByDay[date] = overallCounter.RunsByDay[date] + counter.Count
				}
				questCounters[counter.Key] = counterForQuest
			}
		}
	}
	return overallCounter, questCounters
}

func (s *Server) getTopLaps() []AnniversaryTimes {
	anniversaryTimes := make([]AnniversaryTimes, 0)
	timesByPlayer := make(map[string]map[string]db.QuestSeriesPb)
	if pbs, err := db.GetQuestSeriesPbs("a2022", s.dynamoClient); err == nil {
		for _, pb := range pbs {
			annivTimes, found := timesByPlayer[pb.User]
			if !found {
				annivTimes = make(map[string]db.QuestSeriesPb)
			}
			annivTimes[pb.QuestName] = pb
			timesByPlayer[pb.User] = annivTimes
		}
	}

	for user, times := range timesByPlayer {
		durations := make([]time.Duration, 11)
		formattedDurations := make([]string, 11)
		gameIds := make([]string, 11)
		totalTime := time.Duration(0)
		validTotalTime := true
		for i, questName := range s.anniversaryNamesInOrder {
			formattedTime := "N/A"
			pb, found := times[fmt.Sprintf("%s", questName)]
			if found {
				formattedTime = formatDurationSeconds(pb.Time)
				totalTime += pb.Time
				durations[i] = pb.Time
				gameIds[i] = fmt.Sprintf("%d", pb.Id)
			} else {
				totalTime += time.Minute * 20
				durations[i] = time.Minute * 20
				validTotalTime = false
			}
			formattedDurations[i] = formattedTime
		}
		formattedTotal := "N/A"
		if validTotalTime {
			formattedTotal = formatDurationSeconds(totalTime)
		}
		anniversaryTimes = append(anniversaryTimes, AnniversaryTimes{
			User:            user,
			Total:           formattedTotal,
			totalTime:       totalTime,
			Times:           formattedDurations,
			individualTimes: durations,
			GameIds:         gameIds,
			Colors:          make([]string, 11),
		})
	}

	sort.Slice(anniversaryTimes, func(i, j int) bool {
		return anniversaryTimes[i].totalTime < anniversaryTimes[j].totalTime
	})
	if len(anniversaryTimes) > 10 {
		anniversaryTimes = anniversaryTimes[0:10]
	}
	questBest := make(map[string]time.Duration)
	questWorst := make(map[string]time.Duration)
	for _, times := range anniversaryTimes {
		for i, questName := range s.anniversaryNamesInOrder {
			timeForQuest := times.individualTimes[i]
			if bestTime, found := questBest[questName]; !found || timeForQuest < bestTime {
				questBest[questName] = timeForQuest
			}
			if timeForQuest > questWorst[questName] {
				questWorst[questName] = timeForQuest
			}
		}
	}
	for _, times := range anniversaryTimes {
		for i, questName := range s.anniversaryNamesInOrder {
			worstTimeForQuest := questWorst[questName]
			bestTimeForQuest := questBest[questName]
			questTime := times.individualTimes[i]
			largestDifference := worstTimeForQuest - bestTimeForQuest
			totalDifference := 100 * ((largestDifference) - (worstTimeForQuest - questTime)) / largestDifference
			red := 150 + totalDifference
			green := 250 - totalDifference
			times.Colors[i] = fmt.Sprintf("rgba(%d,%d,120,0.1)", red, green)
		}
	}
	return anniversaryTimes
}

type AnniversaryTimes struct {
	User            string
	Total           string
	totalTime       time.Duration
	Times           []string
	individualTimes []time.Duration
	Colors          []string
	GameIds         []string
}

type QuestCounters struct {
	MesetaCharged      int64
	ClassMesetaCharged map[string]int64
	Runs               int64
	Moving             int64
	Standing           int64
	Attacking          int64
	Casting            int64
	Deaths             int64
	ClassUse           map[string]int64
	SID                map[string]int64
	Players            map[int]int64
	Shifta             map[int]int64
	RunsByDay          map[time.Time]int64
}

type RecordHistoryPoint struct {
	Time time.Time
	P1   *model.Game
	P2   *model.Game
	P3   *model.Game
	P4   *model.Game
}
