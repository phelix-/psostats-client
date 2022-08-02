package pso

import (
	"errors"
	"fmt"
	"github.com/phelix-/psostats/v2/client/internal/pso/constants"
	"log"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/phelix-/psostats/v2/pkg/model"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/client/internal/numbers"
	"github.com/phelix-/psostats/v2/client/internal/pso/inventory"
	"github.com/phelix-/psostats/v2/client/internal/pso/player"
	"github.com/phelix-/psostats/v2/client/internal/pso/quest"
)

const (
	basePlayerArrayAddress   = uintptr(0x00A94254)
	myPlayerIndexAddress     = 0x00A9C4F4
	bPDeRolLeData            = 0x00A43CC8
	monsterDeRolLeHP         = 0x6B4
	monsterDeRolLeHPMax      = 0x6B0
	monsterDeRolLeSkullHP    = 0x6B8
	monsterDeRolLeSkullHPMax = 0x20
	monsterDeRolLeShellHP    = 0x39C
	monsterDeRolLeShellHPMax = 0x1C

	bPBarbaRayData            = 0x00A43CC8
	monsterBarbaRayHP         = 0x704
	monsterBarbaRayHPMax      = 0x700
	monsterBarbaRaySkullHP    = 0x708
	monsterBarbaRaySkullHPMax = 0x20
	monsterBarbaRayShellHP    = 0x7AC
	monsterBarbaRayShellHPMax = 0x1C
)

type BaseGameInfo struct {
	episode      uint16
	difficulty   uint16
	currentMap   uint16
	mapVariation uint16
	currentFloor uint16
}

func (game *BaseGameInfo) DifficultyString() string {
	switch game.difficulty {
	case 0:
		return "Normal"
	case 1:
		return "Hard"
	case 2:
		return "Very Hard"
	case 3:
		return "Ultimate"
	default:
		return "unknown"
	}
}

type Monster struct {
	Name            string
	hp              uint16
	Id              uint16
	UnitxtId        uint32
	SpawnTime       time.Time
	KilledTime      time.Time
	Alive           bool
	Frame1          bool
	Index           int
	LastAttackerIdx uint16
	Location        model.Location
}

type Event struct {
	Second      int
	Description string
}

type QuestRun struct {
	Client                   model.ClientInfo
	Server                   string
	PlayerName               string
	PlayerClass              string
	GuildCard                string
	AllPlayers               []player.BasePlayerInfo
	Difficulty               string
	Episode                  uint16
	QuestName                string
	QuestComplete            bool
	QuestStartTime           time.Time
	QuestStartDate           string
	QuestEndTime             time.Time
	QuestDuration            string
	lastRecordedSecond       int
	lastRecordedHp           uint16
	lastFloor                uint16
	previousMeseta           int
	previousMesetaCharged    int
	DeathCount               int
	HP                       []uint16
	TP                       []uint16
	PB                       []float32
	MesetaCharged            []int
	Room                     []uint16
	maxPartySupplyableShifta int16
	maxPartyPbShifta         int16
	IllegalShifta            bool
	PbCategory               bool
	ShiftaLvl                []int16
	DebandLvl                []int16
	Invincible               []bool
	Events                   []Event
	FastWarps                bool
	Splits                   []model.QuestRunSplit
	Monsters                 map[int]Monster
	PlayerDamage             map[uint16]int64
	LastHits                 map[uint16]int
	Bosses                   map[string]model.BossData
	MonsterCount             []int
	MonstersKilledCount      []int
	MonsterHpPool            []int
	previousDamageDealt      int64
	MonstersDead             int
	Weapons                  map[string]model.Equipment
	equippedWeaponId         string
	FreezeTraps              []uint16
	previousFt               uint16
	previousDt               uint16
	previousCt               uint16
	FTUsed                   uint16
	DTUsed                   uint16
	CTUsed                   uint16
	previousTp               uint16
	TPUsed                   uint16
	previousState            uint16
	TimeByState              map[uint16]uint64
	TechsCast                map[string]int
	TimeStanding             uint64
	TimeMoving               uint64
	TimeAttacking            uint64
	TimeCasting              uint64
	Points                   uint16
	DataFrames               []model.DataFrame
}

