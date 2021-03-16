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

func (pso *PSO) RefreshData() error {
	if !pso.connected {
		log.Fatal("RefreshData: connection to window lost")
		return errors.New("RefreshData: connection to window lost")
	}

	// w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(pso.baseAddress+0x0AA1F1314), 4)

	index, err := pso.getMyPlayerIndex()
	if err != nil {
		log.Fatal("Unable to find player index")
		return err
	}

	address, err := pso.getBaseCharacterAddress(index)
	if err != nil {
		return err
	}

	err = pso.getEpisode()
	if err != nil {
		return err
	}
	err = pso.getDifficulty()
	if err != nil {
		return err
	}

	err = pso.GetMonsterList()
	if err != nil {
		return err
	}

	if address != 0 {
		hp, err := pso.getHP(address)
		if err != nil {
			return err
		}

		maxHp, err := pso.getMaxHP(address)
		if err != nil {
			return err
		}
		err = pso.getTP(address)
		if err != nil {
			return err
		}
		err = pso.getMaxTP(address)
		if err != nil {
			return err
		}
		err = pso.getMeseta(address)
		if err != nil {
			return err
		}
		err = pso.getCharacterName(address)
		if err != nil {
			return err
		}
		err = pso.getClass(address)
		if err != nil {
			return err
		}

		err = pso.getGuildCard(address)
		if err != nil {
			return err
		}

		err = pso.getFloor(address)
		if err != nil {
			return err
		}

		err = pso.getRoom(address)
		if err != nil {
			return err
		}
		err = pso.getInvincibilityFrames(address)
		if err != nil {
			return err
		}
		err = pso.getKillCount(address)
		if err != nil {
			return err
		}

		err = pso.getShiftaLvl(address)
		if err != nil {
			return err
		}

		err = pso.getDebandLvl(address)
		if err != nil {
			return err
		}

		pso.CurrentPlayerData.HP = hp
		pso.CurrentPlayerData.MaxHP = maxHp

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
				quests := pso.questTypes[int(pso.CurrentPlayerData.Episode)]
				questConditions, exists := quests[pso.CurrentPlayerData.QuestName]
				if !exists {
					log.Printf("quest '%v' not found", pso.CurrentPlayerData.QuestName)
				}
				if exists {
					if questConditions.StartingSwitchFloor == pso.CurrentPlayerData.Floor {
						switchSet, err := pso.getFloorSwitch(questConditions.StartingSwitch, pso.CurrentPlayerData.Floor)
						if err != nil {
							return err
						}
						pso.GameState.QuestStarted = switchSet
					}
				} else {
					pso.GameState.QuestStarted = pso.CurrentPlayerData.Floor != 0
				}
				if pso.GameState.QuestStarted {
					pso.GameState.QuestStartTime = time.Now()
				}
			} else if !pso.GameState.QuestComplete {
				quests := pso.questTypes[int(pso.CurrentPlayerData.Episode)]
				questConditions, exists := quests[pso.CurrentPlayerData.QuestName]
				if exists {
					if questConditions.EndingQuestRegister > 0 {
						registerSet, err := pso.isQuestRegisterSet(questConditions.EndingQuestRegister)
						if err != nil {
							return err
						}

						if registerSet {
							pso.GameState.QuestComplete = true
							pso.GameState.QuestEndTime = time.Now()
						}
					} else if questConditions.EndingSwitchFloor == pso.CurrentPlayerData.Floor {
						switched, err := pso.getFloorSwitch(questConditions.EndingSwitch, pso.CurrentPlayerData.Floor)
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
		} else {
			pso.CurrentPlayerData.QuestName = "No Active Quest"
			pso.GameState.QuestComplete = false
			pso.GameState.QuestStarted = false
		}
	}

	// buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(0x0AA01314), 4)
	// if !ok {
	// 	log.Fatal("RefreshData: unable to read process memory")
	// 	return errors.New("RefreshData: unable to read process memory")
	// }
	// log.Printf("maybe? %v", buf[0])

	// byteBuf := bytes.NewBuffer(make([]byte, 0, len(buf)))

	// for _, b := range buf {
	// 	split := make([]byte, 2)
	// 	binary.LittleEndian.PutUint16(split, b)
	// 	byteBuf.Write(split)
	// }

	// err = binary.Read(byteBuf, binary.LittleEndian, pso.dataBlock)
	// if err != nil {
	// 	log.Fatal("RefreshData: unable to encode data block")
	// 	return fmt.Errorf("RefreshData: unable to encode data block: %w", err)
	// }
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

func (pso *PSO) getMeseta(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0xE4C), 4)
	if !ok {
		return errors.New("Unable to getInventoryPointer")
	}
	meseta := uint32(buf[1])<<16 + uint32(buf[0])
	pso.CurrentPlayerData.Meseta = meseta
	return nil
}

