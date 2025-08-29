package weapons

func GetWeaponsUltima() []Weapon {
	vanillaWeapons := GetWeapons()
	weaponsWithModifiedStats := map[string]Weapon{
		"Dark Flow":         {Name: "Dark Flow", MinAtp: 1000, MaxAtp: 1100, Ata: 50, Grind: 0, ComboPreset: Combo{Attack1: "SPECIAL", Attack1Hits: 5, Attack2Hits: 5, Attack3Hits: 5}, Animation: "Sword", Special: "Dark Flow", MaxHit: 100, MaxAttr: 100},
		"Tsumikiri J-Sword": {Name: "Tsumikiri J-Sword", MinAtp: 900, MaxAtp: 950, Ata: 40, Grind: 50, Animation: "Sword", Special: "TJS", MaxHit: 100, MaxAttr: 100},
		"S-Beat's Blade":    {Name: "S-Beat's Blade", MinAtp: 210, MaxAtp: 220, Ata: 35, Grind: 15, MaxHit: 80, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxAttr: 100},
		"P-Arms' Blade":     {Name: "P-Arms' Blade", MinAtp: 250, MaxAtp: 270, Ata: 34, Grind: 25, MaxHit: 80, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxAttr: 100},
		"S-Red's Blade":     {Name: "S-Red's Blade", MinAtp: 340, MaxAtp: 350, Ata: 39, Grind: 15, MaxHit: 80, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxAttr: 100},
		"Two Kamui":         {Name: "Two Kamui", MinAtp: 600, MaxAtp: 650, Ata: 50, Grind: 0, MaxHit: 80, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxAttr: 100},
		"Lavis Blade":       {Name: "Lavis Blade", MinAtp: 380, MaxAtp: 450, Ata: 40, Grind: 0, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxHit: 100, MaxAttr: 100},
		"Girasole":          {Name: "Girasole", MinAtp: 500, MaxAtp: 550, Ata: 50, Grind: 0, MaxHit: 100, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Double Saber", MaxAttr: 100},
		"Inferno Girasole":  {Name: "Inferno Girasole", MinAtp: 700, MaxAtp: 820, Ata: 50, Grind: 20, MaxHit: 100, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Double Saber", MaxAttr: 100},
		"TypeDS/D.Saber":    {Name: "TypeDS/D.Saber", MinAtp: 30, MaxAtp: 30, Ata: 40, Grind: 125, MaxHit: 100, MaxAttr: 100, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Double Saber"},
		"Bringer's Rifle":   {Name: "Bringer's Rifle", MinAtp: 330, MaxAtp: 370, Ata: 63, Grind: 9, Special: "Demon's", MaxAttr: 100, Animation: "Rifle"},
		"Dual Bird":         {Name: "Dual Bird", MinAtp: 200, MaxAtp: 222, Ata: 45, Grind: 50, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		"TypeME/Mechgun":    {Name: "TypeME/Mechgun", MinAtp: 10, MaxAtp: 10, Ata: 20, Grind: 30, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 80, MaxAttr: 100, Animation: "Mechgun"},
		"Rambling May":      {Name: "Rambling May", MinAtp: 360, MaxAtp: 450, Ata: 45, Grind: 50, MaxHit: 100, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, MaxAttr: 100, Animation: "Shot"},
		"Dark Meteor":       {Name: "Dark Meteor", MinAtp: 750, MaxAtp: 880, Ata: 45, Grind: 25, ComboPreset: Combo{Attack2: "NONE", Attack3: "NONE"}, MaxHit: 100, MaxAttr: 100, Animation: "Shot"},
	}
	customWeapons := map[string]Weapon{
		"Power Glove":         {Name: "Power Glove", MinAtp: 1500, MaxAtp: 1500, Ata: 75, Grind: 0, MaxHit: 100, MaxAttr: 100, Animation: "Fist", Special: "Spirit"},
		"Hundred Souls":       {Name: "Hundred Souls", MinAtp: 1000, MaxAtp: 1200, Ata: 50, Grind: 0, Special: "Spirit", MaxHit: 100, MaxAttr: 100, Animation: "Saber"},
		"Blood Sword":         {Name: "Blood Sword", MinAtp: 1000, MaxAtp: 1000, Ata: 70, Grind: 9, Special: "Spirit", MaxHit: 100, MaxAttr: 100, Animation: "Saber"},
		"Fire Rod":            {Name: "Fire Rod", MinAtp: 550, MaxAtp: 650, Ata: 50, Grind: 0, MaxHit: 100, MaxAttr: 100, Special: "Spirit", Animation: "Saber"},
		"Sil Dragon Slayer":   {Name: "Sil Dragon Slayer", MinAtp: 600, MaxAtp: 650, Ata: 70, Grind: 35, Animation: "Sword", Special: "Blizzard", MaxHit: 100, MaxAttr: 100},
		"Crimson Sword":       {Name: "Crimson Sword", MinAtp: 620, MaxAtp: 800, Ata: 80, Grind: 55, Animation: "Sword", Special: "Arrest", MaxHit: 100, MaxAttr: 100},
		"Master Sword":        {Name: "Master Sword", MinAtp: 700, MaxAtp: 780, Ata: 35, Grind: 70, Animation: "Sword", Special: "TJS", MaxHit: 100, MaxAttr: 100},
		"Stealth Sword":       {Name: "Stealth Sword", MinAtp: 300, MaxAtp: 350, Ata: 50, Grind: 0, Animation: "Sword", Special: "Demon's", MaxHit: 100, MaxAttr: 100},
		"Macho Blades":        {Name: "Macho Blades", MinAtp: 400, MaxAtp: 400, Ata: 50, Grind: 25, MaxHit: 80, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxAttr: 100},
		"Blood Tornado":       {Name: "Blood Tornado", MinAtp: 550, MaxAtp: 600, Ata: 70, Grind: 33, Special: "Demon's", ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxHit: 100, MaxAttr: 100},
		"Fury of the Beast":   {Name: "Fury of the Beast", MinAtp: 500, MaxAtp: 600, Ata: 65, Grind: 25, Special: "Charge", ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxHit: 100, MaxAttr: 100},
		"TypeBL/Blade":        {Name: "TypeBL/Blade", MinAtp: 10, MaxAtp: 10, Ata: 35, Grind: 90, MaxHit: 100, MaxAttr: 100, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", Special: "Berserk"},
		"Rico's Parasol":      {Name: "Rico's Parasol", MinAtp: 250, MaxAtp: 300, Ata: 40, Grind: 0, Special: "Charge", Animation: "Partisan", MaxHit: 100, MaxAttr: 100},
		"Great Fairy Sword":   {Name: "Great Fairy Sword", MinAtp: 150, MaxAtp: 150, Ata: 40, Grind: 99, Special: "Charge", Animation: "Partisan", MaxHit: 100, MaxAttr: 100},
		"Sword of Ultima":     {Name: "Sword of Ultima", MinAtp: 300, MaxAtp: 350, Ata: 15, Grind: 55, Special: "Charge", Animation: "Partisan", MaxHit: 100, MaxAttr: 100},
		"Ultima Reaper":       {Name: "Ultima Reaper", MinAtp: 666, MaxAtp: 666, Ata: 45, Grind: 15, Special: "Hell", Animation: "Partisan", MaxHit: 100, MaxAttr: 100},
		"Boomerang":           {Name: "Boomerang", MinAtp: 200, MaxAtp: 200, Ata: 30, Grind: 0, Special: "Arrest", Animation: "Slicer", MaxHit: 100, MaxAttr: 100},
		"Kiss of Death":       {Name: "Kiss of Death", MinAtp: 350, MaxAtp: 350, Ata: 35, Grind: 0, Special: "Hell", Animation: "Slicer", MaxHit: 100, MaxAttr: 100},
		"Slicer of Vengeance": {Name: "Slicer of Vengeance", MinAtp: 470, MaxAtp: 525, Ata: 30, Grind: 0, Special: "Charge", Animation: "Slicer", MaxHit: 100, MaxAttr: 100},
		"Yamigarasu":          {Name: "Yamigarasu", MinAtp: 580, MaxAtp: 650, Ata: 53, Grind: 0, Special: "Hell", MaxHit: 100, MaxAttr: 100, Animation: "Katana"},
		"Ten Years Blades":    {Name: "Ten Years Blades", MinAtp: 600, MaxAtp: 650, Ata: 74, Grind: 200, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 2, Attack3Hits: 2}, MaxHit: 100, MaxAttr: 100, Animation: "Twin Sword"},
		"Tyrfing":             {Name: "Tyrfing", MinAtp: 700, MaxAtp: 740, Ata: 80, Grind: 13, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 2, Attack3Hits: 2}, Special: "Geist", MaxHit: 100, MaxAttr: 100, Animation: "Twin Sword"},
		"Suppressed Gun":      {Name: "Suppressed Gun", MinAtp: 260, MaxAtp: 270, Ata: 47, Grind: 9, Special: "Charge", MaxHit: 100, MaxAttr: 100, Animation: "Handgun"},
		"Serene Swan":         {Name: "Serene Swan", MinAtp: 5, MaxAtp: 15, Ata: 37, Grind: 80, MaxHit: 100, ComboPreset: Combo{Attack1Hits: 4, Attack2Hits: 4, Attack3Hits: 4}, Special: "Devil's", MaxAttr: 100, Animation: "Last Swan"},
		"Asteron Striker":     {Name: "Asteron Striker", MinAtp: 300, MaxAtp: 380, Ata: 60, Grind: 0, Special: "Hell", MaxHit: 100, MaxAttr: 100, Animation: "Handgun"},
		"Hand of Justice":     {Name: "Hand of Justice", MinAtp: 400, MaxAtp: 400, Ata: 65, Grind: 25, Special: "Demon's", MaxHit: 100, MaxAttr: 100, Animation: "Handgun"},
		"Morolian Blaster":    {Name: "Morolian Blaster", MinAtp: 280, MaxAtp: 310, Ata: 60, Grind: 0, Special: "Arrest", MaxHit: 100, MaxAttr: 100, Animation: "Handgun"},
		"Yasminkov 2000H":     {Name: "Yasminkov 2000H", MinAtp: 340, MaxAtp: 340, Ata: 45, Grind: 10, MaxHit: 100, MaxAttr: 100, Animation: "Handgun"},
		"TypeGU/Mechgun":      {Name: "TypeGU/Mechgun", MinAtp: 10, MaxAtp: 10, Ata: 50, Grind: 90, MaxHit: 100, MaxAttr: 100, Animation: "Handgun"},
		"Arrest Needle":       {Name: "Arrest Needle", MinAtp: 300, MaxAtp: 400, Ata: 75, Grind: 60, Special: "Arrest", MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		"Ultima Bringer's":    {Name: "Ultima Bringer's", MinAtp: 330, MaxAtp: 370, Ata: 60, Grind: 70, Special: "Demon's", MaxAttr: 100, Animation: "Rifle"},
		"Egg Blaster MK2":     {Name: "Egg Blaster MK2", MinAtp: 300, MaxAtp: 330, Ata: 50, Grind: 10, MaxHit: 100, MaxAttr: 100, Special: "Berserk", Animation: "Rifle"},
		"Lindcray":            {Name: "Lindcray", MinAtp: 800, MaxAtp: 1200, Ata: 75, Grind: 0, MaxHit: 100, MaxAttr: 100, Special: "Spirit", Animation: "Rifle"},
		"Sacred Bow":          {Name: "Sacred Bow", MinAtp: 1000, MaxAtp: 1400, Ata: 70, Grind: 0, MaxHit: 100, MaxAttr: 100, Special: "Hell", Animation: "Rifle"},
		"Water Gun":           {Name: "Water Gun", MinAtp: 50, MaxAtp: 50, Ata: 50, Grind: 0, MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		"Rianov 303SNR-3":     {Name: "Rianov 303SNR-3", MinAtp: 500, MaxAtp: 500, Ata: 65, Grind: 15, MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		"Rianov 303SNR-4":     {Name: "Rianov 303SNR-4", MinAtp: 350, MaxAtp: 450, Ata: 60, Grind: 60, MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		"Rianov 303SNR-5":     {Name: "Rianov 303SNR-5", MinAtp: 550, MaxAtp: 550, Ata: 70, Grind: 20, MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		"Yasminkov 3000R":     {Name: "Yasminkov 3000R", MinAtp: 370, MaxAtp: 400, Ata: 66, Grind: 60, MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		"Yasminkov 7000R":     {Name: "Yasminkov 7000R", MinAtp: 370, MaxAtp: 450, Ata: 67, Grind: 25, MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		"Rage de Glace":       {Name: "Rage de Glace", MinAtp: 200, MaxAtp: 232, Ata: 85, Grind: 9, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		"Mille Fauciles":      {Name: "Mille Fauciles", MinAtp: 50, MaxAtp: 50, Ata: 65, Grind: 250, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		"Psycho Raven":        {Name: "Psycho Raven", MinAtp: 480, MaxAtp: 480, Ata: 45, Grind: 80, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		"Samba Fiesta":        {Name: "Samba Fiesta", MinAtp: 5, MaxAtp: 10, Ata: 60, Grind: 0, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, Special: "Demon's", MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		"L&K40 Combat":        {Name: "L&K40 Combat", MinAtp: 55, MaxAtp: 66, Ata: 30, Grind: 0, ComboPreset: Combo{Attack1Hits: 5, Attack2Hits: 5, Attack3Hits: 5}, MaxHit: 100, MaxAttr: 100, Animation: "L&K38 Combat"},
		"Crush Cannon":        {Name: "Crush Cannon", MinAtp: 100, MaxAtp: 200, Ata: 50, Grind: 25, MaxHit: 100, MaxAttr: 100, Special: "Gush", Animation: "Shot"},
		"Iron Faust":          {Name: "Iron Faust", MinAtp: 500, MaxAtp: 580, Ata: 42, Grind: 18, MaxHit: 100, MaxAttr: 100, Animation: "Shot"},
		"Arrest Faust":        {Name: "Arrest Faust", MinAtp: 100, MaxAtp: 110, Ata: 50, Grind: 0, MaxHit: 100, MaxAttr: 100, Special: "Arrest", Animation: "Shot"},
		"Frozen Faust":        {Name: "Frozen Faust", MinAtp: 100, MaxAtp: 110, Ata: 45, Grind: 0, MaxHit: 100, MaxAttr: 100, Special: "Blizzard", Animation: "Shot"},

		"ES Launcher":   {Name: "ES Launcher", MinAtp: 10, MaxAtp: 10, Ata: 40, Grind: 180, MaxHit: 0, MaxAttr: 0, Animation: "Launcher", Special: "Berserk"},
		"Outlaw Star":   {Name: "Outlaw Star", MinAtp: 240, MaxAtp: 270, Ata: 30, Grind: 0, MaxHit: 100, MaxAttr: 100, Animation: "Launcher", Special: "Hell"},
		"Banana Cannon": {Name: "Banana Cannon", MinAtp: 450, MaxAtp: 550, Ata: 80, Grind: 50, MaxHit: 100, MaxAttr: 100, Special: "Blizzard", Animation: "Launcher"},

		"Bomb-Chu":     {Name: "Bomb-Chu", MinAtp: 300, MaxAtp: 450, Ata: 250, Grind: 0, ComboPreset: Combo{Attack1Hits: 1, Attack2: "NONE", Attack3: "NONE"}, Special: "Heat", MaxHit: 100, MaxAttr: 100, Animation: "Card"},
		"Whitill Card": {Name: "Whitill Card", MinAtp: 200, MaxAtp: 250, Ata: 46, Grind: 0, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 1, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Card"},
	}
	ultimaWeapons := make([]Weapon, len(vanillaWeapons))
	copy(ultimaWeapons, vanillaWeapons)
	for i, weapon := range ultimaWeapons {
		if modifiedWeapon, found := weaponsWithModifiedStats[weapon.Name]; found {
			ultimaWeapons[i] = modifiedWeapon
		}
	}
	for _, v := range customWeapons {
		ultimaWeapons = append(ultimaWeapons, v)
	}
	return ultimaWeapons
}
