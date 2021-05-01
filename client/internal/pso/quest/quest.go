package quest

import (
	"errors"

	"github.com/TheTitanrain/w32"
	"github.com/phelix-/psostats/v2/client/internal/numbers"
)

func GetQuestPointer(handle w32.HANDLE) uintptr {
	return uintptr(numbers.ReadU32(handle, uintptr(0x00A95AA8)))
}

func GetQuestDataPointer(handle w32.HANDLE, questPtr uintptr) uintptr {
	return uintptr(numbers.ReadU32(handle, questPtr+0x19C))
}

func IsRegisterSet(handle w32.HANDLE, registerId uint16) (bool, error) {
	questRegisterAddress := numbers.ReadU32(handle, uintptr(0x00A954B0))
	registerSet := false
	if questRegisterAddress != 0 {
		buf, _, ok := w32.ReadProcessMemory(handle, uintptr(questRegisterAddress)+(4*uintptr(registerId)), 2)
		if !ok {
			return false, errors.New("unable to isQuestRegisterSet")
		}
		registerSet = buf[0] > 0
	}
	return registerSet, nil
}