func (pso *PSO) StartNewQuest(questConfig quest.Quest) {
	log.Printf("Starting new quest: %v", questConfig.Name)

	allPlayers, err := pso.getOtherPlayers()
	if err != nil {
		log.Panicf("unable to get all players %v", err)
	}
	questStartTime := time.Now()
	pso.GameState.QuestStartTime = questStartTime

	maxPartySupplyableShifta := int16(0)
	for _, p := range allPlayers {
		playerShifta := p.MaxSupplyableShifta()
		if playerShifta > maxPartySupplyableShifta {
			maxPartySupplyableShifta = playerShifta
		}
	}
	maxPartyPbShifta := int16(21 + (20 * (len(allPlayers) - 1)))
	if maxPartySupplyableShifta > maxPartyPbShifta {
		maxPartyPbShifta = maxPartySupplyableShifta
	}
	currentShifta := pso.CurrentPlayerData.ShiftaLvl
	pbShifta := currentShifta > maxPartySupplyableShifta
	pbCharged := pso.CurrentPlayerData.PB > 5.0
	pbCategory := pbShifta || pbCharged
	if !questConfig.TerminalQuest() {
		lowered := pso.CurrentPlayerData.IsLowered()
		shiftaCast := currentShifta > 0
		pbCategory = pbCategory || lowered || shiftaCast
	}

	pso.CurrentQuest = QuestRun{
		Server:                   pso.server,
		PlayerName:               pso.CurrentPlayerData.Name,
		PlayerClass:              pso.CurrentPlayerData.Class,
		GuildCard:                pso.CurrentPlayerData.GuildCard,
		AllPlayers:               allPlayers,
		Difficulty:               pso.GameState.Difficulty,
		Episode:                  pso.GameState.Episode,
		QuestStartTime:           questStartTime,
		QuestStartDate:           questStartTime.Format("15:04 01/02/2006"),
		QuestName:                questConfig.Name,
		lastRecordedSecond:       -1,
		previousMesetaCharged:    0,
		previousMeseta:           -1,
		DeathCount:               0,
		HP:                       make([]uint16, 0),
		TP:                       make([]uint16, 0),
		PB:                       make([]float32, 0),
		MesetaCharged:            make([]int, 0),
		Room:                     make([]uint16, 0),
		maxPartySupplyableShifta: maxPartySupplyableShifta,
		maxPartyPbShifta:         maxPartyPbShifta,
		ShiftaLvl:                make([]int16, 0),
		DebandLvl:                make([]int16, 0),
		IllegalShifta:            false,
		PbCategory:               pbCategory,
		Invincible:               make([]bool, 0),
		Events:                   make([]Event, 0),
		Splits:                   make([]model.QuestRunSplit, len(questConfig.Splits)),
		Monsters:                 make(map[int]Monster),
		Bosses:                   make(map[string]model.BossData),
		LastHits:                 make(map[uint16]int),
		PlayerDamage:             make(map[uint16]int64),
		MonsterCount:             make([]int, 0),
		MonstersKilledCount:      make([]int, 0),
		MonsterHpPool:            make([]int, 0),
		Weapons:                  make(map[string]model.Equipment),
		equippedWeaponId:         "",
		FreezeTraps:              make([]uint16, 0),
		previousDt:               0,
		previousFt:               0,
		previousCt:               0,
		FTUsed:                   0,
		DTUsed:                   0,
		CTUsed:                   0,
		previousTp:               0,
		TPUsed:                   0,
		TimeByState:              make(map[uint16]uint64),
		TechsCast:                make(map[string]int),
		DataFrames:               make([]model.DataFrame, 0),
	}
	pso.startedGame <- pso.CurrentQuest
	pso.GameState.QuestStarted = true
}

