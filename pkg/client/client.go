package client

import (
	"fmt"
	"log"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/phelix-/psostats/v2/pkg/consoleui"
	"github.com/phelix-/psostats/v2/pkg/pso"
)

const (
	defaultTickRate    = time.Second / 5
	defaultUITickRate  = time.Second / 5
	defaultSIOTickRate = time.Second / 5
)

type Client struct {
	pso        *pso.PSO
	version    string
	tickRate   time.Duration
	uiTickRate time.Duration
	ui         *consoleui.ConsoleUI
	uiData     *consoleui.Data
	errChan    chan error
	done       chan struct{}
}

type DataFrame struct {
	Id         uint32
	Timestamp  time.Time
	Deband     int16
	Floor      uint16
	Hp         uint16
	Invincible bool
	Killcount  uint16
	Shifta     int16
	Tp         uint16
	Meseta     uint32
	Monsters   uint32
}

type Game struct {
	Id         uint32
	Timestamp  time.Time
	Quest      string
	Character  string
	Episode    uint16
	Difficulty uint16
	Gc         string
}

func New(version string) (*Client, error) {
	uiData := consoleui.Data{
		Connected: false,
		Status:    "Initializing",
		Version:   version,
	}

	ui, err := consoleui.New(&uiData)
	if err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	// p := widgets.NewParagraph()
	// p.Text = "Hello World!"
	// p.SetRect(0, 0, 25, 5)

	// ui.Render(p)

	// for e := range ui.PollEvents() {
	// 	if e.Type == ui.KeyboardEvent {
	// 		break
	// 	}
	// }

	pso := pso.New()

	return &Client{
		pso:        pso,
		version:    version,
		tickRate:   defaultTickRate,
		uiTickRate: defaultUITickRate,
		ui:         ui,
		uiData:     &uiData,
		errChan:    make(chan error),
		done:       make(chan struct{}),
	}, nil
}

func (c *Client) GetGameInfo() pso.StatsState {
	return c.pso.State
}

func (c *Client) GetFrames() map[int]pso.StatsFrame {
	return c.pso.Frames
}

func (c *Client) Run() error {
	c.ui.DrawScreen(&c.pso.CurrentPlayerData, &c.pso.GameState)
	defer c.ui.Close()

	go c.run()

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>", "<f10>":
				close(c.done)
				return nil
			}
		case err := <-c.errChan:
			close(c.done)
			return fmt.Errorf("Run: error returned on error channel: %w", err)
		}
	}
}

func (c *Client) run() {

	c.pso.StartPersistentConnection(c.errChan)
	go c.runDD()
	go c.runUI()
	// if !c.cfg.OfflineMode {
	// 	go c.runSIO()
	// }
}

func (c *Client) runDD() {
	// gameId := uint32(2)
	gameWritten := false
	gameFile, err := os.Create("game.json")
	if err != nil {
		log.Panic(err)
	}
	defer gameFile.Close()
	file, err := os.Create("frames.json")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	for {
		select {
		case <-time.After(c.tickRate):
			connected, statusString := c.pso.CheckConnection()
			c.uiData.Connected = connected
			c.uiData.Status = statusString
			if !connected {
				c.clearUIData()
				continue
			}

			c.populateUIData()
			if c.pso.CurrentPlayerData.QuestName != "No Active Quest" {

				if !gameWritten {
					// game := Game{
					// 	Id:         gameId,
					// 	Timestamp:  time.Now(),
					// 	Quest:      c.pso.CurrentPlayerData.QuestName,
					// 	Character:  c.pso.CurrentPlayerData.CharacterName,
					// 	Episode:    c.pso.CurrentPlayerData.Episode,
					// 	Difficulty: c.pso.CurrentPlayerData.Difficulty,
					// 	Gc:         c.pso.CurrentPlayerData.Guildcard,
					// }
					// bytes, err := json.Marshal(game)
					// if err != nil {
					// 	log.Panic(err)
					// }
					// gameFile.Write(bytes)
					// gameFile.WriteString("\n")
					// gameWritten = true
				}

				// frame := DataFrame{
				// 	Id:         gameId,
				// 	Timestamp:  time.Now(),
				// 	Deband:     c.pso.CurrentPlayerData.DebandLvl,
				// 	Floor:      c.pso.CurrentPlayerData.Floor,
				// 	Hp:         c.pso.CurrentPlayerData.HP,
				// 	Invincible: c.pso.CurrentPlayerData.InvincibilityFrames > 0,
				// 	Killcount:  c.pso.CurrentPlayerData.KillCount,
				// 	Shifta:     c.pso.CurrentPlayerData.ShiftaLvl,
				// 	Tp:         c.pso.CurrentPlayerData.TP,
				// 	Meseta:     c.pso.CurrentPlayerData.Meseta,
				// 	Monsters:   c.pso.GameState.MonsterCount,
				// }

				// bytes, err := json.Marshal(frame)
				// if err != nil {
				// 	log.Panic(err)
				// }
				// file.Write(bytes)
				// file.WriteString("\n")
			}

			// newStatus := c.pso.GetStatus()

			// when a new quest has started
			// if oldStatus != pso.StatusPlaying && newStatus == pso.StatusPlaying ||
			// 	oldStatus != pso.StatusOtherReplay && newStatus == pso.StatusOtherReplay ||
			// 	oldStatus != pso.StatusOwnReplayFromLeaderboard && newStatus == pso.StatusOwnReplayFromLeaderboard {
			// 	c.statsSent = false
			// }

			// if !c.cfg.OfflineMode {
			// 	if c.dd.GetStatsFinishedLoading() && !c.statsSent {
			// 		if newStatus == devildaggers.StatusDead || newStatus == devildaggers.StatusOtherReplay || newStatus == devildaggers.StatusOwnReplayFromLeaderboard {
			// 			// send stats
			// 			submitGameRequest, err := c.compileGameRequest()
			// 			if err != nil {
			// 				c.errChan <- fmt.Errorf("runGameCapture: could not compile game recording: %w", err)
			// 				return
			// 			}
			// 			gameID, err := c.grpcClient.SubmitGame(submitGameRequest)
			// 			if err != nil {
			// 				c.errChan <- fmt.Errorf("runGameCapture: error submitting game to server: %w", err)
			// 				return
			// 			}
			// 			c.lastSubmittedGameID = gameID
			// 			c.statsSent = true

			// 			if c.cfg.AutoClipboardGame {
			// 				c.copyGameURLToClipboard()
			// 			}

			// 			if (c.cfg.Submit.Stats && !c.dd.GetIsReplay()) ||
			// 				(c.cfg.Submit.ReplayStats && c.dd.GetIsReplay()) {
			// 				if (c.dd.GetLevelHashMD5() == c.v3SurvivalHash) ||
			// 					(!c.cfg.Submit.NonDefaultSpawnsets && c.dd.GetLevelHashMD5() != c.v3SurvivalHash) {
			// 					if c.sioClient.GetStatus() == socketio.StatusLoggedIn {
			// 						err = c.sioClient.SubmitGame(gameID, c.cfg.Discord.NotifyPlayerBest, c.cfg.Discord.NotifyAbove1100)
			// 						if err != nil {
			// 							c.errChan <- fmt.Errorf("runGameCapture: error submitting game to sio: %w", err)
			// 							return
			// 						}
			// 					}
			// 				}
			// 			}
			// 		}
			// 	}
			// }

			// oldStatus = newStatus
		case <-c.done:
			return
		}
	}
}

