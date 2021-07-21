// Draws to the console
package consoleui

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/phelix-/psostats/v2/client/internal/client/config"
	"github.com/phelix-/psostats/v2/pkg/model"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/phelix-/psostats/v2/client/internal/pso"
	"github.com/phelix-/psostats/v2/client/internal/pso/player"
)

type Data struct {
	clientInfo model.ClientInfo
	Connected  bool
	Status     string
}

type ConsoleUI struct {
	data Data
	Motd string
}

func New(clientInfo model.ClientInfo) (*ConsoleUI, error) {
	err := ui.Init()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize termui: %w", err)
	}

	data := Data{
		Connected:  false,
		Status:     "Initializing",
		clientInfo: clientInfo,
	}

	return &ConsoleUI{
		data,
		"",
	}, nil
}

func (cui *ConsoleUI) Close() {
	ui.Close()
}

func (cui *ConsoleUI) ClearScreen() {
	ui.Clear()
}
func (cui *ConsoleUI) DrawScreen(
	playerData *player.BasePlayerInfo,
	gameState *pso.GameState,
	currentQuest *pso.QuestRun,
	config *config.Config,
	floorName string,
) error {
	termWidth, _ := ui.TerminalDimensions()
	cui.drawLogo(termWidth)
	cui.drawMotd(termWidth)
	cui.drawConnection(termWidth)
	cui.drawRecording(gameState, termWidth)
	cui.DrawHP(playerData, termWidth)
	cui.DrawLocation(floorName, playerData, gameState, termWidth)
	showQuestSplits := config.GetQuestSplitsEnabled() && len(currentQuest.Splits) > 0
	cui.drawQuestInfo(gameState, currentQuest, playerData, termWidth)
	cui.drawQuestInfo2(gameState, currentQuest, termWidth)
	if showQuestSplits {
		cui.drawQuestSplits(gameState, currentQuest, termWidth)
	}
	return nil
}

func (cui *ConsoleUI) SetConnectionStatus(connected bool, statusString string) {
	cui.data.Connected = connected
	cui.data.Status = statusString
}

func (cui *ConsoleUI) drawLogo(width int) {
	logo := widgets.NewParagraph()
	logo1 := ` ▄▄▄·.▄▄ ·       .▄▄ · ▄▄▄▄▄ ▄▄▄· ▄▄▄▄▄.▄▄ · 
▐█ ▄█▐█ ▀. ▪     ▐█ ▀. •██  ▐█ ▀█ •██  ▐█ ▀. 
 ██▀·▄▀▀▀█▄ ▄█▀▄ ▄▀▀▀█▄ ▐█.▪▄█▀▀█  ▐█.▪▄▀▀▀█▄
▐█▪·•▐█▄▪▐█▐█▌.▐▌▐█▄▪▐█ ▐█▌·▐█ ▪▐▌ ▐█▌·▐█▄▪▐█
.▀    ▀▀▀▀  ▀█▄▀▪ ▀▀▀▀  ▀▀▀  ▀  ▀  ▀▀▀  ▀▀▀▀ 
                                       v%v.%v.%v`
	logo.Text = fmt.Sprintf(logo1, cui.data.clientInfo.VersionMajor, cui.data.clientInfo.VersionMinor,
		cui.data.clientInfo.VersionPatch)
	logo.Border = false
	offset := (width - 48) / 2
	logo.SetRect(offset, 0, offset+48, 8)
	logo.TextStyle.Fg = ui.ColorCyan
	logo.WrapText = false
	ui.Render(logo)
}

func (cui *ConsoleUI) drawMotd(width int) {
	paragraph := widgets.NewParagraph()
	paragraph.Text = cui.Motd
	if strings.Contains(cui.Motd, "Invalid credentials") {
		paragraph.TextStyle.Fg = ui.ColorRed
	} else if strings.Contains(cui.Motd, "psostats.com/download") {
		paragraph.TextStyle.Fg = ui.ColorYellow
	} else {
		paragraph.TextStyle.Fg = ui.ColorWhite
	}
	paragraph.Border = false
	paragraph.WrapText = false
	paragraph.Text = padToCenter(paragraph.Text, width)

	paragraph.SetRect(0, 8, width, 9)
	ui.Render(paragraph)
}

func padToCenter(str string, width int) string {
	offset := (width - (len(str) + 3)) / 2
	buffer := bytes.NewBufferString("")
	for i := 0; i < offset; i++ {
		buffer.WriteString(" ")
	}
	buffer.WriteString(str)
	return buffer.String()
}