func (pso *PSO) getCharacterName(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x428), 24)
	if !ok {
		return errors.New("Unable to getCharacterName")
	}
	// byteBuf := bytes.NewBuffer(make([]byte, 0, len(buf)))

	endIndex := len(buf)
	for index, b := range buf {
		if b == 0x00 {
			endIndex = index
			break
		}
	}
	something := utf16.Decode(buf[2:endIndex])
	// for _, b := range buf {
	// 	byteBuf.WriteString(fmt.Sprintf("%04X", b))
	// 	// split := make([]byte, 2)
	// 	// binary.BigEndian.PutUint16(split, b)
	// 	// byteBuf.Write(split)
	// }
	aString := string(something)
	pso.CurrentPlayerData.CharacterName = aString
	return nil
}

func (pso *PSO) getGuildCard(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x938), 24)
	if !ok {
		return errors.New("Unable to getGuildCard")
	}

	// endIndex := len(buf)
	// for index, b := range buf {
	// 	if b == 0x00 {
	// 		endIndex = index
	// 		break
	// 	}
	// }
	byteBuf := bytes.NewBuffer(make([]byte, 0, len(buf)))
	for _, b := range buf {
		// byteBuf.WriteString(fmt.Sprintf("%04X", b))
		if b == 0x00 {
			break
		}
		split := make([]byte, 2)
		binary.LittleEndian.PutUint16(split, b)
		byteBuf.Write(split)
	}
	// utf8.DecodeRune()
	// something := utf16.Decode(buf[0:endIndex])
	// aString := string(something)
	aString := byteBuf.String()
	pso.CurrentPlayerData.Guildcard = aString
	return nil
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

func (pso *PSO) getClass(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x961), 1)
	if !ok {
		return errors.New("Unable to getClass")
	}

	class := "Unknown class"
	switch buf[0] {
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
	pso.CurrentPlayerData.Class = class
	return nil
}

func (pso *PSO) getHP(playerAddress int) (uint16, error) {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x334), 2)
	if !ok {
		return 0, errors.New("Unable to getHP")
	}
	hp := buf[0]
	return hp, nil
}

func (pso *PSO) getMaxHP(playerAddress int) (uint16, error) {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x2BC), 2)
	if !ok {
		return 0, errors.New("Unable to getMaxHP")
	}
	hp := buf[0]
	return hp, nil
}
func (pso *PSO) getTP(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x336), 2)
	if !ok {
		return errors.New("Unable to getTP")
	}
	pso.CurrentPlayerData.TP = buf[0]
	return nil
}
func (pso *PSO) getMaxTP(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x2BE), 2)
	if !ok {
		return errors.New("Unable to getMaxTP")
	}
	pso.CurrentPlayerData.MaxTP = buf[0]
	return nil
}
func (pso *PSO) getDifficulty() error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(0xA9CD68), 2)
	if !ok {
		return errors.New("Unable to getEpisode")
	}
	difficulty := buf[0]
	pso.CurrentPlayerData.Difficulty = difficulty
	return nil
}

