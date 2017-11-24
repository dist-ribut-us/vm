package vm

type Op uint16

func (op Op) Bytes(args int) []byte {
	b := make([]byte, args*8+2)
	SetOp(op, &b[0])
	return b
}

type OpFunc func(*VM) error

// control codes
const (
	Stop Op = 0xffff - iota
)

const (
	// Set <R1> <V>: R1 = V
	Set Op = iota
	// Copy <R1> <R2>: R1 = R2
	Copy
	// IAdd <R1> <R2>: R1 += R2 (as integers)
	IAdd
	// FAdd <R1> <R2>: R1 += R2 (as floats)
	FAdd
	// IAdd <R1> <R2>: R1 -= R2 (as integers)
	ISub
	// FAdd <R1> <R2>: R1 -= R2 (as floats)
	FSub
	// IMul <R1> <R2>: R1 *= R2 (as integers)
	IMul
	// FMul <R1> <R2>: R1 *= R2 (as floats)
	FMul
	// Alloc <R1>: R1 holds the number of bytes requests and is set to the Page
	// index of the allocated bytes.
	Alloc
	// Read <R1> <R2> <R3>: Reads 4 bytes into R1 from page R2, position R3
	Read
	// Write <R1> <R2> <R3>: Writes 4 bytes from R1 into page R2 at position R3
	Write
	// Jump <R1> <R2> <R3>: Jumps to page R2 and position to R3 if R1 != 0
	Jump
	// Position <R1> <R2>: Records the current page to R1 and position to R2
	Position
)

var Ops []OpFunc

func init() {
	Ops = make([]OpFunc, 65536)
	base := []OpFunc{
		// Ops must be kept in the same order as the Op codes
		set,
		cp,
		iadd,
		fadd,
		isub,
		fsub,
		imul,
		fmul,
		alloc,
		read,
		write,
		jump,
		position,
	}
	copy(Ops, base)
	Ops[Stop] = func(vm *VM) error {
		vm.Stop = true
		return nil
	}
}

var set = args(2, func(args []uint64, vm *VM) {
	vm.Registers[args[0]] = Register(args[1])
})

var cp = args(2, func(args []uint64, vm *VM) {
	vm.Registers[args[0]] = vm.Registers[args[1]]
})

var iadd = args(2, func(args []uint64, vm *VM) {
	vm.Registers[args[0]] += vm.Registers[args[1]]
})

var isub = args(2, func(args []uint64, vm *VM) {
	vm.Registers[args[0]] -= vm.Registers[args[1]]
})
var imul = args(2, func(args []uint64, vm *VM) {
	vm.Registers[args[0]] *= vm.Registers[args[1]]
})

var fadd = args(2, func(args []uint64, vm *VM) {
	f := vm.Registers[args[0]].GetF() + vm.Registers[args[1]].GetF()
	vm.Registers[args[0]] = RegF(f)
})

var fsub = args(2, func(args []uint64, vm *VM) {
	f := vm.Registers[args[0]].GetF() - vm.Registers[args[1]].GetF()
	vm.Registers[args[0]] = RegF(f)
})
var fmul = args(2, func(args []uint64, vm *VM) {
	f := vm.Registers[args[0]].GetF() * vm.Registers[args[1]].GetF()
	vm.Registers[args[0]] = RegF(f)
})

var alloc = args(1, func(args []uint64, vm *VM) {
	size := vm.Registers[args[0]]
	vm.Registers[args[0]] = Register(len(vm.Pages))
	vm.Pages = append(vm.Pages, make([]byte, size))
})

var read = args(3, func(args []uint64, vm *VM) {
	page := vm.Registers[args[1]]
	pos := vm.Registers[args[2]]
	vm.Registers[args[0]] = GetR(&vm.Pages[page][pos])
})

var write = args(3, func(args []uint64, vm *VM) {
	page := vm.Registers[args[1]]
	pos := vm.Registers[args[2]]
	SetU(vm.Registers[args[0]].GetU(), &vm.Pages[page][pos])
})

func jump(vm *VM) error {
	r1 := GetU(&vm.Pages[vm.Page][vm.Pos+2])
	condition := vm.Registers[r1]
	if condition == 0 {
		vm.Pos += 2 + 8*3
		return nil
	}
	r2 := GetU(&vm.Pages[vm.Page][vm.Pos+10])
	r3 := GetU(&vm.Pages[vm.Page][vm.Pos+18])
	page := vm.Registers[r2].GetU()
	pos := vm.Registers[r3].GetU()
	vm.Page = page
	vm.Pos = pos
	return nil
}

var position = args(2, func(args []uint64, vm *VM) {
	vm.Registers[args[0]] = Register(vm.Page)
	vm.Registers[args[1]] = Register(vm.Pos)
})

func args(ln uint64, fn func([]uint64, *VM)) OpFunc {
	return func(vm *VM) error {
		args := make([]uint64, ln)
		for i := range args {
			args[i] = GetU(&vm.Pages[vm.Page][vm.Pos+2+uint64(i)*8])
		}
		fn(args, vm)
		vm.Pos += 2 + 8*ln
		return nil
	}
}
