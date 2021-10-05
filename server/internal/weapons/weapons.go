package weapons

type Weapon struct {
	Name        string `json:"name"`
	MinAtp      int    `json:"minAtp"`
	MaxAtp      int    `json:"maxAtp"`
	Ata         int    `json:"ata"`
	Grind       int    `json:"grind"`
	MaxHit      int    `json:"maxHit"`
	MaxAttr     int    `json:"maxAttr"`
	ComboPreset Combo  `json:"comboPreset"`
	Special     string `json:"special"`
	Animation   string
}

type Combo struct {
	Attack1     string `json:"attack1"`
	Attack1Hits int    `json:"attack1Hits"`
	Attack2     string `json:"attack2"`
	Attack2Hits int    `json:"attack2Hits"`
	Attack3     string `json:"attack3"`
	Attack3Hits int    `json:"attack3Hits"`
}

type PsoClass struct {
	Name string
	Atp  int
	Ata  int
}

func GetClasses() []PsoClass {
	return []PsoClass{
		{Name: "HUmar", Atp: 1397, Ata: 200},
		{Name: "HUnewearl", Atp: 1237, Ata: 199},
		{Name: "HUcast", Atp: 1639, Ata: 191},
		{Name: "HUcaseal", Atp: 1301, Ata: 218},
		{Name: "RAmar", Atp: 1260, Ata: 249},
		{Name: "RAmarl", Atp: 1145, Ata: 241},
		{Name: "RAcast", Atp: 1350, Ata: 224},
		{Name: "RAcaseal", Atp: 1175, Ata: 231},
		{Name: "FOmar", Atp: 1002, Ata: 163},
		{Name: "FOmarl", Atp: 872, Ata: 170},
		{Name: "FOnewm", Atp: 814, Ata: 180},
		{Name: "FOnewearl", Atp: 583, Ata: 186},
	}
}

