package pso

import (
	"bytes"
	"encoding/binary"
	"errors"
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

func (pso *PSO) ConsolidateFrame() {
	pso.State.QuestName = pso.CurrentPlayerData.QuestName
	pso.State.QuestStarted = pso.GameState.QuestStarted
	pso.State.QuestStartTime = pso.GameState.QuestStartTime
	pso.State.QuestEndTime = pso.GameState.QuestEndTime

	if pso.State.QuestStarted && !pso.State.QuestComplete {
		currentSecond := int(time.Now().Sub(pso.State.QuestStartTime).Seconds())
		if _, exists := pso.Frames[currentSecond]; !exists {
			mesetaCharged := uint32(0)
			if previousFrame, previousFrameExists := pso.Frames[currentSecond-1]; previousFrameExists {
				mesetaCharged = (previousFrame.Meseta - pso.CurrentPlayerData.Meseta) + previousFrame.MesetaCharged
			}
			frame := StatsFrame{
				HP:            pso.CurrentPlayerData.HP,
				TP:            pso.CurrentPlayerData.TP,
				Floor:         pso.CurrentPlayerData.Floor,
				Room:          pso.CurrentPlayerData.Room,
				ShiftaLvl:     pso.CurrentPlayerData.ShiftaLvl,
				DebandLvl:     pso.CurrentPlayerData.DebandLvl,
				Invincible:    pso.CurrentPlayerData.InvincibilityFrames > 0,
				MonsterCount:  pso.GameState.MonsterCount,
				Meseta:        pso.CurrentPlayerData.Meseta,
				MesetaCharged: mesetaCharged,
				Time:          currentSecond,
			}
			pso.Frames[currentSecond] = frame
		}
	}

	pso.State.QuestComplete = pso.GameState.QuestComplete
	if pso.State.QuestComplete {
		pso.State.QuestDuration = pso.State.QuestEndTime.Sub(pso.State.QuestStartTime)
	} else if pso.State.QuestStarted {
		pso.State.QuestDuration = time.Now().Sub(pso.State.QuestStartTime)
	} else {
		pso.State.QuestDuration = time.Duration(0)
	}

}

func (pso *PSO) RefreshData() error {
	if !pso.connected {
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

	err = pso.GetMonsterList()
	if err != nil {
		return err
	}

	if address != 0 {
		err := pso.getPlayerData(address)
		if err != nil {
			return err
		}

		err = pso.getCharacterName(address)
		if err != nil {
			return err
		}

		err = pso.getGuildCard(address)
		if err != nil {
			return err
		}

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
	class := buf[(0x961-base)/2] & 0xF
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
	pso.CurrentPlayerData.Difficulty = difficulty
	return BaseGameInfo{
		episode:    episode,
		difficulty: difficulty,
	}, nil
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
