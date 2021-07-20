package common

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
