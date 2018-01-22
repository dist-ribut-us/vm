package vm

import (
	"fmt"
	"strings"
)

// Op code
type Op uint16

// Bytes returns a byte slice where the first 2 bytes hold the op and sized to
// hold the specified number of args
func (op Op) Bytes(args int) []byte {
	b := make([]byte, args*8+2)
	op.Put(&b[0])
	return b
}

// OpFunc is a function that the VM can invoke for an op
type OpFunc func(*VM) error

// OpDef is a tool for defining ops for the VM. Func must be either an OpFunc,
// ArgFunc or ArgFuncErr. Args takes a slice of bools where true indicates a
// register arg and false indicates a value. The length is used for ArgFunc and
// ArgFuncErr and the boolean values are only used for the description.
type OpDef struct {
	Name string
	Desc string
	Func interface{}
	Args []bool
	Idx  Op
}

// OpFunc produces an OpFunc from an OpDef
func (od OpDef) OpFunc() OpFunc {
	if of, ok := od.Func.(func(*VM) error); ok {
		return of
	}
	if of, ok := od.Func.(OpFunc); ok {
		return of
	}

	if af, ok := od.Func.(func([]Qword, *VM)); ok {
		return argFunc(af, od.Args)
	}
	if af, ok := od.Func.(ArgFunc); ok {
		return argFunc(af, od.Args)
	}

	if afe, ok := od.Func.(func([]Qword, *VM) error); ok {
		return argFuncErr(afe, od.Args)
	}
	if afe, ok := od.Func.(ArgFuncErr); ok {
		return argFuncErr(afe, od.Args)
	}

	panic(od.Name + ": Func must be of type OpFunc, ArgFunc or ArgFuncErr")
}

// ArgFunc is helpful in defining OpFuncs from OpDefs, the args will be passed
// in as a slice of QWords
type ArgFunc func([]Qword, *VM)

func argFunc(fn ArgFunc, boolArgs []bool) OpFunc {
	return func(vm *VM) error {
		args := make([]Qword, len(boolArgs))
		for i := range args {
			args[i] = Get(&vm.Pages[vm.Page][vm.Pos+2+uint64(i)*8])
		}
		fn(args, vm)
		vm.Pos += 2 + 8*uint64(len(boolArgs))
		return nil
	}
}

// ArgFuncErr is helpful in defining OpFuncs from OpDefs, the args will be
// passed in as a slice of QWords. If it returns an error, that will be returned
// from the OpFunc.
type ArgFuncErr func([]Qword, *VM) error

func argFuncErr(fn ArgFuncErr, boolArgs []bool) OpFunc {
	return func(vm *VM) error {
		args := make([]Qword, len(boolArgs))
		for i := range args {
			args[i] = Get(&vm.Pages[vm.Page][vm.Pos+2+uint64(i)*8])
		}
		err := fn(args, vm)
		vm.Pos += 2 + 8*uint64(len(boolArgs))
		return err
	}
}

// Describe uses the OpDef to produce a description of the op.
func (od OpDef) Describe() string {
	if len(od.Args) == 0 {
		if od.Desc == "" {
			return od.Name
		}
		return fmt.Sprintf("%s : %s", od.Name, od.Desc)
	}
	args := make([]string, len(od.Args))
	var rIdx, vIdx int
	for i, isReg := range od.Args {
		if isReg {
			args[i] = fmt.Sprintf("R%d", rIdx)
			rIdx++
		} else {
			args[i] = fmt.Sprintf("V%d", vIdx)
			vIdx++
		}
	}
	argsString := strings.Join(args, " ")
	if od.Desc == "" {
		return fmt.Sprintf("%s %s", od.Name, argsString)
	}
	return fmt.Sprintf("%s %s : %s", od.Name, argsString, od.Desc)
}

// OpList allows bulk operations to be performed on a slice of OpDefs
type OpList []OpDef

// Ops returns a slice of OpFuncs. The slice will always be 65,536 long.
func (os OpList) Ops() []OpFunc {
	ops := make([]OpFunc, 65536)
	var idx Op
	for _, op := range os {
		if op.Idx != 0 {
			idx = op.Idx
		} else {
			idx++
		}
		ops[idx] = op.OpFunc()
	}
	return ops
}

// Describe returns a description of all the ops
func (os OpList) Describe() string {
	ds := make([]string, len(os))
	for i, o := range os {
		ds[i] = o.Describe()
	}
	return strings.Join(ds, "\n")
}
