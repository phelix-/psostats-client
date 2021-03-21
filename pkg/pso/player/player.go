package player

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"unicode/utf16"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/pkg/numbers"
)

type BasePlayerInfo struct {
	Name                string
	Guildcard           string
	Room                uint16
	KillCount           uint16
	ShiftaLvl           int16
	DebandLvl           int16
	MaxHP               uint16
	MaxTP               uint16
	HP                  uint16
	TP                  uint16
	Floor               uint16
	InvincibilityFrames uint32
	Class               string
	Meseta              uint32
}

func GetPlayerData(handle w32.HANDLE, playerAddress int) (BasePlayerInfo, error) {
	base := 0x028
	max := 0xE4E
	buf, _, ok := w32.ReadProcessMemory(handle, uintptr(playerAddress+base), uintptr((max-base)+4))
	if !ok {
		return BasePlayerInfo{}, errors.New("Unable to getPlayerData")
	}
	shiftaMultiplier := numbers.Float32FromU16(buf[(0x278-base)/2], buf[(0x27A-base)/2])
	debandMultiplier := numbers.Float32FromU16(buf[(0x278+12-base)/2], buf[(0x278+14-base)/2])
	class := (buf[(0x961-base)/2] & 0xF00) >> 8
	// log.Printf("%v/%v", read, uintptr(10+(0x3f0-base)/2))
	// for i := 0; i <= len(buf)-8; i += 8 {
	// 	log.Printf("0x%x | 0x%x - %x %x %x %x %x %x %x %x", playerAddress+base+(2*i), base+(2*i), buf[i], buf[i+1], buf[i+2], buf[i+3], buf[i+4], buf[i+5], buf[i+6], buf[i+7])
	// }

	// log.Printf("%v - %v/%v - %v/%v", len(buf), pso.CurrentPlayerData.HP, pso.CurrentPlayerData.MaxHP,
	// 	pso.CurrentPlayerData.TP, pso.CurrentPlayerData.MaxTP)
	name, err := getCharacterName(handle, playerAddress)
	if err != nil {
		return BasePlayerInfo{}, err
	}
	guildcard, err := getGuildCard(handle, playerAddress)

	return BasePlayerInfo{
		Name:                name,
		Guildcard:           guildcard,
		Room:                buf[(0x028-base)/2],
		KillCount:           buf[(0x11A-base)/2],
		ShiftaLvl:           getSDLvlFromMultiplier(shiftaMultiplier),
		DebandLvl:           getSDLvlFromMultiplier(debandMultiplier),
		MaxHP:               buf[(0x2BC-base)/2],
		MaxTP:               buf[(0x2BE-base)/2],
		HP:                  buf[(0x334-base)/2],
		TP:                  buf[(0x336-base)/2],
		Floor:               buf[(0x3F0-base)/2],
		InvincibilityFrames: numbers.Uint32FromU16(buf[(0x720-base)/2], buf[(0x722-base)/2]),
		Class:               getClass(class),
		Meseta:              numbers.Uint32FromU16(buf[(0xE4C-base)/2], buf[(0xE4E-base)/2]),
	}, nil
}

func getCharacterName(handle w32.HANDLE, playerAddress int) (string, error) {
	buf, _, ok := w32.ReadProcessMemory(handle, uintptr(playerAddress+0x428), 24)
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

func getGuildCard(handle w32.HANDLE, playerAddress int) (string, error) {
	buf, _, ok := w32.ReadProcessMemory(handle, uintptr(playerAddress+0x92a), 16)
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

func getClass(classBits uint16) string {
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

func getSDLvlFromMultiplier(multiplier float32) int16 {
	level := int16(0)
	if multiplier != 0 {
		level = int16(1 + math.Round(((math.Abs(float64(multiplier))*100)-10)/1.3))
		if multiplier < 0 {
			level = -level
		}
	}
	return level
}
