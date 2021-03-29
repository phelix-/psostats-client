package player_test

import (
	"testing"

	"github.com/phelix-/psostats/v2/pkg/pso/player"
)

func TestParsePlayerMemory_Class(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assert("HUcast", playerInfo.Class, t)
}

func TestParsePlayerMemory_Room(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertU16(5, playerInfo.Room, t)
}

func TestParsePlayerMemory_ShiftaLvl(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertI16(0, playerInfo.ShiftaLvl, t)
}

func TestParsePlayerMemory_DebandLvl(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertI16(0, playerInfo.DebandLvl, t)
}

func TestParsePlayerMemory_MaxHP(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertU16(2012, playerInfo.MaxHP, t)
}

func TestParsePlayerMemory_MaxTP(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertU16(0, playerInfo.MaxTP, t)
}

func TestParsePlayerMemory_HP(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertU16(2012, playerInfo.HP, t)
}

func TestParsePlayerMemory_TP(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertU16(0, playerInfo.TP, t)
}

func TestParsePlayerMemory_PB(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertF32(0, playerInfo.PB, t)
}

func TestParsePlayerMemory_Floor(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertU16(1, playerInfo.Floor, t)
}

func TestParsePlayerMemory_InvincibilityFrames(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertU32(0, playerInfo.InvincibilityFrames, t)
}

func TestParsePlayerMemory_DamageTraps(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertU16(20, playerInfo.DamageTraps, t)
}

func TestParsePlayerMemory_FreezeTraps(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertU16(20, playerInfo.FreezeTraps, t)
}

func TestParsePlayerMemory_ConfuseTraps(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertU16(20, playerInfo.ConfuseTraps, t)
}

func TestParsePlayerMemory_Meseta(t *testing.T) {
	playerInfo := getExamplePlayerData()
	assertU32(370050, playerInfo.Meseta, t)
}

func TestIsLowered_full(t *testing.T) {
	playerInfo := player.BasePlayerInfo{
		HP:    2012,
		MaxHP: 2012,
	}
	if playerInfo.IsLowered() {
		t.Logf("2012/2012 HP was detected as lowered")
		t.Fail()
	}
}

func TestIsLowered_lowered(t *testing.T) {
	playerInfo := player.BasePlayerInfo{
		HP:    1812,
		MaxHP: 2012,
	}
	if !playerInfo.IsLowered() {
		t.Logf("1812/2012 HP was not detected as lowered")
		t.Fail()
	}
}

func assertF32(expected float32, actual float32, t *testing.T) {
	if expected != actual {
		t.Logf("Expected '%v' but got '%v'", expected, actual)
		t.Fail()
	}
}

func assertU32(expected uint32, actual uint32, t *testing.T) {
	if expected != actual {
		t.Logf("Expected '%v' but got '%v'", expected, actual)
		t.Fail()
	}
}

func assertI16(expected int16, actual int16, t *testing.T) {
	if expected != actual {
		t.Logf("Expected '%v' but got '%v'", expected, actual)
		t.Fail()
	}
}

func assertU16(expected uint16, actual uint16, t *testing.T) {
	if expected != actual {
		t.Logf("Expected '%v' but got '%v'", expected, actual)
		t.Fail()
	}
}

func assert(expected string, actual string, t *testing.T) {
	if expected != actual {
		t.Logf("Expected '%v' but got '%v'", expected, actual)
		t.Fail()
	}
}