func (pso *PSO) consolidateFrame(monsters []Monster) {

	currentQuestRun := pso.CurrentQuest
	currentSecond := int(time.Now().Sub(currentQuestRun.QuestStartTime).Seconds())
	if currentQuestRun.QuestComplete {
		return
	}

	if currentQuestRun.lastRecordedSecond < currentSecond {
		currentQuestRun.lastRecordedSecond = currentSecond
		mesetaCharged := currentQuestRun.previousMesetaCharged
		previousMeseta := currentQuestRun.previousMeseta
		if currentQuestRun.previousCt > pso.CurrentPlayerData.ConfuseTraps {
			currentQuestRun.CTUsed += currentQuestRun.previousCt - pso.CurrentPlayerData.ConfuseTraps
		}
		currentQuestRun.previousCt = pso.CurrentPlayerData.ConfuseTraps

		if currentQuestRun.previousFt > pso.CurrentPlayerData.FreezeTraps {
			currentQuestRun.FTUsed += currentQuestRun.previousFt - pso.CurrentPlayerData.FreezeTraps
		}
		currentQuestRun.previousFt = pso.CurrentPlayerData.FreezeTraps

		if currentQuestRun.previousDt > pso.CurrentPlayerData.DamageTraps {
			currentQuestRun.DTUsed += currentQuestRun.previousDt - pso.CurrentPlayerData.DamageTraps
		}
		currentQuestRun.previousDt = pso.CurrentPlayerData.DamageTraps

		if currentQuestRun.previousTp > pso.CurrentPlayerData.TP {
			currentQuestRun.TPUsed += currentQuestRun.previousTp - pso.CurrentPlayerData.TP
		}
		currentQuestRun.previousTp = pso.CurrentPlayerData.TP

		if previousMeseta != -1 {
			mesetaDifference := previousMeseta - int(pso.CurrentPlayerData.Meseta)
			if mesetaDifference > 0 {
				// negative means meseta picked up, ignoring I guess
				mesetaCharged = mesetaDifference + currentQuestRun.previousMesetaCharged
			}
		}
		currentQuestRun.previousMeseta = int(pso.CurrentPlayerData.Meseta)
		currentQuestRun.previousMesetaCharged = mesetaCharged
		currentQuestRun.HP = append(currentQuestRun.HP, pso.CurrentPlayerData.HP)
		currentQuestRun.TP = append(currentQuestRun.TP, pso.CurrentPlayerData.TP)
		currentQuestRun.PB = append(currentQuestRun.PB, pso.CurrentPlayerData.PB)
		currentQuestRun.Room = append(currentQuestRun.Room, pso.CurrentPlayerData.Room)
		currentQuestRun.ShiftaLvl = append(currentQuestRun.ShiftaLvl, pso.CurrentPlayerData.ShiftaLvl)
		currentQuestRun.DebandLvl = append(currentQuestRun.DebandLvl, pso.CurrentPlayerData.DebandLvl)
		currentQuestRun.MonsterCount = append(currentQuestRun.MonsterCount, pso.GameState.MonsterCount)
		currentQuestRun.MesetaCharged = append(currentQuestRun.MesetaCharged, mesetaCharged)
		currentQuestRun.MonstersKilledCount = append(currentQuestRun.MonstersKilledCount, currentQuestRun.MonstersDead)
		currentQuestRun.FreezeTraps = append(currentQuestRun.FreezeTraps, pso.CurrentPlayerData.FreezeTraps)
		currentQuestRun.Invincible = append(currentQuestRun.Invincible, pso.CurrentPlayerData.InvincibilityFrames > 0)
		damageDealt := currentQuestRun.PlayerDamage[uint16(pso.CurrentPlayerIndex)]
		currentQuestRun.previousDamageDealt = damageDealt
		weaponFound := false
		for _, equipment := range pso.Inventory.Equipment {
			storedEquipment, exists := currentQuestRun.Weapons[equipment.Id]
			if !exists {
				storedEquipment = model.Equipment{
					Id:              equipment.Id,
					Type:            equipment.Type,
					Display:         equipment.Display,
					SecondsEquipped: 0,
				}
			}
			if equipment.Type == model.EquipmentTypeWeapon {
				weaponFound = true
				currentQuestRun.equippedWeaponId = equipment.Id
			}
			storedEquipment.SecondsEquipped += 1
			currentQuestRun.Weapons[equipment.Id] = storedEquipment
		}
		if !weaponFound {
			storedEquipment, exists := currentQuestRun.Weapons[model.WeaponBareHanded]
			if !exists {
				storedEquipment = model.Equipment{
					Id:              model.WeaponBareHanded,
					Type:            model.EquipmentTypeWeapon,
					Display:         model.WeaponBareHanded,
					SecondsEquipped: 0,
					Attacks:         0,
					Techs:           0,
				}
			}
			storedEquipment.SecondsEquipped += 1
			currentQuestRun.equippedWeaponId = model.WeaponBareHanded
			currentQuestRun.Weapons[model.WeaponBareHanded] = storedEquipment
		}
		if pso.CurrentPlayerData.ShiftaLvl > currentQuestRun.maxPartyPbShifta {
			currentQuestRun.IllegalShifta = true
		}
		dataFrame := model.DataFrame{
			Time:               time.Now().Unix(),
			HP:                 pso.CurrentPlayerData.HP,
			TP:                 pso.CurrentPlayerData.TP,
			PB:                 pso.CurrentPlayerData.PB,
			MesetaCharged:      mesetaCharged,
			ShiftaLvl:          pso.CurrentPlayerData.ShiftaLvl,
			DebandLvl:          pso.CurrentPlayerData.DebandLvl,
			Invincible:         pso.CurrentPlayerData.InvincibilityFrames > 0,
			Map:                pso.GameState.Map,
			MapVariation:       pso.GameState.MapVariation,
			FT:                 pso.CurrentPlayerData.FreezeTraps,
			DT:                 pso.CurrentPlayerData.DamageTraps,
			CT:                 pso.CurrentPlayerData.ConfuseTraps,
			DamageDealt:        damageDealt,
			State:              pso.CurrentPlayerData.ActionState,
			Weapon:             currentQuestRun.equippedWeaponId,
			Kills:              pso.CurrentQuest.LastHits[uint16(pso.CurrentPlayerIndex)],
			PlayerByGcLocation: make(map[string]model.Location),
			MonsterLocation:    make(map[int]model.Location),
		}
		for _, monster := range monsters {
			if monster.hp > 0 {
				dataFrame.MonsterLocation[int(monster.Id)] = monster.Location
				dataFrame.MonstersAlive++
			}
		}
		if players, err := pso.getOtherPlayers(); err == nil {
			playerByGcLocation := dataFrame.PlayerByGcLocation
			for _, player := range players {
				if player.Warping && pso.ephineaFastBurstEnabled() {
					currentQuestRun.FastWarps = true
				}
				playerByGcLocation[player.GuildCard] = player.Location
			}
			dataFrame.PlayerByGcLocation = playerByGcLocation
		}
		currentQuestRun.DataFrames = append(currentQuestRun.DataFrames, dataFrame)
	}

	currentState := pso.CurrentPlayerData.ActionState
	currentQuestRun.TimeByState[pso.CurrentPlayerData.ActionState] = currentQuestRun.TimeByState[pso.CurrentPlayerData.ActionState] + 1
	currentWeapon, found := currentQuestRun.Weapons[currentQuestRun.equippedWeaponId]
	if isAttackState(currentState) {
		currentQuestRun.TimeAttacking++
		if found && !isAttackState(currentQuestRun.previousState) {
			currentWeapon.Attacks = currentWeapon.Attacks + 1
			currentQuestRun.Weapons[currentQuestRun.equippedWeaponId] = currentWeapon
		}
	} else if currentState == 2 || currentState == 4 {
		currentQuestRun.TimeMoving++
	} else if currentState == 8 {
		currentQuestRun.TimeCasting++
		if currentQuestRun.previousState != 8 {
			tech := pso.CurrentPlayerData.GetCurrentTech()
			currentQuestRun.TechsCast[tech] = currentQuestRun.TechsCast[tech] + 1
			if found {
				currentWeapon.Techs = currentWeapon.Techs + 1
				currentQuestRun.Weapons[currentQuestRun.equippedWeaponId] = currentWeapon
			}
		}
	} else if currentState == 1 {
		currentQuestRun.TimeStanding++
	}
	currentQuestRun.previousState = currentState

	if players, err := pso.getOtherPlayers(); err == nil {
		for _, player := range players {
			if player.Warping && pso.ephineaFastBurstEnabled() {
				currentQuestRun.FastWarps = true
			}
		}
	}

	if currentQuestRun.lastFloor != pso.GameState.Floor {
		currentQuestRun.Events = append(currentQuestRun.Events, Event{
			Second:      currentSecond,
			Description: pso.GetFloorName(),
		})
		currentQuestRun.lastFloor = pso.GameState.Floor
	}

	if pso.CurrentPlayerData.HP == 0 && currentQuestRun.lastRecordedHp != 0 {
		currentQuestRun.DeathCount++
		currentQuestRun.Events = append(currentQuestRun.Events, Event{
			Second:      currentSecond,
			Description: "Died",
		})
	}
	currentQuestRun.lastRecordedHp = pso.CurrentPlayerData.HP

	if pso.GameState.QuestComplete {
		currentQuestRun.QuestComplete = true
		currentQuestRun.QuestEndTime = pso.GameState.QuestEndTime
		currentQuestRun.QuestDuration = pso.GameState.QuestEndTime.Sub(currentQuestRun.QuestStartTime).String()
		pso.completeGame <- currentQuestRun
	} else {
		currentQuestRun.QuestDuration = time.Now().Sub(currentQuestRun.QuestStartTime).String()
	}
	pso.CurrentQuest = currentQuestRun
}

