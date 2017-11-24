package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSetU(t *testing.T) {
	var x uint64 = 12345
	b := make([]byte, 10)
	SetU(x, &b[0])
	assert.Equal(t, x, GetU(&b[0]))
}

func TestBasic(t *testing.T) {
	p := Programmer{}.
		SetU(1, 123).
		SetF(0, 55.55).
		Copy(2, 1).
		SetF(3, 12.34).
		SetF(4, 11.11).
		FAdd(3, 4).
		Stop()

	v := New(6, p)

	err := v.Run()
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), v.Registers[1].GetU())
	assert.Equal(t, uint64(123), v.Registers[2].GetU())
	assert.Equal(t, float64(55.55), v.Registers[0].GetF())
	assert.Equal(t, float64(23.45), v.Registers[3].GetF())
}

func TestRecover(t *testing.T) {
	p := Programmer{}.
		SetU(100, 123). // Register out of range
		Stop()

	v := New(6, p)

	err := v.Run()
	assert.Error(t, err)
}

func TestPages(t *testing.T) {
	p := Programmer{}.
		SetU(0, 1024).
		Alloc(0). // request 1024 byte page
		SetU(1, 555).
		Write(1, 0, 2). // write 555 to pos 0 of that page
		Read(2, 0, 2).  // read the 555 back out to r2
		Stop()

	v := New(6, p)

	err := v.Run()
	assert.NoError(t, err)
	if assert.Len(t, v.Pages, 2) {
		assert.Len(t, v.Pages[1], 1024)
	}
	assert.Equal(t, Register(555), v.Registers[2])
}

func TestManualMult(t *testing.T) {
	// compute 5x3
	p := Programmer{}.
		SetU(0, 5).
		SetU(1, 3).
		SetU(2, 1).
		Position(3, 4).
		IAdd(5, 0).
		ISub(1, 2).
		Jump(1, 3, 4).
		Stop()

	v := New(6, p)
	v.Panic = true

	err := v.Run()
	assert.NoError(t, err)
	assert.Equal(t, Register(15), v.Registers[5])
}

func TestProgrammer(t *testing.T) {
	ops := OpSet{
		{
			Name: "set",
			ArgFunc: func(args []uint64, vm *VM) {
				vm.Registers[args[0]] = Register(args[1])
			},
			Args: []bool{true, true},
		},
		{
			Name: "stop",
			OpFunc: func(vm *VM) error {
				vm.Stop = true
				return nil
			},
			Idx: Stop,
		},
	}
	assert.Equal(t, "set", ops[0].Name)

	parser := ops.Parser()
	program, err := parser(`
		#def foo 1
		set 1 2
		set 3 4
		stop
	`)
	assert.NoError(t, err)
	t.Error(program)
}
