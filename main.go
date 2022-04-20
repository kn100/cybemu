package main

import (
	"fmt"
	"os"
)

// Inst is a single instruction
type Inst struct {
	// The position of the instruction in the file
	Pos int
	// Either 2, or 4 (or maybe 6, or even 8, unimplemented for now)
	Size int
	// The Opcode of the instruction
	Opcode string
	// The raw bytes that make up the instruction
	Bytes []byte
}

func main() {
	f, err := os.Open("emu_rom.bin")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	finfo, err := os.Stat("emu_rom.bin")
	if err != nil {
		panic(err)
	}
	//
	instructions := []Inst{}
	// Read entire file into byte array
	// TODO: Set array to correct size using os.Stat
	bytes := make([]byte, finfo.Size())
	f.Read(bytes)
	i := 0
	for i < len(bytes) {
		// If we have less than 2 bytes left, we are done
		if i+2 > len(bytes) {
			break
		}
		AH := bytes[i] >> 4
		AL := bytes[i] & 0x0F
		inst := Inst{
			Pos:  i,
			Size: 2,
		}
		switch {
		// noop
		case AH == 0x0 && AL == 0x0:
			inst.Opcode = "nop"
		case AH == 0x0 && AL == 0x1:
			inst.Opcode, inst.Size = Table232(bytes[i : i+8])
		case AH == 0x0 && AL == 0x2:
			BH := bytes[i+1] >> 4
			if BH&0x8 == 0 {
				inst.Opcode = "stc"
			} else {
				inst.Opcode = "???word???"
			}
		case AH == 0x0 && AL == 0x3:
			BH := bytes[i+1] >> 4
			if BH&0x8 == 0 {
				inst.Opcode = "ldc"
			} else {
				inst.Opcode = "???word???"
			}
		case AH == 0x0 && AL == 0x4:
			inst.Opcode = "orc"
		case AH == 0x0 && AL == 0x5:
			inst.Opcode = "xorc"
		case AH == 0x0 && AL == 0x6:
			inst.Opcode = "andc"
		case AH == 0x0 && AL == 0x7:
			inst.Opcode = "ldc"
		case AH == 0x0 && AL == 0x8:
			inst.Opcode = "add"
		case AH == 0x0 && AL == 0x9:
			inst.Opcode = "add"
		case AH == 0x0 && (AL == 0xA || AL == 0xB || AL == 0xF):
			inst.Opcode, inst.Size = Table232(bytes[i : i+8])
		case AH == 0x0 && AL == 0xC:
			inst.Opcode = "mov"
		case AH == 0x0 && AL == 0xD:
			inst.Opcode = "mov"
		case AH == 0x0 && AL == 0xE:
			inst.Opcode = "addx"
		case AH == 0x1 && (AL == 0x0 || AL == 0x1 || AL == 0x2 || AL == 0x3 || AL == 0x7 || AL == 0xA || AL == 0xB || AL == 0xF):
			inst.Opcode, inst.Size = Table232(bytes[i : i+8])
		case AH == 0x1 && AL == 0x4:
			inst.Opcode = "or"
		case AH == 0x1 && AL == 0x5:
			inst.Opcode = "xor"
		case AH == 0x1 && AL == 0x6:
			inst.Opcode = "and"
		case AH == 0x1 && AL == 0x8:
			inst.Opcode = "sub"
		case AH == 0x1 && AL == 0x9:
			inst.Opcode = "sub"
		case AH == 0x1 && AL == 0xC:
			inst.Opcode = "cmp"
		case AH == 0x1 && AL == 0xD:
			inst.Opcode = "cmp"
		case AH == 0x1 && AL == 0xE:
			inst.Opcode = "subx"
		case AH == 0x2 || AH == 0x3:
			inst.Opcode = "mov"
		case AH == 0x4:
			inst.Opcode = Branches(AL)
		case AH == 0x5 && AL == 0x0:
			inst.Opcode = "mulxu"
		case AH == 0x5 && AL == 0x1:
			inst.Opcode = "divxu"
		case AH == 0x5 && AL == 0x2:
			inst.Opcode = "mulxu"
		case AH == 0x5 && AL == 0x3:
			inst.Opcode = "divxu"
		case AH == 0x5 && AL == 0x4:
			inst.Opcode = "rts"
		case AH == 0x5 && AL == 0x5:
			inst.Opcode = "bsr"
		case AH == 0x5 && AL == 0x6:
			inst.Opcode = "rte"
		case AH == 0x5 && AL == 0x7:
			inst.Opcode = "trapa"
		case AH == 0x5 && AL == 0x8:
			inst.Opcode, inst.Size = Table232(bytes[i : i+8])
		case AH == 0x5 && AL == 0x9:
			inst.Opcode = "jmp"
		case AH == 0x5 && AL == 0xA:
			inst.Opcode = "jmp"
		case AH == 0x5 && AL == 0xB:
			inst.Opcode = "jmp"
		case AH == 0x5 && AL == 0xC:
			inst.Opcode = "bsr"
		case AH == 0x5 && AL == 0xD:
			inst.Opcode = "jsr"
		case AH == 0x5 && AL == 0xE:
			inst.Opcode = "jsr"
		case AH == 0x5 && AL == 0xF:
			inst.Opcode = "jsr"
		case AH == 0x6 && AL == 0x0:
			inst.Opcode = "bset"
		case AH == 0x6 && AL == 0x1:
			inst.Opcode = "bnot"
		case AH == 0x6 && AL == 0x2:
			inst.Opcode = "bclr"
		case AH == 0x6 && AL == 0x3:
			inst.Opcode = "btst"
		case AH == 0x6 && AL == 0x4:
			inst.Opcode = "or"
		case AH == 0x6 && AL == 0x5:
			inst.Opcode = "xor"
		case AH == 0x6 && AL == 0x6:
			inst.Opcode = "and"
		case AH == 0x6 && AL == 0x7:
			BH := bytes[i+1] >> 4
			if BH&0x8 == 0 {
				inst.Opcode = "bst"
			} else {
				inst.Opcode = "bist"
			}
		case AH == 0x6 && AL == 0x8:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0x9:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0xA:
			inst.Opcode, inst.Size = Table232(bytes[i : i+8])
		case AH == 0x6 && AL == 0xB:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0xC:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0xD:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0xE:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0xF:
			inst.Opcode = "mov"
		case AH == 0x7 && AL == 0x0:
			inst.Opcode = "bset"
		case AH == 0x7 && AL == 0x1:
			inst.Opcode = "bnot"
		case AH == 0x7 && AL == 0x2:
			inst.Opcode = "bclr"
		case AH == 0x7 && AL == 0x3:
			inst.Opcode = "btst"
		case AH == 0x7 && AL == 0x4:
			BH := bytes[i+1] >> 4
			if BH&0x8 == 0 {
				inst.Opcode = "bor"
			} else {
				inst.Opcode = "bior"
			}
		case AH == 0x7 && AL == 0x5:
			BH := bytes[i+1] >> 4
			if BH&0x8 == 0 {
				inst.Opcode = "bxor"
			} else {
				inst.Opcode = "bixor"
			}
		case AH == 0x7 && AL == 0x6:
			BH := bytes[i+1] >> 4
			if BH&0x8 == 0 {
				inst.Opcode = "band"
			} else {
				inst.Opcode = "biand"
			}
		case AH == 0x7 && AL == 0x7:
			BH := bytes[i+1] >> 4
			if BH&0x8 == 0 {
				inst.Opcode = "bld"
			} else {
				inst.Opcode = "bild"
			}
		case AH == 0x7 && AL == 0x8:
			inst.Opcode = "mov"
		case AH == 0x7 && (AL == 0x9 || AL == 0xA || AL == 0xC || AL == 0xD || AL == 0xE || AL == 0xF):
			inst.Opcode, inst.Size = Table232(bytes[i : i+8])
		case AH == 0x7 && AL == 0xB:
			inst.Opcode = "eepmov"
		case AH == 0x8:
			inst.Opcode = "add"
		case AH == 0x9:
			inst.Opcode = "addx"
		case AH == 0xA:
			inst.Opcode = "cmp"
		case AH == 0xB:
			inst.Opcode = "subx"
		case AH == 0xC:
			inst.Opcode = "or"
		case AH == 0xD:
			inst.Opcode = "xor"
		case AH == 0xE:
			inst.Opcode = "and"
		case AH == 0xF:
			inst.Opcode = "mov"
		default:
			panic("wut")
		}
		for b := 0; b < inst.Size; b++ {
			inst.Bytes = append(inst.Bytes, bytes[i+b])
		}
		instructions = append(instructions, inst)
		i = i + inst.Size
	}
	PrintAssy(instructions)
}