func (pso *PSO) updateCurrentSplit(questConfig quest.Quest) {
	currentQuestRun := pso.CurrentQuest
	currentSplit := pso.GameState.CurrentSplit
	if currentSplit.Index < len(questConfig.Splits) {
		currentSplitCfg := questConfig.Splits[pso.GameState.CurrentSplit.Index]
		switched, err := pso.getFloorSwitch(currentSplitCfg.Trigger.Switch, currentSplitCfg.Trigger.Floor)
		if err == nil && switched {
			currentSplit.End = time.Now()
			currentQuestRun.Splits[currentSplit.Index] = currentSplit
			currentSplit = model.QuestRunSplit{Index: currentSplit.Index + 1}
			pso.GameState.CurrentSplit = currentSplit
		}
	}
	if currentSplit.Start.IsZero() && currentSplit.Index < len(questConfig.Splits) {
		currentSplit.Start = time.Now()
		currentSplit.StartSecond = int(time.Now().Sub(currentQuestRun.QuestStartTime).Seconds())
		currentSplit.Name = questConfig.Splits[pso.GameState.CurrentSplit.Index].Name
		currentQuestRun.Splits[currentSplit.Index] = currentSplit
		pso.GameState.CurrentSplit = currentSplit
	}
}

func (pso *PSO) addExtraQuestInfo(questConfig quest.Quest) {
	if questConfig.Name == "Endless: Episode 1" {
		points := quest.GetRegisterValue(pso.handle, 51, pso.GameState.questRegisterPointer)
		if points > 0 {
			pso.CurrentQuest.Points = points
		}
	}
}

