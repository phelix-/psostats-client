package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/phelix-/psostats/v2/pkg/common"
	"github.com/phelix-/psostats/v2/pkg/model"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	fileNames := []string {
		"./games_1886.json",
		"./games_2000.json",
		"./games_4000.json",
		"./games_6000.json",
		"./games_8000.json",
		"./games_10000.json",
		"./games_11220.json",
		"./games_11700.json",
	}
	sortedGames := make(map[string][]model.Game)
	playerPbs := make(map[string]map[string]time.Duration)
	playerDeaths := make(map[string]int)
	playerGames := make(map[string]int)
	classDeaths := make(map[string]int)
	classGames := make(map[string]int)
	questCounts := make(map[string]int)
	questDeaths := make(map[string]int)
	playerBestDeathlessStreak := make(map[string]int)
	playerCurrentDeathlessStreak := make(map[string]int)
	playerMesetaCharged := make(map[string]int64)
	classMesetaCharged := make(map[string]int64)
	for _,fileName := range fileNames {
		file, _ := os.Open(fileName)
		reader := bufio.NewReaderSize(file, 160000)
		lineNumber := 0
		for true {
			lineNumber++
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatalf("%v", err)
				}
			}

			game := model.Game{}
			err = json.Unmarshal(line, &game)
			if err != nil {
				log.Fatalf("%v", err)
			}
			otherGame, err := parseGameGzip(game.GameGzip)
			if err != nil {
				log.Fatalf("%v", err)
			}
			if !strings.HasPrefix(game.Quest, "Maximum Attack E") ||
				otherGame.PbCategory ||
				otherGame.Difficulty != "Ultimate" ||
				game.Quest == "Maximum Attack E: Spaceship" ||
				otherGame.Server != "ephinea" ||
				game.Quest == "Maximum Attack E: Episode 4" {
				continue
			}
			questName := fmt.Sprintf("%v %v", game.Quest, game.Category)
			gameList := sortedGames[questName]
			if gameList == nil {
				gameList = make([]model.Game, 0)
			}

			gameList = append(gameList, game)
			sortedGames[questName] = gameList

			possibleGames := [][]byte {game.P1Gzip, game.P2Gzip, game.P3Gzip, game.P4Gzip}
			for _,gameGzip := range possibleGames {
				if gameGzip == nil || len(gameGzip) == 0 {
					continue
				}
				subGame, err := parseGameGzip(gameGzip)
				if err != nil {
					log.Fatalf("%v", err)
				}
				userName := subGame.UserName
				pbForGc := playerPbs[userName]
				if pbForGc == nil {
					pbForGc = make(map[string]time.Duration)
				}
				playerPb := pbForGc[game.Quest]
				if playerPb < time.Nanosecond || playerPb > game.Time {
					pbForGc[game.Quest] = game.Time
				}

				playerPbs[userName] = pbForGc
				playerDeaths[userName] = playerDeaths[userName] + subGame.DeathCount
				classDeaths[subGame.PlayerClass] = classDeaths[subGame.PlayerClass] + subGame.DeathCount
				playerGames[userName] = playerGames[userName] + 1
				questDeaths[game.Quest] = questDeaths[game.Quest] + subGame.DeathCount
				questCounts[game.Quest] = questCounts[game.Quest] + 1
				mesetaCharged := subGame.MesetaCharged[len(subGame.MesetaCharged)-1]
				playerMesetaCharged[userName] = playerMesetaCharged[userName] + int64(mesetaCharged)
				classMesetaCharged[subGame.PlayerClass] = classMesetaCharged[subGame.PlayerClass] + int64(mesetaCharged)
				classGames[subGame.PlayerClass] = classGames[subGame.PlayerClass] + 1
				if strings.HasPrefix(subGame.PlayerClass, "HU") {
					if subGame.DeathCount == 0 {
						deathlessStreak := playerCurrentDeathlessStreak[userName]
						deathlessStreak++
						playerCurrentDeathlessStreak[userName] = deathlessStreak
						if deathlessStreak > playerBestDeathlessStreak[userName] {
							playerBestDeathlessStreak[userName] = deathlessStreak
						}
					} else {
						playerCurrentDeathlessStreak[userName] = 0
					}
				}
			}
		}
	}

	classesByName := make(map[string]common.PsoClass)
	for _, class := range common.GetAllClasses() {
		classesByName[class.Name] = class
	}
	for player,games := range playerGames {
		fmt.Printf("%v - %v\n", player, float32(playerDeaths[player]) / float32(games))
	}
	for quest,games := range questCounts {
		fmt.Printf("%v - %v\n", quest, float32(questDeaths[quest]) / float32(games))
	}
	fmt.Println("Deathless Streaks")
	for player,streak := range playerBestDeathlessStreak {
		fmt.Printf("%v - %v\n", player, streak)
	}
	fmt.Println("Meseta")
	for player,meseta := range playerMesetaCharged {
		fmt.Printf("%v - %v - %v\n", player, meseta, meseta / int64(playerGames[player]))
	}
	for class,meseta := range classMesetaCharged {
		fmt.Printf("%v - %v - %v\n", class, meseta, meseta / int64(classGames[class]))
	}
	for class,games := range classGames {
		fmt.Printf("%v - %v\n", class, float32(classDeaths[class]) / float32(games))
	}

	countByDay := make(map[string]map[int]int)
	overallClassCounts := make(map[string]int)
	totalGames := 0
	classCountsTotal := make(map[string]map[string]int)
	shiftaCountsTotal := make(map[string]map[int]int)
	for questName, quests := range sortedGames {
		recordTransitions := make([]TimeAndDuration, 0)

		previousRec := time.Hour
		numGames := len(quests)
		totalGames += numGames
		fmt.Printf("\n[SIZE=5][B]%v[/B][/SIZE] %v total games\n", questName, numGames)

		lastDuration := time.Second
		for _, quest := range quests {
			classCounts := classCountsTotal[quest.Quest]
			shiftaCounts := shiftaCountsTotal[quest.Quest]
			if classCounts == nil {
				classCounts = make(map[string]int)
				shiftaCounts = make(map[int]int)
			}
			day := quest.Timestamp.UTC().YearDay()
			countForQuest := countByDay[quest.Quest]
			if countForQuest == nil {
				countForQuest = make(map[int]int)
			}
			countForQuest[day] = countForQuest[day] + 1
			countByDay[quest.Quest] = countForQuest
			duration := quest.Time

			maxShifta := 0
			for i, player := range quest.PlayerNames {
				if len(player) == 0 {
					continue
				}
				class := quest.PlayerClasses[i]
				if classesByName[class].MaxShifta > maxShifta {
					maxShifta = classesByName[class].MaxShifta
				}
				classCounts[class]++
				overallClassCounts[class]++
			}
			shiftaCounts[maxShifta]++

			if duration < previousRec {
				previousRec = duration
				lastDuration = duration
				//fmt.Printf("\t %v %v %v %v\n", quest.PlayerNames[0], quest.PlayerNames[1], quest.PlayerNames[2], quest.PlayerNames[3])
				players := ""
				for i, player := range quest.PlayerNames {
					if len(player) == 0 {
						continue
					}
					if i > 0 {
						players += ", "
					}
					players += fmt.Sprintf("%v (%v)", player, quest.PlayerClasses[i])
				}
				fmt.Printf("\t[%d, %d],\n", quest.Timestamp.UTC().Unix() * 1000, duration)
				recordTransitions = append(recordTransitions, TimeAndDuration{
					quest.Id,
					quest.Timestamp,
					duration,
					formatDuration(duration.String()),
					players,
				})
			}
			classCountsTotal[quest.Quest] = classCounts
			shiftaCountsTotal[quest.Quest] = shiftaCounts
		}
		fmt.Printf("\t[new Date(\"09/04/2021\").getTime(), %d],\n", lastDuration)
		fmt.Print("---\n")
		for i,record := range recordTransitions {
			fmt.Printf("\"%v\",\n", record.Players)
			if i == len(recordTransitions) - 1 {
				fmt.Printf("\"%v\",\n", record.Players)
			}
		}
		sort.Slice(quests, func(a, b int) bool {
			return quests[a].Time < quests[b].Time
		})
		//jsonBytes, _ := json.MarshalIndent(recordTransitions, "", " ")
		//fmt.Print(string(jsonBytes))

		//fmt.Printf("Current record: [URL='https://psostats.com/game/%v']%v[/URL] %v(%v) %v(%v) %v(%v) %v(%v)\n", quests[0].Id, formatDuration(quests[0].Time), quests[0].PlayerNames[0], quests[0].PlayerClasses[0], quests[0].PlayerNames[1], quests[0].PlayerClasses[1], quests[0].PlayerNames[2], quests[0].PlayerClasses[2], quests[0].PlayerNames[3], quests[0].PlayerClasses[3])
		//fmt.Printf("#2: [URL='https://psostats.com/game/%v']%v[/URL] %v(%v) %v(%v) %v(%v) %v(%v)\n", quests[1].Id, formatDuration(quests[1].Time), quests[1].PlayerNames[0], quests[1].PlayerClasses[0], quests[1].PlayerNames[1], quests[1].PlayerClasses[1], quests[1].PlayerNames[2], quests[1].PlayerClasses[2], quests[1].PlayerNames[3], quests[1].PlayerClasses[3])
		//fmt.Printf("#3: [URL='https://psostats.com/game/%v']%v[/URL] %v(%v) %v(%v) %v(%v) %v(%v)\n", quests[2].Id, formatDuration(quests[2].Time), quests[2].PlayerNames[0], quests[2].PlayerClasses[0], quests[2].PlayerNames[1], quests[2].PlayerClasses[1], quests[2].PlayerNames[2], quests[2].PlayerClasses[2], quests[2].PlayerNames[3], quests[2].PlayerClasses[3])
		//fmt.Printf("#4: [URL='https://psostats.com/game/%v']%v[/URL] %v(%v) %v(%v) %v(%v) %v(%v)\n", quests[3].Id, formatDuration(quests[3].Time), quests[3].PlayerNames[0], quests[3].PlayerClasses[0], quests[3].PlayerNames[1], quests[3].PlayerClasses[1], quests[3].PlayerNames[2], quests[3].PlayerClasses[2], quests[3].PlayerNames[3], quests[3].PlayerClasses[3])
		//fmt.Printf("#5: [URL='https://psostats.com/game/%v']%v[/URL] %v(%v) %v(%v) %v(%v) %v(%v)\n", quests[4].Id, formatDuration(quests[4].Time), quests[4].PlayerNames[0], quests[4].PlayerClasses[0], quests[4].PlayerNames[1], quests[4].PlayerClasses[1], quests[4].PlayerNames[2], quests[4].PlayerClasses[2], quests[4].PlayerNames[3], quests[4].PlayerClasses[3])
		//fmt.Printf("#6: [URL='https://psostats.com/game/%v']%v[/URL] %v(%v) %v(%v) %v(%v) %v(%v)\n", quests[5].Id, formatDuration(quests[5].Time), quests[5].PlayerNames[0], quests[5].PlayerClasses[0], quests[5].PlayerNames[1], quests[5].PlayerClasses[1], quests[5].PlayerNames[2], quests[5].PlayerClasses[2], quests[5].PlayerNames[3], quests[5].PlayerClasses[3])

		firstQuarter := (numGames * 25) / 100
		median := (numGames * 50) / 100
		lastQuarter := median + firstQuarter
		fmt.Printf("\n75th Percentile: %v | Median: %v | 25th Percentile: %v\n", quests[firstQuarter].Time.Milliseconds(), quests[median].Time.Milliseconds(), quests[lastQuarter].Time.Milliseconds())


		//for i, rec := range recordTransitions {
		//	if i > 0 {
		//		rec = recordTransitions[i-1]
		//		fmt.Printf("%v\t%v\t%v\t%v\t%v\n", recordTransitions[i].Time.Add(-1*time.Minute).UTC().Format("2006-01-02 15:04"), formatDuration(rec.Duration.String()), formatDuration(quests[firstQuarter].QuestDuration), formatDuration(quests[median].QuestDuration), formatDuration(quests[lastQuarter].QuestDuration))
		//		rec = recordTransitions[i]
		//	}
		//	//log.Printf("%v %v %v %v", quests[5].PlayerNames[0], quests[5].PlayerNames[1], quests[5].PlayerNames[2], quests[5].PlayerNames[3])
		//}
		//rec := recordTransitions[len(recordTransitions)-1]
		//fmt.Printf("%v\t%v\t%v\t%v\t%v", time.Now().UTC().Format("2006-01-02 15:04"), formatDuration(rec.Duration.String()), formatDuration(quests[firstQuarter].QuestDuration), formatDuration(quests[median].QuestDuration), formatDuration(quests[lastQuarter].QuestDuration))
	}
	for quest,_ := range questCounts {
		fmt.Printf("%v", quest)
		fmt.Printf("Classes - %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d,\n",
			classCountsTotal[quest]["HUmar"],
			classCountsTotal[quest]["HUnewearl"],
			classCountsTotal[quest]["HUcast"],
			classCountsTotal[quest]["HUcaseal"],
			classCountsTotal[quest]["RAmar"],
			classCountsTotal[quest]["RAmarl"],
			classCountsTotal[quest]["RAcast"],
			classCountsTotal[quest]["RAcaseal"],
			classCountsTotal[quest]["FOmar"],
			classCountsTotal[quest]["FOmarl"],
			classCountsTotal[quest]["FOnewm"],
			classCountsTotal[quest]["FOnewearl"],
		)

		fmt.Printf("Shifta Used: %d, %d, %d, %d\n",
			shiftaCountsTotal[quest][3],
			shiftaCountsTotal[quest][15],
			shiftaCountsTotal[quest][20],
			shiftaCountsTotal[quest][30],
		)
	}

	for k,v := range countByDay {
		line := k + " - ["
		for i := 219; i < 248; i++ {
			line += fmt.Sprint(v[i]) + ", "
		}
		line += "]"
		fmt.Println(line)
	}

	sortedPlayerPbs := make([]PlayerAndTotal, 0)
	for player,v := range playerPbs {
		if len(v) == 11 {
			totalDuration := time.Duration(0)
			for _,duration := range v {
				totalDuration += duration
			}
			sortedPlayerPbs = append(sortedPlayerPbs, PlayerAndTotal{
				Player: player,
				Total:  totalDuration,
			})
			//fmt.Printf("%v - %v\n", player, totalDuration.String())
		}
	}

	sort.Slice(sortedPlayerPbs, func(a, b int) bool {
		return sortedPlayerPbs[a].Total < sortedPlayerPbs[b].Total
	})
	for _,a := range sortedPlayerPbs {
		fmt.Printf("%v - %v\n", a.Player, a.Total.String())
	}

	fmt.Printf("Total Games: %v\n", totalGames)
	fmt.Printf("HUmar\t%v(%.0f%%)\nHUnewearl\t%v(%.0f%%)\nHUcast\t%v(%.0f%%)\nHUcaseal\t%v(%.0f%%)\n",
		overallClassCounts["HUmar"], float32(100*overallClassCounts["HUmar"])/float32(4*totalGames),
		overallClassCounts["HUnewearl"], float32(100*overallClassCounts["HUnewearl"])/float32(4*totalGames),
		overallClassCounts["HUcast"], float32(100*overallClassCounts["HUcast"])/float32(4*totalGames),
		overallClassCounts["HUcaseal"], float32(100*overallClassCounts["HUcaseal"])/float32(4*totalGames),
	)
	fmt.Printf("RAmar\t%v(%.0f%%)\nRAmarl\t%v(%.0f%%)\nRAcast\t%v(%.0f%%)\nRAcaseal\t%v(%.0f%%)\n",
		overallClassCounts["RAmar"], float32(100*overallClassCounts["RAmar"])/float32(4*totalGames),
		overallClassCounts["RAmarl"], float32(100*overallClassCounts["RAmarl"])/float32(4*totalGames),
		overallClassCounts["RAcast"], float32(100*overallClassCounts["RAcast"])/float32(4*totalGames),
		overallClassCounts["RAcaseal"], float32(100*overallClassCounts["RAcaseal"])/float32(4*totalGames),
	)
	fmt.Printf("FOmar\t%v(%.0f%%)\nFOmarl\t%v(%.0f%%)\nFOnewm\t%v(%.0f%%)\nFOnewearl\t%v(%.0f%%)\n",
		overallClassCounts["FOmar"], float32(100*overallClassCounts["FOmar"])/float32(4*totalGames),
		overallClassCounts["FOmarl"], float32(100*overallClassCounts["FOmarl"])/float32(4*totalGames),
		overallClassCounts["FOnewm"], float32(100*overallClassCounts["FOnewm"])/float32(4*totalGames),
		overallClassCounts["FOnewearl"], float32(100*overallClassCounts["FOnewearl"])/float32(4*totalGames),
	)

}

type TimeAndDuration struct {
	Id       string
	Time     time.Time
	Duration time.Duration
	FormattedDuration string
	Players  string
}

type PlayerAndTotal struct {
	Player string
	Total time.Duration
}

func formatDuration(dString string) string {
	d, _ := time.ParseDuration(dString)
	d = d.Round(time.Millisecond)
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func parseGameGzip(gameBytes []byte) (*model.QuestRun, error) {
	questRun := model.QuestRun{}
	buffer := bytes.NewBuffer(gameBytes)
	reader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := io.ReadAll(reader)
	if err != io.ErrUnexpectedEOF {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &questRun)
	if err != nil {
		return nil, err
	}

	return &questRun, err
}