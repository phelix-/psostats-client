package player

import (
	"errors"
	"github.com/phelix-/psostats/v2/pkg/model"
	"github.com/phelix-/psostats/v2/pkg/psoclasses"
	"math"
	"strings"
	"unicode/utf16"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/client/internal/numbers"

	constants "github.com/phelix-/psostats/v2/client/internal/pso/constants"
)

const (
	ephineaGuildCardOffset = 0x4
	ephineaGuildCardSize   = 4
)

type BasePlayerInfo struct {
	Name                string
	GuildCard           string
	SectionId           uint8
	Room                uint16
	ShiftaLvl           int16
	DebandLvl           int16
	Level               uint16
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
	actualClass         psoclasses.PsoClass
	Meseta              uint32
	Warping             bool
	AccountMode         constants.EphineaAccountMode
	ActionState         uint16
	currentTech         uint16
	Location            model.Location
}

func (p BasePlayerInfo) GetCurrentTech() string {
	switch p.currentTech {
	case 0x0000:
		return "Foie"
	case 0x0006:
		return "Zonde"
	case 0x0003:
		return "Barta"
	case 0x0001:
		return "Gifoie"
	case 0x0007:
		return "Gizonde"
	case 0x0004:
		return "Gibarta"
	case 0x0002:
		return "Rafoie"
	case 0x0008:
		return "Razonde"
	case 0x0005:
		return "Rabarta"
	case 0x0009:
		return "Grants"
	case 0x0012:
		return "Megid"
	case 0x000F:
		return "Resta"
	case 0x0010:
		return "Anti"
	case 0x0011:
		return "Reverser"
	case 0x000D:
		return "Shifta"
	case 0x000A:
		return "Deband"
	case 0x000B:
		return "Jellen"
	case 0x000C:
		return "Zalure"
	case 0x000E:
		return "Ryuker"
	default:
		return "Unknown"
	}
}

func ParsePlayerMemory(buf []uint16, base uintptr) BasePlayerInfo {
	shiftaMultiplier := numbers.Float32FromU16(buf[(0x278-base)/2], buf[(0x27A-base)/2])
	debandMultiplier := numbers.Float32FromU16(buf[(0x278+12-base)/2], buf[(0x278+14-base)/2])
	class := (buf[(0x961-base)/2] & 0xF00) >> 8

	stateBitfield := buf[(0x33E-base)/2]
	playerWarping := stateBitfield&0x04 > 0

	floor := buf[(0x3F0-base)/2]
	room := buf[(0x028-base)/2]
	playerClass := getClass(class)
	return BasePlayerInfo{
		Room:                room,
		ShiftaLvl:           getSDLvlFromMultiplier(shiftaMultiplier),
		DebandLvl:           getSDLvlFromMultiplier(debandMultiplier),
		SectionId:           uint8(buf[(0x960-base)/2]),
		Level:               buf[(0xE44-base)/2] + 1,
		MaxHP:               buf[(0x2BC-base)/2],
		MaxTP:               buf[(0x2BE-base)/2],
		HP:                  buf[(0x334-base)/2],
		TP:                  buf[(0x336-base)/2],
		Floor:               floor,
		Warping:             playerWarping,
		PB:                  numbers.Float32FromU16(buf[(0x520-base)/2], buf[(0x522-base)/2]),
		InvincibilityFrames: numbers.Uint32FromU16(buf[(0x720-base)/2], buf[(0x722-base)/2]),
		DamageTraps:         buf[(0x89D-base)/2] & 0x00FF,
		FreezeTraps:         (buf[(0x89D-base)/2] & 0xFF00) >> 8,
		ConfuseTraps:        (buf[(0x89F-base)/2] & 0xFF00) >> 8,
		Class:               playerClass.Name,
		actualClass:         playerClass,
		Meseta:              numbers.Uint32FromU16(buf[(0xE4C-base)/2], buf[(0xE4E-base)/2]),
		ActionState:         buf[(0x348-base)/2],
		currentTech:         buf[(0x464-base)/2],
		Location: model.Location{
			Floor:   floor,
			Room:    room,
			Warping: playerWarping,
			Facing:  buf[(0x60-base)/2],
			X:       numbers.Float32FromU16(buf[(0x38-base)/2], buf[(0x3A-base)/2]),
			Y:       numbers.Float32FromU16(buf[(0x3C-base)/2], buf[(0x3E-base)/2]),
			Z:       numbers.Float32FromU16(buf[(0x40-base)/2], buf[(0x42-base)/2]),
		},
	}
}