func (pso *PSO) consolidateMonsterState(monsters []Monster) {
	now := time.Now()
	currentQuestRun := pso.CurrentQuest
	monsterHpPool := 0
	recordThisSecond := len(currentQuestRun.MonsterHpPool)-1 < currentQuestRun.lastRecordedSecond
	for _, monster := range monsters {
		monsterId := int(monster.Id)
		existingMonster, exists := currentQuestRun.Monsters[monsterId]
		monsterHpPool += int(monster.hp)
		if !exists {
			monster.SpawnTime = now
			monster.Alive = true
			currentQuestRun.Monsters[monsterId] = monster
			existingMonster = monster
		} else if existingMonster.Alive && monster.hp <= 0 {
			// We don't allow frame 0 kills because some monsters appear to spawn in with 0 hp.
			// This could be a synchronization issue w/ pso (data is still initializing when we catch it?)
			existingMonster.Alive = false
			currentQuestRun.MonstersDead += 1
			existingMonster.KilledTime = now
			if existingMonster.UnitxtId != 34 &&
				existingMonster.UnitxtId != 45 &&
				existingMonster.UnitxtId != 73 && // barba ray
				existingMonster.UnitxtId != 68 { // recon
				// Excluding DRL and Dark Gunners because they're buggy
				playerName := ""
				if int(monster.LastAttackerIdx) < len(currentQuestRun.AllPlayers) {
					player := currentQuestRun.AllPlayers[monster.LastAttackerIdx]
					playerName = player.Name
				}
				existingMonster.Frame1 = existingMonster.KilledTime.Sub(existingMonster.SpawnTime).Milliseconds() < 60
				if existingMonster.Frame1 {
					log.Printf("frame1? %v(%v) %v - %s", existingMonster.Name, existingMonster.UnitxtId, existingMonster.Id, playerName)
				}
			}
			currentQuestRun.LastHits[monster.LastAttackerIdx] = currentQuestRun.LastHits[monster.LastAttackerIdx] + 1
			currentQuestRun.PlayerDamage[monster.LastAttackerIdx] = currentQuestRun.PlayerDamage[monster.LastAttackerIdx] + int64(existingMonster.hp)
			currentQuestRun.Monsters[monsterId] = existingMonster
		} else if existingMonster.Alive {
			if monster.hp < existingMonster.hp {
				hpLost := int64(existingMonster.hp - monster.hp)
				currentQuestRun.PlayerDamage[monster.LastAttackerIdx] = currentQuestRun.PlayerDamage[monster.LastAttackerIdx] + hpLost
			}
			existingMonster.hp = monster.hp
			currentQuestRun.Monsters[monsterId] = existingMonster
		}
		if isBoss, bossName := isBoss(existingMonster); isBoss {
			idString := fmt.Sprintf("%v", existingMonster.Id)
			if !exists {
				form := 0
				for _, boss := range currentQuestRun.Bosses {
					if boss.UnitxtId == existingMonster.UnitxtId && boss.Id != existingMonster.Id {
						form++
					}
				}
				if form > 0 {
					bossName = fmt.Sprintf("%v (%d)", bossName, form)
				}
				currentQuestRun.Bosses[idString] = model.BossData{
					Name:       bossName,
					Id:         existingMonster.Id,
					UnitxtId:   existingMonster.UnitxtId,
					SpawnTime:  now,
					FirstFrame: currentQuestRun.lastRecordedSecond,
					Hp:         make([]int, 0),
				}
			}
			boss := currentQuestRun.Bosses[idString]
			if recordThisSecond {
				boss.Hp = append(boss.Hp, int(monster.hp))
			}
			boss.KilledTime = existingMonster.KilledTime
			currentQuestRun.Bosses[idString] = boss
		}
	}
	if recordThisSecond {
		currentQuestRun.MonsterHpPool = append(currentQuestRun.MonsterHpPool, monsterHpPool)
	}
	pso.CurrentQuest = currentQuestRun
}

func isBoss(monster Monster) (bool, string) {
	if monster.UnitxtId == 44 {
		return true, "Sil Dragon"
	}
	if monster.UnitxtId == 45 && monster.Index == 0 {
		return true, "Dal Ra Lie"
	}
	if monster.UnitxtId == 46 && monster.Index == 31 {
		return true, "Vol Opt ver. 2 (1)"
	}
	if monster.UnitxtId == 46 && monster.Index == 32 {
		return true, "Vol Opt ver. 2 (2)"
	}
	if monster.UnitxtId == 47 {
		return true, "Dark Falz"
	}
	if monster.UnitxtId == 73 && monster.Index == 0 {
		return true, "Barba Ray"
	}
	if monster.UnitxtId == 76 {
		return true, "Gol Dragon"
	}
	if monster.UnitxtId == 77 {
		return true, "Gal Gryphon"
	}
	if monster.UnitxtId == 78 {
		return true, "Olga Flow"
	}
	if monster.UnitxtId == 106 {
		if monster.Index < 5 {
			return true, fmt.Sprintf("Saint-Million Tail (%v)", monster.Index)
		} else if monster.Index < 9 {
			return true, fmt.Sprintf("Saint-Million Head (%v)", monster.Index-4)
		}
	}
	if monster.UnitxtId == 107 {
		if monster.Index < 5 {
			return true, fmt.Sprintf("Shambertin Tail (%v)", monster.Index)
		} else if monster.Index < 9 {
			return true, fmt.Sprintf("Shambertin Head (%v)", monster.Index-4)
		}
	}
	if monster.UnitxtId == 108 {
		if monster.Index < 5 {
			return true, fmt.Sprintf("Kondrieu Tail (%v)", monster.Index)
		} else if monster.Index < 9 {
			return true, fmt.Sprintf("Kondrieu Head (%v)", monster.Index-4)
		}
	}
	return false, ""
}

