// Seed script: populates quest_records and recent_games_by_player_2 for local testing.
// Usage: go run ./server/cmd/seed/
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/phelix-/psostats/v2/pkg/model"
	db "github.com/phelix-/psostats/v2/server/internal/db"
)

const quest = "Maximum Attack E: Caves"
const player = "testplayer"

func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	}))
	ddb := dynamodb.New(sess)

	if err := seedRecords(ddb); err != nil {
		log.Fatalf("seedRecords: %v", err)
	}
	if err := seedPlayerGames(ddb); err != nil {
		log.Fatalf("seedPlayerGames: %v", err)
	}
	log.Println("Seed complete.")
}

// seedRecords writes one record entry per category into quest_records.
func seedRecords(ddb *dynamodb.DynamoDB) error {
	type record struct {
		category   string
		numPlayers int
		pbRun      bool
		duration   time.Duration
		players    []string
		classes    []string
	}
	records := []record{
		{"1p", 1, true, 4*time.Minute + 2*time.Second + 300*time.Millisecond, []string{player, "", "", ""}, []string{"RAmar", "", "", ""}},
		{"1n", 1, false, 4*time.Minute + 15*time.Second + 100*time.Millisecond, []string{player, "", "", ""}, []string{"HUcast", "", "", ""}},
		{"2p", 2, true, 3*time.Minute + 55*time.Second + 200*time.Millisecond, []string{player, "partner1", "", ""}, []string{"RAmar", "HUcast", "", ""}},
		{"2n", 2, false, 4*time.Minute + 6*time.Second + 700*time.Millisecond, []string{player, "partner1", "", ""}, []string{"HUcast", "RAmar", "", ""}},
		{"3p", 3, true, 2*time.Minute + 55*time.Second + 700*time.Millisecond, []string{player, "partner1", "partner2", ""}, []string{"HUcast", "HUcast", "HUcast", ""}},
		{"3n", 3, false, 3*time.Minute + 10*time.Second + 400*time.Millisecond, []string{player, "partner1", "partner2", ""}, []string{"RAmar", "HUcast", "FOmar", ""}},
		{"4p", 4, true, 2*time.Minute + 42*time.Second + 800*time.Millisecond, []string{player, "partner1", "partner2", "partner3"}, []string{"HUcast", "HUcast", "FOmar", "RAcast"}},
		{"4n", 4, false, 2*time.Minute + 58*time.Second + 100*time.Millisecond, []string{player, "partner1", "partner2", "partner3"}, []string{"HUcast", "RAcast", "FOmar", "RAmar"}},
	}
	for i, r := range records {
		id := fmt.Sprintf("rec-%d", i+1)
		pNames := make([]string, 4)
		pClasses := make([]string, 4)
		copy(pNames, r.players)
		copy(pClasses, r.classes)
		g := model.Game{
			Id:            id,
			Player:        player,
			PlayerNames:   pNames,
			PlayerClasses: pClasses,
			PlayerGcs:     []string{"gc1", "gc2", "gc3", "gc4"},
			Category:      r.category,
			Quest:         quest,
			Episode:       2,
			Time:          r.duration,
			Timestamp:     time.Now().Add(-time.Duration(i+1) * 30 * 24 * time.Hour),
			P1HasStats:    true,
		}
		item, err := dynamodbattribute.MarshalMap(g)
		if err != nil {
			return err
		}
		_, err = ddb.PutItem(&dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String(db.QuestRecordsTable),
		})
		if err != nil {
			return err
		}
		log.Printf("  record %s inserted", r.category)
	}
	return nil
}

// seedPlayerGames writes player run history into recent_games_by_player_2.
func seedPlayerGames(ddb *dynamodb.DynamoDB) error {
	type spec struct {
		numPlayers int
		category   string
		count      int
		baseTime   time.Duration
		variance   time.Duration
	}
	classRotation := []string{"HUmar", "HUnewearl", "HUcast", "HUcaseal", "RAmar", "RAmarl", "RAcast", "RAcaseal", "FOmar", "FOmarl", "FOnewm", "FOnewearl"}
	partnerRotation := []string{"RAcast", "FOmar", "HUcaseal", "FOnewearl", "HUmar", "RAmarl", "FOnewm", "HUnewearl", "RAcaseal", "FOmarl", "HUcast", "RAmar"}
	specs := []spec{
		{1, "1p", 21, 4*time.Minute + 8*time.Second, 2 * time.Minute},
		{2, "2p", 21, 3*time.Minute + 40*time.Second, 90 * time.Second},
		{3, "3p", 21, 3*time.Minute + 0*time.Second, 75 * time.Second},
		{4, "4p", 50, 2*time.Minute + 50*time.Second, 60 * time.Second},
	}

	id := 100
	now := time.Now()

	for _, s := range specs {
		pNames := make([]string, 4)
		pClasses := make([]string, 4)
		pNames[0] = player
		for p := 1; p < s.numPlayers; p++ {
			pNames[p] = fmt.Sprintf("partner%d", p)
			pClasses[p] = partnerRotation[p%len(partnerRotation)]
		}

		for i := 0; i < s.count; i++ {
			pClasses[0] = classRotation[i%len(classRotation)]
			for p := 1; p < s.numPlayers; p++ {
				pClasses[p] = partnerRotation[(i+p)%len(partnerRotation)]
			}
			// Spread runs over the past year
			daysAgo := time.Duration(i*7) * 24 * time.Hour
			ts := now.Add(-daysAgo)

			// Vary the time: best run first in the list, vary by index
			runTime := s.baseTime + time.Duration(i)*s.variance/time.Duration(s.count)

			// alternate p/n across runs so both categories are represented
			pbRun := i%2 == 0
			cat := fmt.Sprintf("%dn", s.numPlayers)
			if pbRun {
				cat = fmt.Sprintf("%dp", s.numPlayers)
			}

			gameId := fmt.Sprintf("g-%d", id)

			g := model.Game{
				Id:            gameId,
				IdInt:         id,
				Player:        player,
				PlayerNames:   append([]string{}, pNames...),
				PlayerClasses: append([]string{}, pClasses...),
				PlayerGcs:     []string{"gc1", "gc2", "gc3", "gc4"},
				Category:      cat,
				Quest:         quest,
				Episode:       2,
				Time:          runTime,
				Timestamp:     ts,
				P1HasStats:    true,
				P2HasStats:    s.numPlayers >= 2,
				P3HasStats:    s.numPlayers >= 3,
				P4HasStats:    s.numPlayers >= 4,
			}

			item, err := dynamodbattribute.MarshalMap(g)
			if err != nil {
				return err
			}
			_, err = ddb.PutItem(&dynamodb.PutItemInput{
				Item:      item,
				TableName: aws.String(db.RecentGamesByPlayerTable),
			})
			if err != nil {
				return err
			}
			id++
		}
		log.Printf("  %dP: %d games inserted", s.numPlayers, s.count)
	}
	return nil
}
