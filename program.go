package vm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type opIdx struct {
	Op
	OpDef
}

func (os OpList) Parser() func(string) ([]byte, error) {
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
		p := programmer{
			byName: byName,
			code:   program,
			vars:   make(map[string]variable),
		}
		if err := p.parse(); err != nil {
			return nil, err
		}
		return p.program, nil
	}
}

type programmer struct {
	byName  map[string]opIdx
	code    string
	program []byte
	vars    map[string]variable
	lexed   []lexedLine
}

type variable struct {
	instance []int
	value    Qword
}

var labelRe = regexp.MustCompile(`\w+:`)

func (p *programmer) parse() error {
	p.lex()
	for _, line := range p.lexed {
		opName := line.word[0]
		if opName == "#def" {
			if err := p.def(line); err != nil {
				return err
			}
			continue
		}
		if labelRe.MatchString(opName) {
			opName = string(opName[:len(opName)-1])
			v := p.vars[opName]
			v.value = Qword(len(p.program))
			p.vars[opName] = v
			continue
		}
		op, ok := p.byName[opName]
		if !ok {
			return line.Error("Op not found")
		}
		if err := p.appendOp(op, line); err != nil {
			return err
		}
	}
	for _, v := range p.vars {
		for _, i := range v.instance {
			v.value.Put(&(p.program[i]))
		}
	}
	return nil
}

func (p *programmer) appendOp(op opIdx, line lexedLine) error {
	if len(line.word)-1 != len(op.Args) {
		return line.Error("Wrong number of arguments")
	}
	pos := len(p.program) + 2
	p.program = append(p.program, op.Bytes(len(op.Args))...)
	for i, arg := range line.word[1:] {
		p.setArg(arg, pos+i*8)
	}
	return nil
}

func (p *programmer) def(line lexedLine) error {
	if len(line.word) != 3 {
		return line.Error("Wrong number of arguments")
	}
	name := line.word[1]
	val, isWord, err := convertArg(line.word[2])
	if err != nil {
		return err
	}
	if isWord {
		return line.Error("definition must be a number")
	}
	v := p.vars[name]
	v.value = val
	p.vars[name] = v
	return nil
}

type lexedLine struct {
	number int
	word   []string
	raw    string
}

func (l lexedLine) Error(ErrorType string) LineError {
	return LineError{
		LineNumber: l.number,
		LineString: l.raw,
		ErrorType:  ErrorType,
	}
}

type LineError struct {
	LineNumber int
	LineString string
	ErrorType  string
}

func (le LineError) Error() string {
	return fmt.Sprintf("%s) %d: %s", le.ErrorType, le.LineNumber, le.LineString)
}

var lineRe = regexp.MustCompile(`^[ \t]*#?(?:([\w\.]+:?)[ \t]*)+`)
var partsRe = regexp.MustCompile(`#?[\w\.]+:?`)

func (p *programmer) lex() {
	for li, lineStr := range strings.Split(p.code, "\n") {
		raw := strings.TrimSpace(lineStr)
		lineStr = lineRe.FindString(raw)
		words := partsRe.FindAllString(lineStr, -1)
		if len(words) == 0 {
			continue
		}
		p.lexed = append(p.lexed, lexedLine{
			number: li,
			word:   words,
			raw:    raw,
		})
	}
}

func (p *programmer) setArg(arg string, pos int) error {
	r, isArg, err := convertArg(arg)
	if err != nil {
		return err
	}
	if !isArg {
		r.Put(&(p.program[pos]))
		return nil
	}

	v := p.vars[arg]
	v.instance = append(v.instance, pos)
	p.vars[arg] = v
	return nil
}

func convertArg(arg string) (Qword, bool, error) {
	if strings.Contains(arg, ".") {
		f, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return 0, false, err
		}
		r := QwordF(f)
		return r, false, nil
	}

	u, err := strconv.ParseUint(arg, 10, 64)
	if err == nil {
		return Qword(u), false, nil
	}

	return 0, true, nil
}
