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

func TestDetermineOperandTypeAndSetDataRegisterDirect(t *testing.T) {
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
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0xD0, 0x51, 0xBC},
				Opcode:         opcode.Divxs,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0B},
			expectedRegDst:      []byte{0x0C},
			expectedOperandType: operand.R8_R16_MULXS_DIVXS,
			expectedString:      "divxs.b r3l, e4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0xD0, 0x53, 0xB2},
				Opcode:         opcode.Divxs,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0B},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.R16_R32_MULXS_DIVXS,
			expectedString:      "divxs.w e3, er2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x51, 0xBC},
				Opcode:         opcode.Divxu,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0B},
			expectedRegDst:      []byte{0x0C},
			expectedOperandType: operand.R8_R16,
			expectedString:      "divxu.b r3l, e4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x53, 0xB2},
				Opcode:         opcode.Divxu,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0B},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.R16_R32,
			expectedString:      "divxu.w e3, er2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x17, 0xD2},
				Opcode:         opcode.Exts,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.R16,
			expectedString:      "exts.w r2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x17, 0xF6},
				Opcode:         opcode.Exts,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.R32,
			expectedString:      "exts.l er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x17, 0x53},
				Opcode:         opcode.Extu,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x03},
			expectedOperandType: operand.R16,
			expectedString:      "extu.w r3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x17, 0x75},
				Opcode:         opcode.Extu,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.R32,
			expectedString:      "extu.l er5",
		},

		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0A, 0x05},
				Opcode:         opcode.Inc,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R8_INC_DEC,
			expectedString:      "inc.b #1, r5h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0B, 0x55},
				Opcode:         opcode.Inc,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R16_INC_DEC,
			expectedString:      "inc.w #1, r5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0B, 0xD5},
				Opcode:         opcode.Inc,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R16_INC_DEC,
			expectedString:      "inc.w #2, r5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0B, 0x75},
				Opcode:         opcode.Inc,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R32_INC_DEC,
			expectedString:      "inc.l #1, er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0B, 0xF5},
				Opcode:         opcode.Inc,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R32_INC_DEC,
			expectedString:      "inc.l #2, er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x03, 0x04},
				Opcode:         opcode.Ldc,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x04},
			expectedOperandType: operand.R8_LDC,
			expectedString:      "ldc.b r4h, ccr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x03, 0x14},
				Opcode:         opcode.Ldc,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x04},
			expectedOperandType: operand.R8_LDC,
			expectedString:      "ldc.b r4h, exr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0C, 0xD4},
				Opcode:         opcode.Mov,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0D},
			expectedRegDst:      []byte{0x04},
			expectedOperandType: operand.R8_R8,
			expectedString:      "mov.b r5l, r4h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0D, 0x4A},
				Opcode:         opcode.Mov,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x04},
			expectedRegDst:      []byte{0x0A},
			expectedOperandType: operand.R16_R16,
			expectedString:      "mov.w r4, e2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x0F, 0x81},
				Opcode:         opcode.Mov,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x08},
			expectedRegDst:      []byte{0x01},
			expectedOperandType: operand.R32_R32_S2,
			expectedString:      "mov.l er0, er1",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0xC0, 0x50, 0x42},
				Opcode:         opcode.Mulxs,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x04},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.R8_R16_MULXS_DIVXS,
			expectedString:      "mulxs.b r4h, r2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0xC0, 0x52, 0x25},
				Opcode:         opcode.Mulxs,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x02},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.R16_R32_MULXS_DIVXS,
			expectedString:      "mulxs.w r2, er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x50, 0x42},
				Opcode:         opcode.Mulxu,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x04},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.R8_R16,
			expectedString:      "mulxu.b r4h, r2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x52, 0x26},
				Opcode:         opcode.Mulxu,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x02},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.R16_R32,
			expectedString:      "mulxu.w r2, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x17, 0x88},
				Opcode:         opcode.Neg,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x08},
			expectedOperandType: operand.R8,
			expectedString:      "neg.b r0l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x17, 0x9C},
				Opcode:         opcode.Neg,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0C},
			expectedOperandType: operand.R16,
			expectedString:      "neg.w e4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x17, 0xB5},
				Opcode:         opcode.Neg,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.R32,
			expectedString:      "neg.l er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x17, 0x08},
				Opcode:         opcode.Not,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x08},
			expectedOperandType: operand.R8,
			expectedString:      "not.b r0l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x17, 0x1C},
				Opcode:         opcode.Not,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0C},
			expectedOperandType: operand.R16,
			expectedString:      "not.w e4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x17, 0x35},
				Opcode:         opcode.Not,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.R32,
			expectedString:      "not.l er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x14, 0x80},
				Opcode:         opcode.Or,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x08},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.R8_R8,
			expectedString:      "or.b r0l, r0h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x14, 0x80},
				Opcode:         opcode.Or,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x08},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.R16_R16,
			expectedString:      "or.w e0, r0",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0xF0, 0x64, 0x05},
				Opcode:         opcode.Or,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x00},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.R32_R32_S4,
			expectedString:      "or.l er0, er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x12, 0x81},
				Opcode:         opcode.Rotl,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x01},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "rotl.b #1, r1h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x12, 0xC8},
				Opcode:         opcode.Rotl,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x08},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "rotl.b #2, r0l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x12, 0x9C},
				Opcode:         opcode.Rotl,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0C},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "rotl.w #1, e4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x12, 0xD9},
				Opcode:         opcode.Rotl,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x09},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "rotl.w #2, e1",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x12, 0xB6},
				Opcode:         opcode.Rotl,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "rotl.l #1, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x12, 0xF6},
				Opcode:         opcode.Rotl,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "rotl.l #2, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x13, 0x81},
				Opcode:         opcode.Rotr,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x01},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "rotr.b #1, r1h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x13, 0xC8},
				Opcode:         opcode.Rotr,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x08},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "rotr.b #2, r0l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x13, 0x9C},
				Opcode:         opcode.Rotr,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0C},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "rotr.w #1, e4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x13, 0xD9},
				Opcode:         opcode.Rotr,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x09},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "rotr.w #2, e1",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x13, 0xB6},
				Opcode:         opcode.Rotr,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "rotr.l #1, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x13, 0xF6},
				Opcode:         opcode.Rotr,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "rotr.l #2, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x12, 0x03},
				Opcode:         opcode.Rotxl,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x03},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "rotxl.b #1, r3h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x12, 0x43},
				Opcode:         opcode.Rotxl,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x03},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "rotxl.b #2, r3h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x12, 0x16},
				Opcode:         opcode.Rotxl,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "rotxl.w #1, r6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x12, 0x56},
				Opcode:         opcode.Rotxl,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "rotxl.w #2, r6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x12, 0x36},
				Opcode:         opcode.Rotxl,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "rotxl.l #1, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x13, 0x03},
				Opcode:         opcode.Rotxr,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x03},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "rotxr.b #1, r3h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x13, 0x43},
				Opcode:         opcode.Rotxr,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x03},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "rotxr.b #2, r3h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x13, 0x16},
				Opcode:         opcode.Rotxr,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "rotxr.w #1, r6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x13, 0x56},
				Opcode:         opcode.Rotxr,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "rotxr.w #2, r6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x13, 0x36},
				Opcode:         opcode.Rotxr,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "rotxr.l #1, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x10, 0x85},
				Opcode:         opcode.Shal,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "shal.b #1, r5h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x10, 0x9B},
				Opcode:         opcode.Shal,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0B},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "shal.w #1, e3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x10, 0xB3},
				Opcode:         opcode.Shal,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x03},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "shal.l #1, er3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x11, 0x85},
				Opcode:         opcode.Shar,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "shar.b #1, r5h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x11, 0x9B},
				Opcode:         opcode.Shar,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0B},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "shar.w #1, e3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x11, 0xB3},
				Opcode:         opcode.Shar,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x03},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "shar.l #1, er3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x10, 0x0E},
				Opcode:         opcode.Shll,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0E},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "shll.b #1, r6l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x10, 0x45},
				Opcode:         opcode.Shll,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "shll.b #2, r5h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x10, 0x14},
				Opcode:         opcode.Shll,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x04},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "shll.w #1, r4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x10, 0x5E},
				Opcode:         opcode.Shll,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0E},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "shll.w #2, e6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x10, 0x32},
				Opcode:         opcode.Shll,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "shll.l #1, er2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x10, 0x75},
				Opcode:         opcode.Shll,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "shll.l #2, er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x11, 0x0E},
				Opcode:         opcode.Shlr,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0E},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "shlr.b #1, r6l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x11, 0x45},
				Opcode:         opcode.Shlr,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R8_SH,
			expectedString:      "shlr.b #2, r5h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x11, 0x14},
				Opcode:         opcode.Shlr,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x04},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "shlr.w #1, r4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x11, 0x5E},
				Opcode:         opcode.Shlr,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x0E},
			expectedOperandType: operand.Ix_R16_SH,
			expectedString:      "shlr.w #2, e6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x11, 0x32},
				Opcode:         opcode.Shlr,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x02},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "shlr.l #1, er2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x11, 0x75},
				Opcode:         opcode.Shlr,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x05},
			expectedOperandType: operand.Ix_R32_SH,
			expectedString:      "shlr.l #2, er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x02, 0x04},
				Opcode:         opcode.Stc,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x04},
			expectedOperandType: operand.R8_STC,
			expectedString:      "stc.b ccr, r4h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x02, 0x14},
				Opcode:         opcode.Stc,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegDst:      []byte{0x04},
			expectedOperandType: operand.R8_STC,
			expectedString:      "stc.b exr, r4h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x18, 0x3E},
				Opcode:         opcode.Sub,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x03},
			expectedRegDst:      []byte{0x0E},
			expectedOperandType: operand.R8_R8,
			expectedString:      "sub.b r3h, r6l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x19, 0x0B},
				Opcode:         opcode.Sub,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x00},
			expectedRegDst:      []byte{0x0B},
			expectedOperandType: operand.R16_R16,
			expectedString:      "sub.w r0, e3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1A, 0x86},
				Opcode:         opcode.Sub,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x08},
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.R32_R32_S2,
			expectedString:      "sub.l er0, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1A, 0xF7},
				Opcode:         opcode.Sub,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0F},
			expectedRegDst:      []byte{0x07},
			expectedOperandType: operand.R32_R32_S2,
			expectedString:      "sub.l sp, sp",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1B, 0x06},
				Opcode:         opcode.Subs,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      nil,
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_ADDS_SUBS,
			expectedString:      "subs #1, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1B, 0x86},
				Opcode:         opcode.Subs,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      nil,
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_ADDS_SUBS,
			expectedString:      "subs #2, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1B, 0x96},
				Opcode:         opcode.Subs,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      nil,
			expectedRegDst:      []byte{0x06},
			expectedOperandType: operand.Ix_R32_ADDS_SUBS,
			expectedString:      "subs #4, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x1E, 0xC0},
				Opcode:         opcode.Subx,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0C},
			expectedRegDst:      []byte{0x00},
			expectedOperandType: operand.R8_R8,
			expectedString:      "subx r4l, r0h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x57, 0x20},
				Opcode:         opcode.Trapa,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      nil,
			expectedRegDst:      nil,
			expectedOperandType: operand.TRAPA_Ix,
			expectedString:      "trapa #2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x15, 0x4C},
				Opcode:         opcode.Xor,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x04},
			expectedRegDst:      []byte{0x0C},
			expectedOperandType: operand.R8_R8,
			expectedString:      "xor.b r4h, r4l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x65, 0xD4},
				Opcode:         opcode.Xor,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x0D},
			expectedRegDst:      []byte{0x04},
			expectedOperandType: operand.R16_R16,
			expectedString:      "xor.w e5, r4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0xF0, 0x65, 0x01},
				Opcode:         opcode.Xor,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterDirect,
			},
			expectedRegSrc:      []byte{0x00},
			expectedRegDst:      []byte{0x01},
			expectedOperandType: operand.R32_R32_S4,
			expectedString:      "xor.l er0, er1",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0xFF, 0x00},
				Opcode:         opcode.Invalid,
				BWL:            size.Unset,
				AddressingMode: addressingmode.None,
			},
			expectedRegSrc:      nil,
			expectedRegDst:      nil,
			expectedOperandType: operand.Unknown,
			expectedString:      "invalid :(",
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