func isAttackState(state uint16) bool {
	return state == 5 || state == 6 || state == 7
}

func (pso *PSO) RefreshData() error {
	if !pso.connected {
		pso.GameState.Clear()
		log.Fatal("RefreshData: connection to window lost")
	}

	index, err := pso.getMyPlayerIndex()
	if err != nil {
		log.Fatal("Unable to find player index")
		return err
	}
	pso.CurrentPlayerIndex = index

	address := pso.getBaseCharacterAddress(index)
	game, err := pso.getBaseGameInfo()
	if err != nil {
		return err
	}
	pso.GameState.Episode = game.episode
	pso.GameState.Map = game.currentMap
	pso.GameState.MapVariation = game.mapVariation
	pso.GameState.Floor = game.currentFloor
	pso.GameState.Difficulty = game.DifficultyString()

	if address != 0 {
		playerData, err := player.GetPlayerData(pso.handle, address, pso.server)
		if err != nil {
			return err
		}
		pso.CurrentPlayerData = playerData

		inventory, err := inventory.ReadInventory(pso.handle, index)
		if err != nil {
			return err
		}
		pso.Inventory = inventory

		monsters, err := pso.GetMonsterList()
		if err != nil {
			return err
		}

		questPtr := quest.GetQuestPointer(pso.handle)
		if questPtr != 0 {
			if questPtr != pso.GameState.questPointer {
				pso.GameState.questRegisterPointer = quest.GetQuestRegisterPointer(pso.handle, questPtr)
				pso.GameState.questPointer = questPtr
			}
			questStartConditionsMet := false
			questDataPtr := quest.GetQuestDataPointer(pso.handle, questPtr)

			questName, err := pso.getQuestName(questDataPtr)
			if err != nil {
				return err
			}
			questNumber := pso.getQuestNumber(questDataPtr)
			questConfig, exists := pso.questTypes.GetQuestConfig(questNumber, int(pso.GameState.Episode), questName)
			if exists {
				questName = questConfig.Name
			} else {
				questName = fmt.Sprintf("%v (Missing Config)", questName)
			}
			pso.GameState.QuestName = questName
			pso.GameState.CmodeStage = questConfig.GetCmodeStage()

			if !pso.GameState.QuestStarted {
				if exists && !questConfig.Ignore {
					quest.GetQuestRegisterPointer(pso.handle, questPtr)
					questStartConditionsMet, err = pso.checkQuestStartConditions(questConfig)
					if err != nil {
						return err
					}
				}
				if questStartConditionsMet && pso.GameState.AllowQuestStart {
					rngSeed := pso.getRngSeed()
					pso.GameState.RngSeed = rngSeed
					pso.StartNewQuest(questConfig)
				}
			} else if !pso.GameState.QuestComplete {
				if exists {
					questEndConditionsMet, err := pso.checkQuestEndConditions(questConfig)
					if err != nil {
						return err
					}
					if questEndConditionsMet {
						pso.GameState.QuestComplete = true
						pso.GameState.QuestEndTime = time.Now()
					} else {
						if pso.GameState.CmodeStage > 0 && pso.GameState.Floor == 0 {
							// Back to pioneer2, cmode failed
							pso.GameState.ClearQuest()
						} else {
							rngSeed := pso.getRngSeed()
							if rngSeed != pso.GameState.RngSeed {
								// unseen quest reset
								pso.GameState.ClearQuest()
							}
						}
					}
				}
			}
			if pso.GameState.QuestStarted {
				pso.consolidateFrame(monsters)
				pso.updateCurrentSplit(questConfig)
				pso.addExtraQuestInfo(questConfig)
				pso.consolidateMonsterState(monsters)
			}
		} else {
			pso.GameState.AllowQuestStart = true
			pso.GameState.ClearQuest()
		}
	} else {
		pso.GameState.Clear()
	}

	return nil
}

func (pso *PSO) getMyPlayerIndex() (uint8, error) {
	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(myPlayerIndexAddress), 4)
	if !ok {
		return 0, errors.New("unable to find player index")
	}
	index := uint8(buf[0])
	return index, nil
}

func (pso *PSO) getBaseCharacterAddress(index uint8) uintptr {
	address := basePlayerArrayAddress + (4 * uintptr(index))
	return uintptr(numbers.ReadU32Unchecked(pso.handle, address))
}