func GetPlayerData(handle w32.HANDLE, playerAddress uintptr, server string) (BasePlayerInfo, error) {
	base := uintptr(0x028)
	max := uintptr(0xE4E)
	buf, _, ok := w32.ReadProcessMemory(handle, playerAddress+base, (max-base)+4)
	if !ok {
		return BasePlayerInfo{}, errors.New("unable to getPlayerData")
	}
	basePlayerInfo := ParsePlayerMemory(buf, base)
	name, err := getCharacterName(handle, playerAddress)
	if err != nil {
		return BasePlayerInfo{}, err
	}
	basePlayerInfo.Name = name

	guildcard, err := getGuildCard(handle, playerAddress, server)
	if err != nil {
		return BasePlayerInfo{}, err
	}
	basePlayerInfo.GuildCard = guildcard

	basePlayerInfo.AccountMode = constants.Normal
	if server == constants.EphineaServerName {
		mode, err := getEphineaAccountMode(handle, playerAddress)
		if err != nil {
			return BasePlayerInfo{}, err
		}
		basePlayerInfo.AccountMode = mode
	} else if server == constants.UnseenServerName {
		basePlayerInfo.AccountMode = constants.Sandbox
	}

	return basePlayerInfo, nil
}

func getCharacterName(handle w32.HANDLE, playerAddress uintptr) (string, error) {
	buf, _, ok := w32.ReadProcessMemory(handle, playerAddress+0x428, 24)
	if !ok {
		return "", errors.New("unable to getCharacterName")
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

func getGuildCard(handle w32.HANDLE, playerAddress uintptr, server string) (string, error) {
	offset := 0
	size := 7
	if server == constants.EphineaServerName {
		offset = ephineaGuildCardOffset
		size = ephineaGuildCardSize
	}

	guildCard, err := numbers.ReadString(handle, playerAddress+0x92c+uintptr(offset), size)
	if err != nil {
		return "", err
	}
	guildCard = strings.Trim(guildCard, "\u0000")
	guildCard = strings.TrimSpace(guildCard)
	index := strings.Index(guildCard, "\u0000")
	if index > 0 {
		guildCard = strings.Split(guildCard, "\u0000")[0]
	}
	return guildCard, nil
}

func getEphineaAccountMode(handle w32.HANDLE, playerAddress uintptr) (constants.EphineaAccountMode, error) {
	buf, _, ok := w32.ReadProcessMemory(handle, playerAddress+0x948, 3)
	if !ok {
		return constants.Normal, errors.New("unable to getEphineaAccountMode")
	}

	// TODO: Add support for hardcore
	//red := buf[1] & 0xFF
	//green := (buf[0] & 0xFF00) >> 8
	blue := buf[0] & 0xFF

	mode := constants.Normal

	if blue == 0x23 {
		mode = constants.Sandbox
	}

	return mode, nil
}

func getClass(classBits uint16) psoclasses.PsoClass {
	class := psoclasses.HUmar
	switch classBits {
	case 0x00:
		class = psoclasses.HUmar
		break
	case 0x01:
		class = psoclasses.HUnewearl
		break
	case 0x02:
		class = psoclasses.HUcast
		break
	case 0x09:
		class = psoclasses.HUcaseal
		break
	case 0x03:
		class = psoclasses.RAmar
		break
	case 0x0B:
		class = psoclasses.RAmarl
		break
	case 0x04:
		class = psoclasses.RAcast
		break
	case 0x05:
		class = psoclasses.RAcaseal
		break
	case 0x0A:
		class = psoclasses.FOmar
		break
	case 0x06:
		class = psoclasses.FOmarl
		break
	case 0x07:
		class = psoclasses.FOnewm
		break
	case 0x08:
		class = psoclasses.FOnewearl
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
	return int16(p.actualClass.MaxShifta)
}
