package weapons

func GetFramesUltima() []Frame {
	return append(GetFrames(),
		Frame{Name: "Sonicteam Armor", Atp: 0, Ata: 0},
		Frame{Name: "Sue's Coat", Atp: 100, Ata: 0},
	)
}