func (pso *PSO) getMonsterUnitxtAddr() (uintptr, error) {
	unitxtAddr, err := numbers.ReadU32(pso.handle, uintptr(0x00a9cd50))
	if err != nil {
		return 0, err
	}
	monsterUnitxtAddr := uint32(0)
	if unitxtAddr != 0 {
		monsterUnitxtAddr, err = numbers.ReadU32(pso.handle, uintptr(unitxtAddr+16))
	}
	return uintptr(monsterUnitxtAddr), err
}

func (pso *PSO) getMonsterName(monsterId uint32) (string, error) {
	if monsterName, exists := pso.MonsterNames[monsterId]; exists {
		return monsterName, nil
	}
	monsterUnitxtAddr, err := pso.getMonsterUnitxtAddr()
	if err != nil {
		return "", err
	} else if monsterUnitxtAddr == 0 {
		return "", errors.New("monsterUnitxtAddr is unset")
	}
	monsterNameAddr, err := numbers.ReadU32(pso.handle, monsterUnitxtAddr+uintptr(4*monsterId))
	if err != nil {
		return "", err
	}
	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(monsterNameAddr), 32)
	if !ok {
		return "", errors.New("unable to getMonsterName")
	}

	endIndex := len(buf)
	for index, b := range buf {
		if b == 0x00 {
			endIndex = index
			break
		}
	}
	something := utf16.Decode(buf[0:endIndex])
	pso.MonsterNames[monsterId] = string(something)
	return string(something), nil
}

func (pso *PSO) getOtherPlayers() ([]player.BasePlayerInfo, error) {
	players := make([]player.BasePlayerInfo, 0)
	for i := 0; i < 12; i++ {
		address := pso.getBaseCharacterAddress(uint8(i))
		if address != 0 {
			playerData, err := player.GetPlayerData(pso.handle, address, pso.server)
			if err != nil {
				return nil, err
			}
			players = append(players, playerData)
		}
	}
	return players, nil
}

func (pso *PSO) getBaseGameInfo() (BaseGameInfo, error) {
	base := 0x00A9B1C8
	max := 0x00A9CD68
	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(base), uintptr((max-base)+2))
	if !ok {
		return BaseGameInfo{}, errors.New("unable to getDifficulty")
	}

	difficulty := buf[(0x00A9CD68-base)/2]
	episode := buf[(0x00A9B1C8-base)/2] + 1
	if episode == 3 {
		episode = 4
	}
	currentMap := numbers.ReadU16(pso.handle, uintptr(0x00AAFC9C))
	currentFloor := numbers.ReadU16(pso.handle, uintptr(0x00AAFCA0))
	mapVariation := numbers.ReadU16(pso.handle, uintptr(0x00AAFC98))
	game := BaseGameInfo{
		episode:      episode,
		difficulty:   difficulty,
		currentMap:   currentMap,
		mapVariation: mapVariation,
		currentFloor: currentFloor,
	}
	return game, nil
}

func (pso *PSO) getFloorSwitch(switchId uint16, floor uint16) (bool, error) {
	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(0xAC9FA0+(32*int(floor))), 32)
	if !ok {
		return false, errors.New("unable to getFloorSwitches")
	}
	var mask uint16
	if switchId%16 >= 8 {
		mask = uint16(0x8000) >> (switchId % 8)
	} else {
		mask = uint16(0x80) >> (switchId % 8)
	}
	switchSet := (buf[switchId/16] & mask) > 0
	return switchSet, nil
}

// -------------- Quest Data Block -------------- //
func (pso *PSO) getQuestNumber(questDataPtr uintptr) uint16 {
	return numbers.ReadU16(pso.handle, questDataPtr+0x10)
}

func (pso *PSO) getQuestName(questDataPtr uintptr) (string, error) {
	buf, _, ok := w32.ReadProcessMemory(pso.handle, questDataPtr+0x18, 64)
	if !ok {
		return "", errors.New("unable to getQuestName")
	}
	endIndex := len(buf)
	for index, b := range buf {
		if b == 0x00 {
			endIndex = index
			break
		}
	}
	questName := string(utf16.Decode(buf[0:endIndex]))
	questName = strings.TrimSpace(questName)
	return questName, nil
}

func (pso *PSO) checkQuestStartConditions(questConfig quest.Quest) (bool, error) {
	questStart := false
	if questConfig.StartsOnRegister() {
		registerSet := quest.IsRegisterSet(pso.handle, *questConfig.Start.Register, pso.GameState.questRegisterPointer)
		if questConfig.GetCmodeStage() > 0 {
			cmodeFailedRegister := quest.IsRegisterSet(pso.handle, 253, pso.GameState.questRegisterPointer)
			questStart = registerSet && !cmodeFailedRegister
		} else {
			questStart = registerSet
		}
	} else if questConfig.TerminalQuest() {
		switchSet, err := pso.getFloorSwitch(questConfig.Start.Switch, questConfig.Start.Floor)
		if err != nil {
			return false, err
		}
		questStart = switchSet
	} else if questConfig.StartsAtWarpIn() {
		allPlayers, err := pso.getOtherPlayers()
		if err != nil {
			log.Panicf("unable to get all players %v", err)
		}
		for _, p := range allPlayers {
			if p.Location.Floor != 0 && !p.Location.Warping {
				questStart = true
				break
			}
		}
	}
	return questStart, nil
}

