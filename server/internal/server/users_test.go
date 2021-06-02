package server_test

import (
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/server/internal/server"
	"log"
	"testing"
)

func Test_hashPassword(t *testing.T) {
	passwordIn := "test"
	log.Printf("%v - %v", passwordIn, server.HashPassword(passwordIn))
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