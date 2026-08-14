// Backfill script: replays anniversary record/PB/stats writes for games that were
// submitted while anniv2026 record-writing was disabled (from the 2026 event's
// opening, 2026-08-12 06:28 PDT, until the fix that re-enabled writes was deployed).
// Regular leaderboard/game storage was never affected - only the anniv_*_2026 tables.
//
// Usage:
//
//	go run ./server/cmd/backfillanniv2026 -end=<deploy-timestamp-RFC3339> [-start=<RFC3339>] [-commit]
//
// -end is required. It MUST be set to the exact timestamp anniversary writes were
// redeployed to production. The stats counters this script writes use DynamoDB ADD
// (increment), which is NOT idempotent - any game the live server already recorded
// after that deploy must NOT also be processed here, or its stats get double-counted.
//
// Without -commit the script only logs what it would do (dry run).
package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"github.com/phelix-/psostats/v2/server/internal/server"
)

// Mirrors the anniversaryQuests set in server/internal/server/server.go.
var anniversaryQuests2026 = map[string]struct{}{
	"Maximum Attack E: Forest": {},
	"Maximum Attack E: Caves":  {},
	"Maximum Attack E: Mines":  {},
	"Maximum Attack E: Ruins":  {},
	"Maximum Attack E: Temple": {},
	"Maximum Attack E: Space":  {},
	"Maximum Attack E: CCA":    {},
	"Maximum Attack E: Seabed": {},
	"Maximum Attack E: Tower":  {},
	"Maximum Attack E: Crater": {},
	"Maximum Attack E: Desert": {},
	"August Atrocity #1":       {},
	"August Atrocity #2":       {},
}

const gamesPerMonthLimit = 5000

func main() {
	startFlag := flag.String("start", "2026-08-12T06:28:00-07:00", "RFC3339 start of the backfill window (inclusive)")
	endFlag := flag.String("end", "", "RFC3339 end of the backfill window (exclusive) - REQUIRED, set to the anniversary-writes deploy timestamp")
	commit := flag.Bool("commit", false, "actually write to DynamoDB (default is dry-run)")
	flag.Parse()

	if *endFlag == "" {
		log.Fatal("-end is required: set it to the timestamp anniversary record writes were redeployed, so this backfill doesn't overlap with games the live server already recorded")
	}
	start, err := time.Parse(time.RFC3339, *startFlag)
	if err != nil {
		log.Fatalf("invalid -start: %v", err)
	}
	end, err := time.Parse(time.RFC3339, *endFlag)
	if err != nil {
		log.Fatalf("invalid -end: %v", err)
	}
	if !end.After(start) {
		log.Fatal("-end must be after -start")
	}

	var awsSession *session.Session
	if _, set := os.LookupEnv("AWS_ACCESS_KEY_ID"); set {
		awsSession = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	} else {
		awsSession = session.Must(session.NewSession(&aws.Config{
			Region:   aws.String("us-west-2"),
			Endpoint: aws.String("http://localhost:8000"),
		}))
	}
	dynamoClient := dynamodb.New(awsSession)

	if *commit {
		log.Printf("COMMIT mode: writes will be made to DynamoDB")
	} else {
		log.Printf("DRY RUN: pass -commit to actually write")
	}
	log.Printf("backfill window: [%s, %s)", start.Format(time.RFC3339), end.Format(time.RFC3339))

	matched, written := 0, 0
	for _, month := range monthsBetween(start, end) {
		games, err := db.GetGamesForMonth(month, gamesPerMonthLimit, dynamoClient)
		if err != nil {
			log.Fatalf("get games for month %s: %v", month, err)
		}
		for _, game := range games {
			if game.Timestamp.Before(start) || !game.Timestamp.Before(end) {
				continue
			}
			if _, ok := anniversaryQuests2026[game.Quest]; !ok {
				continue
			}

			gem := firstStatsGem(game)
			if gem == -1 {
				log.Printf("skip %s (%s): no player POV stats attached", game.Id, game.Quest)
				continue
			}
			questRun, err := db.GetGame(game.Id, gem, dynamoClient)
			if err != nil || questRun == nil {
				log.Printf("skip %s: failed to load full game: %v", game.Id, err)
				continue
			}
			if questRun.Server != "ephinea" || questRun.PbCategory || !server.IsLeaderboardCandidate(*questRun) {
				continue
			}

			matched++
			if *commit {
				if err := backfillGame(*questRun, dynamoClient); err != nil {
					log.Printf("error backfilling %s (%s): %v", questRun.Id, questRun.QuestName, err)
					continue
				}
			}
			written++
			log.Printf("%s: %s %s %dp %s", questRun.QuestName, questRun.Id, questRun.UserName, len(questRun.AllPlayers), questRun.QuestDuration)
		}
	}

	if *commit {
		log.Printf("done: %d games matched, %d written", matched, written)
	} else {
		log.Printf("done: %d games matched, %d would be written (dry run, pass -commit to write)", matched, written)
	}
}

// backfillGame replays the same anniv2026 writes that updateAnniv2026Record performs
// for a live submission. Unlike the live path, it doesn't attach extra teammate POVs
// to an existing record - each game is processed once for the submitting player's POV.
// None of the anniversary quests are score-ranked (Endless-only), so this always
// compares by duration.
func backfillGame(questRun model.QuestRun, dynamoClient *dynamodb.DynamoDB) error {
	numPlayers := len(questRun.AllPlayers)
	questDuration, err := time.ParseDuration(questRun.QuestDuration)
	if err != nil {
		return err
	}

	topRun, err := db.GetAnniv2026Record(questRun.QuestName, numPlayers, questRun.PbCategory, dynamoClient)
	if err != nil {
		return err
	}
	if topRun == nil || questDuration < topRun.Time {
		if err := db.WriteAnniv2026Record(&questRun, dynamoClient); err != nil {
			return err
		}
	}

	pb, err := db.GetQuestSeriesPb("a2026", questRun.UserName, questRun.QuestName, dynamoClient)
	if err != nil {
		return err
	}
	if pb == nil || questDuration < pb.Time {
		if err := db.WriteQuestSeriesPb("a2026", &questRun, dynamoClient); err != nil {
			return err
		}
	}

	db.WriteAnniversaryStats(questRun, dynamoClient)
	return nil
}

func firstStatsGem(game model.Game) int {
	switch {
	case game.P1HasStats:
		return 0
	case game.P2HasStats:
		return 1
	case game.P3HasStats:
		return 2
	case game.P4HasStats:
		return 3
	}
	return -1
}

func monthsBetween(start, end time.Time) []string {
	months := make([]string, 0, 2)
	seen := make(map[string]bool)
	for t := start; !t.After(end); t = t.AddDate(0, 1, 0) {
		m := t.UTC().Format("01/2006")
		if !seen[m] {
			seen[m] = true
			months = append(months, m)
		}
	}
	if endMonth := end.UTC().Format("01/2006"); !seen[endMonth] {
		months = append(months, endMonth)
	}
	return months
}
