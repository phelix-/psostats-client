package inventory

import (
	"errors"
	"fmt"
	"log"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/client/internal/numbers"
)

//goland:noinspection GoUnusedConst
const (
	itemArray          = 0x00A8D81C
	itemArrayCount     = 0x00A8D820
	itemId             = 0xD8
	itemOwnerOffset    = 0xE4
	itemCode           = 0xF2
	itemEquippedOffset = 0x190
	itemTypeOffset     = 0xF2
	itemGroupOffset    = 0xF3
	itemKills          = 0xE8
	itemWepGrind       = 0x1F5
	itemWepSpecial     = 0x1F6
	itemWepStats       = 0x1C8
	itemArmSlots       = 0x1B8
	itemFrameDfp       = 0x1B9
	itemFrameEvp       = 0x1BA
	itemBarrierDfp     = 0x1E4
	itemBarrierEvp     = 0x1E5
	itemUnitMod        = 0x1DC
	itemMagStats       = 0x1C0
	itemMagPBHas       = 0x1C8
	itemMagPB          = 0x1C9
	itemMagColor       = 0x1CA
	itemMagSync        = 0x1BE
	itemMagIQ          = 0x1BC
	itemMagTimer       = 0x1B4
	itemToolCount      = 0x104
	itemTechType       = 0x108
	itemMesetaAmount   = 0x100
)

const (
	WeaponBareHanded     = "Bare Handed"
	EquipmentTypeWeapon  = "Weapon"
	EquipmentTypeFrame   = "Frame"
	EquipmentTypeBarrier = "Barrier"
	EquipmentTypeUnit    = "Unit"
	EquipmentTypeMag     = "Mag"
)

type Equipment struct {
	Type    string
	Display string
}

func ReadInventory(handle w32.HANDLE, playerIndex uint8) ([]Equipment, error) {
	equipment := make([]Equipment, 0)
	buf, _, ok := w32.ReadProcessMemory(handle, uintptr(itemArrayCount), 2)
	if !ok {
		return equipment, errors.New("could not read item count")
	}
	count := numbers.Uint32From16(buf[0:2])
	buf, _, ok = w32.ReadProcessMemory(handle, uintptr(itemArray), 4)
	if !ok {
		return equipment, errors.New("could not read item array")
	}
	address := numbers.Uint32From16(buf[0:2])
	if count != 0 && address != 0 {
		buf, _, ok = w32.ReadProcessMemory(handle, uintptr(address), uintptr(4*count))
		if !ok {
			return equipment, errors.New("could not read item array")
		}

		for i := 0; i < int(count); i++ {
			itemAddr := numbers.Uint32From16(buf[i*2 : (i*2)+2])
			if itemAddr != 0 {
				itemBuffer, _, ok := w32.ReadProcessMemory(handle, uintptr(itemAddr+0xD8), 4)
				if !ok {
					return equipment, errors.New("could not read item")
				}
				itemId := fmt.Sprintf("%04x%04x", itemBuffer[1], itemBuffer[0])
				itemType := numbers.ReadU8(handle, uintptr(itemAddr+itemTypeOffset))
				itemGroup := numbers.ReadU8(handle, uintptr(itemAddr+itemGroupOffset))
				equipped := numbers.ReadU8(handle, uintptr(itemAddr+itemEquippedOffset))&0x01 == 1
				itemOwner := numbers.ReadU8(handle, uintptr(itemAddr+itemOwnerOffset))
				if itemOwner == playerIndex && equipped {
					switch itemType {
					case 0:
						weapon := readWeapon(handle, int(itemAddr), itemId, itemGroup)
						equipment = append(equipment, Equipment{Type: EquipmentTypeWeapon, Display: weapon.String()})
					case 1:
						switch itemGroup {
						case 1:
							frame := readFrame(handle, int(itemAddr), itemId, itemGroup)
							equipment = append(equipment, Equipment{Type: EquipmentTypeFrame, Display: frame.String()})
						case 2:
							barrier := readBarrier(handle, int(itemAddr), itemId, itemGroup)
							equipment = append(equipment, Equipment{Type: EquipmentTypeBarrier, Display: barrier.StringNoSlots()})
						case 3:
							unit := readUnit(handle, int(itemAddr), itemId)
							equipment = append(equipment, Equipment{Type: EquipmentTypeUnit, Display: unit.Name})
						}
					case 2:
						mag := readMag(handle, int(itemAddr), itemId, itemGroup)
						equipment = append(equipment, Equipment{Type: EquipmentTypeMag, Display: mag.String()})
					}
				}
			}
		}
	}

	return equipment, nil
}

