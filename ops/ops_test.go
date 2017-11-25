package ops

import (
	"github.com/dist-ribut-us/vm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOps(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		reg      []vm.Qword
		expected []vm.Qword
		pages    int
	}{
		{
			name: "stop",
			code: `
        stop
      `,
		},
		{
			name: "set",
			code: `
        set 0 111
        set 1 12.34
        stop
      `,
			reg:      []vm.Qword{0, 0},
			expected: []vm.Qword{111, vm.QwordF(12.34)},
		},
		{
			name: "copy",
			code: `
        copy 1 0
        stop
      `,
			reg:      []vm.Qword{10, 0},
			expected: []vm.Qword{10, 10},
		},
		{
			name: "iadd",
			code: `
        iadd 0 0
        stop
      `,
			reg:      []vm.Qword{5},
			expected: []vm.Qword{10},
		},
		{
			name: "iaddv",
			code: `
        iaddv 0 7
        stop
      `,
			reg:      []vm.Qword{12},
			expected: []vm.Qword{19},
		},
		{
			name: "isub",
			code: `
        isub 0 1
        stop
      `,
			reg:      []vm.Qword{7, 5},
			expected: []vm.Qword{2, 5},
		},
		{
			name: "isubv",
			code: `
        isubv 0 5
        stop
      `,
			reg:      []vm.Qword{7},
			expected: []vm.Qword{2},
		},
		{
			name: "imul",
			code: `
        imul 0 1
        stop
      `,
			reg:      []vm.Qword{2, 3},
			expected: []vm.Qword{6, 3},
		},
		{
			name: "imulv",
			code: `
        imulv 0 4
        stop
      `,
			reg:      []vm.Qword{6},
			expected: []vm.Qword{24},
		},
		{
			name: "fadd",
			code: `
        fadd 0 1
        stop
      `,
			reg:      []vm.Qword{vm.QwordF(12.34), vm.QwordF(56.78)},
			expected: []vm.Qword{vm.QwordF(69.12), vm.QwordF(56.78)},
		},
		{
			name: "faddv",
			code: `
        faddv 0 11.11
        stop
      `,
			reg:      []vm.Qword{vm.QwordF(12.34)},
			expected: []vm.Qword{vm.QwordF(23.45)},
		},
		{
			name: "fsub",
			code: `
        fsub 0 1
        stop
      `,
			reg:      []vm.Qword{vm.QwordF(56.78), vm.QwordF(12.34)},
			expected: []vm.Qword{vm.QwordF(44.44), vm.QwordF(12.34)},
		},
		{
			name: "fsubv",
			code: `
        fsubv 0 23.45
        stop
      `,
			reg:      []vm.Qword{vm.QwordF(56.78)},
			expected: []vm.Qword{vm.QwordF(56.78 - 23.45)},
		},
		{
			name: "fmul",
			code: `
        fmul 0 1
        stop
      `,
			reg:      []vm.Qword{vm.QwordF(56.78), vm.QwordF(12.34)},
			expected: []vm.Qword{vm.QwordF(56.78 * 12.34), vm.QwordF(12.34)},
		},
		{
			name: "fmulv",
			code: `
        fmulv 0 11.11
        stop
      `,
			reg:      []vm.Qword{vm.QwordF(56.78)},
			expected: []vm.Qword{vm.QwordF(56.78 * 11.11)},
		},
		{
			name: "alloc",
			code: `
        alloc 0
        stop
      `,
			reg:      []vm.Qword{100},
			expected: []vm.Qword{1},
			pages:    2,
		},
		{
			name: "read",
			code: `
        read 0     0 1
        set  1 12345
        stop
      `,
			reg:      []vm.Qword{0, 36},
			expected: []vm.Qword{12345, 12345},
		},
		{
			name: "write",
			code: `
        alloc 0
        write 1 0 2
        read  0 0 2
        stop
      `,
			reg:      []vm.Qword{100, 123, 0},
			expected: []vm.Qword{123, 123, 0},
		},
		{
			name: "jump",
			code: `
        set 0 end
        jump 0 1 0
        set 100 0 // bad register, will error if we run this line
        end:
        stop
      `,
			reg:      []vm.Qword{0, 0},
			expected: []vm.Qword{62, 0},
		},
		{
			name: "jumpv",
			code: `
        jumpv 0 0 end
        set 100 0 // bad register, will error if we run this line
        end:
        stop
      `,
			reg:      []vm.Qword{1},
			expected: []vm.Qword{1},
		},
		{
			name: "position",
			code: `
        set 0 0 
        position 0 1
        stop
      `,
			reg:      []vm.Qword{0, 0},
			expected: []vm.Qword{0, 36},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := Parser(tc.code)
			assert.NoError(t, err)
			v := vm.New(tc.reg, p, Ops)
			err = v.Run()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, v.Registers)
			if tc.pages > 0 {
				assert.Len(t, v.Pages, tc.pages)
			}
		})
	}
}
