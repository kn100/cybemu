package instruction_test

import (
	"fmt"
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

func TestDetermineOperandTypeAndSetDataRegisterImmediate(t *testing.T) {
	testCases := []struct {
		instruction         instruction.Inst
		expectedOperandType operand.OperandType
		expectedRegSrc      []byte
		expectedRegDst      []byte
		expectedImm         []byte
		expectedString      string
	}{
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x8D, 0x81},
				Opcode:         opcode.Add,
				BWL:            size.Byte,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_R8,
			expectedRegDst:      []byte{0x0D},
			expectedImm:         []byte{0x81},
			expectedString:      "add.b #0x81, r5l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x79, 0x11, 0x30, 0x39},
				Opcode:         opcode.Add,
				BWL:            size.Word,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I16_R16,
			expectedRegDst:      []byte{0x01},
			expectedImm:         []byte{0x30, 0x39},
			expectedString:      "add.w #0x3039, r1",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7A, 0x15, 0x12, 0x34, 0x56, 0x78},
				Opcode:         opcode.Add,
				BWL:            size.Longword,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I32_R32,
			expectedRegDst:      []byte{0x05},
			expectedImm:         []byte{0x12, 0x34, 0x56, 0x78},
			expectedString:      "add.l #0x12345678, er5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x9F, 0x01},
				Opcode:         opcode.Addx,
				BWL:            size.Unset,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_R8,
			expectedRegDst:      []byte{0x0F},
			expectedImm:         []byte{0x01},
			expectedString:      "addx #0x01, r7l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0xE0, 0x7B},
				Opcode:         opcode.And,
				BWL:            size.Byte,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_R8,
			expectedRegDst:      []byte{0x00},
			expectedImm:         []byte{0x7B},
			expectedString:      "and.b #0x7B, r0h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x79, 0x6C, 0x26, 0x94},
				Opcode:         opcode.And,
				BWL:            size.Word,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I16_R16,
			expectedRegDst:      []byte{0x0C},
			expectedImm:         []byte{0x26, 0x94},
			expectedString:      "and.w #0x2694, e4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7A, 0x63, 0x00, 0x0A, 0xBC, 0xDE},
				Opcode:         opcode.And,
				BWL:            size.Longword,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I32_R32,
			expectedRegDst:      []byte{0x03},
			expectedImm:         []byte{0x00, 0x0A, 0xBC, 0xDE},
			expectedString:      "and.l #0x000ABCDE, er3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x06, 0xC0},
				Opcode:         opcode.Andc,
				BWL:            size.Unset,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_CCR,
			expectedImm:         []byte{0xC0},
			expectedString:      "andc #0xC0, ccr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x41, 0x06, 0xA5},
				Opcode:         opcode.Andc,
				BWL:            size.Unset,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_EXR,
			expectedImm:         []byte{0xA5},
			expectedString:      "andc #0xA5, exr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0xA0, 0x00},
				Opcode:         opcode.Cmp,
				BWL:            size.Byte,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_R8,
			expectedRegDst:      []byte{0x00},
			expectedImm:         []byte{0x00},
			expectedString:      "cmp.b #0x00, r0h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x79, 0x2D, 0x1F, 0xFF},
				Opcode:         opcode.Cmp,
				BWL:            size.Word,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I16_R16,
			expectedRegDst:      []byte{0x0D},
			expectedImm:         []byte{0x1F, 0xFF},
			expectedString:      "cmp.w #0x1FFF, e5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7A, 0x24, 0x00, 0x00, 0xFF, 0xFF},
				Opcode:         opcode.Cmp,
				BWL:            size.Longword,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I32_R32,
			expectedRegDst:      []byte{0x04},
			expectedImm:         []byte{0x00, 0x00, 0xFF, 0xFF},
			expectedString:      "cmp.l #0x0000FFFF, er4",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x07, 0xC1},
				Opcode:         opcode.Ldc,
				BWL:            size.Byte,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_CCR,
			expectedImm:         []byte{0xC1},
			expectedString:      "ldc.b #0xC1, ccr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x41, 0x07, 0x6A},
				Opcode:         opcode.Ldc,
				BWL:            size.Byte,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_EXR,
			expectedImm:         []byte{0x6A},
			expectedString:      "ldc.b #0x6A, exr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0xF0, 0x1F},
				Opcode:         opcode.Mov,
				BWL:            size.Byte,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_R8,
			expectedRegDst:      []byte{0x00},
			expectedImm:         []byte{0x1F},
			expectedString:      "mov.b #0x1F, r0h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x79, 0x0B, 0x27, 0x0F},
				Opcode:         opcode.Mov,
				BWL:            size.Word,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I16_R16,
			expectedRegDst:      []byte{0x0B},
			expectedImm:         []byte{0x27, 0x0F},
			expectedString:      "mov.w #0x270F, e3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7A, 0x00, 0x00, 0x00, 0x11, 0xD7},
				Opcode:         opcode.Mov,
				BWL:            size.Longword,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I32_R32,
			expectedRegDst:      []byte{0x00},
			expectedImm:         []byte{0x00, 0x00, 0x11, 0xD7},
			expectedString:      "mov.l #0x000011D7, er0",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0xC1, 0x04},
				Opcode:         opcode.Or,
				BWL:            size.Byte,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_R8,
			expectedRegDst:      []byte{0x01},
			expectedImm:         []byte{0x04},
			expectedString:      "or.b #0x04, r1h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x79, 0x40, 0x00, 0xC0},
				Opcode:         opcode.Or,
				BWL:            size.Word,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I16_R16,
			expectedRegDst:      []byte{0x00},
			expectedImm:         []byte{0x00, 0xC0},
			expectedString:      "or.w #0x00C0, r0",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7A, 0x40, 0x00, 0x00, 0x00, 0xFE},
				Opcode:         opcode.Or,
				BWL:            size.Longword,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I32_R32,
			expectedRegDst:      []byte{0x00},
			expectedImm:         []byte{0x00, 0x00, 0x00, 0xFE},
			expectedString:      "or.l #0x000000FE, er0",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x04, 0x01},
				Opcode:         opcode.Orc,
				BWL:            size.Unset,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_CCR,
			expectedImm:         []byte{0x01},
			expectedString:      "orc #0x01, ccr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x41, 0x04, 0x7B},
				Opcode:         opcode.Orc,
				BWL:            size.Unset,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_EXR,
			expectedImm:         []byte{0x7B},
			expectedString:      "orc #0x7B, exr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x79, 0x3B, 0xFF, 0xF8},
				Opcode:         opcode.Sub,
				BWL:            size.Word,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I16_R16,
			expectedRegDst:      []byte{0x0B},
			expectedImm:         []byte{0xFF, 0xF8},
			expectedString:      "sub.w #0xFFF8, e3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7A, 0x37, 0xFF, 0xFF, 0xFF, 0xF0},
				Opcode:         opcode.Sub,
				BWL:            size.Longword,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I32_R32,
			expectedRegDst:      []byte{0x07},
			expectedImm:         []byte{0xFF, 0xFF, 0xFF, 0xF0},
			expectedString:      "sub.l #0xFFFFFFF0, sp",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0xB5, 0x08},
				Opcode:         opcode.Subx,
				BWL:            size.Unset,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_R8,
			expectedRegDst:      []byte{0x05},
			expectedImm:         []byte{0x08},
			expectedString:      "subx #0x08, r5h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0xD4, 0x80},
				Opcode:         opcode.Xor,
				BWL:            size.Byte,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_R8,
			expectedRegDst:      []byte{0x04},
			expectedImm:         []byte{0x80},
			expectedString:      "xor.b #0x80, r4h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x79, 0x5D, 0x20, 0x00},
				Opcode:         opcode.Xor,
				BWL:            size.Word,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I16_R16,
			expectedRegDst:      []byte{0x0D},
			expectedImm:         []byte{0x20, 0x00},
			expectedString:      "xor.w #0x2000, e5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7A, 0x56, 0x00, 0x00, 0xFF, 0xFF},
				Opcode:         opcode.Xor,
				BWL:            size.Longword,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I32_R32,
			expectedRegDst:      []byte{0x06},
			expectedImm:         []byte{0x00, 0x00, 0xFF, 0xFF},
			expectedString:      "xor.l #0x0000FFFF, er6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x05, 0x40},
				Opcode:         opcode.Xorc,
				BWL:            size.Unset,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_CCR,
			expectedImm:         []byte{0x40},
			expectedString:      "xorc #0x40, ccr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x41, 0x05, 0x9E},
				Opcode:         opcode.Xorc,
				BWL:            size.Unset,
				AddressingMode: addressingmode.Immediate,
			},
			expectedOperandType: operand.I8_EXR,
			expectedImm:         []byte{0x9E},
			expectedString:      "xorc #0x9E, exr",
		},
	}
	for _, tc := range testCases {
		tc.instruction.DetermineOperandTypeAndSetData()
		assert.Equal(t, tc.expectedOperandType, tc.instruction.OperandType, "expected operand type to be %s, got %s", tc.expectedOperandType, tc.instruction.OperandType)
		assert.Equal(t, tc.expectedImm, tc.instruction.Imm, "expected Imm %v, got %v", tc.expectedImm, tc.instruction.Imm)
		assert.Equal(t, tc.expectedRegDst, tc.instruction.RegDst, "expected RegDst %v, got %v", tc.expectedRegDst, tc.instruction.RegDst)
		assert.Equal(t, tc.expectedString, tc.instruction.String())

	}
}

