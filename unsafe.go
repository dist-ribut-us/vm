package vm

import (
	"unsafe"
)

type Register uint64

func RegF(f float64) Register {
	return *(*Register)(unsafe.Pointer(&f))
}

func (r *Register) GetF() float64 {
	return *(*float64)(unsafe.Pointer(r))
}

func (r *Register) GetU() uint64 {
	return uint64(*r)
}

func SetU(i uint64, addr *byte) {
	*(*uint64)(unsafe.Pointer(addr)) = i
}

func SetF(f float64, addr *byte) {
	*(*float64)(unsafe.Pointer(addr)) = f
}

func SetOp(o Op, addr *byte) {
	*(*Op)(unsafe.Pointer(addr)) = o
}

func GetU(addr *byte) uint64 {
	return *(*uint64)(unsafe.Pointer(addr))
}

func GetR(addr *byte) Register {
	return *(*Register)(unsafe.Pointer(addr))
}

func GetF(addr *byte) float64 {
	return *(*float64)(unsafe.Pointer(addr))
}

func GetOp(addr *byte) Op {
	return *(*Op)(unsafe.Pointer(addr))
}