func getWeaponIndex(handle w32.HANDLE, group uint8, index uint8, typeOffset uint8, sizeSomething uint32) uint32 {
	weaponIndex := uint32(0)
	pmtAddress := numbers.ReadU32Unchecked(handle, 0x00a8dc94)
	weaponAddress := numbers.ReadU32Unchecked(handle, uintptr(pmtAddress+uint32(typeOffset)))
	if weaponAddress != 0 {
		groupAddress := weaponAddress + (uint32(group) * 8)
		itemAddress := numbers.ReadU32Unchecked(handle, uintptr(groupAddress+4)) + (sizeSomething * uint32(index))
		weaponIndex = numbers.ReadU32Unchecked(handle, uintptr(itemAddress))
	}
	return weaponIndex
}

func readItemName(handle w32.HANDLE, index int) string {
	unitxtPointer := numbers.ReadU32Unchecked(handle, 0x00a9cd50)
	if unitxtPointer == 0 {
		return "?"
	}
	weaponIndex := numbers.ReadU32Unchecked(handle, uintptr(unitxtPointer+4))
	weaponNameAddress := numbers.ReadU32Unchecked(handle, uintptr(weaponIndex+uint32(4*index)))

	weaponName, err := numbers.ReadNullTerminatedString(handle, uintptr(weaponNameAddress))
	if err != nil {
		log.Fatalf("Error getting weapon name %v", err)
	}
	return fmt.Sprintf("%v", weaponName)
}

func readWeapon(handle w32.HANDLE, itemAddr int, itemId string, itemGroup uint8) Weapon {
	itemIndex := numbers.ReadU8(handle, uintptr(itemAddr+0xF4))
	weaponIndex := getWeaponIndex(handle, itemGroup, itemIndex, 0x00, 44)
	weapon := Weapon{
		Id: itemId,
	}
	weapon.Name = readItemName(handle, int(weaponIndex))
	weapon.Grind = numbers.ReadU8(handle, uintptr(itemAddr+itemWepGrind))
	weapon.Special = numbers.ReadU8(handle, uintptr(itemAddr+itemWepSpecial))
	weapon.SpecialName = getWeaponSpecial(weapon.Special)
	for j := 0; j < 6; j += 2 {
		area := numbers.ReadU8(handle, uintptr(itemAddr+itemWepStats+j))
		percent := numbers.ReadI8(handle, uintptr(itemAddr+itemWepStats+j+1))
		switch area {
		case 1:
			weapon.Native = percent
		case 2:
			weapon.ABeast = percent
		case 3:
			weapon.Machine = percent
		case 4:
			weapon.Dark = percent
		case 5:
			weapon.Hit = percent
		}
	}
	return weapon
}

type Weapon struct {
	Id          string
	Name        string
	Grind       uint8
	Special     uint8
	SpecialName string
	Native      int8
	ABeast      int8
	Machine     int8
	Dark        int8
	Hit         int8
}

func (w Weapon) String() string {
	grindString := ""
	if w.Grind > 0 {
		grindString = fmt.Sprintf(" +%v", w.Grind)
	}
	specialString := ""
	if len(w.SpecialName) > 0 {
		specialString = fmt.Sprintf(" [%v]", w.SpecialName)
	}
	return fmt.Sprintf("%v%v%v [%v/%v/%v/%v|%v]", w.Name, grindString, specialString, w.Native, w.ABeast, w.Machine, w.Dark, w.Hit)
}

func readFrame(handle w32.HANDLE, itemAddr int, itemId string, itemGroup uint8) Frame {
	itemIndex := numbers.ReadU8(handle, uintptr(itemAddr+0xF4))
	weaponIndex := getWeaponIndex(handle, itemGroup-1, itemIndex, 0x04, 32)
	weapon := Frame{
		Id:    itemId,
		Name:  readItemName(handle, int(weaponIndex)),
		Dfp:   numbers.ReadU8(handle, uintptr(itemAddr)+itemFrameDfp),
		Evp:   numbers.ReadU8(handle, uintptr(itemAddr)+itemFrameEvp),
		Slots: numbers.ReadU8(handle, uintptr(itemAddr)+itemArmSlots),
	}
	return weapon
}

type Frame struct {
	Id    string
	Name  string
	Dfp   uint8
	Evp   uint8
	Slots uint8
}

func (f Frame) StringNoSlots() string {
	// lazy reuse of the Frame struct to handle barriers
	return fmt.Sprintf("%v [%v|%v]", f.Name, f.Dfp, f.Evp)
}