// Table 2.3 (2) - returns the instruction, and the size of the instruction (2 or 4 bytes)
func Table232(bytes []byte) (string, int) {
	size := 2
	AH := bytes[0] >> 4
	AL := bytes[0] & 0x0F
	BH := bytes[1] >> 4
	switch {
	case AH == 0x0 && AL == 0x1:
		switch {
		case BH == 0x0:
			return "mov", size
		case BH == 0x1 || BH == 0x2 || BH == 0x3:
			if BH&0x8 == 0 {
				return "ldm", size
			} else {
				return "stm", size
			}
		case BH == 0x4:
			if BH&0x8 == 0 {
				return "ldc", size
			} else {
				return "stc", size
			}
		// 0x5, 0x6, 0x7, are illegal
		case BH == 0x8:
			return "sleep", size
		// 0x9, 0xA, and 0xB are illegal
		case BH == 0xC:
			return Table233(bytes)
		case BH == 0xD:
			return Table233(bytes)
		case BH == 0xE:
			return "tas", size
		case BH == 0xF:
			return Table233(bytes)
		default:
			return "???word???", size
		}
	case AH == 0x0 && AL == 0xA:
		switch {
		case BH == 0x0:
			return "inc", size
		//0x1 to 0x7 are illegal
		case BH == 0x8 || BH == 0x9 || BH == 0xA || BH == 0xB || BH == 0xC || BH == 0xD || BH == 0xE || BH == 0xF:
			return "add", size
		default:
			return "???word???", size
		}
	case AH == 0x0 && AL == 0xB:
		switch {
		case BH == 0x0:
			return "adds", size
		// 0x1 to 0x4 are illegal
		case BH == 0x5:
			return "inc", size
		// 0x6 is illegal
		case BH == 0x7:
			return "inc", size
		case BH == 0x8 || BH == 0x9:
			return "adds", size
		// 0xA, to 0xC are illegal
		case BH == 0xD:
			return "inc", size
		// 0xE is illegal
		case BH == 0xF:
			return "inc", size
		default:
			return "???word???", size
		}
	case AH == 0x0 && AL == 0xF:
		switch {
		case BH == 0x0:
			return "daa", size
		// 0x1 to 0x7 are illegal
		case BH == 0x8 || BH == 0x9 || BH == 0xA || BH == 0xB || BH == 0xC || BH == 0xD || BH == 0xE || BH == 0xF:
			return "mov", size
		default:
			return "???word???", size
		}
	case AH == 0x1 && AL == 0x0:
		switch {
		case BH == 0x0 || BH == 0x1:
			return "shll", size
		// 0x2 is illegal
		case BH == 0x3 || BH == 0x4 || BH == 0x5:
			return "shll", size
		// 0x6 is illegal
		case BH == 0x7:
			return "shll", size
		case BH == 0x8 || BH == 0x9:
			return "shal", size
		// 0xA is illegal
		case BH == 0xB || BH == 0xC || BH == 0xD:
			return "shal", size
		// 0xE is illegal
		case BH == 0xF:
			return "shal", size
		default:
			return "???word???", size
		}
	case AH == 0x1 && AL == 0x1:
		switch {
		case BH == 0x0 || BH == 0x1:
			return "shlr", size
		// 0x2 is illegal
		case BH == 0x3 || BH == 0x4 || BH == 0x5:
			return "shlr", size
		// 0x6 is illegal
		case BH == 0x7:
			return "shlr", size
		case BH == 0x8 || BH == 0x9:
			return "shar", size
		// 0xA is illegal
		case BH == 0xB || BH == 0xC || BH == 0xD:
			return "shar", size
		// 0xE is illegal
		case BH == 0xF:
			return "shar", size
		default:
			return "???word???", size
		}
	case AH == 0x1 && AL == 0x2:
		switch {
		case BH == 0x0 || BH == 0x1:
			return "rotxl", size
		// 0x2 is illegal
		case BH == 0x3 || BH == 0x4 || BH == 0x5:
			return "rotxl", size
		// 0x6 is illegal
		case BH == 0x7:
			return "rotxl", size
		case BH == 0x8 || BH == 0x9:
			return "rotl", size
		// 0xA is illegal
		case BH == 0xB || BH == 0xC || BH == 0xD:
			return "rotl", size
		// 0xE is illegal
		case BH == 0xF:
			return "rotl", size
		default:
			return "???word???", size
		}
	case AH == 0x1 && AL == 0x3:
		switch {
		case BH == 0x0 || BH == 0x1:
			return "rotxr", size
		// 0x2 is illegal
		case BH == 0x3 || BH == 0x4 || BH == 0x5:
			return "rotxr", size
		// 0x6 is illegal
		case BH == 0x7:
			return "rotxr", size
		case BH == 0x8 || BH == 0x9:
			return "rotr", size
		// 0xA is illegal
		case BH == 0xB || BH == 0xC || BH == 0xD:
			return "rotr", size
		// 0xE is illegal
		case BH == 0xF:
			return "rotr", size
		default:
			return "???word???", size
		}
	case AH == 0x1 && AL == 0x7:
		switch {
		case BH == 0x0 || BH == 0x1:
			return "not", size
		// 0x2 is illegal
		case BH == 0x3:
			return "not", size
		// 0x4 is illegal
		case BH == 0x5:
			return "extu", size
		// 0x6 is illegal
		case BH == 0x7:
			return "extu", size
		case BH == 0x8 || BH == 0x9:
			return "neg", size
		// 0xA is illegal
		case BH == 0xB:
			return "neg", size
		// 0xC is illegal
		case BH == 0xD:
			return "exts", size
		// 0xE is illegal
		case BH == 0xF:
			return "exts", size
		default:
			return "???word???", size
		}
	case AH == 0x1 && AL == 0xA:
		switch {
		case BH == 0x0:
			return "dec", size
		// 0x1 to 0x7 is illegal
		case BH == 0x8 || BH == 0x9 || BH == 0xA || BH == 0xB || BH == 0xC || BH == 0xD || BH == 0xE || BH == 0xF:
			return "sub", size
		default:
			return "???word???", size
		}
	case AH == 0x1 && AL == 0xB:
		switch {
		case BH == 0x0:
			return "subs", size
		// 0x1 to 0x4 is illegal
		case BH == 0x5:
			return "dec", size
		// 0x6 is illegal
		case BH == 0x7:
			return "dec", size
		case BH == 0x8 || BH == 0x9:
			return "subs", size
		// 0xA to 0xc is illegal
		case BH == 0xD:
			return "dec", size
		// 0xE is illegal
		case BH == 0xF:
			return "dec", size
		default:
			return "???word???", size
		}
	case AH == 0x1 && AL == 0xF:
		switch {
		case BH == 0x0:
			return "das", size
		// 0x1 to 0x7 is illegal
		case BH == 0x8 || BH == 0x9 || BH == 0xA || BH == 0xB || BH == 0xC || BH == 0xD || BH == 0xE || BH == 0xF:
			return "cmp", size
		default:
			return "???word???", size
		}
	case AH == 0x5 && AL == 0x8:
		return Branches(BH), size
	case AH == 0x6 && AL == 0xA:
		switch {
		case BH == 0x0:
			return "mov", size
		case BH == 0x1:
			return Table234(bytes)
		case BH == 0x2:
			return "mov", size
		case BH == 0x3:
			return Table234(bytes)
		case BH == 0x4:
			return "movfpe", size
		// 0x5 to 0x7 is illegal
		case BH == 0x8:
			return "mov", size
		// 0x9 is illegal
		case BH == 0xA:
			return "mov", size
		// 0xB is illegal
		case BH == 0xC:
			return "movtpe", size
		// 0xD to 0xf is illegal
		default:
			return "???word???", size
		}
	case AH == 0x7 && AL == 0x9:
		switch {
		case BH == 0x0:
			return "mov", size
		case BH == 0x1:
			return "add", size
		case BH == 0x2:
			return "cmp", size
		case BH == 0x3:
			return "sub", size
		case BH == 0x4:
			return "or", size
		case BH == 0x5:
			return "xor", size
		case BH == 0x6:
			return "and", size
		// 0x7 to 0xF is illegal
		default:
			return "???word???", size
		}
	case AH == 0x7 && AL == 0xA:
		switch {
		case BH == 0x0:
			return "mov", size
		case BH == 0x1:
			return "add", size
		case BH == 0x2:
			return "cmp", size
		case BH == 0x3:
			return "sub", size
		case BH == 0x4:
			return "or", size
		case BH == 0x5:
			return "xor", size
		case BH == 0x6:
			return "and", size
		// 0x7 to 0xF is illegal
		default:
			return "???word???", size
		}
	}
	return "???word???", size
}

