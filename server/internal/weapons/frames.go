package weapons

type Frame struct {
	Name string `json:"name"`
	Atp  int    `json:"atp"`
	Ata  int    `json:"ata"`
}

func GetFrames() []Frame {
	return []Frame{
		{Name: "None", Atp: 0, Ata: 0},
		{Name: "Thirteen", Atp: 0, Ata: 0},
		{Name: "D-Parts ver1.01", Atp: 35, Ata: 0},
		{Name: "Crimson Coat", Atp: 0, Ata: 0},
		{Name: "Samurai Armor", Atp: 0, Ata: 0},
		{Name: "Sweetheart (1)", Atp: 0, Ata: 0},
		{Name: "Sweetheart (2)", Atp: 0, Ata: 0},
		{Name: "Sweetheart (3)", Atp: 0, Ata: 0},
	}
}
