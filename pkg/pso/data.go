package pso

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"time"
	"unicode/utf16"

	"github.com/TheTitanrain/w32"
)

const (
	basePlayerArrayAddress = 0x00A94254
	myPlayerIndexAddress   = 0x00A9C4F4
)

type Event struct {
	Second      int
	Description string
}
type Player struct {
	Name  string
	GC    string
	Class string
}

type QuestRun struct {
	PlayerName            string
	PlayerClass           string
	AllPlayers            []Player
	Id                    string
	Difficulty            string
	Episode               uint16
	QuestName             string
	QuestComplete         bool
	QuestStartTime        time.Time
	QuestStartDate        string
	QuestEndTime          time.Time
	QuestDuration         string
	lastRecordedSecond    int
	lastRecordedHp        uint16
	lastFloor             uint16
	previousMeseta        int
	previousMesetaCharged int
	HP                    []uint16
	TP                    []uint16
	Meseta                []uint32
	MesetaCharged         []int
	Room                  []uint16
	ShiftaLvl             []int16
	DebandLvl             []int16
	Invincible            []bool
	Events                []Event
	Monsters              map[int]Monster
	MonsterCount          []int
	MonstersKilledCount   []int
	monstersDead          int
}

func (pso *PSO) StartNewQuest(questName string) {
	log.Printf("Starting new quest")

	allPlayers, err := pso.getOtherPlayers()
	if err != nil {
		log.Panicf("unable to get all players %v", err)
	}
	questStartTime := time.Now()
	pso.GameState.QuestStartTime = questStartTime
	pso.CurrentQuest++
	pso.Quests[pso.CurrentQuest] = QuestRun{
		PlayerName:            pso.CurrentPlayerData.CharacterName,
		PlayerClass:           pso.CurrentPlayerData.Class,
		AllPlayers:            allPlayers,
		Difficulty:            pso.CurrentPlayerData.Difficulty,
		Episode:               pso.CurrentPlayerData.Episode,
		Id:                    fmt.Sprint(pso.CurrentQuest),
		QuestStartTime:        questStartTime,
		QuestStartDate:        questStartTime.Format("15:04 01/02/2006"),
		QuestName:             questName,
		lastRecordedSecond:    -1,
		previousMesetaCharged: 0,
		previousMeseta:        -1,
		HP:                    make([]uint16, 0),
		TP:                    make([]uint16, 0),
		Meseta:                make([]uint32, 0),
		MesetaCharged:         make([]int, 0),
		Room:                  make([]uint16, 0),
		ShiftaLvl:             make([]int16, 0),
		DebandLvl:             make([]int16, 0),
		Invincible:            make([]bool, 0),
		Events:                make([]Event, 0),
		Monsters:              make(map[int]Monster),
		MonsterCount:          make([]int, 0),
		MonstersKilledCount:   make([]int, 0),
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
		mesetaCharged := 0
		previousMeseta := currentQuestRun.previousMeseta
		if previousMeseta != -1 {
			mesetaDifference := (previousMeseta - int(pso.CurrentPlayerData.Meseta))
			if mesetaDifference > 0 {
				// negative means meseta picked up, ignoring i guess
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
		currentQuestRun.MonstersKilledCount = append(currentQuestRun.MonstersKilledCount, currentQuestRun.monstersDead)
	}

	if currentQuestRun.lastFloor != pso.CurrentPlayerData.Floor {
		currentQuestRun.Events = append(currentQuestRun.Events, Event{
			Second:      currentSecond,
			Description: GetFloorName(int(currentQuestRun.Episode), int(pso.CurrentPlayerData.Floor)),
		})
		currentQuestRun.lastFloor = pso.CurrentPlayerData.Floor
	}

	if pso.CurrentPlayerData.HP == 0 && currentQuestRun.lastRecordedHp != 0 {
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
	floorName := fmt.Sprintf("unknown floor e%vf%v", episode, floor)
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
			existingMonster.Alive = false
			currentQuestRun.monstersDead += 1
			existingMonster.KilledTime = now
			existingMonster.Frame1 = existingMonster.KilledTime.Sub(existingMonster.SpawnTime).Milliseconds() < 50
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
		pso.CurrentPlayerData.QuestName = "No Active Quest"
		pso.GameState.QuestComplete = false
		pso.GameState.QuestStarted = false

		log.Fatal("RefreshData: connection to window lost")
		return errors.New("RefreshData: connection to window lost")
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
	_, err = pso.getBaseGameInfo()
	if err != nil {
		return err
	}

	if address != 0 {

		err := pso.getPlayerData(address)
		if err != nil {
			return err
		}
		err = pso.getUnitxtStuff()
		if err != nil {
			return err
		}

		monsters, err := pso.GetMonsterList()
		if err != nil {
			return err
		}

		name, err := pso.getCharacterName(address)
		if err != nil {
			return err
		}
		pso.CurrentPlayerData.CharacterName = name

		gc, err := pso.getGuildCard(address)
		if err != nil {
			return err
		}
		pso.CurrentPlayerData.Guildcard = gc

		questPtr, err := pso.getQuestPointer()
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
				questConditions, exists := pso.questTypes.GetQuest(int(pso.CurrentPlayerData.Episode), pso.CurrentPlayerData.QuestName)
				if exists {
					if questConditions.StartTrigger.Register > 0 {
						registerSet, err := pso.isQuestRegisterSet(uint16(questConditions.StartTrigger.Register))
						if err != nil {
							return err
						}
						pso.GameState.QuestStarted = registerSet
					} else if uint16(questConditions.StartTrigger.Floor) == pso.CurrentPlayerData.Floor {
						switchSet, err := pso.getFloorSwitch(uint16(questConditions.StartTrigger.Switch), pso.CurrentPlayerData.Floor)
						if err != nil {
							return err
						}
						pso.GameState.QuestStarted = switchSet
					}
				} else {
					pso.GameState.QuestStarted = pso.CurrentPlayerData.Floor != 0
				}
				if pso.GameState.QuestStarted {
					pso.StartNewQuest(pso.CurrentPlayerData.QuestName)
				}
			} else if !pso.GameState.QuestComplete {
				questConditions, exists := pso.questTypes.GetQuest(int(pso.CurrentPlayerData.Episode), pso.CurrentPlayerData.QuestName)
				if exists {
					if questConditions.EndTrigger.Register > 0 {
						registerSet, err := pso.isQuestRegisterSet(uint16(questConditions.EndTrigger.Register))
						if err != nil {
							return err
						}

						if registerSet {
							pso.GameState.QuestComplete = true
							pso.GameState.QuestEndTime = time.Now()
						}
					} else if uint16(questConditions.EndTrigger.Floor) == pso.CurrentPlayerData.Floor {
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
				pso.consolidateFrame()
				pso.consolidateMonsterState(monsters)
			}
		} else {
			pso.CurrentPlayerData.QuestName = "No Active Quest"
			pso.GameState.QuestComplete = false
			pso.GameState.QuestStarted = false
		}
	} else {
		pso.CurrentPlayerData.QuestName = "No Active Quest"
		pso.GameState.QuestComplete = false
		pso.GameState.QuestStarted = false
	}
	pso.CurrentPlayerData.Time = time.Now()

	return nil
}

func (pso *PSO) getMyPlayerIndex() (uint8, error) {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(myPlayerIndexAddress), 4)
	if !ok {
		return 0, errors.New("Unable to find player index")
	}
	index := uint8(buf[0])
	return index, nil
}

func (pso *PSO) getBaseCharacterAddress(index uint8) (int, error) {
	address := basePlayerArrayAddress + (4 * int(index))
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(address), 4)
	if !ok {
		return 0, errors.New("Unable to getBaseCharacterAddress")
	}
	baseAddress := int(buf[1])<<16 + int(buf[0])
	// log.Printf("Base address: %v", baseAddress)
	return baseAddress, nil
}

func (pso *PSO) getCharacterName(playerAddress int) (string, error) {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x428), 24)
	if !ok {
		return "", errors.New("Unable to getCharacterName")
	}

	endIndex := len(buf)
	for index, b := range buf {
		if b == 0x00 {
			endIndex = index
			break
		}
	}

	return string(utf16.Decode(buf[2:endIndex])), nil
}

