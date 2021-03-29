package consoleui

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/phelix-/psostats/v2/pkg/pso"
	"github.com/phelix-/psostats/v2/pkg/pso/player"
)

type Data struct {
	Version   string
	Connected bool
	Status    string
}

func (data *Data) Clear() {
	data.Status = "Cleared"
	data.Connected = false
}

type ConsoleUI struct {
	data *Data
}

func New(data *Data) (*ConsoleUI, error) {
	err := ui.Init()
	if err != nil {
		return nil, fmt.Errorf("New: unable to initialize termui: %w", err)
	}

	return &ConsoleUI{
		data,
	}, nil
}

func (cui *ConsoleUI) Close() {
	ui.Close()
}

func (cui *ConsoleUI) ClearScreen() {
	ui.Clear()
}
func (cui *ConsoleUI) DrawScreen(playerData *player.BasePlayerInfo, gameState *pso.GameState) error {
	cui.drawLogo()
	cui.drawConnection()
	cui.drawRecording(gameState)
	cui.DrawHP(playerData)
	cui.DrawLocation(playerData, gameState)
	return nil
}

func (cui *ConsoleUI) drawLogo() {
	logo := widgets.NewParagraph()
	logo.Text = fmt.Sprintf("PSOStats %v", cui.data.Version)
	logo.Border = false
	logo.SetRect(0, 0, 50, 1)
	logo.TextStyle.Fg = ui.ColorCyan
	ui.Render(logo)
}

func (cui *ConsoleUI) drawConnection() {
	connection := widgets.NewParagraph()
	connection.Text = fmt.Sprintf("[[ %v ]] ", cui.data.Status)
	if cui.data.Connected {
		connection.TextStyle.Fg = ui.ColorGreen
	} else {
		connection.TextStyle.Fg = ui.ColorRed
	}
	connection.Border = false
	connection.SetRect(0, 1, 80, 2)
	ui.Render(connection)
}

func (cui *ConsoleUI) drawRecording(gameState *pso.GameState) {
	connection := widgets.NewParagraph()
	if gameState.QuestComplete {
		connection.TextStyle.Fg = ui.ColorGreen
		connection.Text = "[[ Quest Complete ]] " + gameState.QuestEndTime.Sub(gameState.QuestStartTime).String()
	} else if gameState.QuestStarted {
		connection.TextStyle.Fg = ui.ColorGreen
		connection.Text = "[[ Recording ]] " + time.Now().Sub(gameState.QuestStartTime).String()
	} else {
		connection.TextStyle.Fg = ui.ColorRed
		connection.Text = "[[ Waiting for Quest Start ]] "
	}
	connection.Border = false
	connection.SetRect(0, 2, 50, 3)
	ui.Render(connection)
}

func (cui *ConsoleUI) DrawHP(playerData *player.BasePlayerInfo) {
	connection := widgets.NewParagraph()
	connection.Text = fmt.Sprintf("%v - %v (gc: %v)", playerData.Class, playerData.Name, playerData.Guildcard)
	connection.Border = false
	connection.SetRect(0, 3, 80, 4)
	ui.Render(connection)
	//, playerData.HP, playerData.MaxHP,playerData.Meseta, playerData.ShiftaLvl, playerData.DebandLvl

	hpChart := widgets.NewGauge()
	hpChart.Title = "HP"
	// percentHp := float64(playerData.HP) / float64(playerData.MaxHP)
	// hpChart.Percent = int(percentHp * 100)
	percentHp := 0
	if playerData.MaxHP > 0 {
		percentHp = (100 * int(playerData.HP)) / int(playerData.MaxHP)
	}
	hpChart.Percent = percentHp
	hpChart.Label = fmt.Sprintf("%v/%v", playerData.HP, playerData.MaxHP)
	hpChart.Border = false
	hpChart.SetRect(0, 4, 50, 7)
	ui.Render(hpChart)
}

func (cui *ConsoleUI) DrawLocation(playerData *player.BasePlayerInfo, gameState *pso.GameState) {
	floor := widgets.NewParagraph()
	floorName := pso.GetFloorName(int(gameState.Episode), int(playerData.Floor))
	warpingText := ""
	if playerData.Warping {
		warpingText = " (Warping)"
	}
	floor.Text = fmt.Sprintf("%v Episode:%v %v%v Room:%v\n[%v]\nMonsters Alive-%v",
		gameState.Difficulty,
		gameState.Episode, floorName, warpingText, playerData.Room, gameState.QuestName,
		gameState.MonsterCount)
	floor.Border = false
	floor.SetRect(0, 6, 80, 14)
	ui.Render(floor)
}
