// Draws to the console
package consoleui

import (
	"bytes"
	"fmt"
	"github.com/phelix-/psostats/v2/pkg/model"
	"strings"

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
	timeInState *[]pso.TimeDoingAction,
) error {
	termWidth, _ := ui.TerminalDimensions()
	cui.drawLogo(termWidth)
	cui.drawMotd(termWidth)
	cui.drawConnection(termWidth)
	cui.DrawHP(playerData, termWidth)
	cui.drawQuestInfo(playerData, timeInState, termWidth)
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

func (cui *ConsoleUI) drawQuestInfo(
	playerData *player.BasePlayerInfo,
	timeInState *[]pso.TimeDoingAction,
	width int,
) {
	list := widgets.NewList()
	list.Rows = make([]string, 0)

	list.Rows = append(list.Rows, fmt.Sprintf("%v", fmtActionName(playerData.ActionState)))
	for _,timeDoingAction := range *timeInState {
		list.Rows = append(list.Rows, fmt.Sprintf("%v - %2v", fmtActionName(timeDoingAction.Action), timeDoingAction.Time))
	}
	list.WrapText = false
	list.Border = false
	offset := (width / 2) - 42
	list.SetRect(offset, 17, width/2, 40)
	ui.Render(list)
}

func fmtActionName(action uint16) string {
	switch action {
	case 1:
		return "standing"
	case 2:
		return "walking"
	case 4:
		return "running"
	case 5:
		return "combo step 1"
	case 6:
		return "combo step 2"
	case 7:
		return "combo step 3"
	case 8:
		return "casting tech"
	default:
		return fmt.Sprintf("unnamed action %v", action)
	}
}
