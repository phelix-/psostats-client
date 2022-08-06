package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/server/internal/db"
	"sort"
	"strconv"
	"text/template"
	"time"
)

func (s *Server) GamePageV3(c *fiber.Ctx) error {
	gameId := c.Params("gameId")
	gem, err := strconv.Atoi(c.Params("gem"))
	if err != nil {
		gem = -1
	}
	fullGame, err := db.GetFullGame(gameId, s.dynamoClient)
	if err != nil {
		return err
	}
	game, err := db.GetGame(gameId, gem, s.dynamoClient)
	if err != nil {
		return err
	}
	playerDataFrames := make(map[int][]model.DataFrame)
	if fullGame.P1Gzip != nil {
		if dataFrames, err := db.GetDataFrames(gameId, 1, s.dynamoClient); err == nil {
			playerDataFrames[0] = dataFrames
		}
	}
	if fullGame.P2Gzip != nil {
		if dataFrames, err := db.GetDataFrames(gameId, 2, s.dynamoClient); err == nil {
			playerDataFrames[1] = dataFrames
		}
	}
	if fullGame.P3Gzip != nil {
		if dataFrames, err := db.GetDataFrames(gameId, 3, s.dynamoClient); err == nil {
			playerDataFrames[2] = dataFrames
		}
	}
	if fullGame.P4Gzip != nil {
		if dataFrames, err := db.GetDataFrames(gameId, 4, s.dynamoClient); err == nil {
			playerDataFrames[3] = dataFrames
		}
	}

	if game == nil {
		err = s.gameNotFoundTemplate.ExecuteTemplate(c.Response().BodyWriter(), "gameNotFound", nil)
	} else {
		duration, err := time.ParseDuration(game.QuestDuration)
		if err != nil {
			return err
		}

		invincibleRanges := make(map[int]int)
		invincibleStart := -1
		for i, invincible := range game.Invincible {
			if invincible {
				if invincibleStart < 0 {
					invincibleStart = i
				}
			} else {
				if invincibleStart > 0 {
					if i-invincibleStart >= 10 {
						invincibleRanges[invincibleStart] = i
					}
					invincibleStart = -1
				}
			}
		}
		for _, split := range game.Splits {
			if split.StartSecond > 1 {
				game.Events = append(game.Events, model.Event{
					Second:      split.StartSecond,
					Description: split.Name,
				})
			}
		}
		playerIndex := -1
		for i, player := range game.AllPlayers {
			if game.GuildCard == player.GuildCard {
				playerIndex = i
			}
		}
		timeMoving := game.TimeByState[2] + game.TimeByState[4]
		timeStanding := game.TimeByState[1]
		timeAttacking := game.TimeByState[5] + game.TimeByState[6] + game.TimeByState[7]
		timeCasting := game.TimeByState[8]
		totalActions := 0
		for _, weapon := range game.Weapons {
			actions := weapon.Attacks + weapon.Techs
			if actions > totalActions {
				totalActions = actions
			}
		}
		mostTime := 0
		timeByStateMap := make(map[string]TimeAndStateDisplay)
		for _, frame := range game.DataFrames {
			nameForState := getNameForState(frame.State)
			currentValue := timeByStateMap[nameForState.Display]
			nameForState.Time = 1 + currentValue.Time
			timeByStateMap[nameForState.Display] = nameForState
		}
		for _, state := range timeByStateMap {
			if state.Time > mostTime {
				mostTime = state.Time
			}
		}
		timeByState := make([]TimeAndStateDisplay, 0)
		for _, state := range timeByStateMap {
			percentTime := state.Time * 100 / mostTime
			timeByState = append(timeByState, TimeAndStateDisplay{
				Time:        state.Time,
				PercentTime: percentTime,
				PercentRest: 100 - percentTime,
				Display:     state.Display,
				Color:       state.Color,
			})
		}
		sort.Slice(timeByState, func(i, j int) bool {
			return timeByState[i].Time > timeByState[j].Time
		})

		weaponDisplay := make([]WeaponDisplay, 0)
		for _, weapon := range game.Weapons {
			if weapon.Attacks > 0 || weapon.Techs > 0 {
				attacks := weapon.Attacks * 100 / totalActions
				techs := weapon.Techs * 100 / totalActions
				rest := 100 - attacks - techs
				weaponDisplay = append(weaponDisplay, WeaponDisplay{
					Display: weapon.Display,
					Attacks: attacks,
					Techs:   techs,
					Rest:    rest,
				})
			}
		}
		sort.Slice(weaponDisplay, func(i, j int) bool {
			return weaponDisplay[i].Rest < weaponDisplay[j].Rest
		})

		model := struct {
			Game                 model.QuestRun
			SectionId            string
			HasPov               map[int]bool
			FormattedQuestTime   string
			InvincibleRanges     map[int]int
			HpRanges             map[int]uint16
			TpRanges             map[int]uint16
			PbRanges             map[int]int
			MonstersAliveRanges  map[int]int
			MonstersKilledRanges map[int]int
			MesetaChargedRanges  map[int]int
			FreezeTrapRanges     map[int]uint16
			ShiftaRanges         map[int]int16
			DebandRanges         map[int]int16
			HpPoolRanges         map[int]int
			Weapons              []model.Equipment
			Barriers             []model.Equipment
			Frames               []model.Equipment
			Units                []model.Equipment
			Mags                 []model.Equipment
			VideoUrl             string
			TimeMoving           string
			TimeStanding         string
			TimeAttacking        string
			TimeCasting          uint64
			FormattedTimeCasting string
			MapData              []MapData
			PlayerIndex          int
			TechsInOrder         [][]string
			MostActions          int
			TimeByState          []TimeAndStateDisplay
			PlayerDataFrames     map[int][]model.DataFrame
			SortedWeapons        []WeaponDisplay
		}{
			Game:      *game,
			SectionId: getSectionIdForQuest(game),
			HasPov: map[int]bool{
				0: fullGame.P1Gzip != nil,
				1: fullGame.P2Gzip != nil,
				2: fullGame.P3Gzip != nil,
				3: fullGame.P4Gzip != nil,
			},
			FormattedQuestTime:   formatDuration(duration),
			InvincibleRanges:     invincibleRanges,
			HpRanges:             convertU16ToXY(game.HP),
			TpRanges:             convertU16ToXY(game.TP),
			PbRanges:             convertF32ToXY(game.PB),
			MonstersAliveRanges:  convertIntToXY(game.MonsterCount),
			MonstersKilledRanges: convertIntToXY(game.MonstersKilledCount),
			MesetaChargedRanges:  convertIntToXY(game.MesetaCharged),
			FreezeTrapRanges:     convertU16ToXY(game.FreezeTraps),
			ShiftaRanges:         convertToXY(game.ShiftaLvl),
			DebandRanges:         convertToXY(game.DebandLvl),
			HpPoolRanges:         convertIntToXY(game.MonsterHpPool),
			Weapons:              getEquipment(game, model.EquipmentTypeWeapon),
			Barriers:             getEquipment(game, model.EquipmentTypeBarrier),
			Frames:               getEquipment(game, model.EquipmentTypeFrame),
			Units:                getEquipment(game, model.EquipmentTypeUnit),
			Mags:                 getEquipment(game, model.EquipmentTypeMag),
			TimeMoving:           formatDuration(time.Duration(timeMoving) * (time.Second / 30)),
			TimeStanding:         formatDuration(time.Duration(timeStanding) * (time.Second / 30)),
			TimeAttacking:        formatDuration(time.Duration(timeAttacking) * (time.Second / 30)),
			TimeCasting:          timeCasting,
			FormattedTimeCasting: formatDuration(time.Duration(timeCasting) * (time.Second / 30)),
			MapData:              formatMap(game, game.DataFrames),
			PlayerIndex:          playerIndex,
			TechsInOrder: [][]string{
				{"Foie", "Zonde", "Barta"},
				{"Gifoie", "Gizonde", "Gibarta"},
				{"Rafoie", "Razonde", "Rabarta"},
				{"Grants", "Megid"},
				{"Resta", "Anti", "Reverser"},
				{"Shifta", "Deband", "Ryuker"},
				{"Jellen", "Zalure"}},
			MostActions:      totalActions,
			PlayerDataFrames: playerDataFrames,
			TimeByState:      timeByState,
			SortedWeapons:    weaponDisplay,
		}
		funcMap := template.FuncMap{
			"add": func(a, b int) int { return a + b },
		}
		err = s.gameV3Template.Funcs(funcMap).ExecuteTemplate(c.Response().BodyWriter(), "game", model)
	}
	c.Response().Header.Set("Content-Type", "text/html; charset=UTF-8")
	return err
}