// Table2.3 (3) - returns the instruction and the size in bytes (4)
func Table233(bytes []byte) (string, int) {
	size := 4
	AH := bytes[0] >> 4
	AL := bytes[0] & 0x0F
	BH := bytes[1] >> 4
	BL := bytes[1] & 0x0F
	CH := bytes[2] >> 4
	CL := bytes[2] & 0x0F
	DH := bytes[3] >> 4
	switch {
	case AH == 0x0 && AL == 0x1 && BH == 0xC && BL == 0x0 && CH == 0x5:
		switch {
		case CL == 0x0:
			return "mulxs", size
		case CL == 0x2:
			return "mulxs", size
		default:
			return "???word???", 2
		}
	case AH == 0x0 && AL == 0x1 && BH == 0xD && BL == 0x0 && CH == 0x5:
		switch {
		case CL == 0x1:
			return "divxs", size
		case CL == 0x3:
			return "divxs", size
		default:
			return "???word???", 2
		}
	case AH == 0x0 && AL == 0x1 && BH == 0xF && BL == 0x0 && CH == 0x6:
		switch {
		case CL == 0x4:
			return "or", size
		case CL == 0x5:
			return "xor", size
		case CL == 0x6:
			return "and", size
		default:
			return "???word???", 2
		}
	// BH here is actually a register field
	case AH == 0x7 && AL == 0xC && BL == 0x0 && CH == 0x6:
		switch {
		case CL == 0x3:
			return "btst", size
		default:
			return "???word???", 2
		}
	// BH here is actually a register field
	case AH == 0x7 && AL == 0xC && BL == 0x0 && CH == 0x7:
		switch {
		case CL == 0x3:
			return "btst", size
		case CL == 0x4:
			if DH&0x8 == 0 {
				return "bor", size
			} else {
				return "bior", size
			}
		case CL == 0x5:
			if DH&0x8 == 0 {
				return "bxor", size
			} else {
				return "bixor", size
			}
		case CL == 0x6:
			if DH&0x8 == 0 {
				return "band", size
			} else {
				return "biand", size
			}
		case CL == 0x7:
			if DH&0x8 == 0 {
				return "bld", size
			} else {
				return "bild", size
			}
		default:
			return "???word???", 2
		}
	// BH here is actually a register field
	case AH == 0x7 && AL == 0xD && BL == 0x0 && CH == 0x6:
		switch {
		case CL == 0x0:
			return "bset", size
		case CL == 0x1:
			return "bnot", size
		case CL == 0x2:
			return "bclr", size
		case CL == 0x7:
			if DH&0x8 == 0 {
				return "bst", size
			} else {
				return "bist", size
			}
		default:
			return "???word???", 2
		}
	// BH here is actually a register field
	case AH == 0x7 && AL == 0xD && BL == 0x0 && CH == 0x7:
		switch {
		case CL == 0x0:
			return "bset", size
		case CL == 0x1:
			return "bnot", size
		case CL == 0x2:
			return "bclr", size
		default:
			return "???word???", 2
		}
	// BH and BL here are actually an absolute address
	case AH == 0x7 && AL == 0xE && CH == 0x6:
		switch {
		case CL == 0x3:
			return "btst", size
		default:
			return "???word???", 2
		}
	// BH and BL here are actually an absolute address
	case AH == 0x7 && AL == 0xE && CH == 0x7:
		switch {
		case CL == 0x3:
			return "btst", size
		case CL == 0x4:
			if DH&0x8 == 0 {
				return "bor", size
			} else {
				return "bior", size
			}
		case CL == 0x5:
			if DH&0x8 == 0 {
				return "bxor", size
			} else {
				return "bixor", size
			}
		case CL == 0x6:
			if DH&0x8 == 0 {
				return "band", size
			} else {
				return "biand", size
			}
		case CL == 0x7:
			if DH&0x8 == 0 {
				return "bld", size
			} else {
				return "bild", size
			}
		default:
			return "???word???", 2
		}
	// BH and BL here are actually an absolute address
	case AH == 0x7 && AL == 0xF && CH == 0x6:
		switch {
		case CL == 0x0:
			return "bset", size
		case CL == 0x1:
			return "bnot", size
		case CL == 0x2:
			return "bclr", size
		case CL == 0x7:
			if DH&0x8 == 0 {
				return "bst", size
			} else {
				return "bist", size
			}
		default:
			return "???word???", 2
		}
		// BH and BL here are actually an absolute address
	case AH == 0x7 && AL == 0xF && CH == 0x7:
		// 0: bset, 1: bnot, 2: bclr
		switch {
		case CL == 0x0:
			return "bset", size
		case CL == 0x1:
			return "bnot", size
		case CL == 0x2:
			return "bclr", size
		default:
			return "???word???", 2
		}
	}
	return "???word???", 2
}

