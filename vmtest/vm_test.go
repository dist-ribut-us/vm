package vmtest

import (
	"github.com/dist-ribut-us/vm"
	"github.com/dist-ribut-us/vm/ops"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasic(t *testing.T) {
	parser := ops.List.Parser()
	p, err := parser(`
		set 1 123
		set 0 55.55
		copy 2 1
		set 3 12.34
		set 4 11.11
		fadd 3 4
		stop
	`)
	assert.NoError(t, err)

	v := vm.New([]vm.Qword{0, 0, 0, 0, 0}, p, ops.List.Ops())
	v.Panic = true

	err = v.Run()
	assert.NoError(t, err)
	assert.Equal(t, uint64(123), v.Registers[1].GetU())
	assert.Equal(t, uint64(123), v.Registers[2].GetU())
	assert.Equal(t, float64(55.55), v.Registers[0].GetF())
	assert.Equal(t, float64(23.45), v.Registers[3].GetF())
}

func TestRecover(t *testing.T) {
	parser := ops.List.Parser()
	p, err := parser(`
		set 100 123 // Register out of range
		stop
	`)
	assert.NoError(t, err)

	v := vm.New([]vm.Qword{0, 0, 0, 0, 0}, p, ops.List.Ops())

	err = v.Run()
	assert.Error(t, err)
}

func TestPages(t *testing.T) {
	parser := ops.List.Parser()
	p, err := parser(`
		set 0 1024
		alloc 0
		set 1 555
		write 1 0 2
		read 2 0 2
		stop
	`)
	assert.NoError(t, err)
	v := vm.New([]vm.Qword{0, 0, 0, 0, 0}, p, ops.List.Ops())

	err = v.Run()
	assert.NoError(t, err)
	if assert.Len(t, v.Pages, 2) {
		assert.Len(t, v.Pages[1], 1024)
	}
	assert.Equal(t, vm.Qword(555), v.Registers[2])
}

func TestManualMult(t *testing.T) {
	// compute 5x3
	parser := ops.List.Parser()
	p, err := parser(`
		set 0 5
		set 1 3
		set 2 1
		position 3 4
		iadd 5 0
		isub 1 2
		jump 1 3 4
		stop
	`)
	assert.NoError(t, err)
	v := vm.New([]vm.Qword{0, 0, 0, 0, 0, 0}, p, ops.List.Ops())
	v.Panic = true

	err = v.Run()
	assert.NoError(t, err)
	assert.Equal(t, vm.Qword(15), v.Registers[5])
}

func TestVar(t *testing.T) {
	// compute AxB
	parser := ops.List.Parser()
	p, err := parser(`
		#def  A 7
		#def  B 4
		set   0 A
		set   1 B
		loop:
		iadd  2 0
		isubv 1 1
		jumpv 1 0 loop
		stop
	`)
	assert.NoError(t, err)
	v := vm.New([]vm.Qword{0, 0, 0}, p, ops.List.Ops())
	v.Panic = true

	err = v.Run()
	assert.NoError(t, err)
	assert.Equal(t, vm.Qword(28), v.Registers[2])
}
