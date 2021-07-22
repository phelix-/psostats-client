// Draws to the console
package consoleui

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/phelix-/psostats/v2/client/internal/pso"
	"github.com/phelix-/psostats/v2/client/internal/pso/inventory"
	"github.com/phelix-/psostats/v2/client/internal/pso/player"
	"github.com/phelix-/psostats/v2/pkg/model"
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
	equipment *[]inventory.Equipment,
	timeInState *[]pso.TimeDoingAction,
) error {
	termWidth, _ := ui.TerminalDimensions()
	cui.drawConnection(termWidth)
	cui.drawClass(playerData, termWidth)
	cui.drawEquipment(equipment, termWidth)
	cui.drawFrameData(playerData, timeInState, termWidth)
	return nil
}

func (cui *ConsoleUI) SetConnectionStatus(connected bool, statusString string) {
	cui.data.Connected = connected
	cui.data.Status = statusString
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
	connection.SetRect(0, 0, width, 3)
	ui.Render(connection)
}

func (cui *ConsoleUI) drawClass(playerData *player.BasePlayerInfo, width int) {
	playerInfo := widgets.NewParagraph()
	playerInfo.Text = fmt.Sprintf("%v", playerData.Class)
	playerInfo.Border = false
	playerInfo.WrapText = false
	playerInfo.SetRect(0, 3, width, 4)
	ui.Render(playerInfo)
}

func (cui *ConsoleUI) drawEquipment(equipment *[]inventory.Equipment, width int) {
	playerInfo := widgets.NewParagraph()
	for _,equipment := range *equipment {
		if equipment.Type == model.EquipmentTypeWeapon {
			playerInfo.Text = equipment.Display
		}
	}
	playerInfo.Border = false
	playerInfo.WrapText = false
	playerInfo.SetRect(0, 4, width, 5)
	ui.Render(playerInfo)
}

func (cui *ConsoleUI) drawFrameData(
	playerData *player.BasePlayerInfo,
	timeInState *[]pso.TimeDoingAction,
	width int,
) {
	list := widgets.NewList()
	list.Rows = make([]string, 0)

	list.Rows = append(list.Rows, fmt.Sprintf("%v", fmtActionName(playerData.ActionState)))
	total := 0
	for _,timeDoingAction := range *timeInState {
		total += timeDoingAction.Time
		list.Rows = append(list.Rows, fmt.Sprintf("%v - %2v", fmtActionName(timeDoingAction.Action), timeDoingAction.Time))
	}
	if playerData.ActionState == 1 {
		list.Rows = append(list.Rows, fmt.Sprintf("Standing to standing: %vf", total))
	}
	list.WrapText = false
	list.Border = false
	list.SetRect(0, 5, width, 40)
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
	case 23:
		return "emote"
	default:
		return fmt.Sprintf("unnamed action %v", action)
	}
}