func Table234(bytes []byte) (string, int) {
	size := 6
	// Make AH, AL, BH, BL, CH, CL, DH, DL, EH, EL, FH
	AH := bytes[0] >> 4
	AL := bytes[0] & 0x0F
	BH := bytes[1] >> 4
	BL := bytes[1] & 0x0F

	// CH := bytes[2] >> 4      always address bytes
	// CL := bytes[2] & 0x0F    always address bytes
	// DH := bytes[3] >> 4      always address bytes
	// DL := bytes[3] & 0x0F    always address bytes
	EH := bytes[4] >> 4   // sometimes address byte
	EL := bytes[4] & 0x0F // sometimes address byte
	FH := bytes[5] >> 4   // sometimes address byte
	// FL := bytes[5] & 0x0F    always address byte
	GH := bytes[6] >> 4
	// GL appears unused?
	HH := bytes[7] >> 4
	switch {
	case AH == 0x6 && AL == 0xA && BH == 0x1 && BL == 0x0 && EH == 0x6:
		switch EL {
		case 0x3:
			return "btst", size
		default:
			return "???word???", 2
		}
	case AH == 0x6 && AL == 0xA && BH == 0x1 && BL == 0x0 && EH == 0x7:
		switch EL {
		case 0x3:
			return "btst", size
		case 0x4:
			if FH&0x8 == 0 {
				return "bor", size
			} else {
				return "bior", size
			}
		case 0x5:
			if FH&0x8 == 0 {
				return "bxor", size
			} else {
				return "bixor", size
			}
		case 0x6:
			if FH&0x8 == 0 {
				return "band", size
			} else {
				return "biand", size
			}
		case 0x7:
			if FH&0x8 == 0 {
				return "bld", size
			} else {
				return "bild", size
			}
		default:
			return "???word???", 2
		}
	case AH == 0x6 && AL == 0xA && BH == 0x1 && BL == 0x8 && EH == 0x6:
		switch EL {
		case 0x0:
			return "bset", size
		case 0x1:
			return "bnot", size
		case 0x2:
			return "bclr", size
		case 0x7:
			if FH&0x8 == 0 {
				return "bst", size
			} else {
				return "bist", size
			}
		default:
			return "???word???", 2
		}
	case AH == 0x6 && AL == 0xA && BH == 0x1 && BL == 0x8 && EH == 0x7:
		switch EL {
		case 0x0:
			return "bset", size
		case 0x1:
			return "bnot", size
		case 0x2:
			return "bclr", size
		default:
			return "???word???", 2
		}
	// AH 6 AL A BH 3 BL 0 GH 6
	case AH == 0x6 && AL == 0xA && BH == 0x3 && BL == 0x0 && GH == 0x6:
		switch EL {
		case 0x3:
			return "btst", size
		default:
			return "???word???", 2
		}
	case AH == 0x6 && AL == 0xA && BH == 0x3 && BL == 0x0 && GH == 0x7:
		switch EL {
		case 0x3:
			return "btst", size
		case 0x4:
			if HH&0x8 == 0 {
				return "bor", size
			} else {
				return "bior", size
			}
		case 0x5:
			if HH&0x8 == 0 {
				return "bxor", size
			} else {
				return "bixor", size
			}
		case 0x6:
			if HH&0x8 == 0 {
				return "band", size
			} else {
				return "biand", size
			}
		case 0x7:
			if HH&0x8 == 0 {
				return "bld", size
			} else {
				return "bild", size
			}
		default:
			return "???word???", 2
		}
	case AH == 0x6 && AL == 0xA && BH == 0x3 && BL == 0x8 && GH == 0x6:
		switch EL {
		case 0x0:
			return "bset", size
		case 0x1:
			return "bnot", size
		case 0x2:
			return "bclr", size
		case 0x7:
			if HH&0x8 == 0 {
				return "bst", size
			} else {
				return "bist", size
			}
		default:
			return "???word???", 2
		}
	case AH == 0x6 && AL == 0xA && BH == 0x3 && BL == 0x8 && GH == 0x7:
		switch EL {
		case 0x0:
			return "bset", size
		case 0x1:
			return "bnot", size
		case 0x2:
			return "bclr", size
		default:
			return "???word???", 2
		}
	default:
		return "???word???", 2
	}

}

