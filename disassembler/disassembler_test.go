package disassembler_test

import (
	"fmt"
	"testing"

	"github.com/kn100/cybemu/addressingmode"
	"github.com/kn100/cybemu/disassembler"
	"github.com/kn100/cybemu/instruction"
	"github.com/kn100/cybemu/opcode"
	"github.com/kn100/cybemu/operand"
	"github.com/kn100/cybemu/size"
	"github.com/stretchr/testify/assert"
)

// TestDisassembleTimsTestCases Tests the disassembler against Tim's test cases.
// It's a useful test to determine every possible instruction decodes correctly.
func TestDisassemble(t *testing.T) {
	testCases := []struct {
		Bytes         []byte
		ExpectedInsts []instruction.Inst
	}{
		{
			Bytes: []byte{0x00, 0x00, 0x00, 0x00},
			ExpectedInsts: []instruction.Inst{
				{Opcode: opcode.Nop, Bytes: []byte{0x00, 0x00}, TotalBytes: 2, OperandType: operand.None},
				{Opcode: opcode.Nop, Bytes: []byte{0x00, 0x00}, TotalBytes: 2, Pos: 2, OperandType: operand.None},
			},
		},
		{
			Bytes: []byte{
				0x8D, 0x81,
				0x79, 0x6C, 0x26, 0x94,
				0x01, 0x00, 0x78, 0x70, 0x6B, 0x23, 0x00, 0x00, 0x27, 0x0E,
			},
			ExpectedInsts: []instruction.Inst{
				{
					Opcode:         opcode.Add,
					Bytes:          []byte{0x8D, 0x81},
					TotalBytes:     2,
					BWL:            size.Byte,
					AddressingMode: addressingmode.Immediate,
					Pos:            0,
				},
				{
					Opcode:         opcode.And,
					Bytes:          []byte{0x79, 0x6C, 0x26, 0x94},
					TotalBytes:     4,
					BWL:            size.Word,
					AddressingMode: addressingmode.Immediate,
					Pos:            2,
				},
				{
					Opcode:         opcode.Mov,
					Bytes:          []byte{0x01, 0x00, 0x78, 0x70, 0x6B, 0x23, 0x00, 0x00, 0x27, 0x0E},
					TotalBytes:     10,
					BWL:            size.Longword,
					AddressingMode: addressingmode.RegisterIndirectWithDisplacement,
					Pos:            6,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.Bytes), func(t *testing.T) {
			insts := disassembler.Disassemble(tc.Bytes)
			assert.Equal(t, tc.ExpectedInsts, insts)
		})
	}
}

func TestDecodeTimsTestCases(t *testing.T) {
	testCases := []struct {
		Input                  []byte
		ExpectedOpcode         opcode.Opcode
		expectedBWL            size.Size
		expectedAddressingMode addressingmode.AddressingMode
	}{
		{
			Input:          []byte{0x8D, 0x81},
			ExpectedOpcode: opcode.Add,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x08, 0x3E},
			ExpectedOpcode: opcode.Add,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x79, 0x11, 0x30, 0x39},
			ExpectedOpcode: opcode.Add,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x09, 0x0B},
			ExpectedOpcode: opcode.Add,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x7A, 0x15, 0x12, 0x34, 0x56, 0x78},
			ExpectedOpcode: opcode.Add,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x0A, 0x87},
			ExpectedOpcode: opcode.Add,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x0B, 0x00},
			ExpectedOpcode: opcode.Adds,
		},
		{
			Input:          []byte{0x0B, 0x81},
			ExpectedOpcode: opcode.Adds,
		},
		{
			Input:          []byte{0x0B, 0x92},
			ExpectedOpcode: opcode.Adds,
		},
		{
			Input:          []byte{0x95, 0x0A},
			ExpectedOpcode: opcode.Addx,
		},
		{
			Input:          []byte{0x0E, 0xC0},
			ExpectedOpcode: opcode.Addx,
		},
		{
			Input:          []byte{0xE2, 0x2D},
			ExpectedOpcode: opcode.And,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x16, 0x32},
			ExpectedOpcode: opcode.And,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x79, 0x6C, 0x26, 0x94},
			ExpectedOpcode: opcode.And,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x66, 0x2E},
			ExpectedOpcode: opcode.And,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x7A, 0x63, 0x00, 0x0A, 0xBC, 0xDE},
			ExpectedOpcode: opcode.And,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0xF0, 0x66, 0x54},
			ExpectedOpcode: opcode.And,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x06, 0xC0},
			ExpectedOpcode: opcode.Andc,
		},
		{
			Input:          []byte{0x76, 0x25},
			ExpectedOpcode: opcode.Band,
		},
		{
			Input:          []byte{0x7C, 0x20, 0x76, 0x40},
			ExpectedOpcode: opcode.Band,
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x76, 0x50},
			ExpectedOpcode: opcode.Band,
		},
		{
			Input:          []byte{0x40, 0xFE},
			ExpectedOpcode: opcode.Bra,
		},
		{
			Input:          []byte{0x41, 0xFC},
			ExpectedOpcode: opcode.Brn,
		},
		{
			Input:          []byte{0x42, 0xFA},
			ExpectedOpcode: opcode.Bhi,
		},
		{
			Input:          []byte{0x43, 0xF8},
			ExpectedOpcode: opcode.Bls,
		},
		{
			Input:          []byte{0x44, 0xF6},
			ExpectedOpcode: opcode.Bcc,
		},
		{
			Input:          []byte{0x45, 0xF4},
			ExpectedOpcode: opcode.Bcs,
		},
		{
			Input:          []byte{0x46, 0xF2},
			ExpectedOpcode: opcode.Bne,
		},
		{
			Input:          []byte{0x47, 0xF0},
			ExpectedOpcode: opcode.Beq,
		},
		{
			Input:          []byte{0x48, 0xEE},
			ExpectedOpcode: opcode.Bvc,
		},
		{
			Input:          []byte{0x49, 0xEC},
			ExpectedOpcode: opcode.Bvs,
		},
		{
			Input:          []byte{0x4A, 0xEA},
			ExpectedOpcode: opcode.Bpl,
		},
		{
			Input:          []byte{0x4B, 0xE8},
			ExpectedOpcode: opcode.Bmi,
		},
		{
			Input:          []byte{0x4C, 0xE6},
			ExpectedOpcode: opcode.Bge,
		},
		{
			Input:          []byte{0x4D, 0xE4},
			ExpectedOpcode: opcode.Blt,
		},
		{
			Input:          []byte{0x4E, 0xE2},
			ExpectedOpcode: opcode.Bgt,
		},
		{
			Input:          []byte{0x4F, 0xE0},
			ExpectedOpcode: opcode.Ble,
		},
		{
			Input:          []byte{0x58, 0x00, 0x00, 0x3C},
			ExpectedOpcode: opcode.Bra,
		},
		{
			Input:          []byte{0x58, 0x10, 0x00, 0x38},
			ExpectedOpcode: opcode.Brn,
		},
		{
			Input:          []byte{0x58, 0x20, 0x00, 0x34},
			ExpectedOpcode: opcode.Bhi,
		},
		{
			Input:          []byte{0x58, 0x30, 0x00, 0x30},
			ExpectedOpcode: opcode.Bls,
		},
		{
			Input:          []byte{0x58, 0x40, 0x00, 0x2C},
			ExpectedOpcode: opcode.Bcc,
		},
		{
			Input:          []byte{0x58, 0x50, 0x00, 0x28},
			ExpectedOpcode: opcode.Bcs,
		},
		{
			Input:          []byte{0x58, 0x60, 0x00, 0x24},
			ExpectedOpcode: opcode.Bne,
		},
		{
			Input:          []byte{0x58, 0x70, 0x00, 0x20},
			ExpectedOpcode: opcode.Beq,
		},
		{
			Input:          []byte{0x58, 0x80, 0x00, 0x1C},
			ExpectedOpcode: opcode.Bvc,
		},
		{
			Input:          []byte{0x58, 0x90, 0x00, 0x18},
			ExpectedOpcode: opcode.Bvs,
		},
		{
			Input:          []byte{0x58, 0xA0, 0x00, 0x14},
			ExpectedOpcode: opcode.Bpl,
		},
		{
			Input:          []byte{0x58, 0xB0, 0x00, 0x10},
			ExpectedOpcode: opcode.Bmi,
		},
		{
			Input:          []byte{0x58, 0xC0, 0x00, 0x0C},
			ExpectedOpcode: opcode.Bge,
		},
		{
			Input:          []byte{0x58, 0xD0, 0x00, 0x08},
			ExpectedOpcode: opcode.Blt,
		},
		{
			Input:          []byte{0x58, 0xE0, 0x00, 0x04},
			ExpectedOpcode: opcode.Bgt,
		},
		{
			Input:          []byte{0x58, 0xF0, 0x00, 0x00},
			ExpectedOpcode: opcode.Ble,
		},
		{
			Input:          []byte{0x72, 0x4A},
			ExpectedOpcode: opcode.Bclr,
		},
		{
			Input:          []byte{0x7D, 0x40, 0x72, 0x60},
			ExpectedOpcode: opcode.Bclr,
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x72, 0x10},
			ExpectedOpcode: opcode.Bclr,
		},
		{
			Input:          []byte{0x62, 0x93},
			ExpectedOpcode: opcode.Bclr,
		},
		{
			Input:          []byte{0x7D, 0x30, 0x62, 0x40},
			ExpectedOpcode: opcode.Bclr,
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x62, 0x50},
			ExpectedOpcode: opcode.Bclr,
		},
		{
			Input:          []byte{0x76, 0xC2},
			ExpectedOpcode: opcode.Biand,
		},
		{
			Input:          []byte{0x7C, 0x40, 0x76, 0xF0},
			ExpectedOpcode: opcode.Biand,
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x76, 0x80},
			ExpectedOpcode: opcode.Biand,
		},
		{
			Input:          []byte{0x77, 0xC2},
			ExpectedOpcode: opcode.Bild,
		},
		{
			Input:          []byte{0x7C, 0x40, 0x77, 0xF0},
			ExpectedOpcode: opcode.Bild,
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x77, 0x80},
			ExpectedOpcode: opcode.Bild,
		},
		{
			Input:          []byte{0x74, 0xC2},
			ExpectedOpcode: opcode.Bior,
		},
		{
			Input:          []byte{0x7C, 0x40, 0x74, 0xF0},
			ExpectedOpcode: opcode.Bior,
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x74, 0x80},
			ExpectedOpcode: opcode.Bior,
		},
		{
			Input:          []byte{0x67, 0xC2},
			ExpectedOpcode: opcode.Bist,
		},
		{
			Input:          []byte{0x7D, 0x40, 0x67, 0xF0},
			ExpectedOpcode: opcode.Bist,
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x67, 0x80},
			ExpectedOpcode: opcode.Bist,
		},
		{
			Input:          []byte{0x75, 0xC2},
			ExpectedOpcode: opcode.Bixor,
		},
		{
			Input:          []byte{0x7C, 0x40, 0x75, 0xF0},
			ExpectedOpcode: opcode.Bixor,
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x75, 0x80},
			ExpectedOpcode: opcode.Bixor,
		},
		{
			Input:          []byte{0x77, 0x42},
			ExpectedOpcode: opcode.Bld,
		},
		{
			Input:          []byte{0x7C, 0x40, 0x77, 0x70},
			ExpectedOpcode: opcode.Bld,
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x77, 0x00},
			ExpectedOpcode: opcode.Bld,
		},
		{
			Input:          []byte{0x71, 0x42},
			ExpectedOpcode: opcode.Bnot,
		},
		{
			Input:          []byte{0x7D, 0x40, 0x71, 0x70},
			ExpectedOpcode: opcode.Bnot,
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x71, 0x00},
			ExpectedOpcode: opcode.Bnot,
		},
		{
			Input:          []byte{0x61, 0x81},
			ExpectedOpcode: opcode.Bnot,
		},
		{
			Input:          []byte{0x7D, 0x30, 0x61, 0x50},
			ExpectedOpcode: opcode.Bnot,
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x61, 0xE0},
			ExpectedOpcode: opcode.Bnot,
		},
		{
			Input:          []byte{0x74, 0x42},
			ExpectedOpcode: opcode.Bor,
		},
		{
			Input:          []byte{0x7C, 0x40, 0x74, 0x70},
			ExpectedOpcode: opcode.Bor,
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x74, 0x00},
			ExpectedOpcode: opcode.Bor,
		},
		{
			Input:          []byte{0x70, 0x42},
			ExpectedOpcode: opcode.Bset,
		},
		{
			Input:          []byte{0x7D, 0x40, 0x70, 0x70},
			ExpectedOpcode: opcode.Bset,
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x70, 0x00},
			ExpectedOpcode: opcode.Bset,
		},
		{
			Input:          []byte{0x60, 0x81},
			ExpectedOpcode: opcode.Bset,
		},
		{
			Input:          []byte{0x7D, 0x30, 0x60, 0x50},
			ExpectedOpcode: opcode.Bset,
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x60, 0xE0},
			ExpectedOpcode: opcode.Bset,
		},
		{
			Input:          []byte{0x55, 0x00},
			ExpectedOpcode: opcode.Bsr,
		},
		{
			Input:          []byte{0x5C, 0x00, 0xFF, 0x7A},
			ExpectedOpcode: opcode.Bsr,
		},
		{
			Input:          []byte{0x5C, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.Bsr,
		},
		{
			Input:          []byte{0x67, 0x42},
			ExpectedOpcode: opcode.Bst,
		},
		{
			Input:          []byte{0x7D, 0x40, 0x67, 0x70},
			ExpectedOpcode: opcode.Bst,
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x67, 0x00},
			ExpectedOpcode: opcode.Bst,
		},
		{
			Input:          []byte{0x73, 0x42},
			ExpectedOpcode: opcode.Btst,
		},
		{
			Input:          []byte{0x7C, 0x40, 0x73, 0x70},
			ExpectedOpcode: opcode.Btst,
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x73, 0x00},
			ExpectedOpcode: opcode.Btst,
		},
		{
			Input:          []byte{0x63, 0x81},
			ExpectedOpcode: opcode.Btst,
		},
		{
			Input:          []byte{0x7C, 0x30, 0x63, 0x50},
			ExpectedOpcode: opcode.Btst,
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x63, 0xE0},
			ExpectedOpcode: opcode.Btst,
		},
		{
			Input:          []byte{0x75, 0x42},
			ExpectedOpcode: opcode.Bxor,
		},
		{
			Input:          []byte{0x7C, 0x40, 0x75, 0x70},
			ExpectedOpcode: opcode.Bxor,
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x75, 0x00},
			ExpectedOpcode: opcode.Bxor,
		},
		{
			Input:          []byte{0xA5, 0x8F},
			ExpectedOpcode: opcode.Cmp,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x1C, 0x4A},
			ExpectedOpcode: opcode.Cmp,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x79, 0x2D, 0x1F, 0xFF},
			ExpectedOpcode: opcode.Cmp,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x1D, 0xD2},
			ExpectedOpcode: opcode.Cmp,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x7A, 0x24, 0x00, 0x00, 0xFF, 0xFF},
			ExpectedOpcode: opcode.Cmp,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x1F, 0xB5},
			ExpectedOpcode: opcode.Cmp,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x0F, 0x0C},
			ExpectedOpcode: opcode.Daa,
		},
		{
			Input:          []byte{0x1F, 0x05},
			ExpectedOpcode: opcode.Das,
		},
		{
			Input:          []byte{0x1A, 0x05},
			ExpectedOpcode: opcode.Dec,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x1B, 0x54},
			ExpectedOpcode: opcode.Dec,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x1B, 0xDB},
			ExpectedOpcode: opcode.Dec,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x1B, 0x72},
			ExpectedOpcode: opcode.Dec,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x1B, 0xF3},
			ExpectedOpcode: opcode.Dec,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0xD0, 0x51, 0xBC},
			ExpectedOpcode: opcode.Divxs,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x01, 0xD0, 0x53, 0xB2},
			ExpectedOpcode: opcode.Divxs,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x51, 0xBC},
			ExpectedOpcode: opcode.Divxu,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x53, 0xB2},
			ExpectedOpcode: opcode.Divxu,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x7B, 0x5C, 0x59, 0x8F},
			ExpectedOpcode: opcode.Eepmov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x7B, 0xD4, 0x59, 0x8F},
			ExpectedOpcode: opcode.Eepmov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x17, 0xD2},
			ExpectedOpcode: opcode.Exts,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x17, 0xF6},
			ExpectedOpcode: opcode.Exts,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x17, 0x53},
			ExpectedOpcode: opcode.Extu,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x17, 0x75},
			ExpectedOpcode: opcode.Extu,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x0A, 0x04},
			ExpectedOpcode: opcode.Inc,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x0B, 0x5B},
			ExpectedOpcode: opcode.Inc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x0B, 0xD5},
			ExpectedOpcode: opcode.Inc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x0B, 0x72},
			ExpectedOpcode: opcode.Inc,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x0B, 0xF5},
			ExpectedOpcode: opcode.Inc,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x59, 0x60},
			ExpectedOpcode: opcode.Jmp,
		},
		{
			Input:          []byte{0x5A, 0x12, 0x89, 0xDE},
			ExpectedOpcode: opcode.Jmp,
		},
		{
			Input:          []byte{0x5B, 0x3C},
			ExpectedOpcode: opcode.Jmp,
		},
		{
			Input:          []byte{0x5D, 0x60},
			ExpectedOpcode: opcode.Jsr,
		},
		{
			Input:          []byte{0x5E, 0x12, 0x89, 0xDE},
			ExpectedOpcode: opcode.Jsr,
		},
		{
			Input:          []byte{0x5F, 0x3C},
			ExpectedOpcode: opcode.Jsr,
		},
		{
			Input:          []byte{0x07, 0xC1},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x03, 0x04},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x01, 0x40, 0x69, 0x20},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6F, 0x10, 0x1F, 0xFF},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x78, 0x20, 0x6B, 0x20, 0x00, 0x12, 0x34, 0x56},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6D, 0x30},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6B, 0x00, 0x01, 0x26},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6B, 0x20, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0xF6, 0x63},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x0C, 0xD4},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x68, 0x79},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x6E, 0x72, 0xFF, 0xFF},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x78, 0x70, 0x6A, 0x2B, 0x00, 0xFF, 0xFF, 0x9D},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x6C, 0x41},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x26, 0xC0},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x6A, 0x0C, 0x01, 0x26},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x6A, 0x2A, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x68, 0xA0},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x6E, 0xC8, 0x82, 0xFB},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x78, 0x50, 0x6A, 0xA9, 0x00, 0xFE, 0x79, 0x61},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x6C, 0xA3},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x31, 0xC0},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x6A, 0x89, 0x01, 0x26},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x6A, 0xA2, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x79, 0x0B, 0x27, 0x0F},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x0D, 0x4A},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x69, 0x24},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x6F, 0x1A, 0x00, 0xFF},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x78, 0x20, 0x6B, 0x25, 0x00, 0x01, 0x38, 0x80},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x6D, 0x40},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x6D, 0x70},
			ExpectedOpcode: opcode.Pop,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x6B, 0x0E, 0x01, 0x26},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x6B, 0x25, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x69, 0xD8},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x6F, 0xF2, 0x01, 0x01},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x78, 0x60, 0x6B, 0xAD, 0x00, 0x00, 0x27, 0x10},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x6D, 0xD0},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x6D, 0xFE},
			ExpectedOpcode: opcode.Push,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x6B, 0x88, 0x01, 0x26},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x6B, 0xAC, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x7A, 0x00, 0x00, 0x00, 0x11, 0xD7},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x0F, 0x81},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x69, 0x23},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6F, 0x74, 0x00, 0x3E},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x78, 0x70, 0x6B, 0x23, 0x00, 0x00, 0x27, 0x0E},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6D, 0x70},
			ExpectedOpcode: opcode.Pop,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6B, 0x03, 0x01, 0x26},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6B, 0x22, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x69, 0x95},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6F, 0xF4, 0x00, 0x22},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			// Modified this test case to change DH to 0 (previously F) since I
			// can't find a case where it's greater than 7 in the manual.
			Input:          []byte{0x01, 0x00, 0x78, 0x00, 0x6B, 0xA5, 0x00, 0x00, 0x30, 0x0C},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6D, 0xF6},
			ExpectedOpcode: opcode.Push,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6B, 0x81, 0x01, 0x26},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6B, 0xA2, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: opcode.Mov,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x6A, 0x4D, 0xFF, 0xC0},
			ExpectedOpcode: opcode.Movfpe,
		},
		{
			Input:          []byte{0x6A, 0xC5, 0xFF, 0xC0},
			ExpectedOpcode: opcode.Movtpe,
		},
		{
			Input:          []byte{0x01, 0xC0, 0x50, 0x42},
			ExpectedOpcode: opcode.Mulxs,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x01, 0xC0, 0x52, 0x25},
			ExpectedOpcode: opcode.Mulxs,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x50, 0x42},
			ExpectedOpcode: opcode.Mulxu,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x52, 0x26},
			ExpectedOpcode: opcode.Mulxu,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x17, 0x88},
			ExpectedOpcode: opcode.Neg,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x17, 0x9C},
			ExpectedOpcode: opcode.Neg,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x17, 0xB5},
			ExpectedOpcode: opcode.Neg,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x00, 0x00},
			ExpectedOpcode: opcode.Nop,
		},
		{
			Input:          []byte{0x17, 0x0C},
			ExpectedOpcode: opcode.Not,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x17, 0x15},
			ExpectedOpcode: opcode.Not,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x17, 0x31},
			ExpectedOpcode: opcode.Not,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0xC1, 0x04},
			ExpectedOpcode: opcode.Or,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x14, 0x80},
			ExpectedOpcode: opcode.Or,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x79, 0x40, 0x00, 0xC0},
			ExpectedOpcode: opcode.Or,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x64, 0x80},
			ExpectedOpcode: opcode.Or,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x7A, 0x40, 0x00, 0x00, 0x00, 0xFE},
			ExpectedOpcode: opcode.Or,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0xF0, 0x64, 0x05},
			ExpectedOpcode: opcode.Or,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x04, 0x01},
			ExpectedOpcode: opcode.Orc,
		},
		{
			Input:          []byte{0x6D, 0x78},
			ExpectedOpcode: opcode.Pop,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6D, 0x73},
			ExpectedOpcode: opcode.Pop,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x6D, 0xF1},
			ExpectedOpcode: opcode.Push,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6D, 0xF6},
			ExpectedOpcode: opcode.Push,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x12, 0x81},
			ExpectedOpcode: opcode.Rotl,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x12, 0x9C},
			ExpectedOpcode: opcode.Rotl,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x12, 0xB6},
			ExpectedOpcode: opcode.Rotl,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x13, 0x81},
			ExpectedOpcode: opcode.Rotr,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x13, 0x94},
			ExpectedOpcode: opcode.Rotr,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x13, 0xB3},
			ExpectedOpcode: opcode.Rotr,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x12, 0x03},
			ExpectedOpcode: opcode.Rotxl,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x12, 0x16},
			ExpectedOpcode: opcode.Rotxl,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x12, 0x35},
			ExpectedOpcode: opcode.Rotxl,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x13, 0x0C},
			ExpectedOpcode: opcode.Rotxr,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x13, 0x1E},
			ExpectedOpcode: opcode.Rotxr,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x13, 0x34},
			ExpectedOpcode: opcode.Rotxr,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x56, 0x70},
			ExpectedOpcode: opcode.Rte,
		},
		{
			Input:          []byte{0x54, 0x70},
			ExpectedOpcode: opcode.Rts,
		},
		{
			Input:          []byte{0x10, 0x85},
			ExpectedOpcode: opcode.Shal,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x10, 0x9B},
			ExpectedOpcode: opcode.Shal,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x10, 0xB3},
			ExpectedOpcode: opcode.Shal,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x11, 0x84},
			ExpectedOpcode: opcode.Shar,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x11, 0x9E},
			ExpectedOpcode: opcode.Shar,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x11, 0xB5},
			ExpectedOpcode: opcode.Shar,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x10, 0x0E},
			ExpectedOpcode: opcode.Shll,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x10, 0x14},
			ExpectedOpcode: opcode.Shll,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x10, 0x32},
			ExpectedOpcode: opcode.Shll,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x11, 0x03},
			ExpectedOpcode: opcode.Shlr,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x11, 0x1A},
			ExpectedOpcode: opcode.Shlr,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x11, 0x31},
			ExpectedOpcode: opcode.Shlr,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x80},
			ExpectedOpcode: opcode.Sleep,
		},
		{
			Input:          []byte{0x02, 0x04},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x01, 0x40, 0x69, 0xF0},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6F, 0xF0, 0x00, 0x10},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x78, 0x70, 0x6B, 0xA0, 0x00, 0x00, 0x00, 0x64},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6D, 0xE0},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6B, 0x80, 0xFF, 0xC0},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6B, 0xA0, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x18, 0x44},
			ExpectedOpcode: opcode.Sub,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x79, 0x3B, 0xFF, 0xF8},
			ExpectedOpcode: opcode.Sub,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x19, 0x0C},
			ExpectedOpcode: opcode.Sub,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x7A, 0x37, 0xFF, 0xFF, 0xFF, 0xF0},
			ExpectedOpcode: opcode.Sub,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x1A, 0x86},
			ExpectedOpcode: opcode.Sub,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x1B, 0x04},
			ExpectedOpcode: opcode.Subs,
		},
		{
			Input:          []byte{0x1B, 0x85},
			ExpectedOpcode: opcode.Subs,
		},
		{
			Input:          []byte{0x1B, 0x96},
			ExpectedOpcode: opcode.Subs,
		},
		{
			Input:          []byte{0xB5, 0x08},
			ExpectedOpcode: opcode.Subx,
		},
		{
			Input:          []byte{0x1E, 0x09},
			ExpectedOpcode: opcode.Subx,
		},
		{
			Input:          []byte{0x57, 0x20},
			ExpectedOpcode: opcode.Trapa,
		},
		{
			Input:          []byte{0xD4, 0x80},
			ExpectedOpcode: opcode.Xor,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x15, 0x4C},
			ExpectedOpcode: opcode.Xor,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x79, 0x5D, 0x20, 0x00},
			ExpectedOpcode: opcode.Xor,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x65, 0xD4},
			ExpectedOpcode: opcode.Xor,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x7A, 0x56, 0x00, 0x00, 0xFF, 0xFF},
			ExpectedOpcode: opcode.Xor,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0xF0, 0x65, 0x01},
			ExpectedOpcode: opcode.Xor,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x05, 0x40},
			ExpectedOpcode: opcode.Xorc,
		},
		{
			Input:          []byte{0x01, 0x41, 0x06, 0xA5},
			ExpectedOpcode: opcode.Andc,
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x76, 0x30},
			ExpectedOpcode: opcode.Band,
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x76, 0x50},
			ExpectedOpcode: opcode.Band,
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x72, 0x30},
			ExpectedOpcode: opcode.Bclr,
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x72, 0x50},
			ExpectedOpcode: opcode.Bclr,
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x62, 0x70},
			ExpectedOpcode: opcode.Bclr,
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x62, 0xE0},
			ExpectedOpcode: opcode.Bclr,
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x76, 0xB0},
			ExpectedOpcode: opcode.Biand,
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x67, 0x89, 0x76, 0xD0},
			ExpectedOpcode: opcode.Biand,
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x77, 0xB0},
			ExpectedOpcode: opcode.Bild,
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x77, 0xD0},
			ExpectedOpcode: opcode.Bild,
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x74, 0xB0},
			ExpectedOpcode: opcode.Bior,
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x74, 0xD0},
			ExpectedOpcode: opcode.Bior,
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x67, 0xB0},
			ExpectedOpcode: opcode.Bist,
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x67, 0xD0},
			ExpectedOpcode: opcode.Bist,
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x75, 0xB0},
			ExpectedOpcode: opcode.Bixor,
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x75, 0xD0},
			ExpectedOpcode: opcode.Bixor,
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x77, 0x30},
			ExpectedOpcode: opcode.Bld,
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x77, 0x50},
			ExpectedOpcode: opcode.Bld,
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x71, 0x30},
			ExpectedOpcode: opcode.Bnot,
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x71, 0x50},
			ExpectedOpcode: opcode.Bnot,
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x61, 0x70},
			ExpectedOpcode: opcode.Bnot,
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x61, 0xE0},
			ExpectedOpcode: opcode.Bnot,
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x74, 0x30},
			ExpectedOpcode: opcode.Bor,
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x74, 0x50},
			ExpectedOpcode: opcode.Bor,
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x70, 0x30},
			ExpectedOpcode: opcode.Bset,
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x70, 0x50},
			ExpectedOpcode: opcode.Bset,
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x60, 0x70},
			ExpectedOpcode: opcode.Bset,
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x60, 0xE0},
			ExpectedOpcode: opcode.Bset,
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x67, 0x30},
			ExpectedOpcode: opcode.Bst,
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x67, 0x50},
			ExpectedOpcode: opcode.Bst,
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x73, 0x30},
			ExpectedOpcode: opcode.Btst,
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x73, 0x50},
			ExpectedOpcode: opcode.Btst,
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x63, 0x70},
			ExpectedOpcode: opcode.Btst,
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x63, 0xE0},
			ExpectedOpcode: opcode.Btst,
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x75, 0x30},
			ExpectedOpcode: opcode.Bxor,
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x75, 0x50},
			ExpectedOpcode: opcode.Bxor,
		},
		{
			Input:          []byte{0x01, 0x41, 0x07, 0x6A},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x03, 0x14},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x01, 0x41, 0x69, 0x30},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6F, 0x30, 0x01, 0x23},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x78, 0x30, 0x6B, 0x20, 0x12, 0x34, 0x56, 0x78},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6D, 0x30},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6B, 0x00, 0x01, 0x23},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6B, 0x20, 0x12, 0x34, 0x56, 0x78},
			ExpectedOpcode: opcode.Ldc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x10, 0x6D, 0x73},
			ExpectedOpcode: opcode.Ldm,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x20, 0x6D, 0x76},
			ExpectedOpcode: opcode.Ldm,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x30, 0x6D, 0x77},
			ExpectedOpcode: opcode.Ldm,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x41, 0x04, 0x7B},
			ExpectedOpcode: opcode.Orc,
		},
		{
			Input:          []byte{0x12, 0xC8},
			ExpectedOpcode: opcode.Rotl,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x12, 0xD9},
			ExpectedOpcode: opcode.Rotl,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x12, 0xF6},
			ExpectedOpcode: opcode.Rotl,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x13, 0xC8},
			ExpectedOpcode: opcode.Rotr,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x13, 0xD9},
			ExpectedOpcode: opcode.Rotr,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x13, 0xF6},
			ExpectedOpcode: opcode.Rotr,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x12, 0x48},
			ExpectedOpcode: opcode.Rotxl,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x12, 0x59},
			ExpectedOpcode: opcode.Rotxl,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x12, 0x76},
			ExpectedOpcode: opcode.Rotxl,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x13, 0x48},
			ExpectedOpcode: opcode.Rotxr,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x13, 0x59},
			ExpectedOpcode: opcode.Rotxr,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x13, 0x76},
			ExpectedOpcode: opcode.Rotxr,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x10, 0xC5},
			ExpectedOpcode: opcode.Shal,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x10, 0xDE},
			ExpectedOpcode: opcode.Shal,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x10, 0xF5},
			ExpectedOpcode: opcode.Shal,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x11, 0xC5},
			ExpectedOpcode: opcode.Shar,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x11, 0xDE},
			ExpectedOpcode: opcode.Shar,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x11, 0xF5},
			ExpectedOpcode: opcode.Shar,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x10, 0x45},
			ExpectedOpcode: opcode.Shll,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x10, 0x5E},
			ExpectedOpcode: opcode.Shll,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x10, 0x75},
			ExpectedOpcode: opcode.Shll,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x11, 0x45},
			ExpectedOpcode: opcode.Shlr,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x11, 0x5E},
			ExpectedOpcode: opcode.Shlr,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x11, 0x75},
			ExpectedOpcode: opcode.Shlr,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x02, 0x15},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Byte,
		},
		{
			Input:          []byte{0x01, 0x41, 0x69, 0xD0},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6F, 0xD0, 0x01, 0x23},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x78, 0x50, 0x6B, 0xA0, 0x12, 0x34, 0x56, 0x78},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6D, 0xD0},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6B, 0x80, 0x01, 0x23},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6B, 0xA0, 0x12, 0x34, 0x56, 0x78},
			ExpectedOpcode: opcode.Stc,
			expectedBWL:    size.Word,
		},
		{
			Input:          []byte{0x01, 0x10, 0x6D, 0xF2},
			ExpectedOpcode: opcode.Stm,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x20, 0x6D, 0xF4},
			ExpectedOpcode: opcode.Stm,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0x30, 0x6D, 0xF0},
			ExpectedOpcode: opcode.Stm,
			expectedBWL:    size.Longword,
		},
		{
			Input:          []byte{0x01, 0xE0, 0x7B, 0x4C},
			ExpectedOpcode: opcode.Tas,
		},
		{
			Input:          []byte{0x01, 0x41, 0x05, 0x9E},
			ExpectedOpcode: opcode.Xorc,
		},
		{
			Input:          []byte{0x58, 0x00, 0xFB, 0x16},
			ExpectedOpcode: opcode.Bra,
		},
		{
			Input:          []byte{0x58, 0x10, 0xFB, 0x72},
			ExpectedOpcode: opcode.Brn,
		},
		// Below test cases disabled as h8s2000 CPU does not support them.
		// {
		// 	Input:          []byte{0x03, 0x24},
		// 	ExpectedOpcode: "ldmac",
		// },
		// {
		// 	Input:          []byte{0x03, 0x35},
		// 	ExpectedOpcode: "ldmac",
		// },
		// {
		// 	Input:          []byte{0x01, 0x60, 0x6D, 0x35},
		// 	ExpectedOpcode: "mac",
		// },
		// {
		// 	Input:          []byte{0x01, 0xA0},
		// 	ExpectedOpcode: "clrmac",
		// },
		// {
		// 	Input:          []byte{0x02, 0x24},
		// 	ExpectedOpcode: "stmac",
		// },
		// {
		// 	Input:          []byte{0x02, 0x31},
		// 	ExpectedOpcode: "stmac",
		// },
	}
	for _, tc := range testCases {
		// Decode expects the input to potentially be up to 10 bytes, and will
		// panic if provided with less and it's decoding incorrectly. To get
		// nicer errors, we pad here.
		paddedInput := make([]byte, 10)
		copy(paddedInput, tc.Input)
		inst := disassembler.Decode(paddedInput)
		assert.Equal(t, tc.ExpectedOpcode, inst.Opcode, fmt.Sprintf("For byte sequence %x, expected opcode %s, got %s\n", tc.Input, tc.ExpectedOpcode, inst.Opcode))
		assert.Equal(t, tc.expectedBWL, inst.BWL, fmt.Sprintf("For byte sequence %x, expected BWL %s, got %s\n", tc.Input, tc.expectedBWL, inst.BWL))
		assert.Equal(t, len(tc.Input), inst.TotalBytes, fmt.Sprintf("For byte sequence %x, expected %d bytes, got %d bytes\n", tc.Input, len(tc.Input), inst.TotalBytes))
	}
}

