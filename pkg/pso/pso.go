package pso

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"

	"github.com/TheTitanrain/w32"
)

const (
	unseenWindowName             = "PHANTASY STAR ONLINE Blue Burst"
	ephineaWindowName            = "Ephinea: Phantasy Star Online Blue Burst"
	persistentConnectionTickRate = time.Second / 30
	windowsCodeStillActive       = 259
	WarpIn                       = 0
	Switch                       = 1
)

type (
	handle w32.HANDLE
)

type PSO struct {
	questTypes      Quests
	connected       bool
	connectedStatus string
	handle          handle
	baseAddress     uintptr
	// ddstatsBlockAddress address
	CurrentPlayerData PlayerData
	GameState         GameState
	CurrentQuest      int
	Quests            map[int]QuestRun
	Frames            map[int]StatsFrame
	errors            chan error
	done              chan struct{}
	MonsterNames      map[uint32]string
}

type StatsFrame struct {
	HP            uint16
	TP            uint16
	Floor         uint16
	Room          uint16
	ShiftaLvl     int16
	DebandLvl     int16
	Invincible    bool
	MonsterCount  int
	Meseta        uint32
	MesetaCharged uint32
	Time          int
}

type GameState struct {
	MonsterCount      int
	FloorSwitches     bool
	QuestStarted      bool
	QuestComplete     bool
	QuestStartTime    time.Time
	QuestEndTime      time.Time
	monsterUnitxtAddr uint32
}

type PlayerData struct {
	CharacterName       string
	Class               string
	Guildcard           string
	HP                  uint16
	MaxHP               uint16
	TP                  uint16
	MaxTP               uint16
	Difficulty          string
	Episode             uint16
	Floor               uint16
	Room                uint16
	QuestName           string
	KillCount           uint16
	Meseta              uint32
	ShiftaLvl           int16
	DebandLvl           int16
	InvincibilityFrames uint32
	Time                time.Time
}

func New() *PSO {
	return &PSO{
		questTypes:   NewQuests(),
		Quests:       make(map[int]QuestRun),
		Frames:       make(map[int]StatsFrame),
		MonsterNames: make(map[uint32]string),
	}
}

func (pso *PSO) StartPersistentConnection(errors chan error) {
	if pso.done != nil {
		close(pso.done)
	}
	pso.done = make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(persistentConnectionTickRate):
				if !pso.connected {
					connected, connectedStatus, err := pso.Connect()
					pso.connectedStatus = connectedStatus
					pso.connected = connected
					if err != nil {
						errors <- fmt.Errorf("StartPersistentConnection: could not connect to pso: %w", err)
						continue
					}
					if !pso.connected {
						continue
					}
				}
				pso.connected = pso.checkConnection()
				if pso.connected {
					err := pso.RefreshData()
					if err != nil {
						log.Fatal(err)
						errors <- fmt.Errorf("StartPersistentConnection: could not refresh data: %w", err)
						continue
					}
				}
			case <-pso.done:
				pso.Close()
				pso.connected = false
				return
			}
		}
	}()
}

func (pso *PSO) StopPersistentConnection() {
	if pso.done != nil {
		close(pso.done)
	}
}

func (pso *PSO) Connect() (bool, string, error) {
	hwnd := w32.FindWindowW(nil, syscall.StringToUTF16Ptr(unseenWindowName))
	if hwnd == 0 {
		// unseen not found
		hwnd = w32.FindWindowW(nil, syscall.StringToUTF16Ptr(ephineaWindowName))
		if hwnd == 0 {
			return false, "Window not found", nil
		}
	}

	_, pid := w32.GetWindowThreadProcessId(hwnd)
	hndl, err := w32.OpenProcess(w32.PROCESS_ALL_ACCESS, false, uintptr(pid))
	if err != nil {
		return false, "Could not open process", fmt.Errorf("Connect: could not open process with pid %v: %w", pid, err)
	}

	baseAddress, err := getBaseAddress(pid)
	if err != nil {
		return false, "Could not find base address", fmt.Errorf("Connect: could get base address: %w", err)
	}

	pso.handle = handle(hndl)
	pso.baseAddress = baseAddress

	// ddstatsBlockAddress, err := pso.getDevilDaggersBlockBaseAddress()
	// if err != nil {
	// 	pso.connected = false
	// 	return false, fmt.Errorf("Connect: could get ddstats block address: %w", err)
	// }

	// pso.ddstatsBlockAddress = ddstatsBlockAddress

	return true, fmt.Sprintf("Connected to pid %v", pid), nil
}

func (pso *PSO) Close() {
	w32.CloseHandle(w32.HANDLE(pso.handle))
}

func (pso *PSO) CheckConnection() (bool, string) {
	return pso.connected, pso.connectedStatus
}

func (pso *PSO) checkConnection() bool {
	code, err := w32.GetExitCodeProcess(w32.HANDLE(pso.handle))
	return err == nil && code == windowsCodeStillActive
}

func getBaseAddress(pid int) (uintptr, error) {
	var baseAddress uintptr

	snapshot := w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPMODULE|w32.TH32CS_SNAPMODULE32, uint32(pid))
	if snapshot != w32.ERROR_INVALID_HANDLE {
		var me w32.MODULEENTRY32
		me.Size = uint32(unsafe.Sizeof(me))
		if w32.Module32First(snapshot, &me) {
			baseAddress = uintptr(unsafe.Pointer(me.ModBaseAddr))
		}
	}
	defer w32.CloseHandle(snapshot)

	if baseAddress == 0 {
		return 0, fmt.Errorf("getBaseAddress: could not find base address for PID %d", pid)
	}

	return baseAddress, nil
}
