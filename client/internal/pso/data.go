package pso

import (
	"errors"
	"log"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/client/internal/numbers"
	"github.com/phelix-/psostats/v2/client/internal/pso/inventory"
	"github.com/phelix-/psostats/v2/client/internal/pso/player"
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