func (pso *PSO) getEpisode() error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(0x00A9B1C8), 2)
	if !ok {
		return errors.New("Unable to getEpisode")
	}
	episode := buf[0] + 1
	if episode == 3 {
		episode = 4
	}
	pso.CurrentPlayerData.Episode = episode
	return nil
}

func (pso *PSO) getFloor(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x3F0), 2)
	if !ok {
		return errors.New("Unable to getFloor")
	}
	floor := buf[0]
	pso.CurrentPlayerData.Floor = floor
	return nil
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
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(int(questRegisterAddress)+(4*int(registerId))), 1)
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

func (pso *PSO) getRoom(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x28), 2)
	if !ok {
		return errors.New("Unable to getRoom")
	}
	room := buf[0]
	pso.CurrentPlayerData.Room = room
	return nil
}

func (pso *PSO) getKillCount(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x11A), 2)
	if !ok {
		return errors.New("Unable to killCount")
	}
	killCount := buf[0]
	pso.CurrentPlayerData.KillCount = killCount
	return nil
}

func (pso *PSO) getShiftaLvl(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x278), 4)
	if !ok {
		return errors.New("Unable to getShiftaLvl")
	}
	totalValue := uint32(buf[1])<<16 + uint32(buf[0])

	multiplier := math.Float32frombits(totalValue)
	if multiplier != 0 {
		pso.CurrentPlayerData.ShiftaLvl = int16(1 + math.Round(((math.Abs(float64(multiplier))*100)-10)/1.3))
		if multiplier < 0 {
			pso.CurrentPlayerData.ShiftaLvl = -pso.CurrentPlayerData.ShiftaLvl
		}
	} else {
		pso.CurrentPlayerData.ShiftaLvl = 0
	}
	return nil
}

func (pso *PSO) getDebandLvl(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x278+12), 4)
	if !ok {
		return errors.New("Unable to getDebandLvl")
	}
	totalValue := uint32(buf[1])<<16 + uint32(buf[0])

	multiplier := math.Float32frombits(totalValue)
	if multiplier != 0 {
		pso.CurrentPlayerData.DebandLvl = int16(1 + math.Round(((math.Abs(float64(multiplier))*100)-10)/1.3))
		if multiplier < 0 {
			pso.CurrentPlayerData.DebandLvl = -pso.CurrentPlayerData.DebandLvl
		}
	} else {
		pso.CurrentPlayerData.DebandLvl = 0
	}
	return nil
}

func (pso *PSO) getInvincibilityFrames(playerAddress int) error {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(playerAddress+0x720), 4)
	if !ok {
		return errors.New("Unable to iFrames")
	}
	iFrames := uint32(buf[1])<<16 + uint32(buf[0])
	pso.CurrentPlayerData.InvincibilityFrames = iFrames
	return nil
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
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), uintptr(questDataPtr+0x18), 48)
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

func (pso *PSO) GetHP() string {
	return fmt.Sprintf("%v", pso.CurrentPlayerData.HP)
}

func (pso *PSO) GetMaxHP() string {
	return fmt.Sprintf("%v", pso.CurrentPlayerData.MaxHP)
}

func (pso *PSO) GetMonsterList() error {
	entityCount, err := pso.readU32(uintptr(0x00AAE164))
	if err != nil {
		return err
	}

	pso.GameState.MonsterCount = entityCount

	// for i := uint32(0); i < entityCount; i++ {

	// }

	return nil
}

func (pso *PSO) readU32(address uintptr) (uint32, error) {
	buf, _, ok := w32.ReadProcessMemory(w32.HANDLE(pso.handle), address, 4)
	if !ok {
		return 0, errors.New("Unable to getQuestDataPointer")
	}
	return uint32(buf[1])<<16 + uint32(buf[0]), nil
}
