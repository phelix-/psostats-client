// Foundation of the psostats client
package client

import (
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
	if c.config.HostLocalUi != nil && *c.config.HostLocalUi {
		go c.runHttp()
	} else {
		log.Printf("Local UI disabled")
	}

	uiEvents := termui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>", "<f10>":
				close(c.done)
				return nil
			case "w":
				c.writeGameJson()
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

func (c *Client) runHttp() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/game/info", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(c.GetGameInfo())
		if err != nil {
			r.Response.StatusCode = 500
			fmt.Fprintf(w, "Error!")
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(bytes)
	})
	addr := fmt.Sprintf(":%v", c.config.GetUiPort())
	log.Printf("Hosting local ui at localhost%v", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
