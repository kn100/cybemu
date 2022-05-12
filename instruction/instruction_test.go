package instruction_test

import (
	"testing"

	"github.com/kn100/cybemu/addressingmode"
	"github.com/kn100/cybemu/instruction"
	"github.com/kn100/cybemu/opcode"
	"github.com/kn100/cybemu/operand"
	"github.com/kn100/cybemu/size"
	"github.com/stretchr/testify/assert"
)

func TestDetermineOperandTypeAndSetData(t *testing.T) {
	testCases := []struct {
		instruction         instruction.Inst
		expectedOperandType operand.OperandType
		expectedRegSrc      []byte
		expectedRegDst      []byte
		expectedString      string
	}{
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x00, 0x00},
				Opcode:         opcode.Nop,
				BWL:            size.Unset,
				AddressingMode: addressingmode.None,
			},
			expectedRegSrc:      nil,
			expectedRegDst:      nil,
			expectedOperandType: operand.None,
			expectedString:      "nop",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x08, 0x3E},
				Opcode:         opcode.Add,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x03},
			expectedRegDst:      []byte{0x0E},
			expectedOperandType: operand.R8_R8,
			expectedString:      "add.b r3h, r6l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x09, 0x0B},
				Opcode:         opcode.Add,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x00},
			expectedRegDst:      []byte{0x0B},
			expectedOperandType: operand.R16_R16,
			expectedString:      "add.w r0, e3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0A, 0xE0},
				Opcode:         opcode.Add,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0E},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.R32_R32_S2,
			expectedString:      "add.l er6, er0",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0A, 0xF7},
				Opcode:         opcode.Add,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0F},
			expectedRegDst:      []byte{0x07},
			expectedOperandType: operand.R32_R32_S2,
			expectedString:      "add.l sp, sp",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0B, 0x06},
				Opcode:         opcode.Adds,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      nil,
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_ADDS_SUBS,
			expectedString:      "adds #1, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0B, 0x86},
				Opcode:         opcode.Adds,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      nil,
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_ADDS_SUBS,
			expectedString:      "adds #2, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0B, 0x96},
				Opcode:         opcode.Adds,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      nil,
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_ADDS_SUBS,
			expectedString:      "adds #4, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0E, 0xC0},
				Opcode:         opcode.Addx,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0C},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.R8_R8,
			expectedString:      "addx r4l, r0h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x16, 0x32},
				Opcode:         opcode.And,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x03},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.R8_R8,
			expectedString:      "and.b r3h, r2h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x66, 0x2E},
				Opcode:         opcode.And,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x02},
			expectedRegDst:      []byte{0x0E},
			expectedOperandType: operand.R16_R16,
			expectedString:      "and.w r2, e6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0xF0, 0x66, 0x07},
				Opcode:         opcode.And,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x00},
			expectedRegDst:      []byte{0x07},
			expectedOperandType: operand.R32_R32_S4,
			expectedString:      "and.l er0, sp",
		},
		{
			// 76 50
			instruction: instruction.Inst{
				Bytes:          []byte{0x76, 0x50},
				Opcode:         opcode.Band,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "band #5, r0h",
		},
		{
			// 72 50
			instruction: instruction.Inst{
				Bytes:          []byte{0x72, 0x50},
				Opcode:         opcode.Bclr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "bclr #5, r0h",
		},
		{
			// 62 93
			instruction: instruction.Inst{
				Bytes:          []byte{0x62, 0x93},
				Opcode:         opcode.Bclr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x09},
			expectedRegDst:      []byte{0x03},
			expectedOperandType: operand.R8_R8,
			expectedString:      "bclr r1l, r3h",
		},
		{
			// 76 C2
			instruction: instruction.Inst{
				Bytes:          []byte{0x76, 0xC2},
				Opcode:         opcode.Biand,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "biand #4, r2h",
		},
		{
			// 77 C2
			instruction: instruction.Inst{
				Bytes:          []byte{0x77, 0xC2},
				Opcode:         opcode.Bild,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "bild #4, r2h",
		},
		{
			// 74 C2
			instruction: instruction.Inst{
				Bytes:          []byte{0x74, 0xC2},
				Opcode:         opcode.Bior,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "bior #4, r2h",
		},
		{
			// 67 C2
			instruction: instruction.Inst{
				Bytes:          []byte{0x67, 0xC2},
				Opcode:         opcode.Bist,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "bist #4, r2h",
		},
		{
			// 75 C2
			instruction: instruction.Inst{
				Bytes:          []byte{0x75, 0xC2},
				Opcode:         opcode.Bixor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "bixor #4, r2h",
		},
		{
			// 77 42
			instruction: instruction.Inst{
				Bytes:          []byte{0x77, 0x42},
				Opcode:         opcode.Bld,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "bld #4, r2h",
		},
		{
			// 71 50
			instruction: instruction.Inst{
				Bytes:          []byte{0x72, 0x50},
				Opcode:         opcode.Bnot,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "bnot #5, r0h",
		},
		{
			// 61 93
			instruction: instruction.Inst{
				Bytes:          []byte{0x61, 0x93},
				Opcode:         opcode.Bnot,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x09},
			expectedRegDst:      []byte{0x03},
			expectedOperandType: operand.R8_R8,
			expectedString:      "bnot r1l, r3h",
		},
		{
			// 74 50
			instruction: instruction.Inst{
				Bytes:          []byte{0x74, 0x50},
				Opcode:         opcode.Bor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "bor #5, r0h",
		},
		{
			// 70 50
			instruction: instruction.Inst{
				Bytes:          []byte{0x70, 0x50},
				Opcode:         opcode.Bset,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "bset #5, r0h",
		},
		{
			// 60 93
			instruction: instruction.Inst{
				Bytes:          []byte{0x60, 0x93},
				Opcode:         opcode.Bset,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x09},
			expectedRegDst:      []byte{0x03},
			expectedOperandType: operand.R8_R8,
			expectedString:      "bset r1l, r3h",
		},
		{
			// 67 50
			instruction: instruction.Inst{
				Bytes:          []byte{0x67, 0x50},
				Opcode:         opcode.Bst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "bst #5, r0h",
		},
		{
			// 73 50
			instruction: instruction.Inst{
				Bytes:          []byte{0x73, 0x50},
				Opcode:         opcode.Btst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "btst #5, r0h",
		},
		{
			// 63 93
			instruction: instruction.Inst{
				Bytes:          []byte{0x63, 0x93},
				Opcode:         opcode.Btst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x09},
			expectedRegDst:      []byte{0x03},
			expectedOperandType: operand.R8_R8,
			expectedString:      "btst r1l, r3h",
		},
		{
			// 75 50
			instruction: instruction.Inst{
				Bytes:          []byte{0x75, 0x50},
				Opcode:         opcode.Bxor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.Ix_R8,
			expectedString:      "bxor #5, r0h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1C, 0x3E},
				Opcode:         opcode.Cmp,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x03},
			expectedRegDst:      []byte{0x0E},
			expectedOperandType: operand.R8_R8,
			expectedString:      "cmp.b r3h, r6l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1D, 0xD2},
				Opcode:         opcode.Cmp,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0D},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.R16_R16,
			expectedString:      "cmp.w e5, r2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1F, 0xB5},
				Opcode:         opcode.Cmp,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0B},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.R32_R32_S2,
			expectedString:      "cmp.l er3, er5",
		},
		{
			instruction: instruction.Inst{
				// 0F 00
				Bytes:          []byte{0x0F, 0x00},
				Opcode:         opcode.Daa,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.R8,
			expectedString:      "daa r0h",
		},
		{
			instruction: instruction.Inst{
				// 0F 0C
				Bytes:          []byte{0x0F, 0x0C},
				Opcode:         opcode.Daa,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0C},
			expectedOperandType: operand.R8,
			expectedString:      "daa r4l",
		},
		{
			instruction: instruction.Inst{
				// 1F 00
				Bytes:          []byte{0x1F, 0x00},
				Opcode:         opcode.Das,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.R8,
			expectedString:      "das r0h",
		},
		{
			instruction: instruction.Inst{
				// 1F 0C
				Bytes:          []byte{0x1F, 0x0C},
				Opcode:         opcode.Das,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0C},
			expectedOperandType: operand.R8,
			expectedString:      "das r4l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1A, 0x05},
				Opcode:         opcode.Dec,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R8_INC_DEC,
			expectedString:      "dec.b #1, r5h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1B, 0x55},
				Opcode:         opcode.Dec,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R16_INC_DEC,
			expectedString:      "dec.w #1, r5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1B, 0xD5},
				Opcode:         opcode.Dec,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R16_INC_DEC,
			expectedString:      "dec.w #2, r5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1B, 0x75},
				Opcode:         opcode.Dec,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R32_INC_DEC,
			expectedString:      "dec.l #1, er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1B, 0xF5},
				Opcode:         opcode.Dec,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R32_INC_DEC,
			expectedString:      "dec.l #2, er5",
		},
	}
	for _, tc := range testCases {
		tc.instruction.DetermineOperandTypeAndSetData()
		assert.Equal(t, tc.expectedOperandType, tc.instruction.OperandType, "expected operand type to be %s, got %s", tc.expectedOperandType, tc.instruction.OperandType)
		assert.Equal(t, tc.expectedRegSrc, tc.instruction.RegSrc, "expected RegSrc %v, got %v", tc.expectedRegSrc, tc.instruction.RegSrc)
		assert.Equal(t, tc.expectedRegDst, tc.instruction.RegDst, "expected RegDst %v, got %v", tc.expectedRegDst, tc.instruction.RegDst)
		assert.Equal(t, tc.expectedString, tc.instruction.String())

	}
}
