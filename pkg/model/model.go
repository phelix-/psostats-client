package model

import (
	"time"
)

type BasePlayerInfo struct {
	Name      string
	GuildCard string
	Level     uint16
	Class     string
}

type QuestRun struct {
	Server              string
	PlayerName          string
	PlayerClass         string
	GuildCard           string
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
	DeathCount          int
	HP                  []uint16
	TP                  []uint16
	Meseta              []uint32
	MesetaCharged       []int
	Room                []uint16
	IllegalShifta       bool
	PbCategory          bool
	ShiftaLvl           []int16
	DebandLvl           []int16
	Invincible          []bool
	Events              []Event
	Monsters            map[int]Monster
	Bosses              map[string]BossData
	MonsterCount        []int
	MonstersKilledCount []int
	MonsterHpPool       []int
	MonstersDead        int
	WeaponsUsed         map[string]string
	EquipmentUsedTime   map[string]map[string]int
	FreezeTraps         []uint16
	FTUsed              uint16
	DTUsed              uint16
	CTUsed              uint16
	TPUsed              uint16
}

type Monster struct {
	Name       string
	Id         uint16
	UnitxtId   uint32
	SpawnTime  time.Time
	KilledTime time.Time
	Alive      bool
	Frame1     bool
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
	Player           string
	Players          map[string]string
	Category         string
	Episode          int
	Quest            string
	QuestAndCategory string
	Time             time.Duration
	Timestamp        time.Time
	FormattedTime    string
	FormattedDate    string
	GameGzip         []byte
}
