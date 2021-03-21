package numbers

import "math"

func Uint32FromU16(lsb uint16, msb uint16) uint32 {
	return uint32(msb)<<16 + uint32(lsb)
}

func Float32FromU16(lsb uint16, msb uint16) float32 {
	combinedValue := Uint32FromU16(lsb, msb)
	return math.Float32frombits(combinedValue)
}