func getExamplePlayerData() player.BasePlayerInfo {
	playerData := []uint16{0x5, 0x0, 0xffff, 0x0, 0x0, 0x1000, 0x50, 0x14db,
		0xc106, 0x40a4, 0x0, 0xc000, 0x8cef, 0xc49f, 0xc106, 0x40a4,
		0x0, 0xc000, 0x8cef, 0xc49f, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0xe141, 0xffff, 0x0, 0x0,
		0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x3f80, 0xc106, 0x40a4,
		0x8f5c, 0x40e0, 0x8cef, 0xc49f, 0x0, 0x0, 0x0, 0x0,
		0x6598, 0xaa5, 0x2b0, 0x0, 0x6108, 0xaa5, 0xffff, 0x1,
		0x0, 0x0, 0x6080, 0xaa5, 0x0, 0x0, 0x6738, 0xaa5,
		0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0,
		0x8, 0x8, 0xf, 0x0, 0x0, 0x4080, 0x0, 0x41a0,
		0x0, 0x3f80, 0xcccd, 0x3e4c, 0x999a, 0x3f99, 0x0, 0x0,
		0xede0, 0x93, 0xd60c, 0xaa2, 0xa588, 0x1b77, 0xc8a0, 0x1592,
		0xf28, 0x1591, 0x0, 0x3f80, 0x0, 0x0, 0x1, 0x1,
		0xf28, 0x1591, 0xffff, 0xffff, 0xffff, 0xffff, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x6401, 0x6404, 0x6405, 0x4,
		0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0xffff, 0x0,
		0x23dc, 0x1514, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x37, 0x0, 0x44, 0x0, 0x6872, 0x6bf4, 0xfe8e, 0x7359,
		0xcfb, 0x2991, 0xab48, 0xfddb, 0xd867, 0xe753, 0xe5bc, 0xf322,
		0x4bfe, 0x2333, 0xa3e9, 0x7a19, 0x1db1, 0xb444, 0xaea1, 0x9b5a,
		0xcf78, 0xe3d2, 0x638c, 0x129d, 0x3db5, 0x4b74, 0x135d, 0x4a07,
		0x159c, 0x53ca, 0x89bb, 0x379e, 0xbb07, 0xb55a, 0x2f0f, 0xa394,
		0x1d6, 0xf392, 0xb928, 0x30a3, 0xb2c5, 0x707d, 0xfbff, 0x1cc9,
		0x9e94, 0x73f2, 0x7802, 0xdb03, 0x20f9, 0x7b02, 0x2ae1, 0x6ff5,
		0x1dd4, 0xcc7e, 0xdab4, 0xcc9e, 0x55fc, 0xa675, 0x480e, 0x573a,
		0xc936, 0x36dc, 0x8ba2, 0x2574, 0xc720, 0x591e, 0x7c12, 0xe476,
		0xb7a, 0xd28e, 0x853, 0x796a, 0xa9e5, 0x9f01, 0xe159, 0x6dbd,
		0x8c96, 0x87f8, 0xc405, 0xe458, 0x2507, 0x6d4a, 0xa88b, 0x13b3,
		0xec85, 0x1930, 0x3996, 0xdc5b, 0x40cd, 0xdab4, 0xc8e7, 0x4caf,
		0xfc11, 0xdc86, 0x6a53, 0xdf7a, 0x83f4, 0x4e18, 0x5e7b, 0x8645,
		0x211, 0x3c6b, 0x88e6, 0x9eee, 0x95bc, 0xaa25, 0xe1bc, 0xb54f,
		0x396c, 0xdd0d, 0xffff, 0x0, 0x0, 0x0, 0x0, 0x3f80,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6e0c, 0xaa5,
		0x0, 0x0, 0x7dc, 0x0, 0xfa, 0x0, 0x65d, 0x41,
		0x3e9, 0x38d, 0x65d, 0x41, 0x3e9, 0x38d, 0x15b, 0x64,
		0xffff, 0xffff, 0xffff, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20,
		0x20, 0x20, 0x15, 0xd, 0xc106, 0x40a4, 0x8f5c, 0x40e0,
		0x8cef, 0xc49f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x6e30, 0xaa5, 0x6e60, 0xaa5, 0x6e80, 0xaa5,
		0x2, 0x0, 0x1, 0x1, 0x0, 0x4, 0x7dc, 0x0,
		0x0, 0x0, 0x0, 0x161, 0x0, 0x0, 0x0, 0x0,
		0x1, 0x0, 0x9440, 0xb3, 0x0, 0x0, 0x0, 0x0,
		0xffff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0xffff, 0x0, 0x0, 0x0, 0xffff, 0xffff, 0x0, 0x0,
		0x0, 0x0, 0x2, 0x0, 0x9458, 0xb3, 0x6a, 0x1000,
		0x1c, 0x0, 0xd600, 0xaa2, 0x0, 0xffff, 0xc2b0, 0xaa4,
		0xcce0, 0xaa4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0xffff, 0x0, 0x0, 0x0, 0xffff, 0xffff,
		0x0, 0x1200, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x2, 0x0, 0x13, 0x0, 0x1, 0x0, 0x0, 0x0,
		0x0, 0x0, 0xd1bb, 0xffff, 0x0, 0x0, 0xfb09, 0x408a,
		0x0, 0x0, 0x8afb, 0xc49f, 0x5d15, 0xbf2f, 0x0, 0x0,
		0x80eb, 0x3f3a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x9, 0x45, 0x70, 0x68, 0x65, 0x6c, 0x69, 0x78,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3f80, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x6eb0, 0xaa5, 0xff, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3fc0, 0xffff, 0x0, 0xffff, 0x0,
		0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0,
		0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0,
		0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff,
		0xffff, 0xff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff, 0xffff,
		0xffff, 0xffff, 0xffff, 0xaff, 0x1d1d, 0x1d1d, 0x1d1d, 0x1d1d,
		0x1d1d, 0x1d1d, 0x1d1d, 0x1d1d, 0x1d1d, 0x1d, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x3f80, 0x0, 0x0,
		0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0,
		0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0,
		0xffff, 0x0, 0xffff, 0x0, 0x0, 0x0, 0x0, 0x0,
		0xc800, 0xc800, 0xc800, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x1, 0x0, 0x102, 0x0, 0x2, 0x0,
		0x202, 0x0, 0x104, 0x0, 0x804, 0x0, 0x204, 0x0,
		0x904, 0x0, 0xe04, 0x0, 0xd04, 0x0, 0x1004, 0x0,
		0x0, 0x1, 0x4, 0x0, 0x0, 0x1, 0x1, 0x0,
		0x0, 0x0, 0xb04, 0x0, 0x0, 0x1, 0x104, 0x0,
		0x804, 0x0, 0x204, 0x0, 0x904, 0x0, 0xe04, 0x0,
		0xd04, 0x0, 0x1004, 0x0, 0x0, 0x1, 0x4, 0x0,
		0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x105, 0x0, 0x505, 0x0, 0x305, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xcda0, 0xaa4,
		0x0, 0x0, 0xffff, 0xffff, 0xffff, 0xffff, 0xe75c, 0xaa,
		0x0, 0x0, 0xffff, 0xffff, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0xcccd, 0x3ecc, 0xcccd, 0x3ecc,
		0xcccd, 0x3ecc, 0xcccd, 0x3ecc, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x3fc0, 0x0, 0x0, 0x0, 0x3fc0, 0x4, 0x0,
		0x0, 0x0, 0xba16, 0xffff, 0x8, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0xffff, 0xffff,
		0x800, 0x0, 0x6666, 0x3f26, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x3f80,
		0x0, 0x3f80, 0x0, 0xffff, 0x2, 0x0, 0x0, 0x0,
		0xffff, 0xffff, 0x9454, 0xb3, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3f80,
		0x0, 0x3f80, 0x0, 0x3f80, 0xffff, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x47ae, 0x4110, 0x0, 0x0, 0x0, 0x40a0,
		0x0, 0x40a0, 0x0, 0x0, 0x200b, 0x0, 0xc106, 0x40a4,
		0x8f5c, 0x40e0, 0x8cef, 0xc49f, 0x0, 0x0, 0xd, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x4120, 0x2000, 0xc4a1,
		0x0, 0x4120, 0x0, 0xc4a1, 0x0, 0x4120, 0xe000, 0xc4a0,
		0x0, 0x4120, 0xc000, 0xc4a0, 0x0, 0x4120, 0xa000, 0xc4a0,
		0x0, 0x4120, 0x8000, 0xc4a0, 0x0, 0x4120, 0x6000, 0xc4a0,
		0x0, 0x4120, 0x4000, 0xc4a0, 0xafbd, 0x4114, 0x15ab, 0xc4a0,
		0x516, 0x4107, 0xee36, 0xc49f, 0xaffb, 0x40ee, 0xc9dd, 0xc49f,
		0xa94c, 0x40cb, 0xa90b, 0xc49f, 0xc106, 0x40a4, 0x8cef, 0xc49f,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3f80, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0xffff, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x1414, 0x1400, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x3f80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x41a0, 0x0, 0x3f80,
		0x0, 0x3f80, 0x0, 0x0, 0x0, 0x0, 0xffff, 0x0,
		0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0,
		0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0, 0xffff, 0x0,
		0xffff, 0x0, 0x0, 0x0, 0x0, 0x0, 0xffff, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0xffff, 0x0, 0xc106, 0x40a4, 0x22d1, 0x419d,
		0x8cef, 0xc49f, 0x0, 0x0, 0x2020, 0x2020, 0x2020, 0x2020,
		0x3234, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0xffff, 0xffff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x209, 0x300, 0x31, 0x0,
		0x0, 0x17, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0,
		0xf256, 0x3f0b, 0x70a4, 0x3f7d, 0x9, 0x45, 0x70, 0x68,
		0x65, 0x6c, 0x69, 0x78, 0x0, 0x0, 0x0, 0x0,
		0x544b, 0x7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x6070, 0x14db, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x3006, 0xffff, 0x6ad0, 0xaa5, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0xffff, 0xffff, 0xffff, 0xffff,
		0x6401, 0x6404, 0x6405, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0xffff, 0x4, 0x2801, 0x716,
		0xffff, 0xffff, 0x0, 0x0, 0x1640, 0x14de, 0x0, 0x0,
		0x0, 0x4198, 0xffff, 0x0, 0xffff, 0x0, 0x0, 0x3f80,
		0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x3f80,
		0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x0, 0xfca8, 0x5c,
		0x44a8, 0x5e, 0xf840, 0x5c, 0x0, 0x0, 0x0, 0x0,
		0x6da8, 0xaa5, 0x60, 0x0, 0x7a80, 0xb4, 0xa290, 0xac,
		0x20, 0x0, 0x6ad0, 0xaa5, 0x0, 0x0, 0x6880, 0xaa5,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x1, 0x0, 0x7870, 0x4374, 0x5f62, 0x3230, 0x5f78, 0x5f77,
		0x616b, 0x316f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x7870, 0x4374, 0x5f62, 0x3230, 0x5f78, 0x5f77,
		0x616b, 0x326f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x7070, 0xaa5, 0x6de0, 0xaa5, 0xd7e0, 0xaa2,
		0x0, 0x0, 0xffff, 0x0, 0x1, 0x0, 0x0, 0x0,
		0xffff, 0x0, 0xffff, 0x0, 0x0, 0x8000, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xffff, 0x0,
		0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x4040, 0x0, 0x0, 0x0, 0x0, 0x1f0b, 0x3,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x4000, 0x5bd0, 0x2e9,
		0x0, 0x0, 0x0, 0x4000, 0x5be8, 0x2e9, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x300, 0x4300, 0x401, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x6bfc, 0xaa5, 0x2, 0x0, 0x1, 0x0,
		0x7870, 0x5a74, 0x5f63, 0x3130, 0x5f6a, 0x5f74, 0x6f63, 0x6572,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x7870, 0x4374, 0x5f62, 0x3130, 0x5f78, 0x5f77, 0x7568, 0x756b,
		0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x7870, 0x4374, 0x5f62, 0x3130, 0x5f78, 0x5f77, 0x7568, 0x756b,
		0x32, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x7870, 0x4374, 0x5f62, 0x3130, 0x5f78, 0x5f77, 0x7568, 0x756b,
		0x33, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x7870, 0x4374, 0x5f74, 0x3130, 0x5f78, 0x5f77, 0x6574, 0x31,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x2, 0x1, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0xff, 0x1, 0x903, 0x60a,
		0xffff, 0xffff, 0x0, 0x0, 0x200c, 0x14de, 0x0, 0x0,
		0x0, 0x4198, 0xe, 0x0, 0xffff, 0x0, 0x0, 0x3f80,
		0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x3f80,
		0x0, 0x3f80, 0x0, 0x3f80, 0x0, 0x0, 0xf288, 0x5c,
		0x0, 0x0, 0x0, 0x4000, 0x5c18, 0x2e9, 0x0, 0x0,
		0x0, 0x4000, 0x5c30, 0x2e9, 0x0, 0x0, 0x0, 0x4000,
		0x5c48, 0x2e9, 0x0, 0x0, 0x0, 0x4000, 0x5c60, 0x2e9,
		0x0, 0x0, 0x0, 0x4000, 0x5c78, 0x2e9, 0x0, 0x0,
		0x6ff8, 0xaa5, 0x30, 0x0, 0x7a80, 0xb4, 0x43bc, 0xa9,
		0x30, 0x0, 0x7290, 0xaa5, 0x7000, 0xaa5, 0xd7b0, 0xaa2,
		0x6d68, 0xaa5, 0x5, 0x0, 0x0, 0x0, 0x7f80, 0xaa5,
		0x0, 0x0, 0xa810, 0xaa5, 0x9960, 0xaa5, 0x9b90, 0xaa5,
		0x0, 0x0, 0x0, 0x0, 0x8880, 0xaa5, 0x0, 0x0,
		0x0, 0x0, 0x9ed0, 0xaa5, 0xa5c0, 0xaa5, 0xc020, 0xaa5,
		0xc2b0, 0xaa5, 0x65d, 0x0, 0x258, 0x2aa, 0x229, 0x514,
		0x64, 0x0, 0x47ae, 0x4190, 0x0, 0x4120, 0xc7, 0x0,
		0xf498, 0x4f5, 0xa582, 0x5, 0xd7, 0x0, 0x0, 0x0,
	}
	return player.ParsePlayerMemory(playerData, 0x028)
}
