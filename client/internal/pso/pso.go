package pso

import (
	"fmt"
	"github.com/phelix-/psostats/v2/client/internal/pso/inventory"
	"github.com/phelix-/psostats/v2/client/internal/pso/quest"
	"log"
	"syscall"
	"time"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/client/internal/pso/player"
)

const (
	unseenWindowName             = "PHANTASY STAR ONLINE Blue Burst"
	ephineaWindowName            = "Ephinea: Phantasy Star Online Blue Burst"
	persistentConnectionTickRate = time.Second / 30
	windowsCodeStillActive       = 259
)

type PSO struct {
	completeGame      chan QuestRun
	questTypes        quest.Quests
	connected         bool
	connectedStatus   string
	server            string
	handle            w32.HANDLE
	CurrentPlayerData player.BasePlayerInfo
	Equipment         []inventory.Equipment
	GameState         GameState
	CurrentQuest      QuestRun
	errors            chan error
	done              chan struct{}
	MonsterNames      map[uint32]string
}

type GameState struct {
	MonsterCount         int
	QuestName            string
	AllowQuestStart      bool // Guards against starting the client mid-quest
	QuestStarted         bool
	Uploading            bool
	UploadSuccessful     bool
	QuestComplete        bool
	QuestStartTime       time.Time
	QuestEndTime         time.Time
	CmodeStage           int
	RngSeed              uint32
	Difficulty           string
	Episode              uint16
	Map                  uint16
	Floor                uint16
	questPointer         uintptr
	questRegisterPointer uintptr
	CurrentSplit         QuestRunSplit
}

func (state *GameState) ClearQuest() {
	state.MonsterCount = 0
	state.QuestStarted = false
	state.QuestComplete = false
	state.Uploading = false
	state.UploadSuccessful = false
	state.CmodeStage = -1
	state.QuestName = "No Active Quest"
	state.questRegisterPointer = 0
	state.questPointer = 0
	state.CurrentSplit = QuestRunSplit{}
}

func (state *GameState) Clear() {
	state.Difficulty = "Normal"
	state.Episode = 1
	state.ClearQuest()
}

type PlayerData struct {
	CharacterName       string
	Class               string
	Guildcard           string
	HP                  uint16
	MaxHP               uint16
	TP                  uint16
	MaxTP               uint16
	Floor               uint16
	Room                uint16
	Meseta              uint32
	ShiftaLvl           int16
	DebandLvl           int16
	InvincibilityFrames uint32
	Time                time.Time
}

func New(completeGameChannel chan QuestRun) *PSO {
	return &PSO{
		completeGame: completeGameChannel,
		questTypes:   quest.NewQuests(),
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
	server := "unseen"
	hwnd := w32.FindWindowW(nil, syscall.StringToUTF16Ptr(unseenWindowName))
	if hwnd == 0 {
		server = "ephinea"
		// unseen not found
		hwnd = w32.FindWindowW(nil, syscall.StringToUTF16Ptr(ephineaWindowName))
		if hwnd == 0 {
			return false, "Window not found", nil
		}
	}

	_, pid := w32.GetWindowThreadProcessId(hwnd)
	handle, err := w32.OpenProcess(w32.PROCESS_ALL_ACCESS, false, uintptr(pid))
	if err != nil {
		return false, "Could not open process", fmt.Errorf("Connect: could not open process with pid %v: %w", pid, err)
	}

	pso.handle = handle
	pso.server = server

	return true, fmt.Sprintf("Connected to pid %v", pid), nil
}

func (pso *PSO) Close() {
	w32.CloseHandle(pso.handle)
}

func (pso *PSO) CheckConnection() (bool, string) {
	return pso.connected, pso.connectedStatus
}

func (pso *PSO) checkConnection() bool {
	code, err := w32.GetExitCodeProcess(pso.handle)
	return err == nil && code == windowsCodeStillActive
}
