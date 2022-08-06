package psoclasses

import "errors"

type PsoClass struct {
	Name      string
	MaxShifta int
	MinAtp    int
	MaxAtp    int
	Ata       int
}

var HUmar = PsoClass{Name: "HUmar", MaxShifta: 3, MinAtp: 1392, MaxAtp: 1397, Ata: 200}
var HUnewearl = PsoClass{Name: "HUnewearl", MaxShifta: 20, MinAtp: 1232, MaxAtp: 1237, Ata: 199}
var HUcast = PsoClass{Name: "HUcast", MaxShifta: 3, MinAtp: 1634, MaxAtp: 1639, Ata: 191}
var HUcaseal = PsoClass{Name: "HUcaseal", MaxShifta: 3, MinAtp: 1296, MaxAtp: 1301, Ata: 218}
var RAmar = PsoClass{Name: "RAmar", MaxShifta: 15, MinAtp: 1256, MaxAtp: 1260, Ata: 249}
var RAmarl = PsoClass{Name: "RAmarl", MaxShifta: 20, MinAtp: 1141, MaxAtp: 1145, Ata: 241}
var RAcast = PsoClass{Name: "RAcast", MaxShifta: 0, MinAtp: 1346, MaxAtp: 1350, Ata: 224}
var RAcaseal = PsoClass{Name: "RAcaseal", MaxShifta: 0, MinAtp: 1171, MaxAtp: 1175, Ata: 231}
var FOmar = PsoClass{Name: "FOmar", MaxShifta: 30, MinAtp: 1000, MaxAtp: 1002, Ata: 163}
var FOmarl = PsoClass{Name: "FOmarl", MaxShifta: 30, MinAtp: 870, MaxAtp: 872, Ata: 170}
var FOnewm = PsoClass{Name: "FOnewm", MaxShifta: 30, MinAtp: 812, MaxAtp: 814, Ata: 180}
var FOnewearl = PsoClass{Name: "FOnewearl", MaxShifta: 30, MinAtp: 581, MaxAtp: 583, Ata: 186}

func ForName(name string) (PsoClass, error) {
	switch name {
	case "HUmar":
		return HUmar, nil
	case "HUnewearl":
		return HUnewearl, nil
	case "HUcast":
		return HUcast, nil
	case "HUcaseal":
		return HUcaseal, nil
	case "RAmar":
		return RAmar, nil
	case "RAmarl":
		return RAmarl, nil
	case "RAcast":
		return RAcast, nil
	case "RAcaseal":
		return RAcaseal, nil
	case "FOmar":
		return FOmar, nil
	case "FOmarl":
		return FOmarl, nil
	case "FOnewm":
		return FOnewm, nil
	case "FOnewearl":
		return FOnewearl, nil
	}
	return PsoClass{}, errors.New("invalid class")
}

func GetAll() []PsoClass {
	return []PsoClass{
		HUmar,
		HUnewearl,
		HUcast,
		HUcaseal,
		RAmar,
		RAmarl,
		RAcast,
		RAcaseal,
		FOmar,
		FOmarl,
		FOnewm,
		FOnewearl,
	}
}