func (c *Client) runUI() {
	c.ui.ClearScreen()
	for {
		select {
		case <-time.After(c.tickRate):
			err := c.ui.DrawScreen(&c.pso.CurrentPlayerData, &c.pso.GameState)
			if err != nil {
				c.errChan <- fmt.Errorf("runUI: error drawing screen in ui: %w", err)
				return
			}
		case <-c.done:
			c.ui.ClearScreen()
			return
		}
	}
}

func (c *Client) clearUIData() {
	c.uiData.Clear()
}

func (c *Client) populateUIData() {

	// c.uiData.HP = c.pso.CurrentPlayerData.HP
	// c.uiData.MaxHP = c.pso.CurrentPlayerData.MaxHP
	// c.uiData.Status = c.pso.GetStatus()
	// c.uiData.OnlineStatus = c.sioClient.GetStatus()
	// c.uiData.PlayerName = c.pso.GetPlayerName()
	// if c.uiData.PlayerName == "" {
	// 	c.uiData.Status = consoleui.StatusConnecting
	// 	return
	// }
	// c.uiData.LastGameID = c.lastSubmittedGameID
	// status := c.pso.GetStatus()
	// if status == devildaggers.StatusPlaying || status == devildaggers.StatusOtherReplay || status == devildaggers.StatusOwnReplayFromLastRun || status == devildaggers.StatusOwnReplayFromLeaderboard {
	// 	c.uiData.Recording = consoleui.StatusRecording
	// 	if c.statsSent {
	// 		c.uiData.Recording = consoleui.StatusGameSubmitted
	// 	}
	// 	c.uiData.Timer = c.pso.GetTime()
	// 	c.uiData.DaggersHit = c.pso.GetDaggersHit()
	// 	c.uiData.DaggersFired = c.pso.GetDaggersFired()
	// 	c.uiData.Accuracy = c.pso.GetAccuracy()
	// 	c.uiData.GemsCollected = c.pso.GetGemsCollected()
	// 	c.uiData.Homing = c.pso.GetHomingDaggers()
	// 	c.uiData.EnemiesAlive = c.pso.GetEnemiesAlive()
	// 	c.uiData.EnemiesKilled = c.pso.GetKills()
	// 	c.uiData.TotalGems = c.pso.GetTotalGems()
	// 	c.uiData.GemsDespawned = c.pso.GetGemsDespawned()
	// 	c.uiData.GemsEaten = c.pso.GetGemsEaten()
	// 	c.uiData.DaggersEaten = c.pso.GetDaggersEaten()
	// } else {
	// 	c.uiData.Recording = consoleui.StatusNotRecording
	// 	if c.pso.GetStatus() == devildaggers.StatusDead {
	// 		if c.statsSent {
	// 			c.uiData.Recording = consoleui.StatusGameSubmitted
	// 		}
	// 		c.uiData.DeathType = c.pso.GetDeathType()
	// 	}
	// }
}
