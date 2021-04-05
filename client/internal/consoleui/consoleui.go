// Draws to the console
package consoleui

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/phelix-/psostats/v2/client/internal/pso"
	"github.com/phelix-/psostats/v2/client/internal/pso/player"
)

type Data struct {
	Version   string
	Connected bool
	Status    string
}

type ConsoleUI struct {
	data Data
}

func New(version string) (*ConsoleUI, error) {
	err := ui.Init()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize termui: %w", err)
	}

	data := Data{
		Connected: false,
		Status:    "Initializing",
		Version:   version,
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
func (cui *ConsoleUI) DrawScreen(playerData *player.BasePlayerInfo, gameState *pso.GameState, currentQuest *pso.QuestRun) error {
	cui.drawLogo()
	cui.drawConnection()
	cui.drawRecording(gameState)
	cui.DrawHP(playerData)
	cui.DrawLocation(playerData, gameState)
	if gameState.QuestStarted && gameState.AllowQuestStart {
		cui.drawQuestInfo(currentQuest, playerData)
	}
	return nil
}

func (cui *ConsoleUI) SetConnectionStatus(connected bool, statusString string) {
	cui.data.Connected = connected
	cui.data.Status = statusString
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
	recording := widgets.NewParagraph()
	if gameState.QuestComplete {
		recording.TextStyle.Fg = ui.ColorGreen
		recording.Text = "[[ Quest Complete ]] "
	} else if gameState.QuestStarted && gameState.AllowQuestStart {
		recording.TextStyle.Fg = ui.ColorGreen
		recording.Text = "[[ Recording ]] "
	} else {
		recording.TextStyle.Fg = ui.ColorRed
		recording.Text = "[[ Waiting for Quest Start ]] "
	}
	recording.Border = false
	recording.SetRect(0, 2, 50, 3)
	ui.Render(recording)
}

func (cui *ConsoleUI) DrawHP(playerData *player.BasePlayerInfo) {
	playerInfo := widgets.NewParagraph()
	playerInfo.Text = fmt.Sprintf("%v - %v (gc: %v)", playerData.Class, playerData.Name, playerData.GuildCard)
	playerInfo.Border = false
	playerInfo.SetRect(0, 3, 80, 4)
	ui.Render(playerInfo)

	hpChart := widgets.NewGauge()
	hpChart.Title = "HP"
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
	floor.Text = fmt.Sprintf("%v Episode:%v %v%v Room:%v",
		gameState.Difficulty,
		gameState.Episode, floorName, warpingText, playerData.Room)
	floor.Border = false
	floor.SetRect(0, 6, 80, 9)
	ui.Render(floor)
}

func (cui *ConsoleUI) drawQuestInfo(quest *pso.QuestRun, playerData *player.BasePlayerInfo) {
	list := widgets.NewList()
	list.Title = "[[ " + quest.QuestName + " ]]"
	list.Rows = []string{
		formatQuestComplete(quest),
		formatCategory(quest),
		formatQuestTime(quest),
		formatMesetaCharged(quest),
		formatDeaths(quest),
		formatMonstersAlive(quest),
		formatMonstersKilled(quest),
	}
	if playerData.MaxTP > 0 {
		list.Rows = append(list.Rows, formatTpUsed(quest))
	} else {
		list.Rows = append(list.Rows, formatTrapsUsed(quest))
	}
	list.SetRect(0, 9, 80, 22)
	list.Border = false
	ui.Render(list)
}

func formatCategory(quest *pso.QuestRun) string {
	playerCount := len(quest.AllPlayers)
	category := fmt.Sprintf("Category:        %vp ", playerCount)
	if quest.PbCategory {
		category += "PB"
	} else {
		category += "No-PB"
	}
	return category
}

func formatMesetaCharged(quest *pso.QuestRun) string {
	lastFrame := len(quest.MesetaCharged)
	mesetaCharged := 0
	if lastFrame > 0 {
		mesetaCharged = quest.MesetaCharged[lastFrame-1]
	}
	return fmt.Sprintf("Meseta Charged:  %v", mesetaCharged)
}

func formatDeaths(quest *pso.QuestRun) string {
	return fmt.Sprintf("Deaths:          %v", quest.DeathCount)
}

func formatMonstersAlive(quest *pso.QuestRun) string {
	lastFrame := len(quest.MonsterCount)
	monstersAlive := 0
	if lastFrame > 0 {
		monstersAlive = quest.MonsterCount[lastFrame-1]
	}
	return fmt.Sprintf("Enemies Alive:   %v", monstersAlive)
}

func formatMonstersKilled(quest *pso.QuestRun) string {
	return fmt.Sprintf("Enemies Killed:  %v", quest.MonstersDead)
}

func formatTrapsUsed(quest *pso.QuestRun) string {
	return fmt.Sprintf("Traps Used:      %v/%v/%v", quest.FTUsed, quest.DTUsed, quest.CTUsed)
}

func formatTpUsed(quest *pso.QuestRun) string {
	return fmt.Sprintf("TP Used:         %v", quest.TPUsed)
}

func formatQuestTime(quest *pso.QuestRun) string {
	questDuration := time.Duration(0)
	if quest.QuestComplete {
		questDuration = quest.QuestEndTime.Sub(quest.QuestStartTime)
	} else {
		questDuration = time.Now().Sub(quest.QuestStartTime)
	}
	return fmt.Sprintf("Quest Duration:  %v", questDuration.Truncate(time.Millisecond*100))
}

func formatQuestComplete(quest *pso.QuestRun) string {
	if quest.QuestComplete {
		return "Status:          Complete"
	} else {
		return "Status:          Incomplete"
	}
}