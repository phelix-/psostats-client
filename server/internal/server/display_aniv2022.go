package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (s *Server) Anniv2022RecordsPage(c *fiber.Ctx) error {
	overallCounters, questCounters := s.getCounters()

	recordModel := struct {
		QuestNames     []string
		TopLaps        []AnniversaryTimes
		OverallCounter QuestCounters
		QuestCounters  map[string]QuestCounters
	}{
		QuestNames:     s.anniversaryNamesInOrder,
		TopLaps:        s.getTopLaps(),
		OverallCounter: overallCounters,
		QuestCounters:  questCounters,
	}

	s.anniversary2022Template = ensureParsed("./server/internal/templates/anniv2022.gohtml")
	err := s.anniversary2022Template.ExecuteTemplate(c.Response().BodyWriter(), "index", recordModel)
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
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
			if times.individualTimes[i] == questWorst[questName] {
				times.Colors[i] = "rgba(255,0,0,.3)"
			}
			if times.individualTimes[i] == questBest[questName] {
				times.Colors[i] = "rgba(0,255,0,.3)"
			}
		}
	}

	sort.Slice(anniversaryTimes, func(i, j int) bool {
		return anniversaryTimes[i].totalTime < anniversaryTimes[j].totalTime
	})
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
