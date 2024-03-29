package model

import (
	"time"
)

const (
	WeaponBareHanded     = "Bare Handed"
	EquipmentTypeWeapon  = "Weapon"
	EquipmentTypeFrame   = "Frame"
	EquipmentTypeBarrier = "Barrier"
	EquipmentTypeUnit    = "Unit"
	EquipmentTypeMag     = "Mag"
)

type BasePlayerInfo struct {
	Name      string
	GuildCard string
	SectionId uint8
	Level     uint16
	Class     string
}

type QuestRun struct {
	Client              ClientInfo
	Server              string
	PlayerName          string
	PlayerClass         string
	GuildCard           string
	UserName            string
	AllPlayers          []BasePlayerInfo
	Id                  string
	Difficulty          string
	Episode             uint16
	QuestName           string
	QuestComplete       bool
	QuestStartTime      time.Time
	QuestStartDate      string
	QuestEndTime        time.Time
	QuestDuration       string
	SubmittedTime       time.Time
	DeathCount          int
	HP                  []uint16
	TP                  []uint16
	PB                  []float32
	MesetaCharged       []int
	Room                []uint16
	IllegalShifta       bool
	PbCategory          bool
	ShiftaLvl           []int16
	DebandLvl           []int16
	Invincible          []bool
	Events              []Event
	FastWarps           bool
	Splits              []QuestRunSplit
	Monsters            map[int]Monster
	PlayerDamage        map[int]int64
	LastHits            map[int]int
	Bosses              map[string]BossData
	MonsterCount        []int
	MonstersKilledCount []int
	MonsterHpPool       []int
	DamageDealt         []int64
	MonstersDead        int
	Weapons             map[string]Equipment
	WeaponsUsed         map[string]string
	EquipmentUsedTime   map[string]map[string]int
	FreezeTraps         []uint16
	FTUsed              uint16
	DTUsed              uint16
	CTUsed              uint16
	TPUsed              uint16
	MonomateUsed        uint8
	DimateUsed          uint8
	TrimateUsed         uint8
	MonofluidUsed       uint8
	DifluidUsed         uint8
	TrifluidUsed        uint8
	MoonAtomizerUsed    uint8
	StarAtomizerUsed    uint8
	SolAtomizerUsed     uint8
	TelepipeUsed        uint8
	TimeByState         map[int]uint64
	TechsCast           map[string]int
	Points              uint16
	DataFrames          []DataFrame
}

type QuestRunSplit struct {
	Name        string
	Index       int
	StartSecond int
	Start       time.Time
	End         time.Time
}

type Monster struct {
	Name       string
	Id         uint16
	UnitxtId   uint32
	SpawnTime  time.Time
	KilledTime time.Time
	Alive      bool
	Frame1     bool
	Location   Location
}

type Location struct {
	Floor   uint16
	Room    uint16
	X       float32
	Y       float32
	Z       float32
	Warping bool
}

type DataFrame struct {
	HP                 uint16
	TP                 uint16
	PB                 float32
	MesetaCharged      int
	Time               int64
	Map                uint16
	MapVariation       uint16
	ShiftaLvl          int16
	DebandLvl          int16
	Invincible         bool
	FT                 uint16
	DT                 uint16
	CT                 uint16
	MonstersAlive      int
	DamageDealt        int64
	Kills              int
	State              uint16
	Weapon             string
	PlayerByGcLocation map[string]Location
	PlayerLocation     map[int]Location
	MonsterLocation    map[int]Location
}

type BossData struct {
	Name       string
	Id         uint16
	UnitxtId   uint32
	SpawnTime  time.Time
	KilledTime time.Time
	FirstFrame int
	Hp         []int
}

type Event struct {
	Second      int
	Description string
}

type Game struct {
	Id               string
	IdInt            int
	Player           string
	PlayerNames      []string
	PlayerClasses    []string
	PlayerGcs        []string
	Category         string
	Episode          int
	Quest            string
	QuestAndCategory string
	Time             time.Duration
	Timestamp        time.Time
	FormattedTime    string
	FormattedDate    string
	GameGzip         []byte
	P1Gzip           []byte
	P1HasStats       bool
	P1Video          string
	P2Gzip           []byte
	P2HasStats       bool
	P2Video          string
	P3Gzip           []byte
	P3HasStats       bool
	P3Video          string
	P4Gzip           []byte
	P4HasStats       bool
	P4Video          string
	Points           int
}

type FormattedPlayerInfo struct {
	Name      string
	GuildCard string
	HasPov    bool
	Class     string
}

type FormattedGame struct {
	Id           string
	Players      []FormattedPlayerInfo
	PbRun        bool
	NumPlayers   int
	Episode      int
	Quest        string
	Duration     time.Duration
	Time         string
	Date         string
	RelativeDate string
	Pb           bool
	Record       bool
}

type ClientInfo struct {
	VersionMajor int
	VersionMinor int
	VersionPatch int
}

type MessageOfTheDay struct {
	Authorized bool
	Message    string
}

type PostGameResponse struct {
	Pb     bool
	Record bool
	Id     string
}

type Equipment struct {
	Id              string
	UnitxtIndex     string
	Type            string
	Display         string
	SecondsEquipped int
	Attacks         int
	Techs           int
}