func (pso *PSO) checkQuestEndConditions(questConfig quest.Quest) (bool, error) {
	if questConfig.EndsOnRegister() {
		return quest.IsRegisterSet(pso.handle, *questConfig.End.Register, pso.GameState.questRegisterPointer), nil
	} else if questConfig.End.Floor != 0 {
		return pso.getFloorSwitch(questConfig.End.Switch, questConfig.End.Floor)
	} else {
		return false, errors.New(fmt.Sprintf("Quest %v ends on neither switch nor register", questConfig.Name))
	}
}

func (pso *PSO) getRngSeed() uint32 {
	return numbers.ReadU32Unchecked(pso.handle, 0x00A9C22C)
}

func (pso *PSO) getPlayerCount() uint32 {
	return numbers.ReadU32Unchecked(pso.handle, 0x00AAE168)
}

func (pso *PSO) ephineaFastBurstEnabled() bool {
	fastBurst := false
	if pso.server == constants.EphineaServerName {
		a := uintptr(numbers.ReadU32Unchecked(pso.handle, 0x5B92DA))
		if a > 0 {
			a += 0x5B92DF
			slowBurstPtr := uintptr(numbers.ReadU32Unchecked(pso.handle, a))
			if slowBurstPtr > 0 {
				fastBurst = numbers.ReadU16(pso.handle, slowBurstPtr) == 0
			}
		}
	}
	return fastBurst
}

func (pso *PSO) GetMonsterList() ([]Monster, error) {
	npcArrayAddr := uintptr(0x00AAD720)
	npcCount := int(numbers.ReadU32Unchecked(pso.handle, 0x00AAE164))
	playerCount := int(pso.getPlayerCount())
	ephineaMonsters := uintptr(numbers.ReadU32Unchecked(pso.handle, 0x00B5F800))

	buf, _, ok := w32.ReadProcessMemory(pso.handle, npcArrayAddr, uintptr(4*(playerCount+npcCount+1)))
	if !ok {
		return nil, errors.New("unable to GetMonsterList")
	}
	monsterCount := 0
	monsters := make([]Monster, 0)
	for i := playerCount; i < (playerCount + npcCount); i++ {
		monsterAddr := uintptr(numbers.Uint32FromU16(buf[2*i], buf[(2*i)+1]))
		if monsterAddr != 0 {
			monsterId := numbers.ReadU16(pso.handle, monsterAddr+0x1c)
			monsterType, err := numbers.ReadU32(pso.handle, monsterAddr+0x378)
			if err != nil {
				return nil, err
			}
			var hp uint16
			lastAttackerIndex := numbers.ReadU16(pso.handle, monsterAddr+0x2D8)
			if ephineaMonsters != 0 {
				hp = numbers.ReadU16(pso.handle, ephineaMonsters+0x04+(uintptr(monsterId)*32))
			} else {
				hp = numbers.ReadU16(pso.handle, monsterAddr+0x334)
			}
			if monsterType == 45 {
				// DRL
				if i == 0 {
					hp = numbers.ReadU16(pso.handle, monsterAddr+monsterDeRolLeHP)
					// todo Missing skull hp atm
				} else {
					hp = numbers.ReadU16(pso.handle, monsterAddr+monsterDeRolLeShellHP)
				}
			} else if monsterType == 73 {
				// Barba Ray
				if i == 0 {
					hp = numbers.ReadU16(pso.handle, monsterAddr+monsterBarbaRayHP)
					// todo Missing skull hp atm
				} else {
					hp = numbers.ReadU16(pso.handle, monsterAddr+monsterBarbaRayShellHP)
				}
			}
			// underflow seems to be possible
			if hp > 0x8000 {
				hp = 0
			}
			if monsterType != 0 {
				monsterName, err := pso.getMonsterName(monsterType)
				if err != nil {
					log.Printf("cannot read monster name for id %v %v", monsterType, err)
				} else {
					monsters = append(monsters, Monster{
						Name:            monsterName,
						hp:              hp,
						Id:              monsterId,
						Index:           i,
						UnitxtId:        monsterType,
						LastAttackerIdx: lastAttackerIndex,
						Location: model.Location{
							X: numbers.ReadF32(pso.handle, monsterAddr+0x38),
							Y: numbers.ReadF32(pso.handle, monsterAddr+0x3C),
							Z: numbers.ReadF32(pso.handle, monsterAddr+0x40),
						},
					})
					if hp > 0 {
						monsterCount++
					}
				}
			}
		}
	}
	pso.GameState.MonsterCount = monsterCount
	return monsters, nil
}
