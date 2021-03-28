package player

import (
	"errors"
	"log"
	"math"
	"strings"
	"unicode/utf16"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/pkg/numbers"
)

type BasePlayerInfo struct {
	Name                string
	Guildcard           string
	Room                uint16
	ShiftaLvl           int16
	DebandLvl           int16
	MaxHP               uint16
	MaxTP               uint16
	HP                  uint16
	TP                  uint16
	PB                  float32
	Floor               uint16
	InvincibilityFrames uint32
	DamageTraps         uint16
	FreezeTraps         uint16
	ConfuseTraps        uint16
	Class               string
	Meseta              uint32
}

func ParsePlayerMemory(buf []uint16, base int) BasePlayerInfo {
	shiftaMultiplier := numbers.Float32FromU16(buf[(0x278-base)/2], buf[(0x27A-base)/2])
	debandMultiplier := numbers.Float32FromU16(buf[(0x278+12-base)/2], buf[(0x278+14-base)/2])
	class := (buf[(0x961-base)/2] & 0xF00) >> 8

	return BasePlayerInfo{
		Room:                buf[(0x028-base)/2],
		ShiftaLvl:           getSDLvlFromMultiplier(shiftaMultiplier),
		DebandLvl:           getSDLvlFromMultiplier(debandMultiplier),
		MaxHP:               buf[(0x2BC-base)/2],
		MaxTP:               buf[(0x2BE-base)/2],
		HP:                  buf[(0x334-base)/2],
		TP:                  buf[(0x336-base)/2],
		Floor:               buf[(0x3F0-base)/2],
		PB:                  numbers.Float32FromU16(buf[(0x520-base)/2], buf[(0x522-base)/2]),
		InvincibilityFrames: numbers.Uint32FromU16(buf[(0x720-base)/2], buf[(0x722-base)/2]),
		DamageTraps:         buf[(0x89D-base)/2] & 0x00FF,
		FreezeTraps:         (buf[(0x89D-base)/2] & 0xFF00) >> 8,
		ConfuseTraps:        (buf[(0x89F-base)/2] & 0xFF00) >> 8,
		Class:               getClass(class),
		Meseta:              numbers.Uint32FromU16(buf[(0xE4C-base)/2], buf[(0xE4E-base)/2]),
	}
}

func GetPlayerData(handle w32.HANDLE, playerAddress int) (BasePlayerInfo, error) {
	base := 0x028
	max := 0xE4E
	buf, _, ok := w32.ReadProcessMemory(handle, uintptr(playerAddress+base), uintptr((max-base)+4))
	if !ok {
		return BasePlayerInfo{}, errors.New("Unable to getPlayerData")
	}
	// log.Printf("%v/%v", read, uintptr(10+(0x3f0-base)/2))
	// log.Printf("------------")
	// for i := 0; i <= len(buf)-8; i += 8 {
	// 	log.Printf("0x%x, 0x%x, 0x%x, 0x%x, 0x%x, 0x%x, 0x%x, 0x%x", buf[i], buf[i+1], buf[i+2], buf[i+3], buf[i+4], buf[i+5], buf[i+6], buf[i+7])
	// }
	basePlayerInfo := ParsePlayerMemory(buf, base)
	name, err := getCharacterName(handle, playerAddress)
	if err != nil {
		return BasePlayerInfo{}, err
	}
	basePlayerInfo.Name = name

	guildcard, err := getGuildCard(handle, playerAddress)
	if err != nil {
		return BasePlayerInfo{}, err
	}
	basePlayerInfo.Guildcard = guildcard
	return basePlayerInfo, nil

	// log.Printf("%v - %v/%v - %v/%v", len(buf), pso.CurrentPlayerData.HP, pso.CurrentPlayerData.MaxHP,
	// 	pso.CurrentPlayerData.TP, pso.CurrentPlayerData.MaxTP)
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

	name := string(utf16.Decode(buf[0:endIndex]))
	name = strings.TrimPrefix(name, "\tE")
	return name, nil
}

func getGuildCard(handle w32.HANDLE, playerAddress int) (string, error) {
	// return numbers.ReadString(handle, uintptr(playerAddress+0x92a), 8)
	guildcard, err := numbers.ReadString(handle, uintptr(playerAddress+0x92c), 7)
	if err != nil {
		return "", err
	}
	// guildcard = strings.TrimPrefix(guildcard, "A")
	// guildcard = strings.TrimPrefix(guildcard, "B")
	// guildcard = strings.TrimPrefix(guildcard, "C")
	guildcard = strings.Trim(guildcard, "\u0000")
	guildcard = strings.TrimSpace(guildcard)
	index := strings.Index(guildcard, "\u0000")
	if index > 0 {
		guildcard = strings.Split(guildcard, "\u0000")[0]
	}
	return guildcard, nil
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

func (p *BasePlayerInfo) IsLowered() bool {
	maxHp := int(p.MaxHP)
	return int(p.HP) < ((maxHp * 95) / 100)
}

func (p *BasePlayerInfo) MaxSupplyableShifta() int16 {
	switch p.Class {
	case "FOmar":
		fallthrough
	case "FOmarl":
		fallthrough
	case "FOnewm":
		fallthrough
	case "FOnewearl":
		return 30
	case "HUnewearl":
		fallthrough
	case "RAmarl":
		return 20
	case "RAmar":
		return 15
	case "HUmar":
		fallthrough
	case "HUcast":
		fallthrough
	case "HUcaseal":
		return 3
	case "RAcast":
		fallthrough
	case "RAcaseal":
		return 0
	default:
		log.Printf("Unrecongnized class '%v'", p.Class)
		return 0
	}
}