func (pso *PSO) getGuildCard(playerAddress int) (string, error) {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x92a), 16)
	if !ok {
		return "", errors.New("Unable to getGuildCard")
	}

	byteBuf := bytes.NewBuffer(make([]byte, 0, len(buf)))
	for i := 2; i < 8; i++ {
		b := buf[i]
		split := make([]byte, 2)
		binary.LittleEndian.PutUint16(split, b)
		byteBuf.Write(split)
	}
	return byteBuf.String(), nil
}

func (pso *PSO) getUnitxtStuff() error {
	unitxtAddr, err := pso.readU32(uintptr(0x00a9cd50))
	if err != nil {
		return fmt.Errorf("getUnitxtStuff %w", err)
	}
	monsterUnitxtAddr, err := pso.readU32(uintptr(unitxtAddr + 16))
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

	monsterNameAddr, err := pso.readU32(uintptr(pso.GameState.monsterUnitxtAddr + (4 * uint32(monsterId))))
	if err != nil {
		return "", fmt.Errorf("getMonsterName1 %w", err)
	}
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(monsterNameAddr), 32)
	if !ok {
		return "", errors.New("Unable to getMonsterName")
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

func (pso *PSO) getOtherPlayers() ([]Player, error) {
	players := make([]Player, 0)
	for i := 0; i < 12; i++ {
		address, err := pso.getBaseCharacterAddress(uint8(i))
		if err != nil {
			return nil, err
		}
		if address != 0 {
			name, err := pso.getCharacterName(address)
			if err != nil {
				return nil, err
			}
			gc, err := pso.getGuildCard(address)
			if err != nil {
				return nil, err
			}
			players = append(players, Player{
				Name: name,
				GC:   gc,
			})

		}
	}
	// log.Printf("players: %v", players)
	return players, nil
}

// func (pso *PSO) getInventoryAddr(playerAddress int) error {
// 	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0xDF4), 2)
// 	if !ok {
// 		return errors.New("Unable to getInventoryPointer")
// 	}
// 	invAddress := int(buf[1])<<16 + int(buf[0])
// 	if (invAddress != 0) {
// 		buf, _, ok = w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(invAddress), 2)

// 	}
// 	hp := buf[0]
// 	log.Printf("Player HP: %v", hp)
// 	return hp, nil
// }

func (pso *PSO) getClass(classBits uint16) string {
	class := "Unknown class"
	switch classBits {
	case 0x00:
		class = "HUmar"
		break
	case 0x01:
		class = "HUnewearl"
		break
	case 0x02:
		class = "HUcast"
		break
	case 0x09:
		class = "HUcaseal"
		break
	case 0x03:
		class = "RAmar"
		break
	case 0x0B:
		class = "RAmarl"
		break
	case 0x04:
		class = "RAcast"
		break
	case 0x05:
		class = "RAcaseal"
		break
	case 0x0A:
		class = "FOmar"
		break
	case 0x06:
		class = "FOmarl"
		break
	case 0x07:
		class = "FOnewm"
		break
	case 0x08:
		class = "FOnewearl"
		break
	}
	return class
}

func (pso *PSO) getPlayerData(playerAddress int) error {
	base := 0x028
	max := 0xE4E
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+base), uintptr((max-base)+4))
	if !ok {
		return errors.New("Unable to getPlayerData")
	}
	pso.CurrentPlayerData.Room = buf[(0x028-base)/2]
	pso.CurrentPlayerData.KillCount = buf[(0x11A-base)/2]
	shiftaMultiplier := float32FromU16(buf[(0x278-base)/2], buf[(0x27a-base)/2])
	pso.CurrentPlayerData.ShiftaLvl = getShiftaLvlFromMultiplier(shiftaMultiplier)
	debandMultiplier := float32FromU16(buf[(0x278+12-base)/2], buf[(0x278+14-base)/2])
	pso.CurrentPlayerData.DebandLvl = getShiftaLvlFromMultiplier(debandMultiplier)
	pso.CurrentPlayerData.MaxHP = buf[(0x2BC-base)/2]
	pso.CurrentPlayerData.MaxTP = buf[(0x2BE-base)/2]
	pso.CurrentPlayerData.HP = buf[(0x334-base)/2]
	pso.CurrentPlayerData.TP = buf[(0x336-base)/2]
	pso.CurrentPlayerData.Floor = buf[(0x3F0-base)/2]
	pso.CurrentPlayerData.InvincibilityFrames = uint32FromU16(buf[(0x720-base)/2], buf[(0x722-base)/2])
	class := (buf[(0x961-base)/2] & 0xF00) >> 8
	pso.CurrentPlayerData.Class = pso.getClass(class)
	pso.CurrentPlayerData.Meseta = uint32FromU16(buf[(0xE4C-base)/2], buf[(0xE4E-base)/2])
	// log.Printf("%v/%v", read, uintptr(10+(0x3f0-base)/2))
	// for i := 0; i <= len(buf)-8; i += 8 {
	// 	log.Printf("0x%x | 0x%x - %x %x %x %x %x %x %x %x", playerAddress+base+(2*i), base+(2*i), buf[i], buf[i+1], buf[i+2], buf[i+3], buf[i+4], buf[i+5], buf[i+6], buf[i+7])
	// }

	// log.Printf("%v - %v/%v - %v/%v", len(buf), pso.CurrentPlayerData.HP, pso.CurrentPlayerData.MaxHP,
	// 	pso.CurrentPlayerData.TP, pso.CurrentPlayerData.MaxTP)
	return nil
}

