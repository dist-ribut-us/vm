package vm

type VM struct {
	Registers []Qword
	Pages     [][]byte
	Pos, Page uint64
	Ops       []OpFunc
	Panic     bool
	Stop      bool
}

func New(registers int, prog []byte, ops []OpFunc) *VM {
	return &VM{
		Registers: make([]Qword, registers),
		Ops:       ops,
		Pages:     [][]byte{prog},
	}
}

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
		if err != nil {
			return
		}
		if vm.Stop {
			return
		}
	}
}