func (f Frame) String() string {
	return fmt.Sprintf("%v [%v|%v] [%vs]", f.Name, f.Dfp, f.Evp, f.Slots)
}

func readBarrier(handle w32.HANDLE, itemAddr int, itemId string, itemGroup uint8) Frame {
	itemIndex := numbers.ReadU8(handle, uintptr(itemAddr+0xF4))
	weaponIndex := getWeaponIndex(handle, itemGroup-1, itemIndex, 0x04, 32)
	weapon := Frame{
		Id:   itemId,
		Name: readItemName(handle, int(weaponIndex)),
		Dfp:  numbers.ReadU8(handle, uintptr(itemAddr+itemBarrierDfp)),
		Evp:  numbers.ReadU8(handle, uintptr(itemAddr+itemBarrierEvp)),
	}
	return weapon
}

func readUnit(handle w32.HANDLE, itemAddr int, itemId string) Frame {
	itemIndex := numbers.ReadU8(handle, uintptr(itemAddr+0xF4))
	weaponIndex := getWeaponIndex(handle, 0, itemIndex, 0x08, 20)
	weapon := Frame{
		Id:   itemId,
		Name: readItemName(handle, int(weaponIndex)),
	}
	return weapon
}

func readMag(handle w32.HANDLE, itemAddr int, itemId string, itemGroup uint8) Mag {
	weaponIndex := getWeaponIndex(handle, 0, itemGroup, 0x10, 28)
	return Mag{
		Id:   itemId,
		Name: readItemName(handle, int(weaponIndex)),
		Def:  (int(numbers.ReadU8(handle, uintptr(itemAddr+itemMagStats+1)))<<8 + int(numbers.ReadU8(handle, uintptr(itemAddr+itemMagStats+0)))) / 100,
		Pow:  (int(numbers.ReadU8(handle, uintptr(itemAddr+itemMagStats+3)))<<8 + int(numbers.ReadU8(handle, uintptr(itemAddr+itemMagStats+2)))) / 100,
		Dex:  (int(numbers.ReadU8(handle, uintptr(itemAddr+itemMagStats+5)))<<8 + int(numbers.ReadU8(handle, uintptr(itemAddr+itemMagStats+4)))) / 100,
		Mind: (int(numbers.ReadU8(handle, uintptr(itemAddr+itemMagStats+7)))<<8 + int(numbers.ReadU8(handle, uintptr(itemAddr+itemMagStats+6)))) / 100,
	}
}

type Mag struct {
	Id   string
	Name string
	Def  int
	Pow  int
	Dex  int
	Mind int
}

func (m Mag) String() string {
	return fmt.Sprintf("%v [%v/%v/%v/%v]", m.Name, m.Def, m.Pow, m.Dex, m.Mind)
}

func getWeaponSpecial(specialId uint8) string {
	specialName := "?"
	switch specialId {
	case 0:
		specialName = ""
	case 1:
		specialName = "Draw"
	case 2:
		specialName = "Drain"
	case 3:
		specialName = "Fill"
	case 4:
		specialName = "Gush"
	case 5:
		specialName = "Heart"
	case 6:
		specialName = "Mind"
	case 7:
		specialName = "Soul"
	case 8:
		specialName = "Geist"
	case 9:
		specialName = "Master's"
	case 10:
		specialName = "Lord's"
	case 11:
		specialName = "King's"
	case 12:
		specialName = "Charge"
	case 13:
		specialName = "Spirit"
	case 14:
		specialName = "Berserk"
	case 15:
		specialName = "Ice"
	case 16:
		specialName = "Frost"
	case 17:
		specialName = "Freeze"
	case 18:
		specialName = "Blizzard"
	case 19:
		specialName = "Bind"
	case 20:
		specialName = "Hold"
	case 21:
		specialName = "Seize"
	case 22:
		specialName = "Arrest"
	case 23:
		specialName = "Heat"
	case 24:
		specialName = "Fire"
	case 25:
		specialName = "Flame"
	case 26:
		specialName = "Burning"
	case 27:
		specialName = "Shock"
	case 28:
		specialName = "Thunder"
	case 29:
		specialName = "Storm"
	case 30:
		specialName = "Tempest"
	case 31:
		specialName = "Dim"
	case 32:
		specialName = "Shadow"
	case 33:
		specialName = "Dark"
	case 34:
		specialName = "Hell"
	case 35:
		specialName = "Panic"
	case 36:
		specialName = "Riot"
	case 37:
		specialName = "Havoc"
	case 38:
		specialName = "Chaos"
	case 39:
		specialName = "Devil's"
	case 40:
		specialName = "Demon's"
	}
	return specialName
}