func Branches(b byte) string {
	switch b {
	case 0x0:
		return "bra"
	case 0x1:
		return "brn"
	case 0x2:
		return "bh"
	case 0x3:
		return "bls"
	case 0x4:
		return "bcc"
	case 0x5:
		return "bcs"
	case 0x6:
		return "bne"
	case 0x7:
		return "beq"
	case 0x8:
		return "bvc"
	case 0x9:
		return "bvs"
	case 0xA:
		return "bpl"
	case 0xB:
		return "bmi"
	case 0xC:
		return "bge"
	case 0xD:
		return "blt"
	case 0xE:
		return "bgt"
	}
	return "ble"
}

func PrintAssy(instructions []Inst) {
	unks := 0
	decoded := 0
	word := 0
	// Count the number of unks
	for _, inst := range instructions {
		if inst.Opcode == "unk" {
			unks++
		} else if inst.Opcode == "???word???" {
			decoded++
			word++
		} else {
			decoded++
		}
	}
	fmt.Printf("%d instructions decoded, %d unks, %d words\n", decoded, unks, word)
	// display percentage decoded
	fmt.Printf("%d%% decoded\n", (decoded*100)/(decoded+unks))
	nopFlag := false
	for i, inst := range instructions {
		if inst.Size == 2 {
			// First time we see two nops, just print dots
			if !nopFlag && inst.Opcode == "nop" && instructions[i+1].Opcode == "nop" {
				nopFlag = true
				fmt.Println("...")
				continue
			}
			// Do nothing for every future pair of nops
			if nopFlag && inst.Opcode == "nop" && instructions[i+1].Opcode == "nop" {
				continue
			}
			// Once we find real code again, print out the last nop and continue as normal.
			if nopFlag && inst.Opcode == "nop" && instructions[i+1].Opcode != "nop" {
				nopFlag = false
			}
			fmt.Printf("%04X:\t%02X %02X\t%s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Opcode)
		} else if inst.Size == 4 {
			fmt.Printf("%04X:\t%02X %02X %02X %02X\t%s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.Opcode)
		}

	}
}
