package pso

import (
	"errors"
	"fmt"
	"log"
	"time"
	"unicode/utf16"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/client/internal/numbers"
	"github.com/phelix-/psostats/v2/client/internal/pso/inventory"
	"github.com/phelix-/psostats/v2/client/internal/pso/player"
	"github.com/phelix-/psostats/v2/client/internal/pso/quest"
)

const (
	basePlayerArrayAddress = 0x00A94254
	myPlayerIndexAddress   = 0x00A9C4F4
)

type BaseGameInfo struct {
	episode    uint16
	difficulty uint16
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
	Name       string
	hp         uint16
	Id         uint16
	UnitxtId   uint32
	SpawnTime  time.Time
	KilledTime time.Time
	Alive      bool
	Frame1     bool
}

type Event struct {
	Second      int
	Description string
}

type QuestRun struct {
	Server                   string
	PlayerName               string
	PlayerClass              string
	GuildCard                string
	AllPlayers               []player.BasePlayerInfo
	Id                       string
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
	Meseta                   []uint32
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
	Monsters                 map[int]Monster
	MonsterCount             []int
	MonstersKilledCount      []int
	MonstersDead             int
	WeaponsUsed              map[string]string
	FreezeTraps              []uint16
	previousFt               uint16
	previousDt               uint16
	previousCt               uint16
	FTUsed                   uint16
	DTUsed                   uint16
	CTUsed                   uint16
	previousTp               uint16
	TPUsed                   uint16
}

func (pso *PSO) StartNewQuest(questName string, terminalQuest bool) {
	log.Printf("Starting new quest: %v", questName)

	allPlayers, err := pso.getOtherPlayers()
	if err != nil {
		log.Panicf("unable to get all players %v", err)
	}
	questStartTime := time.Now()
	pso.GameState.QuestStartTime = questStartTime
	pso.CurrentQuest++

	maxPartySupplyableShifta := int16(0)
	for _, p := range allPlayers {
		playerShifta := p.MaxSupplyableShifta()
		if playerShifta > maxPartySupplyableShifta {
			maxPartySupplyableShifta = playerShifta
		}
	}
	maxPartyPbShifta := int16(21 + (20 * (len(allPlayers) - 1)))
	currentShifta := pso.CurrentPlayerData.ShiftaLvl
	pbShifta := currentShifta > maxPartySupplyableShifta
	pbCharged := pso.CurrentPlayerData.PB > 5.0
	pbCategory := pbShifta || pbCharged
	if !terminalQuest {
		lowered := pso.CurrentPlayerData.IsLowered()
		shiftaCast := currentShifta > 0
		pbCategory = pbCategory || lowered || shiftaCast
	}

	pso.Quests[pso.CurrentQuest] = QuestRun{
		Server:                   pso.server,
		PlayerName:               pso.CurrentPlayerData.Name,
		PlayerClass:              pso.CurrentPlayerData.Class,
		GuildCard:                pso.CurrentPlayerData.GuildCard,
		AllPlayers:               allPlayers,
		Difficulty:               pso.GameState.Difficulty,
		Episode:                  pso.GameState.Episode,
		Id:                       fmt.Sprint(pso.CurrentQuest),
		QuestStartTime:           questStartTime,
		QuestStartDate:           questStartTime.Format("15:04 01/02/2006"),
		QuestName:                questName,
		lastRecordedSecond:       -1,
		previousMesetaCharged:    0,
		previousMeseta:           -1,
		DeathCount:               0,
		HP:                       make([]uint16, 0),
		TP:                       make([]uint16, 0),
		Meseta:                   make([]uint32, 0),
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
		Monsters:                 make(map[int]Monster),
		MonsterCount:             make([]int, 0),
		MonstersKilledCount:      make([]int, 0),
		WeaponsUsed:              make(map[string]string),
		FreezeTraps:              make([]uint16, 0),
		previousDt:               0,
		previousFt:               0,
		previousCt:               0,
		FTUsed:                   0,
		DTUsed:                   0,
		CTUsed:                   0,
		previousTp:               0,
		TPUsed:                   0,
	}
}

