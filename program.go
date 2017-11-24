package vm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

type Programmer []byte

func (p Programmer) Stop() Programmer {
	return p.Append(Stop)
}

func (p Programmer) Append(op Op, args ...uint64) Programmer {
	b := make([]byte, 2+8*len(args))
	*(*Op)(unsafe.Pointer(&b[0])) = op
	for i, a := range args {
		SetU(a, &b[2+i*8])
	}
	return append(p, b...)
}

func (p Programmer) SetU(r, v uint64) Programmer {
	return p.Append(Set, r, v)
}

func (p Programmer) SetF(r uint64, v float64) Programmer {
	uv := *(*uint64)(unsafe.Pointer(&v))
	return p.Append(Set, r, uv)
}

func (p Programmer) Copy(r1, r2 uint64) Programmer {
	return p.Append(Copy, r1, r2)
}

func (p Programmer) IAdd(r1, r2 uint64) Programmer {
	return p.Append(IAdd, r1, r2)
}

func (p Programmer) FAdd(r1, r2 uint64) Programmer {
	return p.Append(FAdd, r1, r2)
}

func (p Programmer) ISub(r1, r2 uint64) Programmer {
	return p.Append(ISub, r1, r2)
}

func (p Programmer) FSub(r1, r2 uint64) Programmer {
	return p.Append(FSub, r1, r2)
}

func (p Programmer) Alloc(r uint64) Programmer {
	return p.Append(Alloc, r)
}

func (p Programmer) Read(r1, r2, r3 uint64) Programmer {
	return p.Append(Read, r1, r2, r3)
}

func (p Programmer) Write(r1, r2, r3 uint64) Programmer {
	return p.Append(Write, r1, r2, r3)
}

func (p Programmer) Jump(r1, r2, r3 uint64) Programmer {
	return p.Append(Jump, r1, r2, r3)
}

func (p Programmer) Position(r1, r2 uint64) Programmer {
	return p.Append(Position, r1, r2)
}

type OpDef struct {
	OpFunc
	ArgFunc func([]uint64, *VM)
	Name    string
	Args    []bool // T = Register
	Idx     Op
}

func (od *OpDef) Func() OpFunc {
	if od.OpFunc != nil {
		return od.OpFunc
	}
	od.OpFunc = func(vm *VM) error {
		args := make([]uint64, len(od.Args))
		for i := range args {
			args[i] = GetU(&vm.Pages[vm.Page][vm.Pos+2+uint64(i)*8])
		}
		od.ArgFunc(args, vm)
		vm.Pos += 2 + 8*uint64(len(od.Args))
		return nil
	}
	return od.OpFunc
}

type OpSet []OpDef

func (os OpSet) Ops() []OpFunc {
	ops := make([]OpFunc, 65535)
	var idx Op
	for _, op := range os {
		if op.Idx != 0 {
			idx = op.Idx
		} else {
			idx++
		}
		ops[idx] = op.Func()
	}
	return ops
}

type opIdx struct {
	Op
	OpDef
}

func (os OpSet) Parser() func(string) ([]byte, error) {
	byName := make(map[string]opIdx, len(os))
	var idx Op
	for _, op := range os {
		if op.Idx != 0 {
			idx = op.Idx
		} else {
			idx++
		}
		byName[op.Name] = opIdx{
			Op:    idx,
			OpDef: op,
		}
	}

	return func(program string) ([]byte, error) {
		return parse(byName, program)
	}
}

type ParseError struct {
	LineNumber int
	LineString string
	ErrorType  string
}

func (pe ParseError) Error() string {
	return fmt.Sprintf("%s) %d: %s", pe.ErrorType, pe.LineNumber, pe.LineString)
}

func parse(byName map[string]opIdx, code string) ([]byte, error) {
	lexed := lex(code)
	var program []byte
	for _, line := range lexed {
		opName := line.word[0]
		op, ok := byName[opName]
		if !ok {
			return nil, ParseError{
				LineNumber: line.number,
				LineString: line.raw,
				ErrorType:  "Op not found",
			}
		}
		if len(line.word)-1 != len(op.Args) {
			return nil, ParseError{
				LineNumber: line.number,
				LineString: line.raw,
				ErrorType:  "Wrong number of arguments",
			}
		}
		opBytes := op.Bytes(len(op.Args))
		for i, arg := range line.word[1:] {
			setArg(arg, &opBytes[2+8*i])
		}
		program = append(program, opBytes...)
	}
	return program, nil
}

type lexedLine struct {
	number int
	word   []string
	raw    string
}

var lineRe = regexp.MustCompile(`^[ \t]*(?:(\w+)[ \t]*)+`)
var partsRe = regexp.MustCompile(`\w+`)

func lex(program string) []lexedLine {
	var lexed []lexedLine
	for li, lineStr := range strings.Split(program, "\n") {
		raw := strings.TrimSpace(lineStr)
		lineStr = lineRe.FindString(raw)
		words := partsRe.FindAllString(lineStr, -1)
		if len(words) == 0 {
			continue
		}
		lexed = append(lexed, lexedLine{
			number: li,
			word:   words,
			raw:    raw,
		})
	}
	return lexed
}

func setArg(arg string, addr *byte) error {
	if strings.Contains(arg, ".") {
		f, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return err
		}
		SetF(f, addr)
	} else {
		u, err := strconv.ParseUint(arg, 10, 64)
		if err != nil {
			return err
		}
		SetU(u, addr)
	}
	return nil
}
