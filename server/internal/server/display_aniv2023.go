package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/pkg/psoclasses"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"log"
)

func (s *Server) Anniv2023RecordsPage(c *fiber.Ctx) error {
	overallCounters, questCounters := s.getCounters(2023)
	records, err := db.GetQuestRecords(db.Anniv2023RecordsTable, s.dynamoClient)
	if err != nil {
		log.Printf("get recent games %v", err)
		c.Status(500)
		return err
	}
	recordHistory, err := db.GetQuestRecords(db.Anniv2023RecordHistory, s.dynamoClient)
	sortedRecordHistory := s.sortRecordHistory(recordHistory)
	sortedRecs := sortAnnivGames(records)
	recordModel := struct {
		Year            int
		AnnivNumber     int
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
		Year:        2023,
		AnnivNumber: 8,
		QuestNames:  s.anniversaryNamesInOrder,
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
		TopLaps:        s.getTopLaps("a2023"),
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

func (s *Server) Anniv2025RecordsPage(c *fiber.Ctx) error {
	const year = 2025
	overallCounters, questCounters := s.getCounters(year)
	records, err := db.GetQuestRecords(db.Anniv2025RecordsTable, s.dynamoClient)
	if err != nil {
		log.Printf("get recent games %v", err)
		c.Status(500)
		return err
	}
	recordHistory, err := db.GetQuestRecords(db.Anniv2025RecordHistory, s.dynamoClient)
	sortedRecordHistory := s.sortRecordHistory(recordHistory)
	sortedRecs := sortAnnivGames(records)
	recordModel := struct {
		Year            int
		AnnivNumber     int
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
		Year:        year,
		AnnivNumber: 10,
		QuestNames:  s.anniversaryNamesInOrder,
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
		TopLaps:        s.getTopLaps("a2025"),
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