func getNameForState(state uint16) TimeAndStateDisplay {
	switch state {
	case 1:
		return TimeAndStateDisplay{Display: "Standing", Color: "rgba(255,255,255,.3)"}
	case 2:
		return TimeAndStateDisplay{Display: "Walking", Color: "rgba(122,255,122,0.3)"}
	case 4:
		return TimeAndStateDisplay{Display: "Running", Color: "rgba(122, 255, 122, 0.5)"}
	case 5:
		return TimeAndStateDisplay{Display: "Attacking", Color: "rgba(255, 122, 122, 0.5)"}
	case 6:
		return TimeAndStateDisplay{Display: "Attacking", Color: "rgba(255, 122, 122, 0.5)"}
	case 7:
		return TimeAndStateDisplay{Display: "Attacking", Color: "rgba(255, 122, 122, 0.5)"}
	case 8:
		return TimeAndStateDisplay{Display: "Casting", Color: "rgba(122, 122, 255, 0.5)"}
	case 10:
		return TimeAndStateDisplay{Display: "Recoil", Color: "rgba(255,255,255,1)"}
	case 14:
		return TimeAndStateDisplay{Display: "Knocked Down", Color: "rgba(255,255,255,1)"}
	case 15:
		return TimeAndStateDisplay{Display: "Dead", Color: "rgba(0,0,0,0.3)"}
	case 16:
		return TimeAndStateDisplay{Display: "Cutscene", Color: "rgba(255,255,0,.5)"}
	case 18:
		return TimeAndStateDisplay{Display: "Reviving", Color: "rgba(255,255,255,1)"}
	case 20:
		return TimeAndStateDisplay{Display: "Teleporting", Color: "rgba(255,255,255,1)"}
	case 23:
		return TimeAndStateDisplay{Display: "Emoting", Color: "rgba(255,255,255,1)"}
	}
	return TimeAndStateDisplay{Display: fmt.Sprintf("State %d", state), Color: "white"}
}

func getSectionIdForQuest(questRun *model.QuestRun) string {
	sectionId := questRun.AllPlayers[0].SectionId
	return getSectionId(int(sectionId))
}

func getSectionId(index int) string {
	sectionIdString := ""
	switch index {
	case 0:
		sectionIdString = "Viridia"
	case 1:
		sectionIdString = "Greenill"
	case 2:
		sectionIdString = "Skyly"
	case 3:
		sectionIdString = "Bluefull"
	case 4:
		sectionIdString = "Purplenum"
	case 5:
		sectionIdString = "Pinkal"
	case 6:
		sectionIdString = "Redria"
	case 7:
		sectionIdString = "Oran"
	case 8:
		sectionIdString = "Yellowboze"
	case 9:
		sectionIdString = "Whitill"
	}
	return sectionIdString
}

type TimeAndStateDisplay struct {
	Time        int
	PercentTime int
	PercentRest int
	Display     string
	Color       string
}

type WeaponDisplay struct {
	Display string
	Attacks int
	Techs   int
	Rest    int
}