func TestDetermineOperandTypeAndSetDataRegisterAbsoluteAddress(t *testing.T) {
	testCases := []struct {
		instruction         instruction.Inst
		expectedOperandType operand.OperandType
		expectedRegSrc      []byte
		expectedRegDst      []byte
		expectedImm         []byte
		expectedImmL        []byte
		expectedImmR        []byte
		expectedString      string
	}{
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7E, 0xC0, 0x76, 0x50},
				Opcode:         opcode.Band,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x05},
			expectedString:      "band #5, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x10, 0x01, 0x23, 0x76, 0x30},
				Opcode:         opcode.Band,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x01, 0x23},
			expectedImmR:        []byte{0x03},
			expectedString:      "band #3, @0x0123:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x76, 0x50},
				Opcode:         opcode.Band,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x12, 0x34, 0x56, 0x78},
			expectedImmR:        []byte{0x05},
			expectedString:      "band #5, @0x12345678:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7F, 0xC0, 0x72, 0x10},
				Opcode:         opcode.Bclr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x01},
			expectedString:      "bclr #1, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x18, 0x01, 0x23, 0x72, 0x30},
				Opcode:         opcode.Bclr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x01, 0x23},
			expectedImmR:        []byte{0x03},
			expectedString:      "bclr #3, @0x0123:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x72, 0x50},
				Opcode:         opcode.Bclr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x12, 0x34, 0x56, 0x78},
			expectedImmR:        []byte{0x05},
			expectedString:      "bclr #5, @0x12345678:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7F, 0xC0, 0x62, 0x50},
				Opcode:         opcode.Bclr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI8_BCLR,
			expectedRegDst:      []byte{0x05},
			expectedImm:         []byte{0xC0},
			expectedString:      "bclr r5h, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x18, 0x01, 0x23, 0x62, 0x70},
				Opcode:         opcode.Bclr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI16_S6,
			expectedRegDst:      []byte{0x07},
			expectedImm:         []byte{0x01, 0x23},
			expectedString:      "bclr r7h, @0x0123:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x62, 0xE0},
				Opcode:         opcode.Bclr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI32_BCLR,
			expectedRegDst:      []byte{0x0E},
			expectedImm:         []byte{0x12, 0x34, 0x56, 0x78},
			expectedString:      "bclr r6l, @0x12345678:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7E, 0xC0, 0x76, 0x80},
				Opcode:         opcode.Biand,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x08},
			expectedString:      "biand #0, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x10, 0x01, 0x23, 0x76, 0xB0},
				Opcode:         opcode.Biand,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x01, 0x23},
			expectedImmR:        []byte{0x0B},
			expectedString:      "biand #3, @0x0123:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x30, 0x12, 0x34, 0x67, 0x89, 0x76, 0xD0},
				Opcode:         opcode.Biand,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x12, 0x34, 0x67, 0x89},
			expectedImmR:        []byte{0x0D},
			expectedString:      "biand #5, @0x12346789:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7E, 0xC0, 0x77, 0x80},
				Opcode:         opcode.Bild,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x08},
			expectedString:      "bild #0, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x10, 0x01, 0x23, 0x77, 0xB0},
				Opcode:         opcode.Bild,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x01, 0x23},
			expectedImmR:        []byte{0x0B},
			expectedString:      "bild #3, @0x0123:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x77, 0xD0},
				Opcode:         opcode.Bild,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x12, 0x34, 0x56, 0x78},
			expectedImmR:        []byte{0x0D},
			expectedString:      "bild #5, @0x12345678:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7E, 0xC0, 0x74, 0x80},
				Opcode:         opcode.Bior,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x08},
			expectedString:      "bior #0, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x10, 0x00, 0x00, 0x74, 0xF0},
				Opcode:         opcode.Bior,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bior #7, @0x0000:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x74, 0xF0},
				Opcode:         opcode.Bior,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x00, 0x00, 0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bior #7, @0x00000000:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7F, 0xC0, 0x67, 0x80},
				Opcode:         opcode.Bist,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x08},
			expectedString:      "bist #0, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x18, 0x00, 0x00, 0x67, 0xF0},
				Opcode:         opcode.Bist,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bist #7, @0x0000:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x38, 0x00, 0x00, 0x00, 0x00, 0x67, 0xF0},
				Opcode:         opcode.Bist,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x00, 0x00, 0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bist #7, @0x00000000:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7E, 0xC0, 0x75, 0x80},
				Opcode:         opcode.Bixor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x08},
			expectedString:      "bixor #0, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x30, 0x00, 0x00, 0x75, 0xF0},
				Opcode:         opcode.Bixor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bixor #7, @0x0000:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x38, 0x00, 0x00, 0x00, 0x00, 0x75, 0xF0},
				Opcode:         opcode.Bixor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x00, 0x00, 0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bixor #7, @0x00000000:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7E, 0xC0, 0x77, 0x80},
				Opcode:         opcode.Bld,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x08},
			expectedString:      "bld #0, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x30, 0x00, 0x00, 0x77, 0xF0},
				Opcode:         opcode.Bld,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bld #7, @0x0000:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x38, 0x00, 0x00, 0x00, 0x00, 0x77, 0xF0},
				Opcode:         opcode.Bld,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x00, 0x00, 0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bld #7, @0x00000000:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7F, 0xC0, 0x71, 0x10},
				Opcode:         opcode.Bnot,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x01},
			expectedString:      "bnot #1, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x18, 0x01, 0x23, 0x71, 0x30},
				Opcode:         opcode.Bnot,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x01, 0x23},
			expectedImmR:        []byte{0x03},
			expectedString:      "bnot #3, @0x0123:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x71, 0x50},
				Opcode:         opcode.Bnot,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x12, 0x34, 0x56, 0x78},
			expectedImmR:        []byte{0x05},
			expectedString:      "bnot #5, @0x12345678:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7F, 0xC0, 0x61, 0x50},
				Opcode:         opcode.Bnot,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI8_BCLR,
			expectedRegDst:      []byte{0x05},
			expectedImm:         []byte{0xC0},
			expectedString:      "bnot r5h, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x18, 0x01, 0x23, 0x61, 0x70},
				Opcode:         opcode.Bnot,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI16_S6,
			expectedRegDst:      []byte{0x07},
			expectedImm:         []byte{0x01, 0x23},
			expectedString:      "bnot r7h, @0x0123:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x61, 0xE0},
				Opcode:         opcode.Bnot,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			// TODO: Terrible operand type name.
			expectedOperandType: operand.R8_AI32_BCLR,
			expectedRegDst:      []byte{0x0E},
			expectedImm:         []byte{0x12, 0x34, 0x56, 0x78},
			expectedString:      "bnot r6l, @0x12345678:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7E, 0xC0, 0x74, 0x80},
				Opcode:         opcode.Bor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x08},
			expectedString:      "bor #0, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x10, 0x00, 0x00, 0x74, 0xF0},
				Opcode:         opcode.Bor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bor #7, @0x0000:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x74, 0xF0},
				Opcode:         opcode.Bor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x00, 0x00, 0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bor #7, @0x00000000:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7F, 0xC0, 0x70, 0x10},
				Opcode:         opcode.Bset,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x01},
			expectedString:      "bset #1, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x18, 0x01, 0x23, 0x70, 0x30},
				Opcode:         opcode.Bset,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x01, 0x23},
			expectedImmR:        []byte{0x03},
			expectedString:      "bset #3, @0x0123:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x70, 0x50},
				Opcode:         opcode.Bset,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x12, 0x34, 0x56, 0x78},
			expectedImmR:        []byte{0x05},
			expectedString:      "bset #5, @0x12345678:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7F, 0xC0, 0x60, 0x50},
				Opcode:         opcode.Bset,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI8_BCLR,
			expectedRegDst:      []byte{0x05},
			expectedImm:         []byte{0xC0},
			expectedString:      "bset r5h, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x18, 0x01, 0x23, 0x60, 0x70},
				Opcode:         opcode.Bset,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI16_S6,
			expectedRegDst:      []byte{0x07},
			expectedImm:         []byte{0x01, 0x23},
			expectedString:      "bset r7h, @0x0123:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x60, 0xE0},
				Opcode:         opcode.Bset,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			// TODO: Terrible operand type name.
			expectedOperandType: operand.R8_AI32_BCLR,
			expectedRegDst:      []byte{0x0E},
			expectedImm:         []byte{0x12, 0x34, 0x56, 0x78},
			expectedString:      "bset r6l, @0x12345678:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7F, 0xC0, 0x67, 0x80},
				Opcode:         opcode.Bst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x08},
			expectedString:      "bst #0, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x18, 0x00, 0x00, 0x67, 0xF0},
				Opcode:         opcode.Bst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bst #7, @0x0000:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x38, 0x00, 0x00, 0x00, 0x00, 0x67, 0xF0},
				Opcode:         opcode.Bst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x00, 0x00, 0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bst #7, @0x00000000:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7E, 0xC0, 0x73, 0x10},
				Opcode:         opcode.Btst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x01},
			expectedString:      "btst #1, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x10, 0x01, 0x23, 0x73, 0x30},
				Opcode:         opcode.Btst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x01, 0x23},
			expectedImmR:        []byte{0x03},
			expectedString:      "btst #3, @0x0123:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x73, 0x50},
				Opcode:         opcode.Btst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x12, 0x34, 0x56, 0x78},
			expectedImmR:        []byte{0x05},
			expectedString:      "btst #5, @0x12345678:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7E, 0xC0, 0x63, 0x50},
				Opcode:         opcode.Btst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI8_BCLR,
			expectedRegDst:      []byte{0x05},
			expectedImm:         []byte{0xC0},
			expectedString:      "btst r5h, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x10, 0x01, 0x23, 0x63, 0x70},
				Opcode:         opcode.Btst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI16_S6,
			expectedRegDst:      []byte{0x07},
			expectedImm:         []byte{0x01, 0x23},
			expectedString:      "btst r7h, @0x0123:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x63, 0xE0},
				Opcode:         opcode.Btst,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			// TODO: Terrible operand type name.
			expectedOperandType: operand.R8_AI32_BCLR,
			expectedRegDst:      []byte{0x0E},
			expectedImm:         []byte{0x12, 0x34, 0x56, 0x78},
			expectedString:      "btst r6l, @0x12345678:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x7E, 0xC0, 0x75, 0x80},
				Opcode:         opcode.Bxor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI8,
			expectedImmL:        []byte{0xC0},
			expectedImmR:        []byte{0x08},
			expectedString:      "bxor #0, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x10, 0x00, 0x00, 0x75, 0xF0},
				Opcode:         opcode.Bxor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI16,
			expectedImmL:        []byte{0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bxor #7, @0x0000:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x30, 0x00, 0x00, 0x00, 0x00, 0x75, 0xF0},
				Opcode:         opcode.Bxor,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.Ix_AI32,
			expectedImmL:        []byte{0x00, 0x00, 0x00, 0x00},
			expectedImmR:        []byte{0x0F},
			expectedString:      "bxor #7, @0x00000000:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x5A, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Jmp,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.I24,
			expectedImm:         []byte{0x12, 0x89, 0xDE},
			expectedString:      "jmp @0x1289DE:24",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x5E, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Jsr,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.I24,
			expectedImm:         []byte{0x12, 0x89, 0xDE},
			expectedString:      "jsr @0x1289DE:24",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x40, 0x6B, 0x00, 0x01, 0x26},
				Opcode:         opcode.Ldc,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI16_CCR,
			expectedImm:         []byte{0x01, 0x26},
			expectedString:      "ldc.w @0x0126:16, ccr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x40, 0x6B, 0x20, 0x00, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Ldc,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI32_CCR,
			expectedImm:         []byte{0x00, 0x12, 0x89, 0xDE},
			expectedString:      "ldc.w @0x001289DE:32, ccr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x41, 0x6B, 0x00, 0x01, 0x26},
				Opcode:         opcode.Ldc,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI16_CCR,
			expectedImm:         []byte{0x01, 0x26},
			expectedString:      "ldc.w @0x0126:16, exr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x41, 0x6B, 0x20, 0x00, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Ldc,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI32_CCR,
			expectedImm:         []byte{0x00, 0x12, 0x89, 0xDE},
			expectedString:      "ldc.w @0x001289DE:32, exr",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x20, 0xFF},
				Opcode:         opcode.Mov,
				BWL:            size.Byte,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI8_R8,
			expectedImm:         []byte{0xFF},
			expectedRegDst:      []byte{0x00},
			expectedString:      "mov.b @0xFF:8, r0h",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x0C, 0x01, 0x26},
				Opcode:         opcode.Mov,
				BWL:            size.Byte,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI16_R8,
			expectedImm:         []byte{0x01, 0x26},
			expectedRegDst:      []byte{0x0C},
			expectedString:      "mov.b @0x0126:16, r4l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x2A, 0x00, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Mov,
				BWL:            size.Byte,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI32_R8,
			expectedImm:         []byte{0x00, 0x12, 0x89, 0xDE},
			expectedRegDst:      []byte{0x0A},
			expectedString:      "mov.b @0x001289DE:32, r2l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6B, 0x0E, 0x01, 0x26},
				Opcode:         opcode.Mov,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI16_R16,
			expectedImm:         []byte{0x01, 0x26},
			expectedRegDst:      []byte{0x0E},
			expectedString:      "mov.w @0x0126:16, e6",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6B, 0x25, 0x00, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Mov,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI32_R16,
			expectedImm:         []byte{0x00, 0x12, 0x89, 0xDE},
			expectedRegDst:      []byte{0x05},
			expectedString:      "mov.w @0x001289DE:32, r5",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x00, 0x6B, 0x03, 0x01, 0x26},
				Opcode:         opcode.Mov,
				BWL:            size.Longword,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI16_R32,
			expectedImm:         []byte{0x01, 0x26},
			expectedRegDst:      []byte{0x03},
			expectedString:      "mov.l @0x0126:16, er3",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x00, 0x6B, 0x22, 0x00, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Mov,
				BWL:            size.Longword,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI32_R32,
			expectedImm:         []byte{0x00, 0x12, 0x89, 0xDE},
			expectedRegDst:      []byte{0x02},
			expectedString:      "mov.l @0x001289DE:32, er2",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x31, 0xC0},
				Opcode:         opcode.Mov,
				BWL:            size.Byte,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI8,
			expectedImm:         []byte{0xC0},
			expectedRegSrc:      []byte{0x03},
			expectedString:      "mov.b r1h, @0xC0:8",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x89, 0x01, 0x26},
				Opcode:         opcode.Mov,
				BWL:            size.Byte,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI16,
			expectedImm:         []byte{0x01, 0x26},
			expectedRegSrc:      []byte{0x09},
			expectedString:      "mov.b r1l, @0x0126:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0xA2, 0x00, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Mov,
				BWL:            size.Byte,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI32,
			expectedImm:         []byte{0x00, 0x12, 0x89, 0xDE},
			expectedRegSrc:      []byte{0x02},
			expectedString:      "mov.b r2h, @0x001289DE:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6B, 0x88, 0x01, 0x26},
				Opcode:         opcode.Mov,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R16_AI16,
			expectedImm:         []byte{0x01, 0x26},
			expectedRegSrc:      []byte{0x08},
			expectedString:      "mov.w e0, @0x0126:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6B, 0xAC, 0x00, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Mov,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R16_AI32,
			expectedImm:         []byte{0x00, 0x12, 0x89, 0xDE},
			expectedRegSrc:      []byte{0x0C},
			expectedString:      "mov.w e4, @0x001289DE:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x00, 0x6B, 0x81, 0x01, 0x26},
				Opcode:         opcode.Mov,
				BWL:            size.Longword,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R32_AI16,
			expectedImm:         []byte{0x01, 0x26},
			expectedRegSrc:      []byte{0x01},
			expectedString:      "mov.l er1, @0x0126:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x00, 0x6B, 0xA2, 0x00, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Mov,
				BWL:            size.Longword,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R32_AI32,
			expectedImm:         []byte{0x00, 0x12, 0x89, 0xDE},
			expectedRegSrc:      []byte{0x02},
			expectedString:      "mov.l er2, @0x001289DE:32",
		},
		//
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0x4D, 0xFF, 0xC0},
				Opcode:         opcode.Movfpe,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.AI16_R8,
			expectedImm:         []byte{0xFF, 0xC0},
			expectedRegSrc:      []byte{0x0D},
			expectedString:      "movfpe @0xFFC0:16, r5l",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x6A, 0xC5, 0xFF, 0xC0},
				Opcode:         opcode.Movtpe,
				BWL:            size.Unset,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.R8_AI16,
			expectedImm:         []byte{0xFF, 0xC0},
			expectedRegSrc:      []byte{0x05},
			expectedString:      "movtpe r5h, @0xFFC0:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x40, 0x6B, 0x80, 0x01, 0x26},
				Opcode:         opcode.Stc,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.CCR_AI16,
			expectedImm:         []byte{0x01, 0x26},
			expectedString:      "stc.w ccr, @0x0126:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x40, 0x6B, 0xA0, 0x00, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Stc,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.CCR_AI32,
			expectedImm:         []byte{0x00, 0x12, 0x89, 0xDE},
			expectedString:      "stc.w ccr, @0x001289DE:32",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x41, 0x6B, 0x80, 0x01, 0x26},
				Opcode:         opcode.Stc,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.CCR_AI16,
			expectedImm:         []byte{0x01, 0x26},
			expectedString:      "stc.w exr, @0x0126:16",
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x01, 0x41, 0x6B, 0xA0, 0x00, 0x12, 0x89, 0xDE},
				Opcode:         opcode.Stc,
				BWL:            size.Word,
				AddressingMode: addressingmode.AbsoluteAddress,
			},
			expectedOperandType: operand.CCR_AI32,
			expectedImm:         []byte{0x00, 0x12, 0x89, 0xDE},
			expectedString:      "stc.w exr, @0x001289DE:32",
		},
	}
	for _, tc := range testCases {
		tc.instruction.DetermineOperandTypeAndSetData()
		fmt.Printf("%+v", tc.instruction)
		assert.Equal(t, tc.expectedOperandType, tc.instruction.OperandType, "expected operand type to be %s, got %s", tc.expectedOperandType, tc.instruction.OperandType)
		assert.Equal(t, tc.expectedImmL, tc.instruction.ImmL, "expected ImmL %v, got %v", tc.expectedImmL, tc.instruction.ImmL)
		assert.Equal(t, tc.expectedImmR, tc.instruction.ImmR, "expected ImmR %v, got %v", tc.expectedImmR, tc.instruction.ImmR)
		assert.Equal(t, tc.expectedImm, tc.instruction.Imm, "expected Imm %v, got %v", tc.expectedImm, tc.instruction.Imm)
		assert.Equal(t, tc.expectedString, tc.instruction.String())
	}
}

