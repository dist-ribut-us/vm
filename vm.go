package vm

type VM struct {
	Registers []Qword
	Pages     [][]byte
	Pos, Page uint64
	Ops       []OpFunc
	Cost      []int64
	Supply    int64
	Panic     bool
	Stop      bool
}

func New(registers []Qword, prog []byte, ops []OpFunc) *VM {
	return &VM{
		Registers: registers,
		Ops:       ops,
		Pages:     [][]byte{prog},
	}
}

type SupplyDepleted struct{}

func (SupplyDepleted) Error() string {
	return "Supply Depleted"
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
		if int(op) < len(vm.Cost) {
			vm.Supply -= vm.Cost[op]
			if vm.Supply < 0 {
				err = SupplyDepleted{}
				return
			}
		}
		if vm.Stop {
			return
		}
	}
}