func TestDetermineOperandTypeAndSetDataRegisterIndirect(t *testing.T) {
	testCases := []struct {
		instruction         instruction.Inst
		expectedOperandType operand.OperandType
		expectedRegSrc      []byte
		expectedRegDst      []byte
		expectedString      string
	}{
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7C, 0x20, 0x76, 0x40},
				Opcode:         opcode.Band,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x02},
			expectedString:      "band #4, @er2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7D, 0x40, 0x72, 0x60},
				Opcode:         opcode.Bclr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "bclr #6, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7D, 0x30, 0x62, 0x40},
				Opcode:         opcode.Bclr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.R8_AR32_BCLR,
			expectedRegSrc:      []byte{0x04},
			expectedRegDst:      []byte{0x03},
			expectedString:      "bclr r4h, @er3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7C, 0x40, 0x76, 0xF0},
				Opcode:         opcode.Biand,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "biand #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7C, 0x40, 0x77, 0xF0},
				Opcode:         opcode.Bild,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "bild #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7C, 0x40, 0x74, 0xF0},
				Opcode:         opcode.Bior,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "bior #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7D, 0x40, 0x67, 0xF0},
				Opcode:         opcode.Bist,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "bist #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7C, 0x40, 0x75, 0xF0},
				Opcode:         opcode.Bixor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "bixor #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7C, 0x40, 0x77, 0x70},
				Opcode:         opcode.Bld,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "bld #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7D, 0x40, 0x71, 0x70},
				Opcode:         opcode.Bnot,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "bnot #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7D, 0x30, 0x61, 0x50},
				Opcode:         opcode.Bnot,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.R8_AR32,
			expectedRegSrc:      []byte{0x05},
			expectedRegDst:      []byte{0x03},
			expectedString:      "bnot r5h, @er3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7C, 0x40, 0x74, 0x70},
				Opcode:         opcode.Bor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "bor #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7D, 0x40, 0x70, 0x70},
				Opcode:         opcode.Bset,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "bset #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7D, 0x30, 0x60, 0x50},
				Opcode:         opcode.Bset,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.R8_AR32,
			expectedRegSrc:      []byte{0x05},
			expectedRegDst:      []byte{0x03},
			expectedString:      "bset r5h, @er3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7D, 0x40, 0x67, 0x70},
				Opcode:         opcode.Bst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "bst #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7C, 0x40, 0x73, 0x70},
				Opcode:         opcode.Btst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "btst #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7C, 0x30, 0x63, 0x50},
				Opcode:         opcode.Btst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.R8_AR32,
			expectedRegSrc:      []byte{0x05},
			expectedRegDst:      []byte{0x03},
			expectedString:      "btst r5h, @er3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7C, 0x40, 0x73, 0x70},
				Opcode:         opcode.Bxor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.Ix_AR32,
			expectedRegDst:      []byte{0x04},
			expectedString:      "bxor #7, @er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x59, 0x60},
				Opcode:         opcode.Jmp,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.AR32_S2,
			expectedRegDst:      []byte{0x06},
			expectedString:      "jmp @er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x5D, 0x60},
				Opcode:         opcode.Jsr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.AR32_S2,
			expectedRegDst:      []byte{0x06},
			expectedString:      "jsr @er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x40, 0x69, 0x20},
				Opcode:         opcode.Ldc,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.R8_LDC, // This one is very strange... Check TODO
			expectedRegSrc:      []byte{0x02},
			expectedString:      "ldc.w @er2, ccr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x41, 0x69, 0x20},
				Opcode:         opcode.Ldc,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.R8_LDC, // This one is very strange... Check TODO
			expectedRegSrc:      []byte{0x02},
			expectedString:      "ldc.w @er2, exr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x68, 0x79},
				Opcode:         opcode.Mov,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.AR32_R8,
			expectedRegSrc:      []byte{0x07},
			expectedRegDst:      []byte{0x09},
			expectedString:      "mov.b @sp, r1l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x69, 0x24},
				Opcode:         opcode.Mov,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.AR32_R16,
			expectedRegSrc:      []byte{0x02},
			expectedRegDst:      []byte{0x04},
			expectedString:      "mov.w @er2, r4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x00, 0x69, 0x23},
				Opcode:         opcode.Mov,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.AR32_R32,
			expectedRegSrc:      []byte{0x02},
			expectedRegDst:      []byte{0x03},
			expectedString:      "mov.l @er2, er3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x68, 0xA0},
				Opcode:         opcode.Mov,
				BWL:            size.Byte,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.R8_AR32,
			expectedRegSrc:      []byte{0x00},
			expectedRegDst:      []byte{0x0A},
			expectedString:      "mov.b r0h, @er2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x69, 0xD8},
				Opcode:         opcode.Mov,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.R16_AR32,
			expectedRegSrc:      []byte{0x08},
			expectedRegDst:      []byte{0x0D},
			expectedString:      "mov.w e0, @er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x00, 0x69, 0x95},
				Opcode:         opcode.Mov,
				BWL:            size.Longword,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.R32_AR32,
			expectedRegSrc:      []byte{0x05},
			expectedRegDst:      []byte{0x09},
			expectedString:      "mov.l er5, @er1",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x40, 0x69, 0xF0},
				Opcode:         opcode.Stc,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.R8_STC,
			expectedRegDst:      []byte{0x0F},
			expectedString:      "stc.w ccr, @sp",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x41, 0x69, 0xD0},
				Opcode:         opcode.Stc,
				BWL:            size.Word,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.R8_STC, // This one is very strange... Check TODO
			expectedRegDst:      []byte{0x0D},
			expectedString:      "stc.w exr, @er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0xE0, 0x7B, 0x4C},
				Opcode:         opcode.Tas,
				BWL:            size.Unset,
				AddressingMode: addressingmode.RegisterIndirect,
			},
			expectedOperandType: operand.S4_R32,
			expectedString:      "tas @er4",
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