func uint32FromU16(lsb uint16, msb uint16) uint32 {
	return uint32(msb)<<16 + uint32(lsb)
}

func float32FromU16(lsb uint16, msb uint16) float32 {
	combinedValue := uint32FromU16(lsb, msb)
	return math.Float32frombits(combinedValue)
}

func getShiftaLvlFromMultiplier(multiplier float32) int16 {
	level := int16(0)
	if multiplier != 0 {
		level = int16(1 + math.Round(((math.Abs(float64(multiplier))*100)-10)/1.3))
		if multiplier < 0 {
			level = -level
		}
	}
	return level
}

type BaseGameInfo struct {
	episode    uint16
	difficulty uint16
}

func (pso *PSO) getBaseGameInfo() (BaseGameInfo, error) {
	base := 0x00A9B1C8
	max := 0x00A9CD68
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(base), uintptr((max-base)+2))
	if !ok {
		return BaseGameInfo{}, errors.New("Unable to getDifficulty")
	}

	difficulty := buf[(0x00A9CD68-base)/2]
	episode := buf[(0x00A9B1C8-base)/2] + 1
	if episode == 3 {
		episode = 4
	}
	pso.CurrentPlayerData.Episode = episode
	pso.CurrentPlayerData.Difficulty = renderDifficulty(difficulty)
	return BaseGameInfo{
		episode:    episode,
		difficulty: difficulty,
	}, nil
}

