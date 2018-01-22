package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnsafe(t *testing.T) {
	b := make([]byte, 10)
	var r Qword = 12345
	r.Put(&b[3])
	assert.Equal(t, r, Get(&b[3]))
}
