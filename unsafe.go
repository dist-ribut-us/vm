package vm

import (
	"unsafe"
)

// Qword is used to represent values within the VM.
type Qword uint64

// QwordF converts a float64 to a Qword
func QwordF(f float64) Qword {
	return *(*Qword)(unsafe.Pointer(&f))
}

// GetF returns the Qword as a float64
func (r Qword) GetF() float64 {
	return *(*float64)(unsafe.Pointer(&r))
}

// GetU returns the Qword as a uint64
func (r Qword) GetU() uint64 {
	return uint64(r)
}

// Put takes an address as a byte pointer and sets the QWord to the value stored
// there.
func (r Qword) Put(addr *byte) {
	*(*Qword)(unsafe.Pointer(addr)) = r
}

// Get takes an address as a byte pointer and returns the value stored there as
// a Qword
func Get(addr *byte) Qword {
	return *(*Qword)(unsafe.Pointer(addr))
}

// Put takes and address a byte pointer and sets the op to the value stored
// there
func (o Op) Put(addr *byte) {
	*(*Op)(unsafe.Pointer(addr)) = o
}

// GetOp takes a byte pointer and returns the value stored there as an Op
func GetOp(addr *byte) Op {
	return *(*Op)(unsafe.Pointer(addr))
}
