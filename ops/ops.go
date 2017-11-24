package ops

import (
	"github.com/dist-ribut-us/vm"
)

var List = vm.OpList{
	{
		Name: "set",
		Desc: "set the value of a register",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] = args[1]
		},
		Args: []bool{true, false},
	},
	{
		Name: "copy",
		Desc: "copy the value from the second register into the first",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] = v.Registers[args[1]]
		},
		Args: []bool{true, true},
	},
	{
		Name: "iadd",
		Desc: "add the value from the second register into the first, interpreting the values as integers",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] += v.Registers[args[1]]
		},
		Args: []bool{true, true},
	},
	{
		Name: "isub",
		Desc: "subtract the value of the second register from the first, interpreting the values as integers",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] -= v.Registers[args[1]]
		},
		Args: []bool{true, true},
	},
	{
		Name: "isubv",
		Desc: "subtract the value of the second argument from the first register, interpreting the values as integers",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] -= args[1]
		},
		Args: []bool{true, true},
	},
	{
		Name: "imul",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] *= v.Registers[args[1]]
		},
		Args: []bool{true, true},
	},
	{
		Name: "fadd",
		Func: func(args []vm.Qword, v *vm.VM) {
			f := v.Registers[args[0]].GetF() + v.Registers[args[1]].GetF()
			v.Registers[args[0]] = vm.QwordF(f)
		},
		Args: []bool{true, true},
	},
	{
		Name: "fsub",
		Func: func(args []vm.Qword, v *vm.VM) {
			f := v.Registers[args[0]].GetF() - v.Registers[args[1]].GetF()
			v.Registers[args[0]] = vm.Qword(f)
		},
		Args: []bool{true, true},
	},
	{
		Name: "fmul",
		Func: func(args []vm.Qword, v *vm.VM) {
			f := v.Registers[args[0]].GetF() * v.Registers[args[1]].GetF()
			v.Registers[args[0]] = vm.Qword(f)
		},
		Args: []bool{true, true},
	},
	{
		Name: "alloc",
		Func: func(args []vm.Qword, v *vm.VM) {
			size := v.Registers[args[0]]
			v.Registers[args[0]] = vm.Qword(len(v.Pages))
			v.Pages = append(v.Pages, make([]byte, size))
		},
		Args: []bool{true},
	},
	{
		Name: "read",
		Func: func(args []vm.Qword, v *vm.VM) {
			page := v.Registers[args[1]]
			pos := v.Registers[args[2]]
			v.Registers[args[0]] = vm.Get(&v.Pages[page][pos])
		},
		Args: []bool{true, true, true},
	},
	{
		Name: "write",
		Func: func(args []vm.Qword, v *vm.VM) {
			page := v.Registers[args[1]]
			pos := v.Registers[args[2]]
			v.Registers[args[0]].Put(&v.Pages[page][pos])
		},
		Args: []bool{true, true, true},
	},
	{
		Name: "jump",
		Func: func(v *vm.VM) error {
			r1 := vm.Get(&v.Pages[v.Page][v.Pos+2])
			condition := v.Registers[r1]
			if condition == 0 {
				v.Pos += 2 + 8*3
				return nil
			}
			r2 := vm.Get(&v.Pages[v.Page][v.Pos+10])
			r3 := vm.Get(&v.Pages[v.Page][v.Pos+18])
			page := v.Registers[r2].GetU()
			pos := v.Registers[r3].GetU()
			v.Page = page
			v.Pos = pos
			return nil
		},
		Args: []bool{true, true, true},
	},
	{
		Name: "jumpv",
		Func: func(v *vm.VM) error {
			r1 := vm.Get(&v.Pages[v.Page][v.Pos+2])
			condition := v.Registers[r1]
			if condition == 0 {
				v.Pos += 2 + 8*3
				return nil
			}
			page := vm.Get(&v.Pages[v.Page][v.Pos+10]).GetU()
			pos := vm.Get(&v.Pages[v.Page][v.Pos+18]).GetU()
			v.Page = page
			v.Pos = pos
			return nil
		},
		Args: []bool{true, false, false},
	},
	{
		Name: "position",
		Func: func(args []vm.Qword, v *vm.VM) {
			v.Registers[args[0]] = vm.Qword(v.Page)
			v.Registers[args[1]] = vm.Qword(v.Pos)
		},
		Args: []bool{true, true},
	},
	// Keep this at the end
	{
		Name: "stop",
		Func: func(v *vm.VM) error {
			v.Stop = true
			return nil
		},
		Idx: 65535,
	},
}

var Ops = List.Ops()
var Parser = List.Parser()
