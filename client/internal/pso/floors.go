package pso

import "fmt"

func GetFloorName(episode, floor, cmodeStage int) string {
	floorName := fmt.Sprintf("Unknown Floor E%vF%v", episode, floor)
	switch episode {
	case 1:
		if cmodeStage > 0 && floor > 0 && floor < 11 {
			area := 0
			switch cmodeStage {
			case 1:
				area = floor
			case 2:
				area = 3 + floor
			case 3:
				area = 8 + floor
			case 4:
				area = 13 + floor
			case 5:
				area = 19 + floor
			case 6:
				area = 24 + floor
			case 7:
				area = 30 + floor
			case 8:
				area = 36 + floor
			case 9:
				area = 40 + floor
			}
			floorName = fmt.Sprintf("Area %v", area)
		} else {
			switch floor {
			case 0:
				floorName = "Pioneer II"
			case 1:
				floorName = "Forest 1"
			case 2:
				floorName = "Forest 2"
			case 3:
				floorName = "Cave 1"
			case 4:
				floorName = "Cave 2"
			case 5:
				floorName = "Cave 3"
			case 6:
				floorName = "Mine 1"
			case 7:
				floorName = "Mine 2"
			case 8:
				floorName = "Ruins 1"
			case 9:
				floorName = "Ruins 2"
			case 10:
				floorName = "Ruins 3"
			case 11:
				floorName = "Dragon"
			case 12:
				floorName = "De Rol Le"
			case 13:
				floorName = "Vol Opt"
			case 14:
				floorName = "Dark Falz"
			case 15:
				floorName = "Lobby"
			}
		}
	case 2:
		if cmodeStage > 0 && floor > 0 && floor < 12 {
			area := 0
			switch cmodeStage {
			case 1:
				area = floor
			case 2:
				area = 7 + floor
			case 3:
				area = 14 + floor
			case 4:
				area = 20 + floor
			case 5:
				area = 28 + floor
			}
			floorName = fmt.Sprintf("Area %v", area)
		} else {
			switch floor {
			case 0:
				floorName = "Lab"
			case 1:
				floorName = "Temple Alpha"
			case 2:
				floorName = "Temple Beta 2"
			case 3:
				floorName = "Spaceship Alpha"
			case 4:
				floorName = "Spaceship Beta"
			case 5:
				floorName = "CCA"
			case 6:
				floorName = "Jungle North"
			case 7:
				floorName = "Jungle East"
			case 8:
				floorName = "Mountain"
			case 9:
				floorName = "Seaside"
			case 10:
				floorName = "Seabed Upper"
			case 11:
				floorName = "Seabed Lower"
			case 12:
				floorName = "Gal Gryphon"
			case 13:
				floorName = "Olga Flow"
			case 14:
				floorName = "Barba Ray"
			case 15:
				floorName = "Gol Dragon"
			case 16:
				floorName = "Seaside at Night"
			case 17:
				floorName = "Control Tower"
			}
		}
	case 4:
		switch floor {
		case 0:
			floorName = "Pioneer 2"
		case 1:
			floorName = "Crater East"
		case 2:
			floorName = "Crater West"
		case 3:
			floorName = "Crater South"
		case 4:
			floorName = "Crater North"
		case 5:
			floorName = "Crater Interior"
		case 6:
			floorName = "Desert 1"
		case 7:
			floorName = "Desert 2"
		case 8:
			floorName = "Desert 3"
		case 9:
			floorName = "Saint-Milion"
		case 15:
			floorName = "Lobby"
		}
	}
	return floorName
}
