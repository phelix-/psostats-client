package server_test

import (
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/server/internal/server"
	"log"
	"testing"
	"time"
)

func Test_hashPassword(t *testing.T) {
	passwordIn := "test"
	log.Printf("%v - %v", passwordIn, server.HashPassword(passwordIn))
}

func Test_gamesMatch(t *testing.T) {
	a := model.QuestRun{GuildCard: "u1", UserName: "u1", QuestStartTime: time.Now()}
	b := model.QuestRun{GuildCard: "u2", UserName: "u2", QuestStartTime: time.Now()}

	if !server.GamesMatch(a, b) {
		t.Error("match")
	}
	b.QuestStartTime = time.Now().Add(45 * time.Second)
	if server.GamesMatch(a, b) {
		t.Error("match")
	}
	b.QuestStartTime = time.Now().Add(-45 * time.Second)
	if server.GamesMatch(a, b) {
		t.Error("match")
	}
}

func Test_cmodeRegex(t *testing.T) {
	questRun := model.QuestRun{QuestName: "1c3", Difficulty: "Normal", QuestComplete: true}
	if !server.IsLeaderboardCandidate(questRun) {
		t.Error("1c3 normal")
	}
	questRun = model.QuestRun{QuestName: "2c5", Difficulty: "Normal", QuestComplete: true}
	if !server.IsLeaderboardCandidate(questRun) {
		t.Error("2c5 normal")
	}
	questRun = model.QuestRun{QuestName: "ma1c", Difficulty: "Normal", QuestComplete: true}
	if server.IsLeaderboardCandidate(questRun) {
		t.Error("ma1c normal")
	}
	questRun = model.QuestRun{QuestName: "ma1c", Difficulty: "Ultimate", QuestComplete: true}
	if !server.IsLeaderboardCandidate(questRun) {
		t.Error("ma1c ult")
	}
}

func Test_sendWebhook(t *testing.T) {
	s := server.New(nil)
	duration := (1 * time.Minute) + (25 * time.Second) + (31 * time.Millisecond)
	questRun := model.QuestRun{QuestName: "Maximum Attack E: CCA", QuestDuration: duration.String(), Id: "1", AllPlayers: []model.BasePlayerInfo{
		{
			Name: "phelix",
			Class: "HUmar",
		},{
			Name: "shoebert",
			Class: "FOnewearl",
		},{
			Name: "Clio",
			Class: "RAmarl",
		},{
			Name: "Krang",
			Class: "HUcast",
		},
	}}
	//previousRecord := model.Game{
	//	Time: (1 * time.Minute) + (45 * time.Second),
	//}
	s.QuestRecordWebhook(questRun, nil)
}