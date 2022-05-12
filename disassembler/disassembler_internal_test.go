package disassembler

import (
	"testing"

	"github.com/kn100/cybemu/opcode"
	"github.com/stretchr/testify/assert"
)

// mostly here to verify that the table correctly returns a word on an invalid instruction
func TestTable234(t *testing.T) {
	testCases := []struct {
		Input          []byte
		ExpectedOpcode opcode.Opcode
	}{
		{
			// 6A 10 00 00 63 00 00 00
			Input:          []byte{0x6A, 0x10, 0x00, 0x00, 0x63, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.Btst,
		},
		{
			// FF FF FF FF FF FF FF FF
			Input:          []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			ExpectedOpcode: opcode.Invalid,
		},
	}
	for _, tc := range testCases {
		inst := table234(tc.Input)
		assert.Equal(t, tc.ExpectedOpcode, inst.Opcode)
	}
}
