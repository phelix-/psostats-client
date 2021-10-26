// Foundation of the psostats client
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gizak/termui/v3/widgets"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"syscall"
	"time"

	"github.com/hishboy/gocommons/lang"

	"github.com/phelix-/psostats/v2/pkg/model"

	termui "github.com/gizak/termui/v3"
	"github.com/phelix-/psostats/v2/client/internal/client/config"
	"github.com/phelix-/psostats/v2/client/internal/consoleui"
	"github.com/phelix-/psostats/v2/client/internal/pso"
)

type Client struct {
	pso           *pso.PSO
	clientInfo    model.ClientInfo
	httpClient    http.Client
	config        *config.Config
	uiRefreshRate time.Duration
	ui            *consoleui.ConsoleUI
	currentGameId int
	errChan       chan error
	done          chan struct{}
	startedGame   chan pso.QuestRun
	completeGame  chan pso.QuestRun
	gameQueue     *lang.Queue
}

func New(clientInfo model.ClientInfo) (*Client, error) {
	ui, err := consoleui.New(clientInfo)
	if err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	startedGameChannel := make(chan pso.QuestRun)
	completeGameChannel := make(chan pso.QuestRun)
	pso := pso.New(startedGameChannel, completeGameChannel)
	clientConfig, err := config.ReadFromFile("./config.yaml")
	if err != nil {
		if pathErr, ok := err.(*fs.PathError); ok && pathErr.Err == syscall.ERROR_FILE_NOT_FOUND {
			showMissingConfigUi(ui)
			log.Fatalf("Missing config file, shutting down")
		} else {
			log.Fatalf("Unable to start client %v", err)
		}
	}

	return &Client{
		pso:           pso,
		clientInfo:    clientInfo,
		httpClient:    http.Client{},
		config:        clientConfig,
		uiRefreshRate: clientConfig.GetUiRefreshRate(),
		ui:            ui,
		errChan:       make(chan error),
		done:          make(chan struct{}),
		startedGame:   startedGameChannel,
		completeGame:  completeGameChannel,
		gameQueue:     lang.NewQueue(),
	}, nil
}

func (c *Client) Run() error {
	defer c.ui.Close()

	if err := c.getMotd(); err != nil {
		c.ui.Motd = fmt.Sprintf("Error getting message of the day %v", err)
	}

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
				for c.gameQueue.Len() > 0 {
					game, ok := c.gameQueue.Poll().(*pso.QuestRun)
					if ok {
						c.uploadGame(*game)
					} else {
						log.Printf("Manual upload failed")
						break
					}

				}
				c.pso.GameState.AwaitingUpload = false
			case "<Resize>":
				c.ui.ClearScreen()
			}
		case game := <-c.startedGame:
			if c.config.GetQuestSplitsEnabled() && c.config.GetQuestSplitsCompareTo() != "none" {
				err := c.getQuestSplits(game.QuestName, len(game.AllPlayers), game.PbCategory)
				if err != nil {
					log.Printf("Error getting quest splits %v", err)
				}
			}
		case game := <-c.completeGame:
			if c.config.AutoUploadEnabled() {
				c.uploadGame(game)
			} else {
				c.gameQueue.Push(&game)
				c.pso.GameState.AwaitingUpload = true
			}
		case err := <-c.errChan:
			close(c.done)
			return fmt.Errorf("run: error returned on error channel %w", err)
		}
	}
}

func showMissingConfigUi(cui *consoleui.ConsoleUI) {
	width, _ := termui.TerminalDimensions()
	cui.DrawLogo(width)
	paragraph := widgets.NewParagraph()
	paragraph.Text = "Config file missing. Press any key to quit"
	offset := (width - 44) / 2
	paragraph.SetRect(offset, 10, offset+48, 13)
	paragraph.Border = false
	termui.Render(paragraph)
	uiEvents := termui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			if e.Type == termui.KeyboardEvent {
				return
			}
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
	jsonBytes, err := json.Marshal(c.pso.CurrentQuest)
	if err != nil {
		log.Printf("Unable to generate json")
	}
	file.Write(jsonBytes)
}

func (c *Client) uploadGame(game pso.QuestRun) {
	if c.pso.GameState.Uploading || c.pso.GameState.UploadSuccessful {
		// Prevent resubmission
		return
	}
	game.Client = c.clientInfo
	c.pso.GameState.Uploading = true
	jsonBytes, err := json.Marshal(game)
	if err != nil {
		log.Printf("Unable to generate json")
	}
	buf := bytes.NewBuffer(jsonBytes)
	request, err := http.NewRequest("POST", c.config.GetServerBaseUrl()+"/api/game", buf)
	if err != nil {
		log.Printf("Failed to build request %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(*c.config.User, *c.config.Password)
	response, err := c.httpClient.Do(request)
	if err != nil {
		log.Printf("Unable to upload game %v", err)
	}
	if response == nil || response.StatusCode != 200 {
		log.Printf("Got response status %v: %v", response.StatusCode, response.Body)
	} else {
		c.pso.GameState.UploadSuccessful = true
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("Error reading response from server %v", err)
		}
		postResponse := model.PostGameResponse{}
		err = json.Unmarshal(responseBytes, &postResponse)
		if err != nil {
			log.Printf("Error reading response from server %v", err)
		} else {
			gameUrl := fmt.Sprintf("Last game: %v/%v", c.config.GetServerBaseUrl(), postResponse.Id)
			if postResponse.Record {
				gameUrl = fmt.Sprintf("%v - RECORD", gameUrl)
			} else if postResponse.Pb {
				gameUrl = fmt.Sprintf("%v - PB", gameUrl)
			}
			c.ui.Motd = gameUrl
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

			currentQuest := c.pso.CurrentQuest
			floorName := c.pso.GetFloorName()
			err := c.ui.DrawScreen(&c.pso.CurrentPlayerData, &c.pso.GameState, &currentQuest, c.config, floorName)
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

func (c *Client) getMotd() error {
	jsonBytes, err := json.Marshal(c.clientInfo)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(jsonBytes)
	request, err := http.NewRequest("POST", c.config.GetServerBaseUrl()+"/api/motd", buf)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(*c.config.User, *c.config.Password)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	motd := model.MessageOfTheDay{}
	if err := json.Unmarshal(responseBytes, &motd); err != nil {
		return err
	}

	if !motd.Authorized {
		c.ui.Motd = "Invalid credentials"
	} else {
		c.ui.Motd = motd.Message
	}
	return nil
}

func (c *Client) getQuestSplits(questName string, players int, pbCategory bool) error {
	c.ui.QuestSplits = nil
	compareTo := c.config.GetQuestSplitsCompareTo()
	urlEncodedQuestName := url.PathEscape(questName)
	var path string
	if compareTo == "pb" {
		path = fmt.Sprintf("%v/api/pb-splits/%v?players=%d&pb=%v", c.config.GetServerBaseUrl(), urlEncodedQuestName, players, pbCategory)
	} else {
		path = fmt.Sprintf("%v/api/record-splits/%v?players=%d&pb=%v", c.config.GetServerBaseUrl(), urlEncodedQuestName, players, pbCategory)
	}
	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	request.SetBasicAuth(*c.config.User, *c.config.Password)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		splits := make([]model.QuestRunSplit, 0)
		if err := json.Unmarshal(responseBytes, &splits); err != nil {
			return err
		}

		c.ui.QuestSplits = splits
	}
	return nil
}