func (pso *PSO) consolidateFrame() {

	currentQuestRun := pso.Quests[pso.CurrentQuest]
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
		currentQuestRun.Room = append(currentQuestRun.Room, pso.CurrentPlayerData.Room)
		currentQuestRun.ShiftaLvl = append(currentQuestRun.ShiftaLvl, pso.CurrentPlayerData.ShiftaLvl)
		currentQuestRun.DebandLvl = append(currentQuestRun.DebandLvl, pso.CurrentPlayerData.DebandLvl)
		currentQuestRun.MonsterCount = append(currentQuestRun.MonsterCount, pso.GameState.MonsterCount)
		currentQuestRun.Meseta = append(currentQuestRun.Meseta, pso.CurrentPlayerData.Meseta)
		currentQuestRun.MesetaCharged = append(currentQuestRun.MesetaCharged, mesetaCharged)
		currentQuestRun.MonstersKilledCount = append(currentQuestRun.MonstersKilledCount, currentQuestRun.MonstersDead)
		currentQuestRun.FreezeTraps = append(currentQuestRun.FreezeTraps, pso.CurrentPlayerData.FreezeTraps)
		currentQuestRun.Invincible = append(currentQuestRun.Invincible, pso.CurrentPlayerData.InvincibilityFrames > 0)
		for id, display := range pso.Equipment {
			currentQuestRun.WeaponsUsed[id] = display
		}
		if pso.CurrentPlayerData.ShiftaLvl > currentQuestRun.maxPartyPbShifta {
			currentQuestRun.IllegalShifta = true
		}
	}

	if currentQuestRun.lastFloor != pso.CurrentPlayerData.Floor {
		currentQuestRun.Events = append(currentQuestRun.Events, Event{
			Second:      currentSecond,
			Description: GetFloorName(int(currentQuestRun.Episode), int(pso.CurrentPlayerData.Floor)),
		})
		currentQuestRun.lastFloor = pso.CurrentPlayerData.Floor
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
	} else {
		currentQuestRun.QuestDuration = time.Now().Sub(currentQuestRun.QuestStartTime).String()
	}
	pso.Quests[pso.CurrentQuest] = currentQuestRun
}

func GetFloorName(episode int, floor int) string {
	floorName := fmt.Sprintf("Unknown Floor E%vF%v", episode, floor)
	switch episode {
	case 1:
		switch floor {
		case 0:
			floorName = "Pioneer II"
		case 1:
			floorName = "Forest 1"
		case 2:
			floorName = "Forest 2"
		case 3:
			floorName = "Cave 1"
		case 4:
			floorName = "Cave 2"
		case 5:
			floorName = "Cave 3"
		case 6:
			floorName = "Mine 1"
		case 7:
			floorName = "Mine 2"
		case 8:
			floorName = "Ruins 1"
		case 9:
			floorName = "Ruins 2"
		case 10:
			floorName = "Ruins 3"
		case 11:
			floorName = "Dragon"
		case 12:
			floorName = "De Rol Le"
		case 13:
			floorName = "Vol Opt"
		case 14:
			floorName = "Dark Falz"
		case 15:
			floorName = "Lobby"
		}
	case 2:
		switch floor {
		case 0:
			floorName = "Lab"
		case 1:
			floorName = "Temple Alpha"
		case 2:
			floorName = "Temple Beta 2"
		case 3:
			floorName = "Spaceship Alpha"
		case 4:
			floorName = "Spaceship Beta"
		case 5:
			floorName = "CCA"
		case 6:
			floorName = "Jungle North"
		case 7:
			floorName = "Jungle East"
		case 8:
			floorName = "Mountain"
		case 9:
			floorName = "Seaside"
		case 10:
			floorName = "Seabed Upper"
		case 11:
			floorName = "Seabed Lower"
		case 12:
			floorName = "Gal Gryphon"
		case 13:
			floorName = "Olga Flow"
		case 14:
			floorName = "Barba Ray"
		case 15:
			floorName = "Gol Dragon"
		case 16:
			floorName = "Seaside at Night"
		case 17:
			floorName = "Control Tower"

		}
	case 4:
		switch floor {
		case 0:
			floorName = "Pioneer 2"
		case 1:
			floorName = "Crater East"
		case 2:
			floorName = "Crater West"
		case 3:
			floorName = "Crater South"
		case 4:
			floorName = "Crater North"
		case 5:
			floorName = "Crater Interior"
		case 6:
			floorName = "Desert 1"
		case 7:
			floorName = "Desert 2"
		case 8:
			floorName = "Desert 3"
		case 9:
			floorName = "Saint-Milion"
		case 15:
			floorName = "Lobby"
		}
	}
	return floorName
}

