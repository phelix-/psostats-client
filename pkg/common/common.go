package common

import "fmt"

type PsoClass struct {
	Name string
	MaxShifta int
}

func GetAllClasses() []PsoClass {
	return []PsoClass {
		{Name: "HUmar", MaxShifta: 3},
		{Name: "HUnewearl", MaxShifta: 20},
		{Name: "HUcast", MaxShifta: 3},
		{Name: "HUcaseal", MaxShifta: 3},
		{Name: "RAmar", MaxShifta: 15},
		{Name: "RAmarl", MaxShifta: 20},
		{Name: "RAcast", MaxShifta: 0},
		{Name: "RAcaseal", MaxShifta: 0},
		{Name: "FOmar", MaxShifta: 30},
		{Name: "FOmarl", MaxShifta: 30},
		{Name: "FOnewm", MaxShifta: 30},
		{Name: "FOnewearl", MaxShifta: 30},
	}
}

func GetFloorName(mapNum uint16) string {
	floorName := fmt.Sprintf("Unknown Map %v", mapNum)
	switch mapNum {
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
	case 16:
		floorName = "BA Spaceship"
	case 17:
		floorName = "BA Temple"
	case 18:
		floorName = "Lab"
	case 19:
		floorName = "Temple Alpha"
	case 20:
		floorName = "Temple Beta 2"
	case 21:
		floorName = "Spaceship Alpha"
	case 22:
		floorName = "Spaceship Beta"
	case 23:
		floorName = "CCA"
	case 24:
		floorName = "Jungle North"
	case 25:
		floorName = "Jungle East"
	case 26:
		floorName = "Mountain"
	case 27:
		floorName = "Seaside"
	case 28:
		floorName = "Seabed Upper"
	case 29:
		floorName = "Seabed Lower"
	case 30:
		floorName = "Gal Gryphon"
	case 31:
		floorName = "Olga Flow"
	case 32:
		floorName = "Barba Ray"
	case 33:
		floorName = "Gol Dragon"
	case 34:
		floorName = "Seaside at Night"
	case 35:
		floorName = "Control Tower"
	case 36:
		floorName = "Pioneer 2"
	case 37:
		floorName = "Crater East"
	case 38:
		floorName = "Crater West"
	case 39:
		floorName = "Crater South"
	case 40:
		floorName = "Crater North"
	case 41:
		floorName = "Crater Interior"
	case 42:
		floorName = "Desert 1"
	case 43:
		floorName = "Desert 2"
	case 44:
		floorName = "Desert 3"
	case 45:
		floorName = "Saint-Milion"
	}


	return floorName
}