package pso

import (
	"fmt"
	"github.com/phelix-/psostats/v2/pkg/common"
)

func (pso *PSO) GetFloorName() string {
	cmodeStage := pso.GameState.CmodeStage
	if cmodeStage < 1 {
		return common.GetFloorName(pso.GameState.Map)
	} else {
		episode := pso.GameState.Episode
		floor := int(pso.GameState.Floor)
		floorName := fmt.Sprintf("Unknown Floor E%vF%v", episode, floor)
		switch episode {
		case 1:
			if cmodeStage > 0 {
				area := 0
				switch cmodeStage {
				case 1:
					if floor == 11 {
						area = 3
					} else {
						area = floor
					}
				case 2:
					area = 1 + floor
				case 3:
					area = 5 + floor
				case 4:
					if floor == 12 {
						area = 19
					} else {
						area = 9 + floor
					}
				case 5:
					area = 14 + floor
				case 6:
					if floor == 13 {
						area = 30
					} else {
						area = 18 + floor
					}
				case 7:
					area = 23 + floor
				case 8:
					area = 27 + floor
				case 9:
					if floor == 14 {
						area = 46
					} else {
						area = 40 + floor
					}
				}
				floorName = fmt.Sprintf("Area %v", area)
			}
		case 2:
			if cmodeStage > 0 {
				area := 0
				switch cmodeStage {
				case 1:
					if floor == 14 {
						area = 7
					} else {
						area = floor
					}
				case 2:
					if floor == 15 {
						area = 14
					} else {
						area = 7 + floor
					}
				case 3:
					if floor == 12 {
						area = 20
					} else {
						area = 14 + floor
					}
				case 4:
					if floor == 13 {
						area = 27
					} else {
						area = 20 + floor
					}
				case 5:
					area = 27 + floor
				}
				floorName = fmt.Sprintf("Area %v", area)
			}
		}
		return floorName
	}
}