// HLToBVA decides whether to range the 4 most significant bits (h) or the 4
// least significant bits (l) of the particular byte.
type HLToBVA int

const (
	H HLToBVA = iota
	L
)

type Range int

const (
	zeroToSeven Range = iota
	eightToF
	zeroToF
	zeroToSix
	zeroTo3
	eightToE
)

// These tests ensure that instruction sequences take into account that
// it is sometimes only valid if the registers in question are in a range.
func TestDecodeRangeCases(t *testing.T) {
	ranges := map[Range][]int{
		zeroToSeven: {0x0, 0x7},
		eightToF:    {0x8, 0xF},
		zeroToF:     {0x0, 0xF},
		zeroToSix:   {0x0, 0x6}, // For mov that isn't pop
		eightToE:    {0x8, 0xE}, // For mov that isn't push
		zeroTo3:     {0x0, 0x3}, // Trapa
	}
	testCases := []struct {
		Input          []byte
		ExpectedOpcode opcode.Opcode
		ByteToBVA      int
		HLToBVA        HLToBVA
		Range          []int
	}{
		{
			Input:          []byte{0x7A, 0x1F, 0x00, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.Add,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x0A, 0xF0},
			ExpectedOpcode: opcode.Add,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			Input:          []byte{0x0A, 0xF0},
			ExpectedOpcode: opcode.Add,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x0B, 0x0F},
			ExpectedOpcode: opcode.Adds,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x0B, 0x8F},
			ExpectedOpcode: opcode.Adds,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x0B, 0x9F},
			ExpectedOpcode: opcode.Adds,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x7a, 0x6F, 0x00, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.And,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x01, 0xF0, 0x66, 0xF0},
			ExpectedOpcode: opcode.And,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x01, 0xF0, 0x66, 0x0F},
			ExpectedOpcode: opcode.And,
			ByteToBVA:      3,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x76, 0xF0},
			ExpectedOpcode: opcode.Band,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x7C, 0xF0, 0x76, 0x00},
			ExpectedOpcode: opcode.Band,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x7C, 0x00, 0x76, 0x00},
			ExpectedOpcode: opcode.Band,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x7E, 0x00, 0x76, 0x00},
			ExpectedOpcode: opcode.Band,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x6A, 0x10, 0x00, 0x00, 0x76, 0xF0},
			ExpectedOpcode: opcode.Band,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x76, 0xF0},
			ExpectedOpcode: opcode.Band,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x72, 0xF0},
			ExpectedOpcode: opcode.Bclr,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x7D, 0xF0, 0x72, 0x00},
			ExpectedOpcode: opcode.Bclr,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x7D, 0x00, 0x72, 0xF0},
			ExpectedOpcode: opcode.Bclr,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x7F, 0x00, 0x72, 0xF0},
			ExpectedOpcode: opcode.Bclr,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 18 00 00 72 F0
			Input:          []byte{0x6A, 0x18, 0x00, 0x00, 0x72, 0xF0},
			ExpectedOpcode: opcode.Bclr,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 38 00 00 00 00 72 F0
			Input:          []byte{0x6A, 0x38, 0x00, 0x00, 0x00, 0x00, 0x72, 0xF0},
			ExpectedOpcode: opcode.Bclr,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			Input:          []byte{0x7D, 0xF0, 0x62, 0x00},
			ExpectedOpcode: opcode.Bclr,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 76 F0
			Input:          []byte{0x76, 0xF0},
			ExpectedOpcode: opcode.Biand,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 7C F0 76 D0
			Input:          []byte{0x7C, 0xF0, 0x76, 0xD0},
			ExpectedOpcode: opcode.Biand,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C 00 76 F0
			Input:          []byte{0x7C, 0x00, 0x76, 0xF0},
			ExpectedOpcode: opcode.Biand,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 7E 00 76 F0
			Input:          []byte{0x7E, 0x00, 0x76, 0xF0},
			ExpectedOpcode: opcode.Biand,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 6A 10 00 00 76 F0
			Input:          []byte{0x6A, 0x10, 0x00, 0x00, 0x76, 0xF0},
			ExpectedOpcode: opcode.Biand,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 6A 30 00 00 00 00 76 F0
			Input:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x76, 0xF0},
			ExpectedOpcode: opcode.Biand,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 77 F0
			Input:          []byte{0x77, 0xF0},
			ExpectedOpcode: opcode.Bild,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 7C F0 77 F0
			Input:          []byte{0x7C, 0xF0, 0x77, 0xF0},
			ExpectedOpcode: opcode.Bild,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C 00 77 F0
			Input:          []byte{0x7C, 0x00, 0x77, 0xF0},
			ExpectedOpcode: opcode.Bild,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 6A 10 00 00 77 F0
			Input:          []byte{0x6A, 0x10, 0x00, 0x00, 0x77, 0xF0},
			ExpectedOpcode: opcode.Bild,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 6A 30 00 00 00 00 77 F0
			Input:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x77, 0xF0},
			ExpectedOpcode: opcode.Bild,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 74 F0
			Input:          []byte{0x74, 0xF0},
			ExpectedOpcode: opcode.Bior,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 7C F0 74 F0
			Input:          []byte{0x7C, 0xF0, 0x74, 0xF0},
			ExpectedOpcode: opcode.Bior,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C 00 74 F0
			Input:          []byte{0x7C, 0x00, 0x74, 0xF0},
			ExpectedOpcode: opcode.Bior,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 7E 00 74 F0
			Input:          []byte{0x7E, 0x00, 0x74, 0xF0},
			ExpectedOpcode: opcode.Bior,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 6A 10 00 00 74 F0
			Input:          []byte{0x6A, 0x10, 0x00, 0x00, 0x74, 0xF0},
			ExpectedOpcode: opcode.Bior,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 6A 30 00 00 00 00 74 F0
			Input:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x74, 0xF0},
			ExpectedOpcode: opcode.Bior,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 67 F0
			Input:          []byte{0x67, 0xF0},
			ExpectedOpcode: opcode.Bist,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 7D F0 67 F0
			Input:          []byte{0x7D, 0xF0, 0x67, 0xF0},
			ExpectedOpcode: opcode.Bist,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7D 00 67 F0
			Input:          []byte{0x7D, 0x00, 0x67, 0xF0},
			ExpectedOpcode: opcode.Bist,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 7F 00 67 F0
			Input:          []byte{0x7F, 0x00, 0x67, 0xF0},
			ExpectedOpcode: opcode.Bist,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 6A 18 00 00 67 F0
			Input:          []byte{0x6A, 0x18, 0x00, 0x00, 0x67, 0xF0},
			ExpectedOpcode: opcode.Bist,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 6A 38 00 00 00 00 67 F0
			Input:          []byte{0x6A, 0x38, 0x00, 0x00, 0x00, 0x00, 0x67, 0xF0},
			ExpectedOpcode: opcode.Bist,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 75 F0
			Input:          []byte{0x75, 0xF0},
			ExpectedOpcode: opcode.Bixor,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 7C F0 75 F0
			Input:          []byte{0x7C, 0xF0, 0x75, 0xF0},
			ExpectedOpcode: opcode.Bixor,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C 00 75 F0
			Input:          []byte{0x7C, 0x00, 0x75, 0xF0},
			ExpectedOpcode: opcode.Bixor,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 7E 00 75 F0
			Input:          []byte{0x7E, 0x00, 0x75, 0xF0},
			ExpectedOpcode: opcode.Bixor,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 6A 10 00 00 75 F0
			Input:          []byte{0x6A, 0x10, 0x00, 0x00, 0x75, 0xF0},
			ExpectedOpcode: opcode.Bixor,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 6A 30 00 00 00 00 75 F0
			Input:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x75, 0xF0},
			ExpectedOpcode: opcode.Bixor,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 77 ?0
			Input:          []byte{0x77, 0x00},
			ExpectedOpcode: opcode.Bld,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C F0 77 00
			Input:          []byte{0x7C, 0xF0, 0x77, 0x00},
			ExpectedOpcode: opcode.Bld,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C 00 77 F0
			Input:          []byte{0x7C, 0x00, 0x77, 0xF0},
			ExpectedOpcode: opcode.Bld,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7E 00 77 F0
			Input:          []byte{0x7E, 0x00, 0x77, 0xF0},
			ExpectedOpcode: opcode.Bld,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 10 00 00 77 F0
			Input:          []byte{0x6A, 0x10, 0x00, 0x00, 0x77, 0xF0},
			ExpectedOpcode: opcode.Bld,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 30 00 00 00 00 77 F0
			Input:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x77, 0xF0},
			ExpectedOpcode: opcode.Bld,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 71 F0
			Input:          []byte{0x71, 0xF0},
			ExpectedOpcode: opcode.Bnot,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7D F0 71 00
			Input:          []byte{0x7D, 0xF0, 0x71, 0x00},
			ExpectedOpcode: opcode.Bnot,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7D 00 71 F0
			Input:          []byte{0x7D, 0x00, 0x71, 0xF0},
			ExpectedOpcode: opcode.Bnot,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7F 00 71 F0
			Input:          []byte{0x7F, 0x00, 0x71, 0xF0},
			ExpectedOpcode: opcode.Bnot,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 18 00 00 71 F0
			Input:          []byte{0x6A, 0x18, 0x00, 0x00, 0x71, 0xF0},
			ExpectedOpcode: opcode.Bnot,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 38 00 00 00 00 71 F0
			Input:          []byte{0x6A, 0x38, 0x00, 0x00, 0x00, 0x00, 0x71, 0xF0},
			ExpectedOpcode: opcode.Bnot,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7D F0 61 00
			Input:          []byte{0x7D, 0xF0, 0x61, 0x00},
			ExpectedOpcode: opcode.Bnot,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 74 F0
			Input:          []byte{0x74, 0xF0},
			ExpectedOpcode: opcode.Bor,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C F0 74 00
			Input:          []byte{0x7C, 0xF0, 0x74, 0x00},
			ExpectedOpcode: opcode.Bor,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C 00 74 F0
			Input:          []byte{0x7C, 0x00, 0x74, 0xF0},
			ExpectedOpcode: opcode.Bor,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7E 00 74 F0
			Input:          []byte{0x7E, 0x00, 0x74, 0xF0},
			ExpectedOpcode: opcode.Bor,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 10 00 00 74 F0
			Input:          []byte{0x6A, 0x10, 0x00, 0x00, 0x74, 0xF0},
			ExpectedOpcode: opcode.Bor,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 30 00 00 00 00 74 F0
			Input:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x74, 0xF0},
			ExpectedOpcode: opcode.Bor,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 70 F0
			Input:          []byte{0x70, 0xF0},
			ExpectedOpcode: opcode.Bset,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7D F0 70 00
			Input:          []byte{0x7D, 0xF0, 0x70, 0x00},
			ExpectedOpcode: opcode.Bset,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7D 00 70 F0
			Input:          []byte{0x7D, 0x00, 0x70, 0xF0},
			ExpectedOpcode: opcode.Bset,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7F 00 70 F0
			Input:          []byte{0x7F, 0x00, 0x70, 0xF0},
			ExpectedOpcode: opcode.Bset,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 18 00 00 70 F0
			Input:          []byte{0x6A, 0x18, 0x00, 0x00, 0x70, 0xF0},
			ExpectedOpcode: opcode.Bset,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 38 00 00 00 00 70 F0
			Input:          []byte{0x6A, 0x38, 0x00, 0x00, 0x00, 0x00, 0x70, 0xF0},
			ExpectedOpcode: opcode.Bset,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7D F0 60 00
			Input:          []byte{0x7D, 0xF0, 0x60, 0x00},
			ExpectedOpcode: opcode.Bset,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 67 F0
			Input:          []byte{0x67, 0xF0},
			ExpectedOpcode: opcode.Bst,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7D F0 67 00
			Input:          []byte{0x7D, 0xF0, 0x67, 0x00},
			ExpectedOpcode: opcode.Bst,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7D 00 67 F0
			Input:          []byte{0x7D, 0x00, 0x67, 0xF0},
			ExpectedOpcode: opcode.Bst,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7F 00 67 F0
			Input:          []byte{0x7F, 0x00, 0x67, 0xF0},
			ExpectedOpcode: opcode.Bst,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 18 00 00 67 F0
			Input:          []byte{0x6A, 0x18, 0x00, 0x00, 0x67, 0xF0},
			ExpectedOpcode: opcode.Bst,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 38 00 00 00 00 67 F0
			Input:          []byte{0x6A, 0x38, 0x00, 0x00, 0x00, 0x00, 0x67, 0xF0},
			ExpectedOpcode: opcode.Bst,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 73 F0
			Input:          []byte{0x73, 0xF0},
			ExpectedOpcode: opcode.Btst,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C F0 73 00
			Input:          []byte{0x7C, 0xF0, 0x73, 0x00},
			ExpectedOpcode: opcode.Btst,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C 00 73 F0
			Input:          []byte{0x7C, 0x00, 0x73, 0xF0},
			ExpectedOpcode: opcode.Btst,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7E 00 73 F0
			Input:          []byte{0x7E, 0x00, 0x73, 0xF0},
			ExpectedOpcode: opcode.Btst,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 10 00 00 73 F0
			Input:          []byte{0x6A, 0x10, 0x00, 0x00, 0x73, 0xF0},
			ExpectedOpcode: opcode.Btst,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 30 00 00 00 00 73 F0
			Input:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x73, 0xF0},
			ExpectedOpcode: opcode.Btst,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C F0 63 00
			Input:          []byte{0x7C, 0xF0, 0x63, 0x00},
			ExpectedOpcode: opcode.Btst,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 75 F0
			Input:          []byte{0x75, 0xF0},
			ExpectedOpcode: opcode.Bxor,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C F0 75 00
			Input:          []byte{0x7C, 0xF0, 0x75, 0x00},
			ExpectedOpcode: opcode.Bxor,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7C 00 75 F0
			Input:          []byte{0x7C, 0x00, 0x75, 0xF0},
			ExpectedOpcode: opcode.Bxor,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7E 00 75 F0
			Input:          []byte{0x7E, 0x00, 0x75, 0xF0},
			ExpectedOpcode: opcode.Bxor,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 10 00 00 75 F0
			Input:          []byte{0x6A, 0x10, 0x00, 0x00, 0x75, 0xF0},
			ExpectedOpcode: opcode.Bxor,
			ByteToBVA:      5,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6A 30 00 00 00 00 75 F0
			Input:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x75, 0xF0},
			ExpectedOpcode: opcode.Bxor,
			ByteToBVA:      7,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7A 2F 00 00 00 00
			Input:          []byte{0x7A, 0x2F, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.Cmp,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 1F F0
			Input:          []byte{0x1F, 0xF0},
			ExpectedOpcode: opcode.Cmp,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 1F FF
			Input:          []byte{0x1F, 0xFF},
			ExpectedOpcode: opcode.Cmp,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 1B 7F
			Input:          []byte{0x1B, 0x7F},
			ExpectedOpcode: opcode.Dec,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 1B FF
			Input:          []byte{0x1B, 0xFF},
			ExpectedOpcode: opcode.Dec,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 D0 53 0F
			Input:          []byte{0x01, 0xD0, 0x53, 0x0F},
			ExpectedOpcode: opcode.Divxs,
			ByteToBVA:      3,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 53 0F
			Input:          []byte{0x53, 0x0F},
			ExpectedOpcode: opcode.Divxu,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 17 FF
			Input:          []byte{0x17, 0xFF},
			ExpectedOpcode: opcode.Exts,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 17 7F
			Input:          []byte{0x17, 0x7F},
			ExpectedOpcode: opcode.Extu,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 0B 7F
			Input:          []byte{0x0B, 0x7F},
			ExpectedOpcode: opcode.Inc,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 59 F0
			Input:          []byte{0x59, 0xF0},
			ExpectedOpcode: opcode.Jmp,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 5D F0
			Input:          []byte{0x5D, 0xF0},
			ExpectedOpcode: opcode.Jsr,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 40 69 F0
			Input:          []byte{0x01, 0x40, 0x69, 0xF0},
			ExpectedOpcode: opcode.Ldc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 40 6F F0 00 00
			Input:          []byte{0x01, 0x40, 0x6F, 0xF0, 0x00, 0x00},
			ExpectedOpcode: opcode.Ldc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 40 78 F0 6B 20 FF FF FF FF
			Input:          []byte{0x01, 0x40, 0x78, 0xF0, 0x6B, 0x20, 0xFF, 0xFF, 0xFF, 0xFF},
			ExpectedOpcode: opcode.Ldc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 40 6D F0
			Input:          []byte{0x01, 0x40, 0x6D, 0xF0},
			ExpectedOpcode: opcode.Ldc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		// exr cases
		{
			// 01 41 69 F0
			Input:          []byte{0x01, 0x41, 0x69, 0xF0},
			ExpectedOpcode: opcode.Ldc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 41 6F F0 00 00
			Input:          []byte{0x01, 0x41, 0x6F, 0xF0, 0x00, 0x00},
			ExpectedOpcode: opcode.Ldc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 41 78 F0 6B 20 FF FF FF FF
			Input:          []byte{0x01, 0x41, 0x78, 0xF0, 0x6B, 0x20, 0xFF, 0xFF, 0xFF, 0xFF},
			ExpectedOpcode: opcode.Ldc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 41 6D F0
			Input:          []byte{0x01, 0x41, 0x6D, 0xF0},
			ExpectedOpcode: opcode.Ldc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 10 6D 7F
			Input:          []byte{0x01, 0x10, 0x6D, 0x7F},
			ExpectedOpcode: opcode.Ldm,
			ByteToBVA:      3,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 0F F0
			Input:          []byte{0x0F, 0xF0},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 0F F0
			Input:          []byte{0x0F, 0xF0},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 68 F0
			Input:          []byte{0x68, 0xF0},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToF],
		},
		{
			// 6E F0 00 00
			Input:          []byte{0x6E, 0xF0, 0x00, 0x00},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToF],
		},
		{
			// 78 F0 6A 20 00 00 00 00
			Input:          []byte{0x78, 0xF0, 0x6A, 0x20, 0x00, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 6C F0
			Input:          []byte{0x6C, 0xF0},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToF],
		},
		{
			// 69 F0
			Input:          []byte{0x69, 0xF0},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToF],
		},
		{
			// 6F F0 00 00
			Input:          []byte{0x6F, 0xF0, 0x00, 0x00},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToF],
		},
		{
			// 78 F0 6B 20 00 00 00 00
			Input:          []byte{0x78, 0xF0, 0x6B, 0x20, 0x00, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		// Below cases are difficult to test with this structure. TODO.
		// {
		// 	// 6D F0
		// 	Input:          []byte{0x6D, 0xF0},
		// 	ExpectedOpcode: opcode.Mov, // rly mov
		// 	ByteToBVA:      1,
		// 	HLToBVA:        H,
		// 	Range:          ranges[zeroToSix],
		// },
		// {
		// 	// 6D F0
		// 	Input:          []byte{0x6D, 0xF0},
		// 	ExpectedOpcode: opcode.Push, // rly mov
		// 	ByteToBVA:      1,
		// 	HLToBVA:        H,
		// 	Range:          ranges[eightToE],
		// },
		{
			// 7A 0F
			Input:          []byte{0x7A, 0x0F},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 00 69 F0
			Input:          []byte{0x01, 0x00, 0x69, 0xF0},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToF],
		},
		{
			// 01 00 6F F0
			Input:          []byte{0x01, 0x00, 0x6F, 0xF0},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToF],
		},
		{
			// 01 00 69 0F
			Input:          []byte{0x01, 0x00, 0x69, 0x0F},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      3,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 00 78 F0 6B A0 00 00 00 00
			Input:          []byte{0x01, 0x00, 0x78, 0xF0, 0x6B, 0xA0, 0x00, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 00 78 00 6B AF 00 00 00 00
			Input:          []byte{0x01, 0x00, 0x78, 0x00, 0x6B, 0xAF, 0x00, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      5,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 00 6D F0
			Input:          []byte{0x01, 0x00, 0x6D, 0xF0},
			ExpectedOpcode: opcode.Mov,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSix],
		},
		{
			// 01 00 6D FF
			Input:          []byte{0x01, 0x00, 0x6D, 0xFF},
			ExpectedOpcode: opcode.Push,
			ByteToBVA:      3,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 C0 52 0F
			Input:          []byte{0x01, 0xC0, 0x52, 0x0F},
			ExpectedOpcode: opcode.Mulxs,
			ByteToBVA:      3,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 52 0F
			Input:          []byte{0x52, 0x0F},
			ExpectedOpcode: opcode.Mulxu,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 17 BF
			Input:          []byte{0x17, 0xBF},
			ExpectedOpcode: opcode.Neg,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 17 3F
			Input:          []byte{0x17, 0x3F},
			ExpectedOpcode: opcode.Not,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7A 4F 00 00 00 00
			Input:          []byte{0x7A, 0x4F, 0x00, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.Or,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 F0 64 F0
			Input:          []byte{0x01, 0xF0, 0x64, 0xF0},
			ExpectedOpcode: opcode.Or,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 F0 64 0F
			Input:          []byte{0x01, 0xF0, 0x64, 0x0F},
			ExpectedOpcode: opcode.Or,
			ByteToBVA:      3,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 12 BF
			Input:          []byte{0x12, 0xBF},
			ExpectedOpcode: opcode.Rotl,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 12 FF
			Input:          []byte{0x12, 0xFF},
			ExpectedOpcode: opcode.Rotl,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 13 BF
			Input:          []byte{0x13, 0xBF},
			ExpectedOpcode: opcode.Rotr,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 13 FF
			Input:          []byte{0x13, 0xFF},
			ExpectedOpcode: opcode.Rotr,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 12 3F
			Input:          []byte{0x12, 0x3F},
			ExpectedOpcode: opcode.Rotxl,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 12 7F
			Input:          []byte{0x12, 0x7F},
			ExpectedOpcode: opcode.Rotxl,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 13 3F
			Input:          []byte{0x13, 0x3F},
			ExpectedOpcode: opcode.Rotxr,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 13 7F
			Input:          []byte{0x13, 0x7F},
			ExpectedOpcode: opcode.Rotxr,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 10 BF
			Input:          []byte{0x10, 0xBF},
			ExpectedOpcode: opcode.Shal,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 10 FF
			Input:          []byte{0x10, 0xFF},
			ExpectedOpcode: opcode.Shal,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 11 BF
			Input:          []byte{0x11, 0xBF},
			ExpectedOpcode: opcode.Shar,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 11 FF
			Input:          []byte{0x11, 0xFF},
			ExpectedOpcode: opcode.Shar,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		// START
		{
			// 10 3F
			Input:          []byte{0x10, 0x3F},
			ExpectedOpcode: opcode.Shll,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 10 7F
			Input:          []byte{0x10, 0x7F},
			ExpectedOpcode: opcode.Shll,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 11 3F
			Input:          []byte{0x11, 0x3F},
			ExpectedOpcode: opcode.Shlr,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 11 7F
			Input:          []byte{0x11, 0x7F},
			ExpectedOpcode: opcode.Shlr,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 40 69 F0
			Input:          []byte{0x01, 0x40, 0x69, 0xF0},
			ExpectedOpcode: opcode.Stc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 01 40 6F F0 00 00
			Input:          []byte{0x01, 0x40, 0x6F, 0xF0, 0x00, 0x00},
			ExpectedOpcode: opcode.Stc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 01 40 78 F0 6B A0 FF FF FF FF
			Input:          []byte{0x01, 0x40, 0x78, 0xF0, 0x6B, 0xA0, 0xFF, 0xFF, 0xFF, 0xFF},
			ExpectedOpcode: opcode.Stc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 40 6D F0
			Input:          []byte{0x01, 0x40, 0x6D, 0xF0},
			ExpectedOpcode: opcode.Stc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		// exr cases
		{
			// 01 41 69 F0
			Input:          []byte{0x01, 0x41, 0x69, 0xF0},
			ExpectedOpcode: opcode.Stc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 01 41 6F F0 00 00
			Input:          []byte{0x01, 0x41, 0x6F, 0xF0, 0x00, 0x00},
			ExpectedOpcode: opcode.Stc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 01 41 78 F0 6B A0 FF FF FF FF
			Input:          []byte{0x01, 0x41, 0x78, 0xF0, 0x6B, 0xA0, 0xFF, 0xFF, 0xFF, 0xFF},
			ExpectedOpcode: opcode.Stc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 41 6D F0
			Input:          []byte{0x01, 0x41, 0x6D, 0xF0},
			ExpectedOpcode: opcode.Stc,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 01 10 6D FF
			Input:          []byte{0x01, 0x10, 0x6D, 0xFF},
			ExpectedOpcode: opcode.Stm,
			ByteToBVA:      3,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 20 6D FF
			Input:          []byte{0x01, 0x20, 0x6D, 0xFF},
			ExpectedOpcode: opcode.Stm,
			ByteToBVA:      3,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 30 6D FF
			Input:          []byte{0x01, 0x30, 0x6D, 0xFF},
			ExpectedOpcode: opcode.Stm,
			ByteToBVA:      3,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 7A 3F 00 00 00 00
			Input:          []byte{0x7A, 0x3F, 0x00, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.Sub,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 1A F0
			Input:          []byte{0x1A, 0xF0},
			ExpectedOpcode: opcode.Sub,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[eightToF],
		},
		{
			// 1A FF
			Input:          []byte{0x1A, 0xFF},
			ExpectedOpcode: opcode.Sub,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 1B 0F
			Input:          []byte{0x1B, 0x0F},
			ExpectedOpcode: opcode.Subs,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 1B 8F
			Input:          []byte{0x1B, 0x8F},
			ExpectedOpcode: opcode.Subs,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 1B 9F
			Input:          []byte{0x1B, 0x9F},
			ExpectedOpcode: opcode.Subs,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 E0 7B FC
			Input:          []byte{0x01, 0xE0, 0x7B, 0xFC},
			ExpectedOpcode: opcode.Tas,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 57 F0
			Input:          []byte{0x57, 0xF0},
			ExpectedOpcode: opcode.Trapa,
			ByteToBVA:      1,
			HLToBVA:        H,
			Range:          ranges[zeroTo3],
		},
		{
			// 7A 5F 00 00 00 00
			Input:          []byte{0x7A, 0x5F, 0x00, 0x00, 0x00, 0x00},
			ExpectedOpcode: opcode.Xor,
			ByteToBVA:      1,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 F0 65 F0
			Input:          []byte{0x01, 0xF0, 0x65, 0xF0},
			ExpectedOpcode: opcode.Xor,
			ByteToBVA:      3,
			HLToBVA:        H,
			Range:          ranges[zeroToSeven],
		},
		{
			// 01 F0 65 0F
			Input:          []byte{0x01, 0xF0, 0x65, 0x0F},
			ExpectedOpcode: opcode.Xor,
			ByteToBVA:      3,
			HLToBVA:        L,
			Range:          ranges[zeroToSeven],
		},
	}
	for _, tc := range testCases {
		// Decode expects the input to potentially be up to 10 bytes, and will
		// panic if provided with less and it's decoding incorrectly. To get
		// nicer errors, we pad here.
		paddedInput := make([]byte, 10)
		copy(paddedInput, tc.Input)

		valsToTest := []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF}

		for _, b := range valsToTest {
			if tc.HLToBVA == L {
				// Replace last 4 bits of paddedInput[ByteToBVA] with b
				paddedInput[tc.ByteToBVA] = paddedInput[tc.ByteToBVA]&0xF0 | b
			} else {
				// Replace first 4 bits of paddedInput[ByteToBVA] with b
				paddedInput[tc.ByteToBVA] = paddedInput[tc.ByteToBVA]&0x0F | b<<4
			}
			inst := disassembler.Decode(paddedInput)
			if int(b) >= tc.Range[0] && int(b) <= tc.Range[1] {
				// If it's in the range we expect to be successful, test it and check opcode is what we wanted
				assert.Equal(t, tc.ExpectedOpcode, inst.Opcode, fmt.Sprintf("For byte sequence %x, expected opcode %s, got %s\n", paddedInput, tc.ExpectedOpcode, inst.Opcode))
			} else {
				assert.NotEqual(t, tc.ExpectedOpcode, inst.Opcode, fmt.Sprintf("For byte sequence %x, expected opcode not to be %s\n", paddedInput, tc.ExpectedOpcode))
			}
		}
	}
}
