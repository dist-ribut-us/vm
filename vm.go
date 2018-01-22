package vm

// VM is a register machine with pages of memory and a slice of operations that
// can execute a program.
type VM struct {
	Registers []Qword
	Pages     [][]byte
	Pos, Page uint64
	Ops       []OpFunc
	Panic     bool
	Stop      bool
	Extend    interface{}
}

// New creates a VM with the specified register values, program and ops
func New(registers []Qword, prog []byte, ops []OpFunc) *VM {
	return &VM{
		Registers: registers,
		Ops:       ops,
		Pages:     [][]byte{prog},
	}
}

// Run the VM
func (vm *VM) Run() (err error) {
	if !vm.Panic {
		defer func() {
			if r := recover(); r != nil {
				if rerr, ok := r.(error); ok {
					err = rerr
				} else {
					panic(r)
				}
			}
		}()
	}
	for {
		op := GetOp(&vm.Pages[vm.Page][vm.Pos])
		err = vm.Ops[op](vm)
		if err != nil || vm.Stop {
			return
		}
	}
}
