package pso

import (
	"errors"
	"log"
	"unicode/utf16"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/client/internal/numbers"
	"github.com/phelix-/psostats/v2/client/internal/pso/inventory"
	"github.com/phelix-/psostats/v2/client/internal/pso/player"
	"github.com/phelix-/psostats/v2/client/internal/pso/quest"
)

const (
	basePlayerArrayAddress   = uintptr(0x00A94254)
	myPlayerIndexAddress     = 0x00A9C4F4
)

func (pso *PSO) RefreshData() error {
	if !pso.connected {
		log.Fatal("RefreshData: connection to window lost")
	}

	index, err := pso.getMyPlayerIndex()
	if err != nil {
		log.Fatal("Unable to find player index")
		return err
	}

	address := pso.getBaseCharacterAddress(index)

	if address != 0 {
		previousAction := pso.CurrentPlayerData.ActionState
		playerData, err := player.GetPlayerData(pso.handle, address, pso.server)
		if err != nil {
			return err
		}
		pso.CurrentPlayerData = playerData

		if playerData.ActionState != 1 && previousAction == 1 {
			pso.TimeInState = make([]TimeDoingAction, 0)
		}

		if playerData.ActionState != 1 {
			if len(pso.TimeInState) == 0 {
				action := TimeDoingAction{
					Action: playerData.ActionState,
					Time:   0,
				}
				action.Time++
				pso.TimeInState = append(pso.TimeInState, action)
			}
			if len(pso.TimeInState) > 0 {
				action := pso.TimeInState[len(pso.TimeInState) - 1]
				if action.Action != playerData.ActionState {
					action = TimeDoingAction{
						Action: playerData.ActionState,
						Time:   0,
					}
					pso.TimeInState = append(pso.TimeInState, action)
				}
				action.Time++
				pso.TimeInState[len(pso.TimeInState) - 1] = action
			}
		}

		equipment, err := inventory.ReadInventory(pso.handle, index)
		if err != nil {
			return err
		}
		pso.Equipment = equipment

	}

	return nil
}

func (pso *PSO) getMyPlayerIndex() (uint8, error) {
	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(myPlayerIndexAddress), 4)
	if !ok {
		return 0, errors.New("unable to find player index")
	}
	index := uint8(buf[0])
	return index, nil
}

func (pso *PSO) getBaseCharacterAddress(index uint8) uintptr {
	address := basePlayerArrayAddress + (4 * uintptr(index))
	return uintptr(numbers.ReadU32Unchecked(pso.handle, address))
}

func (pso *PSO) getMonsterUnitxtAddr() (uintptr, error) {
	unitxtAddr, err := numbers.ReadU32(pso.handle, uintptr(0x00a9cd50))
	if err != nil {
		return 0, err
	}
	monsterUnitxtAddr := uint32(0)
	if unitxtAddr != 0 {
		monsterUnitxtAddr, err = numbers.ReadU32(pso.handle, uintptr(unitxtAddr+16))
	}
	return uintptr(monsterUnitxtAddr), err
}

func (pso *PSO) getMonsterName(monsterId uint32) (string, error) {
	if monsterName, exists := pso.MonsterNames[monsterId]; exists {
		return monsterName, nil
	}
	monsterUnitxtAddr, err := pso.getMonsterUnitxtAddr()
	if err != nil {
		return "", err
	} else if monsterUnitxtAddr == 0 {
		return "", errors.New("monsterUnitxtAddr is unset")
	}
	monsterNameAddr, err := numbers.ReadU32(pso.handle, monsterUnitxtAddr+uintptr(4*monsterId))
	if err != nil {
		return "", err
	}
	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(monsterNameAddr), 32)
	if !ok {
		return "", errors.New("unable to getMonsterName")
	}

	endIndex := len(buf)
	for index, b := range buf {
		if b == 0x00 {
			endIndex = index
			break
		}
	}
	something := utf16.Decode(buf[0:endIndex])
	pso.MonsterNames[monsterId] = string(something)
	return string(something), nil
}

func (pso *PSO) getOtherPlayers() ([]player.BasePlayerInfo, error) {
	players := make([]player.BasePlayerInfo, 0)
	for i := 0; i < 12; i++ {
		address := pso.getBaseCharacterAddress(uint8(i))
		if address != 0 {
			playerData, err := player.GetPlayerData(pso.handle, address, pso.server)
			if err != nil {
				return nil, err
			}
			players = append(players, playerData)
		}
	}
	return players, nil
}

func (pso *PSO) getFloorSwitch(switchId uint16, floor uint16) (bool, error) {
	buf, _, ok := w32.ReadProcessMemory(pso.handle, uintptr(0xAC9FA0+(32*int(floor))), 32)
	if !ok {
		return false, errors.New("unable to getFloorSwitches")
	}
	var mask uint16
	if switchId%16 >= 8 {
		mask = uint16(0x8000) >> (switchId % 8)
	} else {
		mask = uint16(0x80) >> (switchId % 8)
	}
	switchSet := (buf[switchId/16] & mask) > 0
	return switchSet, nil
}

// -------------- Quest Data Block -------------- //
func (pso *PSO) getQuestName(questDataPtr uintptr) (string, error) {
	buf, _, ok := w32.ReadProcessMemory(pso.handle, questDataPtr+0x18, 64)
	if !ok {
		return "", errors.New("unable to getQuestName")
	}
	endIndex := len(buf)
	for index, b := range buf {
		if b == 0x00 {
			endIndex = index
			break
		}
	}
	questName := string(utf16.Decode(buf[0:endIndex]))
	return questName, nil
}

func (pso *PSO) checkQuestStartConditions(questConfig quest.Quest) (bool, error) {
	questStart := false
	if questConfig.StartsOnRegister() {
		registerPointer := quest.GetQuestRegisterPointer(pso.handle)
		registerSet, err := quest.IsRegisterSet(pso.handle, *questConfig.Start.Register, registerPointer)
		if err != nil {
			return false, err
		}
		if questConfig.GetCmodeStage() > 0 {
			cmodeFailedRegister, err := quest.IsRegisterSet(pso.handle, 253, registerPointer)
			if err != nil {
				return false, err
			}
			questStart = registerSet && !cmodeFailedRegister
		} else {
			questStart = registerSet
		}
	} else if questConfig.TerminalQuest() {
		switchSet, err := pso.getFloorSwitch(questConfig.Start.Switch, questConfig.Start.Floor)
		if err != nil {
			return false, err
		}
		questStart = switchSet
	} else if questConfig.StartsAtWarpIn() {
		allPlayers, err := pso.getOtherPlayers()
		if err != nil {
			log.Panicf("unable to get all players %v", err)
		}
		for _, p := range allPlayers {
			if p.Floor != 0 && !p.Warping {
				questStart = true
				break
			}
		}
	}
	return questStart, nil
}
