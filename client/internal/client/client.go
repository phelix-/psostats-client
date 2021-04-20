// Foundation of the psostats client
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	termui "github.com/gizak/termui/v3"
	"github.com/phelix-/psostats/v2/client/internal/client/config"
	"github.com/phelix-/psostats/v2/client/internal/consoleui"
	"github.com/phelix-/psostats/v2/client/internal/pso"
)

type Client struct {
	pso           *pso.PSO
	version       string
	config        *config.Config
	uiRefreshRate time.Duration
	ui            *consoleui.ConsoleUI
	currentGameId int
	errChan       chan error
	done          chan struct{}
}

func New(version string) (*Client, error) {
	ui, err := consoleui.New(version)
	if err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	pso := pso.New()
	clientConfig := config.ReadFromFile("./config.yaml")

	return &Client{
		pso:           pso,
		version:       version,
		config:        clientConfig,
		uiRefreshRate: clientConfig.GetUiRefreshRate(),
		ui:            ui,
		errChan:       make(chan error),
		done:          make(chan struct{}),
	}, nil
}

func (c *Client) GetGameInfo() pso.QuestRun {
	return c.pso.Quests[c.pso.CurrentQuest]
}

func (c *Client) Run() error {
	defer c.ui.Close()

	c.pso.StartPersistentConnection(c.errChan)
	go c.runUI()

	uiEvents := termui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q":
				close(c.done)
				return nil
			case "w":
				c.writeGameJson()
			case "u":
				c.uploadGame()
			}
		case err := <-c.errChan:
			close(c.done)
			return fmt.Errorf("run: error returned on error channel %w", err)
		}
	}
}

func (c *Client) GetNextGameId() string {
	c.currentGameId++
	return fmt.Sprint(c.currentGameId)
}

func (c *Client) writeGameJson() {
	filename := fmt.Sprintf("./game-%v.json", time.Now().Format("2006_01_02-1504"))
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Unable to write to %v, %v", filename, err)
	}
	defer file.Close()
	jsonBytes, err := json.Marshal(c.pso.Quests[c.pso.CurrentQuest])
	if err != nil {
		log.Printf("Unable to generate json")
	}
	file.Write(jsonBytes)
}

func (c *Client) uploadGame() {
	jsonBytes, err := json.Marshal(c.pso.Quests[c.pso.CurrentQuest])
	if err != nil {
		log.Printf("Unable to generate json")
	}
	buf := bytes.NewBuffer(jsonBytes)
	response, err := http.Post(c.config.GetServerBaseUrl()+"/api/game", "application/json", buf)
	if err != nil {
		log.Printf("Unable to upload game %v", err)
	}
	if response.StatusCode != 200 {
		log.Printf("Got response status %v: %v", response.StatusCode, response.Body)
	}
}

func (c *Client) runUI() {
	c.ui.ClearScreen()
	for {
		select {
		case <-time.After(c.uiRefreshRate):
			connected, statusString := c.pso.CheckConnection()
			c.ui.SetConnectionStatus(connected, statusString)

			currentQuest := c.pso.Quests[c.pso.CurrentQuest]
			err := c.ui.DrawScreen(&c.pso.CurrentPlayerData, &c.pso.GameState, &currentQuest)
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