func TestDetermineOperandAndTypeAndSetDataRegisterProgramCounterRelative(t *testing.T) {
	testCases := []struct {
		instruction         instruction.Inst
		expectedOperandType operand.OperandType
		expectedImm         []byte
		expectedString      string
	}{
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x40, 0xFE},
				Opcode:         opcode.Bra,
				BWL:            size.Unset,
				AddressingMode: addressingmode.ProgramCounterRelative,
				Pos:            60, //0x3C
			},
			expectedOperandType: operand.O8,
			expectedString:      "bra 0x0000003C:8",
			expectedImm:         []byte{0xFE},
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x41, 0xFC},
				Opcode:         opcode.Brn,
				BWL:            size.Unset,
				AddressingMode: addressingmode.ProgramCounterRelative,
				Pos:            62, //0x3E
			},
			expectedOperandType: operand.O8,
			expectedString:      "brn 0x0000003C:8",
			expectedImm:         []byte{0xFC},
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x42, 0xFA},
				Opcode:         opcode.Bhi,
				BWL:            size.Unset,
				AddressingMode: addressingmode.ProgramCounterRelative,
				Pos:            64, //0x40
			},
			expectedOperandType: operand.O8,
			expectedString:      "bhi 0x0000003C:8",
			expectedImm:         []byte{0xFA},
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x58, 0x00, 0x00, 0x3C},
				Opcode:         opcode.Bra,
				BWL:            size.Unset,
				AddressingMode: addressingmode.ProgramCounterRelative,
				Pos:            92, //0x5C
			},
			expectedOperandType: operand.O16,
			expectedString:      "bra 0x0000009C:16",
			expectedImm:         []byte{0x00, 0x3C},
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x58, 0x10, 0x00, 0x38},
				Opcode:         opcode.Brn,
				BWL:            size.Unset,
				AddressingMode: addressingmode.ProgramCounterRelative,
				Pos:            96, //0x60
			},
			expectedOperandType: operand.O16,
			expectedString:      "brn 0x0000009C:16",
			expectedImm:         []byte{0x00, 0x38},
		},
		{
			instruction: instruction.Inst{
				Bytes:          []byte{0x58, 0x20, 0x00, 0x34},
				Opcode:         opcode.Bhi,
				BWL:            size.Unset,
				AddressingMode: addressingmode.ProgramCounterRelative,
				Pos:            100, //0x64
			},
			expectedOperandType: operand.O16,
			expectedString:      "bhi 0x0000009C:16",
			expectedImm:         []byte{0x00, 0x34},
		},
	}
	for _, tc := range testCases {
		tc.instruction.DetermineOperandTypeAndSetData()
		fmt.Printf("%+v", tc.instruction)
		assert.Equal(t, tc.expectedOperandType, tc.instruction.OperandType, "expected operand type to be %s, got %s", tc.expectedOperandType, tc.instruction.OperandType)
		assert.Equal(t, tc.expectedImm, tc.instruction.Imm, "expected Imm %v, got %v", tc.expectedImm, tc.instruction.Imm)
		assert.Equal(t, tc.expectedString, tc.instruction.String())
	}
}
