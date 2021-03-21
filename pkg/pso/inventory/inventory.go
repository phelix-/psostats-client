package inventory

import (
	"errors"
	"fmt"
	"log"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/pkg/numbers"
)

const (
	itemArray        = 0x00A8D81C
	itemArrayCount   = 0x00A8D820
	itemId           = 0xD8
	itemOwner        = 0xE4
	itemCode         = 0xF2
	itemEquipped     = 0x190
	itemKills        = 0xE8
	itemWepGrind     = 0x1F5
	itemWepSpecial   = 0x1F6
	itemWepStats     = 0x1C8
	itemArmSlots     = 0x1B8
	itemFrameDfp     = 0x1B9
	itemFrameEvp     = 0x1BA
	itemBarrierDfp   = 0x1E4
	itemBarrierEvp   = 0x1E5
	itemUnitMod      = 0x1DC
	itemMagStats     = 0x1C0
	itemMagPBHas     = 0x1C8
	itemMagPB        = 0x1C9
	itemMagColor     = 0x1CA
	itemMagSync      = 0x1BE
	itemMagIQ        = 0x1BC
	itemMagTimer     = 0x1B4
	itemToolCount    = 0x104
	itemTechType     = 0x108
	itemMesetaAmount = 0x100
)

func ReadInventory(handle w32.HANDLE, playerAddress int) (Weapon, error) {
	weapon := Weapon{}
	buf, _, ok := w32.ReadProcessMemory(handle, uintptr(itemArrayCount), 2)
	if !ok {
		return weapon, errors.New("Could not read item count")
	}
	count := numbers.Uint32From16(buf[0:2])
	buf, _, ok = w32.ReadProcessMemory(handle, uintptr(itemArray), 4)
	address := numbers.Uint32From16(buf[0:2])
	buf, _, ok = w32.ReadProcessMemory(handle, uintptr(address), uintptr(4*count))

	for i := 0; i < int(count); i++ {
		itemAddr := numbers.Uint32From16(buf[i*2 : (i*2)+2])
		if itemAddr != 0 {
			itemBuffer, _, ok := w32.ReadProcessMemory(handle, uintptr(itemAddr+0xD8), 4)
			if !ok {
				return weapon, errors.New("Could not read item")
			}
			itemId := fmt.Sprintf("%04x%04x", itemBuffer[1], itemBuffer[0])
			// itemDataBuffer, _, ok := w32.ReadProcessMemory(handle, uintptr(itemAddr+0xF2), 8)
			// itemType := itemDataBuffer[0] & 0xFF00
			itemType := readU8(handle, uintptr(itemAddr+0xF2))
			equipped := readU8(handle, uintptr(itemAddr+itemEquipped))&0x01 == 1
			// log.Printf("%v equipped=%v", itemId, equipped)
			if itemType == 0 && equipped {
				weapon = readWeapon(handle, int(itemAddr), itemId)
			}
			// log.Printf("%04x%04x %04x%04x", itemDataBuffer[1], itemDataBuffer[0], itemDataBuffer[3], itemDataBuffer[2])

		}
	}

	// buf, _, ok := w32.ReadProcessMemory(handle, uintptr(playerAddress+base), uintptr((max-base)+4))
	return weapon, nil
}

func readU8(handle w32.HANDLE, address uintptr) uint8 {
	buf, _, ok := w32.ReadProcessMemory(handle, address, 1)
	if !ok {
		log.Fatalf("Error reading 0x%08x", address)
	}
	return uint8(buf[0])
}

func readU32(handle w32.HANDLE, address uintptr) uint32 {
	buf, _, ok := w32.ReadProcessMemory(handle, address, 4)
	if !ok {
		log.Fatalf("Error reading 0x%08x", address)
	}
	return numbers.Uint32From16(buf[0:2])
}

func getWeaponIndex(handle w32.HANDLE, group uint8, index uint8) uint32 {
	weaponIndex := uint32(0)
	pmtAddress := readU32(handle, 0x00a8dc94)
	weaponAddress := readU32(handle, uintptr(pmtAddress))
	if weaponAddress != 0 {
		groupAddress := weaponAddress + (uint32(group) * 8)
		itemAddress := readU32(handle, uintptr(groupAddress+4)) + (44 * uint32(index))
		weaponIndex = readU32(handle, uintptr(itemAddress))
	}
	return weaponIndex
}

func readItemName(handle w32.HANDLE, index int) string {
	unitxtPointer := readU32(handle, 0x00a9cd50)
	if unitxtPointer == 0 {
		return "?"
	}
	weaponIndex := readU32(handle, uintptr(unitxtPointer+4))
	weaponNameAddress := readU32(handle, uintptr(weaponIndex+uint32(4*index)))

	// log.Printf("0x%08x", weaponNameAddress)
	weaponName, err := numbers.ReadNullTerminatedString(handle, uintptr(weaponNameAddress))
	if err != nil {
		log.Fatalf("Error getting weapon name %v", err)
	}
	return fmt.Sprintf("%v", weaponName)
}

func readWeapon(handle w32.HANDLE, itemAddr int, itemId string) Weapon {
	itemGroup := readU8(handle, uintptr(itemAddr+0xF3))
	itemIndex := readU8(handle, uintptr(itemAddr+0xF4))
	weaponIndex := getWeaponIndex(handle, itemGroup, itemIndex)
	weapon := Weapon{
		Id: itemId,
	}
	weapon.Name = readItemName(handle, int(weaponIndex))
	weapon.Grind = readU8(handle, uintptr(itemAddr+itemWepGrind))
	weapon.Special = readU8(handle, uintptr(itemAddr+itemWepSpecial))
	// log.Printf("i=%2v itemAddr=%08x - itemId= type=%x grind=%v special=%x",
	// 	i, itemAddr, , itemType, grind, special)
	for j := 0; j < 6; j += 2 {
		area := readU8(handle, uintptr(itemAddr+itemWepStats+j))
		percent := readU8(handle, uintptr(itemAddr+itemWepStats+j+1))
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
	// log.Printf("%v", weapon.String())
	return weapon
}

type Weapon struct {
	Id      string
	Name    string
	Grind   uint8
	Special uint8
	Native  uint8
	ABeast  uint8
	Machine uint8
	Dark    uint8
	Hit     uint8
}

func (w *Weapon) String() string {
	return fmt.Sprintf("%v +%v [%v] [%v/%v/%v/%v|%v]", w.Name, w.Grind, w.Special, w.Native, w.ABeast, w.Machine, w.Dark, w.Hit)
}
