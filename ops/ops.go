package ops

import (
	"github.com/dist-ribut-us/vm"
)

var List = vm.OpList{
	{
		Name: "set",
		Desc: "set R0 to V0",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] = args[1]
		},
		Args: []bool{true, false},
	},
	{
		Name: "copy",
		Desc: "copy the value of R1 into R0",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] = v.Registers[args[1]]
		},
		Args: []bool{true, true},
	},
	{
		Name: "iadd",
		Desc: "set R0 to R0+R1 treating both as integers",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] += v.Registers[args[1]]
		},
		Args: []bool{true, true},
	},
	{
		Name: "iaddv",
		Desc: "set R0 to R0+V0 treating both as integers",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] += args[1]
		},
		Args: []bool{true, true},
	},
	{
		Name: "isub",
		Desc: "set R0 to R0-R1 treating both as integers",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] -= v.Registers[args[1]]
		},
		Args: []bool{true, true},
	},
	{
		Name: "isubv",
		Desc: "set R0 to R0-V0 treating both as integers",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] -= args[1]
		},
		Args: []bool{true, true},
	},
	{
		Name: "imul",
		Desc: "set R0 to R0*R1 treating both as integers",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] *= v.Registers[args[1]]
		},
		Args: []bool{true, true},
	},
	{
		Name: "imulv",
		Desc: "set R0 to R0*V0 treating both as integers",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] *= args[1]
		},
		Args: []bool{true, true},
	},
	{
		Name: "fadd",
		Desc: "set R0 to R0+R1 treating both as floating point numbers",
		Func: func(args []vm.Qword, v *vm.VM) {
			f := v.Registers[args[0]].GetF() + v.Registers[args[1]].GetF()
			v.Registers[args[0]] = vm.QwordF(f)
		},
		Args: []bool{true, true},
	},
	{
		Name: "faddv",
		Desc: "set R0 to R0+V0 treating both as floating point numbers",
		Func: func(args []vm.Qword, v *vm.VM) {
			f := v.Registers[args[0]].GetF() + args[1].GetF()
			v.Registers[args[0]] = vm.QwordF(f)
		},
		Args: []bool{true, true},
	},
	{
		Name: "fsub",
		Desc: "set R0 to R0-R1 treating both as floating point numbers",
		Func: func(args []vm.Qword, v *vm.VM) {
			f := v.Registers[args[0]].GetF() - v.Registers[args[1]].GetF()
			v.Registers[args[0]] = vm.QwordF(f)
		},
		Args: []bool{true, true},
	},
	{
		Name: "fsubv",
		Desc: "set R0 to R0-V0 treating both as floating point numbers",
		Func: func(args []vm.Qword, v *vm.VM) {
			f := v.Registers[args[0]].GetF() - args[1].GetF()
			v.Registers[args[0]] = vm.QwordF(f)
		},
		Args: []bool{true, true},
	},
	{
		Name: "fmul",
		Desc: "set R0 to R0*R1 treating both as floating point numbers",
		Func: func(args []vm.Qword, v *vm.VM) {
			f := v.Registers[args[0]].GetF() * v.Registers[args[1]].GetF()
			v.Registers[args[0]] = vm.QwordF(f)
		},
		Args: []bool{true, true},
	},
	{
		Name: "fmulv",
		Desc: "set R0 to R0*V0 treating both as floating point numbers",
		Func: func(args []vm.Qword, v *vm.VM) {
			f := v.Registers[args[0]].GetF() * args[1].GetF()
			v.Registers[args[0]] = vm.QwordF(f)
		},
		Args: []bool{true, true},
	},
	{
		Name: "alloc",
		Desc: "allocates a new page with a size of R0 then sets R0 to the page number",
		Func: func(args []vm.Qword, v *vm.VM) {
			size := v.Registers[args[0]]
			v.Registers[args[0]] = vm.Qword(len(v.Pages))
			v.Pages = append(v.Pages, make([]byte, size))
		},
		Args: []bool{true},
	},
	{
		Name: "read",
		Desc: "sets R0 to the value at page R1, position R2",
		Func: func(args []vm.Qword, v *vm.VM) {
			page := v.Registers[args[1]]
			pos := v.Registers[args[2]]
			v.Registers[args[0]] = vm.Get(&v.Pages[page][pos])
		},
		Args: []bool{true, true, true},
	},
	{
		Name: "write",
		Desc: "writes the value in R0 to page R1, position R2",
		Func: func(args []vm.Qword, v *vm.VM) {
			page := v.Registers[args[1]]
			pos := v.Registers[args[2]]
			v.Registers[args[0]].Put(&v.Pages[page][pos])
		},
		Args: []bool{true, true, true},
	},
	{
		Name: "jump",
		Desc: "if R0 is not 0, it will jump to page R1, position R2",
		Func: func(v *vm.VM) error {
			r1 := vm.Get(&v.Pages[v.Page][v.Pos+2])
			condition := v.Registers[r1]
			if condition == 0 {
				v.Pos += 2 + 8*3
				return nil
			}
			r2 := vm.Get(&v.Pages[v.Page][v.Pos+10])
			r3 := vm.Get(&v.Pages[v.Page][v.Pos+18])
			v.Page = v.Registers[r2].GetU()
			v.Pos = v.Registers[r3].GetU()
			return nil
		},
		Args: []bool{true, true, true},
	},
	{
		Name: "jumpv",
		Desc: "if R0 is not 0, it will jump to page V0, position V1",
		Func: func(v *vm.VM) error {
			r1 := vm.Get(&v.Pages[v.Page][v.Pos+2])
			condition := v.Registers[r1]
			if condition == 0 {
				v.Pos += 2 + 8*3
				return nil
			}
			v.Page = vm.Get(&v.Pages[v.Page][v.Pos+10]).GetU()
			v.Pos = vm.Get(&v.Pages[v.Page][v.Pos+18]).GetU()
			return nil
		},
		Args: []bool{true, false, false},
	},
	{
		Name: "position",
		Desc: "sets the first register to the current page and the second register to the position of the next instruction",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] = vm.Qword(v.Page)
			v.Registers[args[1]] = vm.Qword(v.Pos + 2 + 8*2)
		},
		Args: []bool{true, true},
	},
	// Keep this at the end
	{
		Name: "stop",
		Desc: "stops the VM",
		Func: func(v *vm.VM) error {
			v.Stop = true
			return nil
		},
		Idx: 65535,
	},
}

var Ops = List.Ops()
var Parser = List.Parser()
