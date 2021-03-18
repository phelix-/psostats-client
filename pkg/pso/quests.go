package pso

type Quest struct {
	StartingSwitchFloor uint16
	StartingSwitch      uint16
	EndingQuestRegister uint16
	EndingSwitchFloor   uint16
	EndingSwitch        uint16
}

func Ep1Quests() map[string]Quest {
	var quests = map[string]Quest{
		"Maximum Attack 4 -1A-": {
			StartingSwitchFloor: 4,
			StartingSwitch:      99,
			EndingSwitchFloor:   10,
			EndingSwitch:        31,
		},
	}
	// "Sweep-Up Operation 3":  WarpIn,
	// "": Switch,
	// "Maximum Attack 4 -1C-": Switch,
	// "Maximum Attack 4 -4C-": Switch,

	// quests[]
	return quests
}

func Ep2Quests() map[string]Quest {
	var quests = map[string]Quest{
		"Maximum Attack E: GDV": {
			StartingSwitchFloor: 5,
			StartingSwitch:      0,
			EndingQuestRegister: 254,
			EndingSwitchFloor:   10,
			EndingSwitch:        31,
		},
	}
	// "Sweep-Up Operation 3":  WarpIn,
	// "": Switch,
	// "Maximum Attack 4 -1C-": Switch,
	// "Maximum Attack 4 -4C-": Switch,

	// quests[]
	return quests
}

func Ep4Quests() map[string]Quest {
	var quests = map[string]Quest{
		"Maximum Attack 4 -4C-": {
			StartingSwitchFloor: 5,
			StartingSwitch:      66,
			EndingSwitchFloor:   8,
			EndingSwitch:        192,
		},
	}
	// "Sweep-Up Operation 3":  WarpIn,
	// "": Switch,
	// "Maximum Attack 4 -1C-": Switch,
	// "Maximum Attack 4 -4C-": Switch,

	// quests[]
	return quests
}
