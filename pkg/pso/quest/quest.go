package quest

import (
	"errors"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/pkg/numbers"
)

func GetQuestPointer(handle w32.HANDLE) (uint32, error) {
	buf, _, ok := w32.ReadProcessMemory(handle, uintptr(0x00A95AA8), 4)
	if !ok {
		return 0, errors.New("Unable to getQuestPointer")
	}
	questPointer := numbers.Uint32From16(buf[0:2])
	return questPointer, nil
}

func IsRegisterSet(handle w32.HANDLE, registerId uint16) (bool, error) {
	questRegisterAddress, err := numbers.ReadU32(handle, uintptr(0x00A954B0))
	if err != nil {
		return false, err
	}
	registerSet := false
	if questRegisterAddress != 0 {
		buf, _, ok := w32.ReadProcessMemory(handle, uintptr(int(questRegisterAddress)+(4*int(registerId))), 2)
		if !ok {
			return false, errors.New("Unable to isQuestRegisterSet")
		}
		registerSet = buf[0] > 0
	}
	// log.Printf("R[%v]@%x = %v", registerId, uintptr(0x00A954B0+(4*int(registerId))), registerSet)
	// log.Print("\n")
	return registerSet, nil
}