func GetWeapons() []Weapon {
	return []Weapon{
		{Name: "Unarmed", MinAtp: 0, MaxAtp: 0, Ata: 0, Grind: 0, MaxHit: 0, MaxAttr: 0, Animation: "Fist", Special: ""},

		{Name: "Saber", MinAtp: 40, MaxAtp: 55, Ata: 30, Grind: 35, MaxHit: 100, MaxAttr: 100, Animation: "Saber", Special: "Charge"},
		{Name: "Brand", MinAtp: 80, MaxAtp: 100, Ata: 33, Grind: 32, MaxHit: 100, MaxAttr: 100, Animation: "Saber", Special: "Charge"},
		{Name: "Buster", MinAtp: 120, MaxAtp: 160, Ata: 35, Grind: 30, MaxHit: 100, MaxAttr: 100, Animation: "Saber", Special: "Charge"},
		{Name: "Pallasch", MinAtp: 170, MaxAtp: 220, Ata: 38, Grind: 26, MaxHit: 100, MaxAttr: 100, Animation: "Saber", Special: "Charge"},
		{Name: "Gladius", MinAtp: 240, MaxAtp: 280, Ata: 40, Grind: 18, MaxHit: 100, MaxAttr: 100, Animation: "Saber", Special: "Charge"},
		{Name: "Battledore", MinAtp: 1, MaxAtp: 1, Ata: 1, Grind: 0, MaxHit: 100, MaxAttr: 100, Animation: "Saber"},
		{Name: "Red Saber", MinAtp: 450, MaxAtp: 489, Ata: 51, Grind: 78, MaxHit: 100, MaxAttr: 100, Animation: "Saber", Special: "Fill"},
		//{Name: "Lame d'Argent", MinAtp: 430, MaxAtp: 465, Ata: 40, Grind: 35, MaxHit: 100, MaxAttr: 100, Animation: "Saber"},
		{Name: "Lavis Cannon", MinAtp: 730, MaxAtp: 750, Ata: 54, Grind: 0, MaxHit: 100, MaxAttr: 100, Animation: "Saber", Special: "Lavis"},
		{Name: "Excalibur", MinAtp: 900, MaxAtp: 950, Ata: 60, Grind: 0, Special: "Berserk", MaxHit: 100, MaxAttr: 100, Animation: "Saber"},
		{Name: "Galatine", MinAtp: 990, MaxAtp: 1260, Ata: 77, Grind: 9, Special: "Spirit", MaxHit: 100, MaxAttr: 100, Animation: "Saber"},
		{Name: "ES Saber", MinAtp: 150, MaxAtp: 150, Ata: 50, Grind: 250, MaxHit: 0, MaxAttr: 0, Animation: "Saber", Special: "Berserk"},
		{Name: "ES Axe", MinAtp: 200, MaxAtp: 200, Ata: 50, Grind: 250, MaxHit: 0, MaxAttr: 0, Animation: "Saber", Special: "Berserk"},

		{Name: "Sword", MinAtp: 25, MaxAtp: 60, Ata: 15, Grind: 46, Animation: "Sword", Special: "Charge", MaxHit: 100, MaxAttr: 100},
		{Name: "Gigush", MinAtp: 55, MaxAtp: 100, Ata: 18, Grind: 32, Animation: "Sword", Special: "Charge", MaxHit: 100, MaxAttr: 100},
		{Name: "Breaker", MinAtp: 100, MaxAtp: 150, Ata: 20, Grind: 18, Animation: "Sword", Special: "Charge", MaxHit: 100, MaxAttr: 100},
		{Name: "Claymore", MinAtp: 150, MaxAtp: 200, Ata: 23, Grind: 16, Animation: "Sword", Special: "Charge", MaxHit: 100, MaxAttr: 100},
		{Name: "Calibur", MinAtp: 210, MaxAtp: 255, Ata: 25, Grind: 10, Animation: "Sword", Special: "Charge", MaxHit: 100, MaxAttr: 100},
		{Name: "Flowen's Sword (3084)", MinAtp: 300, MaxAtp: 320, Ata: 34, Grind: 85, Animation: "Sword", Special: "Spirit", MaxHit: 100, MaxAttr: 100},
		{Name: "Red Sword", MinAtp: 400, MaxAtp: 611, Ata: 37, Grind: 52, Animation: "Sword", Special: "Arrest", MaxHit: 100, MaxAttr: 100},
		{Name: "Chain Sawd", MinAtp: 500, MaxAtp: 525, Ata: 36, Grind: 15, Animation: "Sword", Special: "Gush", MaxHit: 100, MaxAttr: 100},
		{Name: "Zanba", MinAtp: 310, MaxAtp: 438, Ata: 38, Grind: 38, Special: "Berserk", Animation: "Sword", MaxHit: 100, MaxAttr: 100},
		{Name: "Sealed J-Sword", MinAtp: 420, MaxAtp: 525, Ata: 35, Grind: 0, Special: "Hell", Animation: "Sword", MaxHit: 100, MaxAttr: 100},
		{Name: "Laconium Axe", MinAtp: 700, MaxAtp: 750, Ata: 40, Grind: 25, Animation: "Sword", Special: "Berserk", MaxHit: 100, MaxAttr: 100},
		{Name: "Dark Flow", MinAtp: 756, MaxAtp: 900, Ata: 50, Grind: 0, ComboPreset: Combo{Attack1: "S", Attack1Hits: 5, Attack2: "NONE", Attack3: "NONE"}, Animation: "Sword", Special: "Dark Flow", MaxHit: 100, MaxAttr: 100},
		{Name: "Tsumikiri J-Sword", MinAtp: 700, MaxAtp: 756, Ata: 40, Grind: 50, Animation: "Sword", Special: "TJS", MaxHit: 100, MaxAttr: 100},
		{Name: "TypeSW/J-Sword", MinAtp: 100, MaxAtp: 150, Ata: 40, Grind: 125, Animation: "Sword", Special: "Spirit", MaxHit: 100, MaxAttr: 100},
		{Name: "ES Sword", MinAtp: 200, MaxAtp: 200, Ata: 35, Grind: 250, MaxHit: 0, MaxAttr: 0, Animation: "Sword", Special: "Berserk"},

		{Name: "Dagger", MinAtp: 25, MaxAtp: 40, Ata: 20, Grind: 65, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxHit: 100, MaxAttr: 100},
		{Name: "Knife", MinAtp: 50, MaxAtp: 70, Ata: 22, Grind: 50, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxHit: 100, MaxAttr: 100},
		{Name: "Blade", MinAtp: 80, MaxAtp: 100, Ata: 24, Grind: 35, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxHit: 100, MaxAttr: 100},
		{Name: "Edge", MinAtp: 105, MaxAtp: 130, Ata: 26, Grind: 25, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxHit: 100, MaxAttr: 100},
		{Name: "Ripper", MinAtp: 125, MaxAtp: 160, Ata: 28, Grind: 15, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxHit: 100, MaxAttr: 100},
		{Name: "S-Beat's Blade", MinAtp: 210, MaxAtp: 220, Ata: 35, Grind: 15, MaxHit: 50, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxAttr: 100},
		{Name: "P-Arms' Blade", MinAtp: 250, MaxAtp: 270, Ata: 34, Grind: 25, MaxHit: 50, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxAttr: 100},
		{Name: "Red Dagger", MinAtp: 245, MaxAtp: 280, Ata: 35, Grind: 65, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxHit: 100, MaxAttr: 100},
		{Name: "S-Red's Blade", MinAtp: 340, MaxAtp: 350, Ata: 39, Grind: 15, MaxHit: 50, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxAttr: 100},
		{Name: "Two Kamui", MinAtp: 600, MaxAtp: 650, Ata: 50, Grind: 0, MaxHit: 0, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxAttr: 100},
		{Name: "Lavis Blade", MinAtp: 380, MaxAtp: 450, Ata: 40, Grind: 0, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxHit: 100, MaxAttr: 100},
		{Name: "Daylight Scar", MinAtp: 500, MaxAtp: 550, Ata: 48, Grind: 25, Special: "Berserk", ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", MaxHit: 100, MaxAttr: 100},
		{Name: "ES Blade", MinAtp: 10, MaxAtp: 10, Ata: 35, Grind: 200, MaxHit: 0, MaxAttr: 0, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Dagger", Special: "Berserk"},

		{Name: "Gungnir", MinAtp: 150, MaxAtp: 180, Ata: 32, Grind: 10, Animation: "Partisan", MaxHit: 100, MaxAttr: 100},
		{Name: "Vjaya", MinAtp: 160, MaxAtp: 220, Ata: 36, Grind: 15, Special: "Vjaya", Animation: "Partisan", MaxHit: 100, MaxAttr: 100},
		{Name: "Tyrell's Parasol", MinAtp: 250, MaxAtp: 300, Ata: 40, Grind: 0, Animation: "Partisan", MaxHit: 100, MaxAttr: 100},
		{Name: "Madam's Umbrella", MinAtp: 210, MaxAtp: 280, Ata: 40, Grind: 0, Animation: "Partisan", MaxHit: 100, MaxAttr: 100},
		{Name: "Plantain Huge Fan", MinAtp: 265, MaxAtp: 300, Ata: 38, Grind: 9, Animation: "Partisan", MaxHit: 100, MaxAttr: 100},
		{Name: "Asteron Belt", MinAtp: 380, MaxAtp: 400, Ata: 55, Grind: 9, Animation: "Partisan", MaxHit: 100, MaxAttr: 100},
		{Name: "Yunchang", MinAtp: 300, MaxAtp: 350, Ata: 49, Grind: 25, Animation: "Partisan", MaxHit: 100, MaxAttr: 100},
		{Name: "ES Partisan", MinAtp: 10, MaxAtp: 10, Ata: 40, Grind: 200, MaxHit: 0, MaxAttr: 0, Animation: "Partisan"},
		{Name: "ES Scythe", MinAtp: 10, MaxAtp: 10, Ata: 40, Grind: 180, MaxHit: 0, MaxAttr: 0, Animation: "Partisan"},

		{Name: "Diska", MinAtp: 85, MaxAtp: 105, Ata: 25, Grind: 10, Special: "Charge", Animation: "Slicer", MaxHit: 100, MaxAttr: 100},
		{Name: "Diska of Braveman", MinAtp: 150, MaxAtp: 167, Ata: 31, Grind: 9, Special: "Berserk", Animation: "Slicer", MaxHit: 100, MaxAttr: 100},
		{Name: "Slicer of Fanatic", MinAtp: 340, MaxAtp: 360, Ata: 40, Grind: 30, Special: "Demon's", Animation: "Slicer", MaxHit: 100, MaxAttr: 100},
		{Name: "Red Slicer", MinAtp: 190, MaxAtp: 200, Ata: 38, Grind: 45, Animation: "Slicer", MaxHit: 100, MaxAttr: 100},
		{Name: "Rainbow Baton", MinAtp: 300, MaxAtp: 320, Ata: 40, Grind: 24, Animation: "Slicer", MaxHit: 100, MaxAttr: 100},
		{Name: "ES Slicer", MinAtp: 10, MaxAtp: 10, Ata: 35, Grind: 140, MaxHit: 0, MaxAttr: 0, Special: "Berserk", Animation: "Slicer"},
		{Name: "ES J-Cutter", MinAtp: 25, MaxAtp: 25, Ata: 35, Grind: 150, MaxHit: 0, MaxAttr: 0, Special: "Berserk", Animation: "Slicer"},

		{Name: "Demolition Comet", MinAtp: 530, MaxAtp: 530, Ata: 38, Grind: 25, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Double Saber", MaxAttr: 100},
		{Name: "Girasole", MinAtp: 500, MaxAtp: 550, Ata: 50, Grind: 0, MaxHit: 0, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Double Saber", MaxAttr: 100},
		{Name: "Twin Blaze", MinAtp: 300, MaxAtp: 520, Ata: 40, Grind: 9, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Double Saber", MaxAttr: 100},
		{Name: "Meteor Cudgel", MinAtp: 300, MaxAtp: 560, Ata: 42, Grind: 15, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Double Saber", MaxAttr: 100},
		{Name: "Vivienne", MinAtp: 575, MaxAtp: 590, Ata: 49, Grind: 50, MaxHit: 0, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Double Saber", MaxAttr: 100},
		{Name: "Black King Bar", MinAtp: 590, MaxAtp: 600, Ata: 43, Grind: 80, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Double Saber", MaxAttr: 100},
		{Name: "Double Cannon", MinAtp: 620, MaxAtp: 650, Ata: 45, Grind: 0, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Double Saber", MaxAttr: 100},
		{Name: "ES Twin", MinAtp: 50, MaxAtp: 50, Ata: 40, Grind: 250, MaxHit: 0, MaxAttr: 0, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Double Saber"},

		{Name: "Toy Hammer", MinAtp: 1, MaxAtp: 400, Ata: 53, Grind: 0, MaxHit: 100, MaxAttr: 100, Animation: "Katana"},
		{Name: "Raikiri", MinAtp: 550, MaxAtp: 560, Ata: 30, Grind: 0, MaxHit: 100, MaxAttr: 100, Animation: "Katana"},
		{Name: "Orotiagito", MinAtp: 750, MaxAtp: 800, Ata: 55, Grind: 0, MaxHit: 0, MaxAttr: 100, Animation: "Katana"},

		{Name: "Musashi", MinAtp: 330, MaxAtp: 350, Ata: 35, Grind: 40, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 2, Attack3Hits: 2}, MaxHit: 100, MaxAttr: 100, Animation: "Twin Sword"},
		{Name: "Yamato", MinAtp: 380, MaxAtp: 390, Ata: 40, Grind: 60, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 2, Attack3Hits: 2}, MaxHit: 100, MaxAttr: 100, Animation: "Twin Sword"},
		{Name: "G-Assassin's Sabers", MinAtp: 350, MaxAtp: 360, Ata: 35, Grind: 25, MaxHit: 50, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 2, Attack3Hits: 2}, MaxAttr: 100, Animation: "Twin Sword"},
		{Name: "Asuka", MinAtp: 560, MaxAtp: 570, Ata: 50, Grind: 30, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 2, Attack3Hits: 2}, MaxHit: 100, MaxAttr: 100, Animation: "Twin Sword"},
		{Name: "Sange & Yasha", MinAtp: 640, MaxAtp: 650, Ata: 50, Grind: 30, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 2, Attack3Hits: 2}, MaxHit: 100, MaxAttr: 100, Animation: "Twin Sword"},
		{Name: "Jizai", MinAtp: 800, MaxAtp: 810, Ata: 55, Grind: 40, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 2, Attack3Hits: 2}, MaxHit: 100, MaxAttr: 100, Animation: "Twin Sword"},
		{Name: "TypeSS/Swords", MinAtp: 150, MaxAtp: 150, Ata: 45, Grind: 125, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 2, Attack3Hits: 2}, MaxHit: 100, MaxAttr: 100, Animation: "Twin Sword"},
		{Name: "ES Swords", MinAtp: 180, MaxAtp: 180, Ata: 45, Grind: 250, MaxHit: 0, MaxAttr: 0, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 2, Attack3Hits: 2}, Animation: "Twin Sword"},

		{Name: "Raygun", MinAtp: 150, MaxAtp: 180, Ata: 35, Grind: 15, Special: "Charge", MaxHit: 100, MaxAttr: 100, Animation: "Handgun"},
		{Name: "Master Raven", MinAtp: 350, MaxAtp: 380, Ata: 52, Grind: 9, MaxHit: 0, ComboPreset: Combo{Attack1Hits: 3, Attack2: "NONE", Attack3: "NONE"}, MaxAttr: 100, Animation: "Master Raven"},
		{Name: "Last Swan", MinAtp: 80, MaxAtp: 90, Ata: 32, Grind: 9, MaxHit: 0, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxAttr: 100, Animation: "Last Swan"},
		{Name: "Heaven Striker", MinAtp: 550, MaxAtp: 600, Ata: 55, Grind: 20, MaxHit: 100, MaxAttr: 100, Animation: "Handgun"},

		{Name: "Laser", MinAtp: 200, MaxAtp: 210, Ata: 50, Grind: 25, MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		{Name: "Spread Needle", MinAtp: 1, MaxAtp: 110, Ata: 40, Grind: 40, MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		{Name: "Bringer's Rifle", MinAtp: 330, MaxAtp: 370, Ata: 63, Grind: 9, MaxHit: 50, Special: "Demon's", MaxAttr: 100, Animation: "Rifle"},
		{Name: "Frozen Shooter", MinAtp: 240, MaxAtp: 250, Ata: 60, Grind: 9, MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		{Name: "Snow Queen", MinAtp: 330, MaxAtp: 350, Ata: 60, Grind: 18, ComboPreset: Combo{Attack2: "NONE", Attack3: "NONE"}, MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		{Name: "Holy Ray", MinAtp: 290, MaxAtp: 300, Ata: 70, Grind: 40, MaxHit: 100, MaxAttr: 100, Animation: "Rifle"},
		{Name: "ES Rifle", MinAtp: 10, MaxAtp: 10, Ata: 60, Grind: 220, MaxHit: 0, MaxAttr: 0, Special: "Berserk", Animation: "Rifle"},
		{Name: "ES Needle", MinAtp: 10, MaxAtp: 10, Ata: 40, Grind: 70, MaxHit: 0, MaxAttr: 0, Special: "Berserk", Animation: "Rifle"},

		{Name: "Mechgun", MinAtp: 2, MaxAtp: 4, Ata: 0, Grind: 9, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "Assault", MinAtp: 5, MaxAtp: 8, Ata: 3, Grind: 9, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "Repeater", MinAtp: 5, MaxAtp: 12, Ata: 6, Grind: 9, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "Gatling", MinAtp: 5, MaxAtp: 16, Ata: 9, Grind: 9, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "Vulcan", MinAtp: 5, MaxAtp: 20, Ata: 12, Grind: 9, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		//{Name: "Samba Maracas", MinAtp: 5, MaxAtp: 10, Ata: 10, Grind: 0, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		//{Name: "Rocket Punch", MinAtp: 50, MaxAtp: 300, Ata: 10, Grind: 50, ComboPreset: Combo{Attack1Hits: 3, Attack2: "NONE", Attack3: "NONE"}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "M&A60 Vise", MinAtp: 15, MaxAtp: 25, Ata: 15, Grind: 9, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		//{Name: "H&S25 Justice", MinAtp: 15, MaxAtp: 30, Ata: 18, Grind: 9, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		//{Name: "L&K14 Combat", MinAtp: 15, MaxAtp: 30, Ata: 18, Grind: 20, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		//{Name: "Twin Psychogun", MinAtp: 35, MaxAtp: 40, Ata: 23, Grind: 0, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		//{Name: "Red Mechgun", MinAtp: 50, MaxAtp: 50, Ata: 25, Grind: 30, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "Yasminkov 9000M", MinAtp: 40, MaxAtp: 80, Ata: 27, Grind: 10, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "Rage de Feu", MinAtp: 175, MaxAtp: 185, Ata: 40, Grind: 9, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "Guld Milla", MinAtp: 180, MaxAtp: 200, Ata: 30, Grind: 9, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "Mille Marteaux", MinAtp: 200, MaxAtp: 220, Ata: 45, Grind: 12, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "Dual Bird", MinAtp: 200, MaxAtp: 210, Ata: 22, Grind: 21, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 0, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "TypeME/Mechgun", MinAtp: 10, MaxAtp: 10, Ata: 20, Grind: 30, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, MaxHit: 40, MaxAttr: 100, Animation: "Mechgun"},
		{Name: "ES Mechgun", MinAtp: 10, MaxAtp: 10, Ata: 20, Grind: 50, MaxHit: 0, MaxAttr: 0, ComboPreset: Combo{Attack1Hits: 3, Attack2Hits: 3, Attack3Hits: 3}, Animation: "Mechgun"},
		{Name: "ES Psychogun", MinAtp: 10, MaxAtp: 10, Ata: 20, Grind: 50, MaxHit: 0, MaxAttr: 0, Animation: "Mechgun"},
		{Name: "ES Punch", MinAtp: 10, MaxAtp: 10, Ata: 40, Grind: 250, MaxHit: 0, MaxAttr: 0, ComboPreset: Combo{Attack1Hits: 3, Attack2: "NONE", Attack3: "NONE"}, Animation: "Mechgun"},

		{Name: "Shot", MinAtp: 20, MaxAtp: 25, Ata: 27, Grind: 20, MaxHit: 100, MaxAttr: 100, Animation: "Shot"},
		{Name: "Spread", MinAtp: 30, MaxAtp: 50, Ata: 28, Grind: 20, MaxHit: 100, MaxAttr: 100, Animation: "Shot"},
		{Name: "Cannon", MinAtp: 40, MaxAtp: 80, Ata: 30, Grind: 15, MaxHit: 100, MaxAttr: 100, Animation: "Shot"},
		{Name: "Launcher", MinAtp: 50, MaxAtp: 110, Ata: 31, Grind: 15, MaxHit: 100, MaxAttr: 100, Animation: "Shot"},
		{Name: "Arms", MinAtp: 60, MaxAtp: 140, Ata: 33, Grind: 10, MaxHit: 100, MaxAttr: 100, Animation: "Shot"},
		{Name: "L&K38 Combat", MinAtp: 150, MaxAtp: 250, Ata: 40, Grind: 25, ComboPreset: Combo{Attack1Hits: 5, Attack2: "NONE", Attack3: "NONE"}, MaxHit: 100, MaxAttr: 100, Animation: "L&K38 Combat"},
		{Name: "Rambling May", MinAtp: 360, MaxAtp: 450, Ata: 45, Grind: 0, MaxHit: 0, ComboPreset: Combo{Attack1Hits: 2, Attack2Hits: 2, Attack3Hits: 2}, MaxAttr: 100, Animation: "Shot"},
		{Name: "Baranz Launcher", MinAtp: 230, MaxAtp: 240, Ata: 40, Grind: 30, MaxHit: 50, MaxAttr: 100, Animation: "Shot"},
		{Name: "Dark Meteor", MinAtp: 150, MaxAtp: 280, Ata: 45, Grind: 25, ComboPreset: Combo{Attack2: "NONE", Attack3: "NONE"}, MaxHit: 100, MaxAttr: 100, Animation: "Shot"},
		{Name: "TypeSH/Shot", MinAtp: 10, MaxAtp: 10, Ata: 40, Grind: 60, MaxHit: 100, MaxAttr: 100, Animation: "Shot"},
		{Name: "ES Shot", MinAtp: 10, MaxAtp: 10, Ata: 40, Grind: 125, MaxHit: 0, MaxAttr: 0, Animation: "Shot"},
		{Name: "ES Bazooka", MinAtp: 10, MaxAtp: 10, Ata: 40, Grind: 250, MaxHit: 0, MaxAttr: 0, Animation: "Shot"},

		{Name: "ES Launcher", MinAtp: 10, MaxAtp: 10, Ata: 40, Grind: 180, MaxHit: 0, MaxAttr: 0, Animation: "Launcher", Special: "Berserk"},
		{Name: "Cannon Rouge", MinAtp: 600, MaxAtp: 750, Ata: 45, Grind: 30, ComboPreset: Combo{Attack1Hits: 1, Attack2: "NONE", Attack3: "NONE"}, MaxHit: 100, MaxAttr: 100, Animation: "Launcher"},

		{Name: "Gal Wind", MinAtp: 270, MaxAtp: 310, Ata: 40, Grind: 15, MaxHit: 50, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 1, Attack3Hits: 3}, MaxAttr: 100, Animation: "Card"},
		{Name: "Guardianna", MinAtp: 200, MaxAtp: 280, Ata: 40, Grind: 9, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 1, Attack3Hits: 3}, MaxHit: 100, MaxAttr: 100, Animation: "Card"},
		{Name: "ES Cards", MinAtp: 10, MaxAtp: 10, Ata: 45, Grind: 150, MaxHit: 0, MaxAttr: 0, ComboPreset: Combo{Attack1Hits: 1, Attack2Hits: 1, Attack3Hits: 3}, Animation: "Card"},
	}

}
