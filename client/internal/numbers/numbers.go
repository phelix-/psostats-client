package numbers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/TheTitanrain/w32"
)

func Uint32From16(slice []uint16) uint32 {
	return Uint32FromU16(slice[0], slice[1])
}

func Uint32FromU16(lsb uint16, msb uint16) uint32 {
	return uint32(msb)<<16 + uint32(lsb)
}

func Float32FromU16(lsb uint16, msb uint16) float32 {
	combinedValue := Uint32FromU16(lsb, msb)
	return math.Float32frombits(combinedValue)
}

func ReadString(handle w32.HANDLE, address uintptr, length int) (string, error) {
	buf, _, ok := w32.ReadProcessMemory(handle, uintptr(address), uintptr(length*2))
	if !ok {
		return "", errors.New(fmt.Sprintf("Unable to read string at 0x%08x", address))
	}

	byteBuf := bytes.NewBuffer(make([]byte, 0, len(buf)))
	for i := 0; i < length; i++ {
		b := buf[i]
		split := make([]byte, 2)
		binary.LittleEndian.PutUint16(split, b)
		byteBuf.Write(split)
	}
	return byteBuf.String(), nil
}

func ReadNullTerminatedString(handle w32.HANDLE, address uintptr) (string, error) {
	buf, _, ok := w32.ReadProcessMemory(handle, uintptr(address), 48)
	if !ok {
		return "", errors.New(fmt.Sprintf("Unable to read string at 0x%08x", address))
	}

	dataFound := false
	byteBuf := bytes.NewBuffer(make([]byte, 0, len(buf)))
	for i := 0; i < 24; i++ {
		b := buf[i]
		if b > 0 {
			dataFound = true
		} else if dataFound {
			break
		}
		split := make([]byte, 1)
		split[0] = byte(b)
		byteBuf.Write(split)
	}
	return byteBuf.String(), nil
}

func ReadU16(handle w32.HANDLE, address uintptr) uint16 {
	buf, _, ok := w32.ReadProcessMemory(handle, address, 2)
	if !ok {
		log.Fatalf("Unable to readU16 @0x%08x", address)
	}
	return buf[0]
}

func ReadU32Unchecked(handle w32.HANDLE, address uintptr) uint32 {
	ret, err := ReadU32(handle, address)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return ret
}

func ReadU32(handle w32.HANDLE, address uintptr) (uint32, error) {
	buf, _, ok := w32.ReadProcessMemory(handle, address, 4)
	if !ok {
		return 0, errors.New(fmt.Sprintf("Unable to readU32 0x%08x", address))
	}
	return Uint32From16(buf[0:2]), nil
}