func renderDifficulty(difficulty uint16) string {
	switch difficulty {
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

func (pso *PSO) getFloorSwitches(floor uint16) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(0xAC9FA0+(32*int(floor))), 32)
	if !ok {
		return errors.New("Unable to getFloorSwitches")
	}
	anySet := false
	for _, byte := range buf {
		if byte > 0 {
			anySet = true
			break
		}
	}
	pso.GameState.FloorSwitches = anySet
	return nil
}

func (pso *PSO) isQuestRegisterSet(registerId uint16) (bool, error) {
	questRegisterAddress, err := pso.readU32(uintptr(0x00A954B0))
	if err != nil {
		return false, err
	}
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(int(questRegisterAddress)+(4*int(registerId))), 2)
	if !ok {
		return false, errors.New("Unable to isQuestRegisterSet")
	}
	registerSet := buf[0] > 0
	// log.Printf("R[%v]@%x = %v", registerId, uintptr(0x00A954B0+(4*int(registerId))), registerSet)
	// log.Print("\n")
	return registerSet, nil
}

func (pso *PSO) getFloorSwitch(switchId uint16, floor uint16) (bool, error) {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(0xAC9FA0+(32*int(floor))), 32)
	if !ok {
		return false, errors.New("Unable to getFloorSwitches")
	}
	mask := uint16(0x80) >> (switchId % 16)
	switchSet := (buf[switchId/16] & mask) > 0
	// log.Printf("switch[%v] = %v", switchId, switchSet)
	// log.Print("\n")
	return switchSet, nil
}