func (pso *PSO) consolidateMonsterState(monsters []Monster) {
	now := time.Now()
	currentQuestRun := pso.Quests[pso.CurrentQuest]
	for _, monster := range monsters {
		monsterId := int(monster.Id)
		existingMonster, exists := currentQuestRun.Monsters[monsterId]
		if !exists {
			monster.SpawnTime = now
			monster.Alive = true
			currentQuestRun.Monsters[monsterId] = monster
			existingMonster = monster
		} else if existingMonster.Alive && monster.hp <= 0 {
			// We don't allow frame 0 kills because some of the monsters appear to spawn in with 0 hp.
			// This could be a synchronization issue w/ pso (data is still initializing when we catch it?)
			existingMonster.Alive = false
			currentQuestRun.MonstersDead += 1
			existingMonster.KilledTime = now
			existingMonster.Frame1 = existingMonster.KilledTime.Sub(existingMonster.SpawnTime).Milliseconds() < 60
			if existingMonster.Frame1 {
				log.Printf("frame1? %v | %v", existingMonster.Id, existingMonster.UnitxtId)
			}
			currentQuestRun.Monsters[monsterId] = existingMonster

		}
	}
	pso.Quests[pso.CurrentQuest] = currentQuestRun
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

	address, err := pso.getBaseCharacterAddress(index)
	if err != nil {
		return err
	}
	game, err := pso.getBaseGameInfo()
	if err != nil {
		return err
	}
	pso.GameState.Episode = game.episode
	pso.GameState.Difficulty = game.DifficultyString()

	if address != 0 {

		playerData, err := player.GetPlayerData(pso.handle, address)
		if err != nil {
			return err
		}
		pso.CurrentPlayerData = playerData
		err = pso.getUnitxtStuff()
		if err != nil {
			return err
		}

		equipment, err := inventory.ReadInventory(pso.handle, index)
		if err != nil {
			return err
		}
		pso.Equipment = equipment

		monsters, err := pso.GetMonsterList()
		if err != nil {
			return err
		}

		questPtr, err := quest.GetQuestPointer(pso.handle)
		if err != nil {
			return err
		}
		if questPtr != 0 {
			questDataPtr, err := pso.getQuestDataPointer(questPtr)
			if err != nil {
				return err
			}

			err = pso.getQuestName(questDataPtr)
			if err != nil {
				return err
			}

			if !pso.GameState.QuestStarted {
				questConditions, exists := pso.questTypes.GetQuest(int(pso.GameState.Episode), pso.GameState.QuestName)
				if exists && !questConditions.Ignore {
					if questConditions.StartsOnRegister() {
						registerSet, err := quest.IsRegisterSet(pso.handle, uint16(*questConditions.StartTrigger.Register))
						if err != nil {
							return err
						}
						pso.GameState.QuestStarted = registerSet
					} else if questConditions.TerminalQuest() {
						switchSet, err := pso.getFloorSwitch(uint16(questConditions.StartTrigger.Switch), questConditions.StartTrigger.Floor)
						if err != nil {
							return err
						}
						pso.GameState.QuestStarted = switchSet
					} else if questConditions.StartsAtWarpIn() {
						allPlayers, err := pso.getOtherPlayers()
						if err != nil {
							log.Panicf("unable to get all players %v", err)
						}
						for playerIndex, p := range allPlayers {
							if p.Floor != 0 && !p.Warping {
								if len(pso.GameState.PlayerArray) > playerIndex {
									previousPlayerState := pso.GameState.PlayerArray[playerIndex]
									if previousPlayerState.GuildCard == p.GuildCard && previousPlayerState.Warping {
										pso.GameState.QuestStarted = true
									}
								}
							}
						}
						pso.GameState.PlayerArray = allPlayers
					}
				}
				if pso.GameState.QuestStarted && pso.GameState.AllowQuestStart {
					pso.StartNewQuest(pso.GameState.QuestName, exists && questConditions.TerminalQuest())
				}
			} else if !pso.GameState.QuestComplete {
				questConditions, exists := pso.questTypes.GetQuest(int(pso.GameState.Episode), pso.GameState.QuestName)
				if exists {
					if questConditions.EndsOnRegister() {
						registerSet, err := quest.IsRegisterSet(pso.handle, uint16(*questConditions.EndTrigger.Register))
						if err != nil {
							return err
						}

						if registerSet {
							pso.GameState.QuestComplete = true
							pso.GameState.QuestEndTime = time.Now()
						}
					} else if questConditions.EndTrigger.Floor == pso.CurrentPlayerData.Floor {
						switched, err := pso.getFloorSwitch(uint16(questConditions.EndTrigger.Switch), pso.CurrentPlayerData.Floor)
						if err != nil {
							return err
						}
						if switched {
							pso.GameState.QuestComplete = true
							pso.GameState.QuestEndTime = time.Now()
						}
					}
				}
			}
			if pso.GameState.QuestStarted {
				if pso.GameState.AllowQuestStart {
					pso.consolidateFrame()
					pso.consolidateMonsterState(monsters)
				}
			} else {
				pso.GameState.AllowQuestStart = true
			}
		} else {
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

func (pso *PSO) getBaseCharacterAddress(index uint8) (int, error) {
	address := basePlayerArrayAddress + (4 * int(index))
	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(address), 4)
	if !ok {
		return 0, errors.New("unable to getBaseCharacterAddress")
	}
	baseAddress := int(buf[1])<<16 + int(buf[0])
	// log.Printf("Base address: %v", baseAddress)
	return baseAddress, nil
}

func (pso *PSO) getUnitxtStuff() error {
	unitxtAddr, err := numbers.ReadU32(pso.handle, uintptr(0x00a9cd50))
	if err != nil {
		return fmt.Errorf("getUnitxtStuff %w", err)
	}
	monsterUnitxtAddr, err := numbers.ReadU32(pso.handle, uintptr(unitxtAddr+16))
	if err != nil {
		return fmt.Errorf("getUnitxtStuff2 %w", err)
	}
	pso.GameState.monsterUnitxtAddr = monsterUnitxtAddr

	return nil
}

func (pso *PSO) getMonsterName(monsterId uint32) (string, error) {
	monsterName, exists := pso.MonsterNames[monsterId]
	if exists {
		return monsterName, nil
	}

	monsterNameAddr, err := numbers.ReadU32(pso.handle, uintptr(pso.GameState.monsterUnitxtAddr+(4*monsterId)))
	if err != nil {
		return "", fmt.Errorf("getMonsterName1 %v %w", monsterId, err)
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
		address, err := pso.getBaseCharacterAddress(uint8(i))
		if err != nil {
			return nil, err
		}
		if address != 0 {
			playerData, err := player.GetPlayerData(pso.handle, address)
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
	game := BaseGameInfo{
		episode:    episode,
		difficulty: difficulty,
	}
	return game, nil
}

func (pso *PSO) getFloorSwitch(switchId uint16, floor uint16) (bool, error) {
	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(0xAC9FA0+(32*int(floor))), 32)
	if !ok {
		return false, errors.New("unable to getFloorSwitches")
	}
	mask := uint16(0x00)
	if switchId%16 > 8 {
		mask = uint16(0x8000) >> (switchId % 8)
	} else {
		mask = uint16(0x80) >> (switchId % 8)
	}
	// log.Printf("%v | 0x%04x", switchId%16, mask)
	switchSet := (buf[switchId/16] & mask) > 0
	// log.Printf("switch[%v] = %v | 0x%04x 0x%04x", switchId, switchSet, buf[switchId/16], mask)
	// log.Print("\n")
	return switchSet, nil
}

// -------------- Quest Data Block -------------- //
func (pso *PSO) getQuestDataPointer(questPtr uint32) (uint32, error) {
	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(questPtr+0x19C), 4)
	if !ok {
		return 0, errors.New("unable to getQuestDataPointer")
	}
	questDataPtr := uint32(buf[1])<<16 + uint32(buf[0])
	return questDataPtr, nil
}

func (pso *PSO) getQuestName(questDataPtr uint32) error {
	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(questDataPtr+0x18), 64)
	if !ok {
		return errors.New("unable to getQuestName")
	}
	endIndex := len(buf)
	for index, b := range buf {
		if b == 0x00 {
			endIndex = index
			break
		}
	}

	something := utf16.Decode(buf[0:endIndex])
	aString := string(something)
	pso.GameState.QuestName = aString
	return nil
}

func (pso *PSO) getPlayerCount() (uint32, error) {
	playerCountAddr := 0x00AAE168
	return numbers.ReadU32(pso.handle, uintptr(playerCountAddr))
}

func (pso *PSO) GetMonsterList() ([]Monster, error) {
	npcCountAddr := 0x00AAE164
	npcArrayAddr := 0x00AAD720
	npcCount, err := numbers.ReadU32(pso.handle, uintptr(npcCountAddr))
	if err != nil {
		return nil, err
	}
	playerCount, err := pso.getPlayerCount()
	pCountInt := int(playerCount)

	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(npcArrayAddr), uintptr(4*(playerCount+npcCount+1)))
	if !ok {
		return nil, errors.New("unable to GetMonsterList")
	}
	monsterCount := 0
	monsters := make([]Monster, 0)
	for i := 0; i < int(npcCount); i++ {
		effectiveI := i + pCountInt
		// log.Printf("npc[%v] 0x%08x | 0x%08x 0x%08x", i, npcArrayAddr+(2*effectiveI), buf[2*effectiveI], buf[(2*effectiveI)+1])
		monsterAddr := numbers.Uint32FromU16(buf[2*effectiveI], buf[(2*effectiveI)+1])
		monsterId, err := numbers.ReadU16(pso.handle, uintptr(monsterAddr+0x1c))
		if err != nil {
			return nil, err
		}
		monsterType, err := numbers.ReadU32(pso.handle, uintptr(monsterAddr+0x378))
		if err != nil {
			return nil, err
		}
		hp, err := numbers.ReadU16(pso.handle, uintptr(monsterAddr+0x334))
		if err != nil {
			return nil, err
		}
		// log.Printf("npc[@0x%08x] id = 0x%04x - 0x%08x", monsterAddr, monsterId, monsterType)
		if monsterType != 0 {
			monsterName, err := pso.getMonsterName(monsterType)
			if err != nil {
				return nil, err
			}

			monsters = append(monsters, Monster{
				Name:     monsterName,
				hp:       hp,
				Id:       monsterId,
				UnitxtId: monsterType,
			})
			if hp > 0 {
				monsterCount++
			}
		}
	}
	pso.GameState.MonsterCount = monsterCount
	return monsters, nil
}
