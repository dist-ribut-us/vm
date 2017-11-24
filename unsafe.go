package vm

import (
	"unsafe"
)

type Qword uint64

func QwordF(f float64) Qword {
	return *(*Qword)(unsafe.Pointer(&f))
}

func (r Qword) GetF() float64 {
	return *(*float64)(unsafe.Pointer(&r))
}

func (r Qword) GetU() uint64 {
	return uint64(r)
}

func (r Qword) Put(addr *byte) {
	*(*Qword)(unsafe.Pointer(addr)) = r
}

func Get(addr *byte) Qword {
	return *(*Qword)(unsafe.Pointer(addr))
}

func (o Op) Put(addr *byte) {
	*(*Op)(unsafe.Pointer(addr)) = o
}

func GetOp(addr *byte) Op {
	return *(*Op)(unsafe.Pointer(addr))
}