func (pso *PSO) getQuestPointer() (uint32, error) {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(0xA95AA8), 4)
	if !ok {
		return 0, errors.New("Unable to getQuestPointer")
	}
	questPointer := uint32(buf[1])<<16 + uint32(buf[0])
	return questPointer, nil
}

func (pso *PSO) getQuestDataPointer(questPtr uint32) (uint32, error) {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(questPtr+0x19C), 4)
	if !ok {
		return 0, errors.New("Unable to getQuestDataPointer")
	}
	questDataPtr := uint32(buf[1])<<16 + uint32(buf[0])
	return questDataPtr, nil
}

func (pso *PSO) getQuestName(questDataPtr uint32) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(questDataPtr+0x18), 64)
	if !ok {
		return errors.New("Unable to getQuestName")
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
	pso.CurrentPlayerData.QuestName = aString
	return nil
}

func (pso *PSO) GetPlayerData() PlayerData {
	return pso.CurrentPlayerData
}

func (pso *PSO) getPlayerCount() (uint32, error) {
	playerCountAddr := 0x00AAE168
	return pso.readU32(uintptr(playerCountAddr))
}

func (pso *PSO) GetMonsterList() ([]Monster, error) {
	npcCountAddr := 0x00AAE164
	npcArrayAddr := 0x00AAD720
	npcCount, err := pso.readU32(uintptr(npcCountAddr))
	if err != nil {
		return nil, err
	}
	playerCount, err := pso.getPlayerCount()
	pCountInt := int(playerCount)

	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(npcArrayAddr), uintptr(4*(playerCount+npcCount+1)))
	if !ok {
		return nil, errors.New("Unable to GetMonsterList")
	}
	// for i := 0; i < pCountInt; i++ {
	// 	log.Printf("player[%v] 0x%08x | 0x%08x 0x%08x", i, npcArrayAddr+(2*i), buf[2*i], buf[(2*i)+1])
	// }
	// log.Printf("%v npcs", npcCount)
	monsterCount := 0
	monsters := make([]Monster, 0)
	for i := 0; i < int(npcCount); i++ {
		effectiveI := (i + pCountInt)
		// log.Printf("npc[%v] 0x%08x | 0x%08x 0x%08x", i, npcArrayAddr+(2*effectiveI), buf[2*effectiveI], buf[(2*effectiveI)+1])
		monsterAddr := uint32FromU16(buf[2*effectiveI], buf[(2*effectiveI)+1])
		monsterId, err := pso.readU16(uintptr(monsterAddr + 0x1c))
		if err != nil {
			return nil, err
		}
		monsterType, err := pso.readU32(uintptr(monsterAddr + 0x378))
		if err != nil {
			return nil, err
		}
		hp, err := pso.readU16(uintptr(monsterAddr + 0x334))
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

	// for i := 0; i < int(npcCount); i++ {
	// 	monsterPtr, err := pso.readU32(uintptr(npcArrayAddr + (4 * i)))
	// }

	// for i := uint32(0); i < entityCount; i++ {

	// }

	return monsters, nil
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

func (pso *PSO) readU16(address uintptr) (uint16, error) {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), address, 2)
	if !ok {
		return 0, errors.New(fmt.Sprintf("Unable to readU16 @0x%08x", address))
	}
	return buf[0], nil
}

func (pso *PSO) readU32(address uintptr) (uint32, error) {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), address, 4)
	if !ok {
		return 0, errors.New(fmt.Sprintf("Unable to read 0x%08x", address))
	}
	return uint32(buf[1])<<16 + uint32(buf[0]), nil
}
