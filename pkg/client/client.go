package client

import (
	"encoding/json"
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
	pso           *pso.PSO
	version       string
	tickRate      time.Duration
	uiTickRate    time.Duration
	ui            *consoleui.ConsoleUI
	uiData        *consoleui.Data
	currentGameId int
	errChan       chan error
	done          chan struct{}
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

func (c *Client) GetGameInfo() pso.QuestRun {
	return c.pso.Quests[c.pso.CurrentQuest]
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
			case "w":
				filename := fmt.Sprintf("./game-%v.json", time.Now().Format("2006_01_02-1504"))
				file, err := os.Create(filename)
				if err != nil {
					log.Printf("Unable to write to %v, %v", filename, err)
				}
				defer file.Close()
				json, err := json.Marshal(c.pso.Quests[c.pso.CurrentQuest])
				if err != nil {
					log.Printf("Unable to generate json")
				}
				file.Write(json)
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

func (c *Client) GetNextGameId() string {
	c.currentGameId++
	return fmt.Sprint(c.currentGameId)
}

func (c *Client) runDD() {
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