func (cui *ConsoleUI) drawConnection(width int) {
	connection := widgets.NewParagraph()
	connection.Text = fmt.Sprintf("[[ %v ]] ", cui.data.Status)
	if cui.data.Connected {
		connection.TextStyle.Fg = ui.ColorGreen
	} else {
		connection.TextStyle.Fg = ui.ColorRed
	}
	connection.Border = false
	connection.WrapText = false
	connection.Text = padToCenter(connection.Text, width+1)
	connection.SetRect(0, 9, width, 10)
	ui.Render(connection)
}

func (cui *ConsoleUI) drawRecording(gameState *pso.GameState, width int) {
	recording := widgets.NewParagraph()
	if gameState.UploadSuccessful {
		recording.TextStyle.Fg = ui.ColorGreen
		recording.Text = "[[ Uploaded ]] "
	} else if gameState.Uploading {
		recording.TextStyle.Fg = ui.ColorGreen
		recording.Text = "[[ Uploading ]] "
	} else if gameState.AwaitingUpload {
		recording.TextStyle.Fg = ui.ColorYellow
		recording.Text = "[[ Quest Ready For Upload: Press u to Upload ]] "
	} else if gameState.QuestComplete {
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
	recording.WrapText = false
	recording.Text = padToCenter(recording.Text, width+1)
	recording.SetRect(0, 10, width, 11)
	ui.Render(recording)
}

func (cui *ConsoleUI) DrawHP(playerData *player.BasePlayerInfo, width int) {
	playerInfo := widgets.NewParagraph()
	playerInfo.Text = fmt.Sprintf("%v - %v (gc: %v)", playerData.Class, playerData.Name, playerData.GuildCard)
	playerInfo.Border = false
	playerInfo.WrapText = false
	playerInfo.Text = padToCenter(playerInfo.Text, width)
	playerInfo.SetRect(0, 11, width, 12)
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
	hpOffset := (width - 50) / 2
	hpChart.SetRect(hpOffset, 12, hpOffset+50, 15)
	ui.Render(hpChart)
}

func (cui *ConsoleUI) DrawLocation(floorName string, playerData *player.BasePlayerInfo, gameState *pso.GameState, width int) {
	floor := widgets.NewParagraph()
	warpingText := ""
	if playerData.Warping {
		warpingText = " (Warping)"
	}
	floor.Text = fmt.Sprintf("%v Episode:%v %v%v Room:%v",
		gameState.Difficulty,
		gameState.Episode, floorName, warpingText, playerData.Room)
	floor.Border = false
	floor.WrapText = false

	floor.Text = padToCenter(floor.Text, width)
	floor.SetRect(0, 15, width, 16)
	ui.Render(floor)
}

func (cui *ConsoleUI) drawQuestInfo(
	gameState *pso.GameState,
	quest *pso.QuestRun,
	playerData *player.BasePlayerInfo,
	width int,
) {
	list := widgets.NewList()
	if len(gameState.QuestName) > 0 && gameState.QuestName != "No Active Quest" {
		list.Title = fmt.Sprintf("%40v", "[[ "+gameState.QuestName+" ]]")
	}
	if gameState.QuestStarted {
		list.Rows = []string{
			fmt.Sprintf("%40v", formatQuestComplete(quest)),
			fmt.Sprintf("%40v", formatCategory(quest)),
			fmt.Sprintf("%40v", formatQuestTime(quest)),
		}
		list.Rows = append(list.Rows, fmt.Sprintf("%40v", formatMesetaCharged(quest)))
		list.Rows = append(list.Rows, fmt.Sprintf("%40v", formatDeaths(quest)))
		list.Rows = append(list.Rows, fmt.Sprintf("%40v", formatMonstersAlive(quest)))
		list.Rows = append(list.Rows, fmt.Sprintf("%40v", formatMonstersKilled(quest)))
		if playerData.MaxTP > 0 {
			list.Rows = append(list.Rows, fmt.Sprintf("%40v", formatTpUsed(quest)))
		} else {
			list.Rows = append(list.Rows, fmt.Sprintf("%40v", formatTrapsUsed(quest)))
		}
	}
	list.WrapText = false
	list.Border = false
	offset := (width / 2) - 42
	list.SetRect(offset, 17, width/2, 40)
	ui.Render(list)
}


func (cui *ConsoleUI) drawQuestInfo2(
	gameState *pso.GameState,
	quest *pso.QuestRun,
	width int,
) {
	list := widgets.NewList()
	if gameState.QuestStarted {
		list.Rows = []string{
			formatTimeStanding(quest),
			formatTimeMoving(quest),
			formatTimeAttacking(quest),
		}
		if quest.TimeByState[8] > 0 {
			list.Rows = append(list.Rows, formatTimeCasting(quest))
		}
	}
	list.WrapText = false
	list.Border = false
	list.SetRect(width/2, 17, width, 23)
	ui.Render(list)
}

func (cui *ConsoleUI) drawQuestSplits(
	gameState *pso.GameState,
	quest *pso.QuestRun,
	width int,
) {
	list := widgets.NewList()
	if gameState.QuestStarted {
		list.Rows = make([]string, 0)
		for _, split := range quest.Splits {
			if !split.Start.IsZero() {
				list.Rows = append(list.Rows, formatSplitTime(split, quest))
			}
		}
	}

	list.WrapText = false
	list.Border = false
	list.SetRect(width/2, 22, width, 40)
	ui.Render(list)
}

func formatCategory(quest *pso.QuestRun) string {
	playerCount := len(quest.AllPlayers)
	category := fmt.Sprintf("Category:%4vp ", playerCount)
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
	return fmt.Sprintf("Meseta Charged:%11v", mesetaCharged)
}

func formatDeaths(quest *pso.QuestRun) string {
	return fmt.Sprintf("Deaths:%11v", quest.DeathCount)
}

func formatMonstersAlive(quest *pso.QuestRun) string {
	lastFrame := len(quest.MonsterCount)
	monstersAlive := 0
	if lastFrame > 0 {
		monstersAlive = quest.MonsterCount[lastFrame-1]
	}
	return fmt.Sprintf("Enemies Alive:%11v", monstersAlive)
}

func formatMonstersKilled(quest *pso.QuestRun) string {
	return fmt.Sprintf("Enemies Killed:%11v", quest.MonstersDead)
}

func formatTrapsUsed(quest *pso.QuestRun) string {
	return fmt.Sprintf("Traps Used:%5v/%2v/%2v", quest.FTUsed, quest.DTUsed, quest.CTUsed)
}

func formatTpUsed(quest *pso.QuestRun) string {
	return fmt.Sprintf("TP Used:%11v", quest.TPUsed)
}

func formatQuestTime(quest *pso.QuestRun) string {
	questDuration := time.Duration(0)
	if quest.QuestComplete {
		questDuration = quest.QuestEndTime.Sub(quest.QuestStartTime)
	} else {
		questDuration = time.Now().Sub(quest.QuestStartTime)
	}
	return fmt.Sprintf("Quest Duration:%11v", questDuration.Truncate(time.Millisecond*100))
}

func formatSplitTime(split model.QuestRunSplit, quest *pso.QuestRun) string {
	questDuration := time.Duration(0)
	if !split.Start.IsZero() {
		if !split.End.IsZero() {
			questDuration = split.End.Sub(split.Start)
		} else if quest.QuestComplete {
			questDuration = quest.QuestEndTime.Sub(split.Start)
		} else {
			questDuration = time.Now().Sub(split.Start)
		}
	}
	return fmt.Sprintf("%12v: %v", split.Name, questDuration.Truncate(time.Millisecond*100))
}

func formatQuestComplete(quest *pso.QuestRun) string {
	if quest.QuestComplete {
		return "Status:   Complete"
	} else {
		return "Status: Incomplete"
	}
}

func formatTimeMoving(quest *pso.QuestRun) string {
	timeMoving := quest.TimeByState[2] + quest.TimeByState[4]
	movingDuration := (time.Second / 30) * time.Duration(timeMoving)
	return fmt.Sprintf("%15v: %v", "Time Moving", movingDuration.Truncate(time.Millisecond*100))
}

func formatTimeStanding(quest *pso.QuestRun) string {
	timeMoving := quest.TimeByState[1]
	movingDuration := (time.Second / 30) * time.Duration(timeMoving)
	return fmt.Sprintf("%15v: %v", "Time Standing", movingDuration.Truncate(time.Millisecond*100))
}

func formatTimeAttacking(quest *pso.QuestRun) string {
	timeMoving := quest.TimeByState[5] + quest.TimeByState[6] + quest.TimeByState[7]
	movingDuration := (time.Second / 30) * time.Duration(timeMoving)
	return fmt.Sprintf("%15v: %v", "Time Attacking", movingDuration.Truncate(time.Millisecond*100))
}

func formatTimeCasting(quest *pso.QuestRun) string {
	timeMoving := quest.TimeByState[8]
	movingDuration := (time.Second / 30) * time.Duration(timeMoving)
	return fmt.Sprintf("%15v: %v", "Time Casting", movingDuration.Truncate(time.Millisecond*100))
}