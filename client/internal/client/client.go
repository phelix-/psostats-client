// Foundation of the psostats client
package client

import (
	"fmt"
	"log"
	"time"

	"github.com/phelix-/psostats/v2/pkg/model"

	termui "github.com/gizak/termui/v3"
	"github.com/phelix-/psostats/v2/client/internal/client/config"
	"github.com/phelix-/psostats/v2/client/internal/consoleui"
	"github.com/phelix-/psostats/v2/client/internal/pso"
)

type Client struct {
	pso           *pso.PSO
	clientInfo    model.ClientInfo
	config        *config.Config
	uiRefreshRate time.Duration
	ui            *consoleui.ConsoleUI
	errChan       chan error
	done          chan struct{}
}

func New(clientInfo model.ClientInfo) (*Client, error) {
	ui, err := consoleui.New(clientInfo)
	if err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	pso := pso.New()
	clientConfig := config.ReadFromFile("./config.yaml")

	return &Client{
		pso:           pso,
		clientInfo:    clientInfo,
		config:        clientConfig,
		uiRefreshRate: clientConfig.GetUiRefreshRate(),
		ui:            ui,
		errChan:       make(chan error),
		done:          make(chan struct{}),
	}, nil
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

			case "<Resize>":
				c.ui.ClearScreen()
			}
		case err := <-c.errChan:
			close(c.done)
			return fmt.Errorf("run: error returned on error channel %w", err)
		}
	}
}

func (c *Client) runUI() {
	c.ui.ClearScreen()
	for {
		select {
		case <-time.After(c.uiRefreshRate):
			connected, statusString := c.pso.CheckConnection()
			c.ui.SetConnectionStatus(connected, statusString)

			err := c.ui.DrawScreen(&c.pso.CurrentPlayerData, &c.pso.TimeInState)
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

