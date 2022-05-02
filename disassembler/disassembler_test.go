package disassembler_test

import (
	"fmt"
	"testing"

	"github.com/kn100/cybemu/disassembler"
	"github.com/stretchr/testify/assert"
)

func TestDisassembleTimsTestCases(t *testing.T) {
	testCases := []struct {
		Input                  []byte
		ExpectedOpcode         string
		expectedBWL            disassembler.Size
		expectedAddressingMode disassembler.AddressingMode
	}{
		{
			Input:          []byte{0x8D, 0x81},
			ExpectedOpcode: "add",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x08, 0x3E},
			ExpectedOpcode: "add",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x79, 0x11, 0x30, 0x39},
			ExpectedOpcode: "add",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x09, 0x0B},
			ExpectedOpcode: "add",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x7A, 0x15, 0x12, 0x34, 0x56, 0x78},
			ExpectedOpcode: "add",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x0A, 0x87},
			ExpectedOpcode: "add",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x0B, 0x00},
			ExpectedOpcode: "adds",
		},
		{
			Input:          []byte{0x0B, 0x81},
			ExpectedOpcode: "adds",
		},
		{
			Input:          []byte{0x0B, 0x92},
			ExpectedOpcode: "adds",
		},
		{
			Input:          []byte{0x95, 0x0A},
			ExpectedOpcode: "addx",
		},
		{
			Input:          []byte{0x0E, 0xC0},
			ExpectedOpcode: "addx",
		},
		{
			Input:          []byte{0xE2, 0x2D},
			ExpectedOpcode: "and",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x16, 0x32},
			ExpectedOpcode: "and",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x79, 0x6C, 0x26, 0x94},
			ExpectedOpcode: "and",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x66, 0x2E},
			ExpectedOpcode: "and",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x7A, 0x63, 0x00, 0x0A, 0xBC, 0xDE},
			ExpectedOpcode: "and",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0xF0, 0x66, 0x54},
			ExpectedOpcode: "and",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x06, 0xC0},
			ExpectedOpcode: "andc",
		},
		{
			Input:          []byte{0x76, 0x25},
			ExpectedOpcode: "band",
		},
		{
			Input:          []byte{0x7C, 0x20, 0x76, 0x40},
			ExpectedOpcode: "band",
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x76, 0x50},
			ExpectedOpcode: "band",
		},
		{
			Input:          []byte{0x40, 0xFE},
			ExpectedOpcode: "bra",
		},
		{
			Input:          []byte{0x41, 0xFC},
			ExpectedOpcode: "brn",
		},
		{
			Input:          []byte{0x42, 0xFA},
			ExpectedOpcode: "bhi",
		},
		{
			Input:          []byte{0x43, 0xF8},
			ExpectedOpcode: "bls",
		},
		{
			Input:          []byte{0x44, 0xF6},
			ExpectedOpcode: "bcc",
		},
		{
			Input:          []byte{0x45, 0xF4},
			ExpectedOpcode: "bcs",
		},
		{
			Input:          []byte{0x46, 0xF2},
			ExpectedOpcode: "bne",
		},
		{
			Input:          []byte{0x47, 0xF0},
			ExpectedOpcode: "beq",
		},
		{
			Input:          []byte{0x48, 0xEE},
			ExpectedOpcode: "bvc",
		},
		{
			Input:          []byte{0x49, 0xEC},
			ExpectedOpcode: "bvs",
		},
		{
			Input:          []byte{0x4A, 0xEA},
			ExpectedOpcode: "bpl",
		},
		{
			Input:          []byte{0x4B, 0xE8},
			ExpectedOpcode: "bmi",
		},
		{
			Input:          []byte{0x4C, 0xE6},
			ExpectedOpcode: "bge",
		},
		{
			Input:          []byte{0x4D, 0xE4},
			ExpectedOpcode: "blt",
		},
		{
			Input:          []byte{0x4E, 0xE2},
			ExpectedOpcode: "bgt",
		},
		{
			Input:          []byte{0x4F, 0xE0},
			ExpectedOpcode: "ble",
		},
		{
			Input:          []byte{0x58, 0x00, 0x00, 0x3C},
			ExpectedOpcode: "bra",
		},
		{
			Input:          []byte{0x58, 0x10, 0x00, 0x38},
			ExpectedOpcode: "brn",
		},
		{
			Input:          []byte{0x58, 0x20, 0x00, 0x34},
			ExpectedOpcode: "bhi",
		},
		{
			Input:          []byte{0x58, 0x30, 0x00, 0x30},
			ExpectedOpcode: "bls",
		},
		{
			Input:          []byte{0x58, 0x40, 0x00, 0x2C},
			ExpectedOpcode: "bcc",
		},
		{
			Input:          []byte{0x58, 0x50, 0x00, 0x28},
			ExpectedOpcode: "bcs",
		},
		{
			Input:          []byte{0x58, 0x60, 0x00, 0x24},
			ExpectedOpcode: "bne",
		},
		{
			Input:          []byte{0x58, 0x70, 0x00, 0x20},
			ExpectedOpcode: "beq",
		},
		{
			Input:          []byte{0x58, 0x80, 0x00, 0x1C},
			ExpectedOpcode: "bvc",
		},
		{
			Input:          []byte{0x58, 0x90, 0x00, 0x18},
			ExpectedOpcode: "bvs",
		},
		{
			Input:          []byte{0x58, 0xA0, 0x00, 0x14},
			ExpectedOpcode: "bpl",
		},
		{
			Input:          []byte{0x58, 0xB0, 0x00, 0x10},
			ExpectedOpcode: "bmi",
		},
		{
			Input:          []byte{0x58, 0xC0, 0x00, 0x0C},
			ExpectedOpcode: "bge",
		},
		{
			Input:          []byte{0x58, 0xD0, 0x00, 0x08},
			ExpectedOpcode: "blt",
		},
		{
			Input:          []byte{0x58, 0xE0, 0x00, 0x04},
			ExpectedOpcode: "bgt",
		},
		{
			Input:          []byte{0x58, 0xF0, 0x00, 0x00},
			ExpectedOpcode: "ble",
		},
		{
			Input:          []byte{0x72, 0x4A},
			ExpectedOpcode: "bclr",
		},
		{
			Input:          []byte{0x7D, 0x40, 0x72, 0x60},
			ExpectedOpcode: "bclr",
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x72, 0x10},
			ExpectedOpcode: "bclr",
		},
		{
			Input:          []byte{0x62, 0x93},
			ExpectedOpcode: "bclr",
		},
		{
			Input:          []byte{0x7D, 0x30, 0x62, 0x40},
			ExpectedOpcode: "bclr",
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x62, 0x50},
			ExpectedOpcode: "bclr",
		},
		{
			Input:          []byte{0x76, 0xC2},
			ExpectedOpcode: "biand",
		},
		{
			Input:          []byte{0x7C, 0x40, 0x76, 0xF0},
			ExpectedOpcode: "biand",
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x76, 0x80},
			ExpectedOpcode: "biand",
		},
		{
			Input:          []byte{0x77, 0xC2},
			ExpectedOpcode: "bild",
		},
		{
			Input:          []byte{0x7C, 0x40, 0x77, 0xF0},
			ExpectedOpcode: "bild",
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x77, 0x80},
			ExpectedOpcode: "bild",
		},
		{
			Input:          []byte{0x74, 0xC2},
			ExpectedOpcode: "bior",
		},
		{
			Input:          []byte{0x7C, 0x40, 0x74, 0xF0},
			ExpectedOpcode: "bior",
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x74, 0x80},
			ExpectedOpcode: "bior",
		},
		{
			Input:          []byte{0x67, 0xC2},
			ExpectedOpcode: "bist",
		},
		{
			Input:          []byte{0x7D, 0x40, 0x67, 0xF0},
			ExpectedOpcode: "bist",
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x67, 0x80},
			ExpectedOpcode: "bist",
		},
		{
			Input:          []byte{0x75, 0xC2},
			ExpectedOpcode: "bixor",
		},
		{
			Input:          []byte{0x7C, 0x40, 0x75, 0xF0},
			ExpectedOpcode: "bixor",
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x75, 0x80},
			ExpectedOpcode: "bixor",
		},
		{
			Input:          []byte{0x77, 0x42},
			ExpectedOpcode: "bld",
		},
		{
			Input:          []byte{0x7C, 0x40, 0x77, 0x70},
			ExpectedOpcode: "bld",
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x77, 0x00},
			ExpectedOpcode: "bld",
		},
		{
			Input:          []byte{0x71, 0x42},
			ExpectedOpcode: "bnot",
		},
		{
			Input:          []byte{0x7D, 0x40, 0x71, 0x70},
			ExpectedOpcode: "bnot",
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x71, 0x00},
			ExpectedOpcode: "bnot",
		},
		{
			Input:          []byte{0x61, 0x81},
			ExpectedOpcode: "bnot",
		},
		{
			Input:          []byte{0x7D, 0x30, 0x61, 0x50},
			ExpectedOpcode: "bnot",
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x61, 0xE0},
			ExpectedOpcode: "bnot",
		},
		{
			Input:          []byte{0x74, 0x42},
			ExpectedOpcode: "bor",
		},
		{
			Input:          []byte{0x7C, 0x40, 0x74, 0x70},
			ExpectedOpcode: "bor",
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x74, 0x00},
			ExpectedOpcode: "bor",
		},
		{
			Input:          []byte{0x70, 0x42},
			ExpectedOpcode: "bset",
		},
		{
			Input:          []byte{0x7D, 0x40, 0x70, 0x70},
			ExpectedOpcode: "bset",
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x70, 0x00},
			ExpectedOpcode: "bset",
		},
		{
			Input:          []byte{0x60, 0x81},
			ExpectedOpcode: "bset",
		},
		{
			Input:          []byte{0x7D, 0x30, 0x60, 0x50},
			ExpectedOpcode: "bset",
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x60, 0xE0},
			ExpectedOpcode: "bset",
		},
		{
			Input:          []byte{0x5C, 0x00, 0xFF, 0x7A},
			ExpectedOpcode: "bsr",
		},
		{
			Input:          []byte{0x5C, 0x00, 0x00, 0x00},
			ExpectedOpcode: "bsr",
		},
		{
			Input:          []byte{0x67, 0x42},
			ExpectedOpcode: "bst",
		},
		{
			Input:          []byte{0x7D, 0x40, 0x67, 0x70},
			ExpectedOpcode: "bst",
		},
		{
			Input:          []byte{0x7F, 0xC0, 0x67, 0x00},
			ExpectedOpcode: "bst",
		},
		{
			Input:          []byte{0x73, 0x42},
			ExpectedOpcode: "btst",
		},
		{
			Input:          []byte{0x7C, 0x40, 0x73, 0x70},
			ExpectedOpcode: "btst",
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x73, 0x00},
			ExpectedOpcode: "btst",
		},
		{
			Input:          []byte{0x63, 0x81},
			ExpectedOpcode: "btst",
		},
		{
			Input:          []byte{0x7C, 0x30, 0x63, 0x50},
			ExpectedOpcode: "btst",
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x63, 0xE0},
			ExpectedOpcode: "btst",
		},
		{
			Input:          []byte{0x75, 0x42},
			ExpectedOpcode: "bxor",
		},
		{
			Input:          []byte{0x7C, 0x40, 0x75, 0x70},
			ExpectedOpcode: "bxor",
		},
		{
			Input:          []byte{0x7E, 0xC0, 0x75, 0x00},
			ExpectedOpcode: "bxor",
		},
		{
			Input:          []byte{0xA5, 0x8F},
			ExpectedOpcode: "cmp",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x1C, 0x4A},
			ExpectedOpcode: "cmp",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x79, 0x2D, 0x1F, 0xFF},
			ExpectedOpcode: "cmp",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x1D, 0xD2},
			ExpectedOpcode: "cmp",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x7A, 0x24, 0x00, 0x00, 0xFF, 0xFF},
			ExpectedOpcode: "cmp",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x1F, 0xB5},
			ExpectedOpcode: "cmp",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x0F, 0x0C},
			ExpectedOpcode: "daa",
		},
		{
			Input:          []byte{0x1F, 0x05},
			ExpectedOpcode: "das",
		},
		{
			Input:          []byte{0x1A, 0x05},
			ExpectedOpcode: "dec",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x1B, 0x54},
			ExpectedOpcode: "dec",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x1B, 0xDB},
			ExpectedOpcode: "dec",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x1B, 0x72},
			ExpectedOpcode: "dec",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x1B, 0xF3},
			ExpectedOpcode: "dec",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0xD0, 0x51, 0xBC},
			ExpectedOpcode: "divxs",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x01, 0xD0, 0x53, 0xB2},
			ExpectedOpcode: "divxs",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x51, 0xBC},
			ExpectedOpcode: "divxu",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x53, 0xB2},
			ExpectedOpcode: "divxu",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x7B, 0x5C, 0x59, 0x8F},
			ExpectedOpcode: "eepmov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x7B, 0xD4, 0x59, 0x8F},
			ExpectedOpcode: "eepmov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x17, 0xD2},
			ExpectedOpcode: "exts",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x17, 0xF6},
			ExpectedOpcode: "exts",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x17, 0x53},
			ExpectedOpcode: "extu",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x17, 0x75},
			ExpectedOpcode: "extu",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x0A, 0x04},
			ExpectedOpcode: "inc",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x0B, 0x5B},
			ExpectedOpcode: "inc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x0B, 0xD5},
			ExpectedOpcode: "inc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x0B, 0x72},
			ExpectedOpcode: "inc",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x0B, 0xF5},
			ExpectedOpcode: "inc",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x59, 0x60},
			ExpectedOpcode: "jmp",
		},
		{
			Input:          []byte{0x5A, 0x12, 0x89, 0xDE},
			ExpectedOpcode: "jmp",
		},
		{
			Input:          []byte{0x5B, 0x3C},
			ExpectedOpcode: "jmp",
		},
		{
			Input:          []byte{0x5D, 0x60},
			ExpectedOpcode: "jsr",
		},
		{
			Input:          []byte{0x5E, 0x12, 0x89, 0xDE},
			ExpectedOpcode: "jsr",
		},
		{
			Input:          []byte{0x5F, 0x3C},
			ExpectedOpcode: "jsr",
		},
		{
			Input:          []byte{0x07, 0xC1},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x03, 0x04},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x01, 0x40, 0x69, 0x20},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6F, 0x10, 0x1F, 0xFF},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x78, 0x20, 0x6B, 0x20, 0x00, 0x12, 0x34, 0x56},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6D, 0x30},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6B, 0x00, 0x01, 0x26},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6B, 0x20, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0xF6, 0x63},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x0C, 0xD4},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x68, 0x79},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x6E, 0x72, 0xFF, 0xFF},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x78, 0x70, 0x6A, 0x2B, 0x00, 0xFF, 0xFF, 0x9D},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x6C, 0x41},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x26, 0xC0},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x6A, 0x0C, 0x01, 0x26},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x6A, 0x2A, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x68, 0xA0},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x6E, 0xC8, 0x82, 0xFB},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x78, 0x50, 0x6A, 0xA9, 0x00, 0xFE, 0x79, 0x61},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x6C, 0xA3},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x31, 0xC0},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x6A, 0x89, 0x01, 0x26},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x6A, 0xA2, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x79, 0x0B, 0x27, 0x0F},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x0D, 0x4A},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x69, 0x24},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x6F, 0x1A, 0x00, 0xFF},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x78, 0x20, 0x6B, 0x25, 0x00, 0x01, 0x38, 0x80},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x6D, 0x70},
			ExpectedOpcode: "pop",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x6B, 0x0E, 0x01, 0x26},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x6B, 0x25, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x69, 0xD8},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x6F, 0xF2, 0x01, 0x01},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x78, 0x60, 0x6B, 0xAD, 0x00, 0x00, 0x27, 0x10},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x6D, 0xFE},
			ExpectedOpcode: "push",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x6B, 0x88, 0x01, 0x26},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x6B, 0xAC, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x7A, 0x00, 0x00, 0x00, 0x11, 0xD7},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x0F, 0x81},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x69, 0x23},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6F, 0x74, 0x00, 0x3E},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x78, 0x70, 0x6B, 0x23, 0x00, 0x00, 0x27, 0x0E},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6D, 0x70},
			ExpectedOpcode: "pop",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6B, 0x03, 0x01, 0x26},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6B, 0x22, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x69, 0x95},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6F, 0xF4, 0x00, 0x22},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x78, 0xF0, 0x6B, 0xA5, 0x00, 0x00, 0x30, 0x0C},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6D, 0xF6},
			ExpectedOpcode: "push",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6B, 0x81, 0x01, 0x26},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6B, 0xA2, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: "mov",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x6A, 0x4D, 0xFF, 0xC0},
			ExpectedOpcode: "movfpe",
		},
		{
			Input:          []byte{0x6A, 0xC5, 0xFF, 0xC0},
			ExpectedOpcode: "movtpe",
		},
		{
			Input:          []byte{0x01, 0xC0, 0x50, 0x42},
			ExpectedOpcode: "mulxs",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x01, 0xC0, 0x52, 0x25},
			ExpectedOpcode: "mulxs",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x50, 0x42},
			ExpectedOpcode: "mulxu",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x52, 0x26},
			ExpectedOpcode: "mulxu",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x17, 0x88},
			ExpectedOpcode: "neg",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x17, 0x9C},
			ExpectedOpcode: "neg",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x17, 0xB5},
			ExpectedOpcode: "neg",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x00, 0x00},
			ExpectedOpcode: "nop",
		},
		{
			Input:          []byte{0x17, 0x0C},
			ExpectedOpcode: "not",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x17, 0x15},
			ExpectedOpcode: "not",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x17, 0x31},
			ExpectedOpcode: "not",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0xC1, 0x04},
			ExpectedOpcode: "or",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x14, 0x80},
			ExpectedOpcode: "or",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x79, 0x40, 0x00, 0xC0},
			ExpectedOpcode: "or",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x64, 0x80},
			ExpectedOpcode: "or",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x7A, 0x40, 0x00, 0x00, 0x00, 0xFE},
			ExpectedOpcode: "or",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0xF0, 0x64, 0x05},
			ExpectedOpcode: "or",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x04, 0x01},
			ExpectedOpcode: "orc",
		},
		{
			Input:          []byte{0x6D, 0x78},
			ExpectedOpcode: "pop",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6D, 0x73},
			ExpectedOpcode: "pop",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x6D, 0xF1},
			ExpectedOpcode: "push",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x00, 0x6D, 0xF6},
			ExpectedOpcode: "push",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x12, 0x81},
			ExpectedOpcode: "rotl",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x12, 0x9C},
			ExpectedOpcode: "rotl",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x12, 0xB6},
			ExpectedOpcode: "rotl",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x13, 0x81},
			ExpectedOpcode: "rotr",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x13, 0x94},
			ExpectedOpcode: "rotr",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x13, 0xB3},
			ExpectedOpcode: "rotr",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x12, 0x03},
			ExpectedOpcode: "rotxl",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x12, 0x16},
			ExpectedOpcode: "rotxl",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x12, 0x35},
			ExpectedOpcode: "rotxl",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x13, 0x0C},
			ExpectedOpcode: "rotxr",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x13, 0x1E},
			ExpectedOpcode: "rotxr",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x13, 0x34},
			ExpectedOpcode: "rotxr",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x56, 0x70},
			ExpectedOpcode: "rte",
		},
		{
			Input:          []byte{0x54, 0x70},
			ExpectedOpcode: "rts",
		},
		{
			Input:          []byte{0x10, 0x85},
			ExpectedOpcode: "shal",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x10, 0x9B},
			ExpectedOpcode: "shal",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x10, 0xB3},
			ExpectedOpcode: "shal",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x11, 0x84},
			ExpectedOpcode: "shar",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x11, 0x9E},
			ExpectedOpcode: "shar",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x11, 0xB5},
			ExpectedOpcode: "shar",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x10, 0x0E},
			ExpectedOpcode: "shll",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x10, 0x14},
			ExpectedOpcode: "shll",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x10, 0x32},
			ExpectedOpcode: "shll",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x11, 0x03},
			ExpectedOpcode: "shlr",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x11, 0x1A},
			ExpectedOpcode: "shlr",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x11, 0x31},
			ExpectedOpcode: "shlr",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x80},
			ExpectedOpcode: "sleep",
		},
		{
			Input:          []byte{0x02, 0x04},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x01, 0x40, 0x69, 0xF0},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6F, 0xF0, 0x00, 0x10},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x78, 0x70, 0x6B, 0xA0, 0x00, 0x00, 0x00, 0x64},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6D, 0xE0},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6B, 0x80, 0xFF, 0xC0},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x40, 0x6B, 0xA0, 0x00, 0x12, 0x89, 0xDE},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x18, 0x44},
			ExpectedOpcode: "sub",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x79, 0x3B, 0xFF, 0xF8},
			ExpectedOpcode: "sub",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x19, 0x0C},
			ExpectedOpcode: "sub",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x7A, 0x37, 0xFF, 0xFF, 0xFF, 0xF0},
			ExpectedOpcode: "sub",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x1A, 0x86},
			ExpectedOpcode: "sub",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x1B, 0x04},
			ExpectedOpcode: "subs",
		},
		{
			Input:          []byte{0x1B, 0x85},
			ExpectedOpcode: "subs",
		},
		{
			Input:          []byte{0x1B, 0x96},
			ExpectedOpcode: "subs",
		},
		{
			Input:          []byte{0xB5, 0x08},
			ExpectedOpcode: "subx",
		},
		{
			Input:          []byte{0x1E, 0x09},
			ExpectedOpcode: "subx",
		},
		{
			Input:          []byte{0x57, 0x20},
			ExpectedOpcode: "trapa",
		},
		{
			Input:          []byte{0xD4, 0x80},
			ExpectedOpcode: "xor",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x15, 0x4C},
			ExpectedOpcode: "xor",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x79, 0x5D, 0x20, 0x00},
			ExpectedOpcode: "xor",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x65, 0xD4},
			ExpectedOpcode: "xor",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x7A, 0x56, 0x00, 0x00, 0xFF, 0xFF},
			ExpectedOpcode: "xor",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0xF0, 0x65, 0x01},
			ExpectedOpcode: "xor",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x05, 0x40},
			ExpectedOpcode: "xorc",
		},
		{
			Input:          []byte{0x01, 0x41, 0x06, 0xA5},
			ExpectedOpcode: "andc",
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x76, 0x30},
			ExpectedOpcode: "band",
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x76, 0x50},
			ExpectedOpcode: "band",
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x72, 0x30},
			ExpectedOpcode: "bclr",
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x72, 0x50},
			ExpectedOpcode: "bclr",
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x62, 0x70},
			ExpectedOpcode: "bclr",
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x62, 0xE0},
			ExpectedOpcode: "bclr",
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x76, 0xB0},
			ExpectedOpcode: "biand",
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x67, 0x89, 0x76, 0xD0},
			ExpectedOpcode: "biand",
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x77, 0xB0},
			ExpectedOpcode: "bild",
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x77, 0xD0},
			ExpectedOpcode: "bild",
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x74, 0xB0},
			ExpectedOpcode: "bior",
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x74, 0xD0},
			ExpectedOpcode: "bior",
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x67, 0xB0},
			ExpectedOpcode: "bist",
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x67, 0xD0},
			ExpectedOpcode: "bist",
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x75, 0xB0},
			ExpectedOpcode: "bixor",
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x75, 0xD0},
			ExpectedOpcode: "bixor",
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x77, 0x30},
			ExpectedOpcode: "bld",
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x77, 0x50},
			ExpectedOpcode: "bld",
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x71, 0x30},
			ExpectedOpcode: "bnot",
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x71, 0x50},
			ExpectedOpcode: "bnot",
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x61, 0x70},
			ExpectedOpcode: "bnot",
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x61, 0xE0},
			ExpectedOpcode: "bnot",
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x74, 0x30},
			ExpectedOpcode: "bor",
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x74, 0x50},
			ExpectedOpcode: "bor",
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x70, 0x30},
			ExpectedOpcode: "bset",
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x70, 0x50},
			ExpectedOpcode: "bset",
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x60, 0x70},
			ExpectedOpcode: "bset",
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x60, 0xE0},
			ExpectedOpcode: "bset",
		},
		{
			Input:          []byte{0x6A, 0x18, 0x01, 0x23, 0x67, 0x30},
			ExpectedOpcode: "bst",
		},
		{
			Input:          []byte{0x6A, 0x38, 0x12, 0x34, 0x56, 0x78, 0x67, 0x50},
			ExpectedOpcode: "bst",
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x73, 0x30},
			ExpectedOpcode: "btst",
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x73, 0x50},
			ExpectedOpcode: "btst",
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x63, 0x70},
			ExpectedOpcode: "btst",
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x63, 0xE0},
			ExpectedOpcode: "btst",
		},
		{
			Input:          []byte{0x6A, 0x10, 0x01, 0x23, 0x75, 0x30},
			ExpectedOpcode: "bxor",
		},
		{
			Input:          []byte{0x6A, 0x30, 0x12, 0x34, 0x56, 0x78, 0x75, 0x50},
			ExpectedOpcode: "bxor",
		},
		{
			Input:          []byte{0x01, 0x41, 0x07, 0x6A},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x03, 0x14},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x01, 0x41, 0x69, 0x30},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6F, 0x30, 0x01, 0x23},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x78, 0x30, 0x6B, 0x20, 0x12, 0x34, 0x56, 0x78},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6D, 0x30},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6B, 0x00, 0x01, 0x23},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6B, 0x20, 0x12, 0x34, 0x56, 0x78},
			ExpectedOpcode: "ldc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x10, 0x6D, 0x73},
			ExpectedOpcode: "ldm",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x20, 0x6D, 0x76},
			ExpectedOpcode: "ldm",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x30, 0x6D, 0x77},
			ExpectedOpcode: "ldm",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x41, 0x04, 0x7B},
			ExpectedOpcode: "orc",
		},
		{
			Input:          []byte{0x12, 0xC8},
			ExpectedOpcode: "rotl",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x12, 0xD9},
			ExpectedOpcode: "rotl",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x12, 0xF6},
			ExpectedOpcode: "rotl",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x13, 0xC8},
			ExpectedOpcode: "rotr",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x13, 0xD9},
			ExpectedOpcode: "rotr",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x13, 0xF6},
			ExpectedOpcode: "rotr",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x12, 0x48},
			ExpectedOpcode: "rotxl",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x12, 0x59},
			ExpectedOpcode: "rotxl",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x12, 0x76},
			ExpectedOpcode: "rotxl",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x13, 0x48},
			ExpectedOpcode: "rotxr",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x13, 0x59},
			ExpectedOpcode: "rotxr",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x13, 0x76},
			ExpectedOpcode: "rotxr",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x10, 0xC5},
			ExpectedOpcode: "shal",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x10, 0xDE},
			ExpectedOpcode: "shal",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x10, 0xF5},
			ExpectedOpcode: "shal",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x11, 0xC5},
			ExpectedOpcode: "shar",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x11, 0xDE},
			ExpectedOpcode: "shar",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x11, 0xF5},
			ExpectedOpcode: "shar",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x10, 0x45},
			ExpectedOpcode: "shll",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x10, 0x5E},
			ExpectedOpcode: "shll",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x10, 0x75},
			ExpectedOpcode: "shll",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x11, 0x45},
			ExpectedOpcode: "shlr",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x11, 0x5E},
			ExpectedOpcode: "shlr",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x11, 0x75},
			ExpectedOpcode: "shlr",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x02, 0x15},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Byte,
		},
		{
			Input:          []byte{0x01, 0x41, 0x69, 0xD0},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6F, 0xD0, 0x01, 0x23},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x78, 0x50, 0x6B, 0xA0, 0x12, 0x34, 0x56, 0x78},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6D, 0xD0},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6B, 0x80, 0x01, 0x23},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x41, 0x6B, 0xA0, 0x12, 0x34, 0x56, 0x78},
			ExpectedOpcode: "stc",
			expectedBWL:    disassembler.Word,
		},
		{
			Input:          []byte{0x01, 0x10, 0x6D, 0xF2},
			ExpectedOpcode: "stm",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x20, 0x6D, 0xF4},
			ExpectedOpcode: "stm",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0x30, 0x6D, 0xF0},
			ExpectedOpcode: "stm",
			expectedBWL:    disassembler.Longword,
		},
		{
			Input:          []byte{0x01, 0xE0, 0x7B, 0x4C},
			ExpectedOpcode: "tas",
		},
		{
			Input:          []byte{0x01, 0x41, 0x05, 0x9E},
			ExpectedOpcode: "xorc",
		},
		{
			Input:          []byte{0x58, 0x00, 0xFB, 0x16},
			ExpectedOpcode: "bra",
		},
		{
			Input:          []byte{0x58, 0x10, 0xFB, 0x72},
			ExpectedOpcode: "brn",
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
		assert.Equal(t, tc.expectedBWL, inst.BWL, fmt.Sprintf("For byte sequence %x, expected BWL %s, got %s\n", tc.Input, BWLToString(tc.expectedBWL), BWLToString(inst.BWL)))
		assert.Equal(t, len(tc.Input), inst.TotalBytes, fmt.Sprintf("For byte sequence %x, expected %d bytes, got %d bytes\n", tc.Input, len(tc.Input), inst.TotalBytes))
	}
}

func BWLToString(bwl disassembler.Size) string {
	switch bwl {
	case disassembler.Byte:
		return "Byte"
	case disassembler.Word:
		return "Word"
	case disassembler.Longword:
		return "Long"
	case disassembler.Unset:
		return "Unset"
	}
	return ""
}
