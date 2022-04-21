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
		switch AH {
		case 0x0:
			switch AL {
			case 0x0:
				inst.Opcode = "nop"
			case 0x1:
				inst.Opcode, inst.Size = Table232(bytes[i : i+8])
			case 0x2:
				BH := bytes[i+1] >> 4
				if BH&0x8 == 0 {
					inst.Opcode = "stc"
				} else {
					inst.Opcode = ".word"
				}
			case 0x3:
				BH := bytes[i+1] >> 4
				if BH&0x8 == 0 {
					inst.Opcode = "ldc"
				} else {
					inst.Opcode = ".word"
				}
			case 0x4:
				inst.Opcode = "orc"
			case 0x5:
				inst.Opcode = "xorc"
			case 0x6:
				inst.Opcode = "andc"
			case 0x7:
				inst.Opcode = "ldc"
			case 0x8, 0x9:
				inst.Opcode = "add"
			case 0xA, 0xB, 0xF:
				inst.Opcode, inst.Size = Table232(bytes[i : i+8])
			case 0xC, 0xD:
				inst.Opcode = "mov"
			case 0xE:
				inst.Opcode = "addx"
			}
		case 0x1:
			switch AL {
			case 0x0, 0x1, 0x2, 0x3, 0x7, 0xA, 0xB, 0xF:
				inst.Opcode, inst.Size = Table232(bytes[i : i+8])
			case 0x4, 0x5, 0x6:
				inst.Opcode = OrXorAnd(AL, false)
			case 0x8, 0x9:
				inst.Opcode = "sub"
			case 0xC, 0xD:
				inst.Opcode = "cmp"
			case 0xE:
				inst.Opcode = "subx"
			}
		case 0x2, 0x3:
			inst.Opcode = "mov"
		case 0x4:
			inst.Opcode = Branches(AL)
		case 0x5:
			switch AL {
			case 0x0, 0x2:
				inst.Opcode = "mulxu"
			case 0x1, 0x3:
				inst.Opcode = "divxu"
			case 0x4:
				inst.Opcode = "rts"
			case 0x5:
				inst.Opcode = "bsr"
			case 0x6:
				inst.Opcode = "rte"
			case 0x7:
				inst.Opcode = "trapa"
			case 0x8:
				inst.Opcode, inst.Size = Table232(bytes[i : i+8])
			case 0x9, 0xA, 0xB:
				inst.Opcode = "jmp"
			case 0xC:
				inst.Opcode = "bsr"
			case 0xD, 0xE, 0xF:
				inst.Opcode = "jsr"
			}
		case 0x6:
			switch AL {
			case 0x0, 0x1, 0x2:
				inst.Opcode = BSetBNotBClr(AL)
			case 0x3:
				inst.Opcode = "btst"
			case 0x4, 0x5, 0x6:
				inst.Opcode = OrXorAnd(AL, false)
			case 0x7:
				BH := bytes[i+1] >> 4
				if BH&0x8 == 0 {
					inst.Opcode = "bst"
				} else {
					inst.Opcode = "bist"
				}
			case 0x8, 0x9, 0xB, 0xC, 0xD, 0xE, 0xF:
				inst.Opcode = "mov"
			case 0xA:
				inst.Opcode, inst.Size = Table232(bytes[i : i+8])
			}
		case 0x7:
			switch AL {
			case 0x0, 0x1, 0x2:
				inst.Opcode = BSetBNotBClr(AL)
			case 0x3, 0x4, 0x5, 0x6, 0x7:
				BH := bytes[i+1] >> 4
				inst.Opcode = BorBxorBandBld(AL, BH)
			case 0x8:
				inst.Opcode = "mov"
			case 0x9, 0xA, 0xC, 0xD, 0xE, 0xF:
				inst.Opcode, inst.Size = Table232(bytes[i : i+8])
			case 0xB:
				inst.Opcode = "eepmov"
			}

		case 0x8:
			inst.Opcode = "add"
		case 0x9:
			inst.Opcode = "addx"
		case 0xA:
			inst.Opcode = "cmp"
		case 0xB:
			inst.Opcode = "subx"
		case 0xC, 0xD, 0xE:
			inst.Opcode = OrXorAnd(AH, true)
		case 0xF:
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
	switch AH {
	case 0x0:
		switch AL {
		case 0x1:
			switch BH {
			case 0x0:
				return "mov", size
			case 0x1, 0x2, 0x3:
				if BH&0x8 == 0 {
					return "ldm", size
				} else {
					return "stm", size
				}
			case 0x4:
				if BH&0x8 == 0 {
					return "ldc", size
				} else {
					return "stc", size
				}
			case 0x8:
				return "sleep", size
			case 0xC:
				return Table233(bytes)
			case 0xD:
				return Table233(bytes)
			case 0xE:
				return "tas", size
			case 0xF:
				return Table233(bytes)
			default:
				return ".word", 2
			}
		case 0xA:
			switch BH {
			case 0x0:
				return "inc", size
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				return "add", size
			default:
				return ".word", 2
			}
		case 0xB:
			switch BH {
			case 0x0, 0x8, 0x9:
				return "adds", size
			case 0x5, 0x7, 0xD, 0xF:
				return "inc", size
			default:
				return ".word", 2
			}
		case 0xF:
			switch BH {
			case 0x0:
				return "daa", size
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				return "mov", size
			default:
				return ".word", 2
			}
		}
	case 0x1:
		switch AL {
		case 0x0:
			switch BH {
			case 0x0, 0x1, 0x3, 0x4, 0x5, 0x7:
				return "shll", size
			case 0x8, 0x9, 0xB, 0xC, 0xD, 0xF:
				return "shal", size
			default:
				return ".word", 2
			}
		case 0x1:
			switch BH {
			case 0x0, 0x1, 0x3, 0x4, 0x5, 0x7:
				return "shlr", size
			case 0x8, 0x9, 0xB, 0xC, 0xD, 0xF:
				return "shar", size
			default:
				return ".word", 2
			}
		case 0x2:
			switch BH {
			case 0x0, 0x1, 0x3, 0x4, 0x5, 0x7:
				return "rotxl", size
			case 0x8, 0x9, 0xB, 0xC, 0xD, 0xF:
				return "rotl", size
			default:
				return ".word", 2
			}
		case 0x3:
			switch BH {
			case 0x0, 0x1, 0x3, 0x4, 0x5, 0x7:
				return "rotxr", size
			case 0x8, 0x9, 0xB, 0xC, 0xD, 0xF:
				return "rotr", size
			default:
				return ".word", 2
			}
		case 0x7:
			switch BH {
			case 0x0, 0x1, 0x3:
				return "not", size
			case 0x5, 0x7:
				return "extu", size
			case 0x8, 0x9, 0xB:
				return "neg", size
			case 0xD, 0xF:
				return "exts", size
			default:
				return ".word", 2
			}
		case 0xA:
			switch BH {
			case 0x0:
				return "dec", size
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				return "sub", size
			default:
				return ".word", 2
			}
		case 0xB:
			switch BH {
			case 0x0, 0x8, 0x9:
				return "subs", size
			case 0x5, 0x7, 0xD, 0xF:
				return "dec", size
			default:
				return ".word", 2
			}
		case 0xF:
			switch BH {
			case 0x0:
				return "das", size
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				return "cmp", size
			default:
				return ".word", 2
			}
		}
	case 0x5:
		switch AL {
		case 0x8:
			return Branches(BH), size
		}
	case 0x6:
		switch AL {
		case 0xA:
			switch BH {
			case 0x0, 0x2, 0xA, 0x8:
				return "mov", size
			case 0x1, 0x3:
				return Table234(bytes)
			case 0x4:
				return "movfpe", size
			case 0xC:
				return "movtpe", size
			default:
				return ".word", 2
			}
		}
	case 0x7:
		switch AL {
		case 0x9, 0xA:
			switch BH {
			case 0x0:
				return "mov", size
			case 0x1:
				return "add", size
			case 0x2:
				return "cmp", size
			case 0x3:
				return "sub", size
			case 0x4, 0x5, 0x6:
				return OrXorAnd(BH, false), size
			default:
				return ".word", 2
			}
		}
	}
	return ".word", 2
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
	switch AH {
	case 0x0:
		if AL != 0x1 {
			return ".word", 2
		}

		switch BH {
		case 0xC:
			if !(BL == 0x0 && CH == 0x5) ||
				(CL != 0x0 && CL != 0x2) {
				return ".word", 2
			}
			return "mulxs", size

		case 0xD:
			if !(BL == 0x0 && CH == 0x5 && CL == 0x1) &&
				!(BL == 0x0 && CH == 0x5 && CL == 0x3) {
				return ".word", 2
			}
			return "divxs", size

		case 0xF:
			if !(BL == 0x0 && CH == 0x6) ||
				(CL != 0x4 && CL != 0x5 && CL != 0x6) {
				return ".word", 2
			}
			return OrXorAnd(CL, false), size
		}

	case 0x7:
		switch AL {
		case 0xC:
			if BL != 0x0 {
				return ".word", 2
			}
			switch CH {
			case 0x6:
				if CL != 0x3 {
					return ".word", 2
				}
				return "btst", size
			case 0x7:
				switch CL {
				case 0x3, 0x4, 0x5, 0x6, 0x7:
					return BorBxorBandBld(CL, DH), size
				default:
					return ".word", 2
				}
			}
		case 0xD:
			if BL != 0x0 {
				return ".word", 2
			}
			switch CH {
			case 0x6:
				switch CL {
				case 0x0, 0x1, 0x2:
					return BSetBNotBClr(CL), 2
				case 0x7:
					if DH&0x8 == 0 {
						return "bst", size
					} else {
						return "bist", size
					}
				default:
					return ".word", 2
				}
			case 0x7:
				switch CL {
				case 0x0, 0x1, 0x2:
					return BSetBNotBClr(CL), 2
				default:
					return ".word", 2
				}

			}
		case 0xE:
			switch CH {
			case 0x6:
				if CL != 0x3 {
					return ".word", 2
				}
				return "btst", size
			case 0x7:
				switch CL {
				case 0x3, 0x4, 0x5, 0x6, 0x7:
					return BorBxorBandBld(CL, DH), size
				default:
					return ".word", 2
				}
			}
		case 0xF:
			switch CH {
			case 0x6:
				switch CL {
				case 0x0, 0x1, 0x2:
					return BSetBNotBClr(CL), 2
				case 0x7:
					if DH&0x8 == 0 {
						return "bst", size
					} else {
						return "bist", size
					}
				default:
					return ".word", 2
				}
			case 0x7:
				switch CL {
				case 0x0, 0x1, 0x2:
					return BSetBNotBClr(CL), 2
				default:
					return ".word", 2
				}
			}
		}
	}
	return ".word", 2
}

func Table234(bytes []byte) (string, int) {
	size := 6
	AH := bytes[0] >> 4
	AL := bytes[0] & 0x0F

	if !(AH == 0x6 && AL == 0xA) {
		panic("Attempted to lookup instruction with wrong prefix. Only 0x6A is valid for the first byte." + fmt.Sprintf("First byte was %x", bytes[0]))
	}

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

	switch BH {
	case 0x1:
		switch BL {
		case 0x0:
			switch EH {
			case 0x6:
				switch EL {
				case 0x3:
					return "btst", size
				default:
					return ".word", 2
				}
			case 0x7:
				switch EL {
				case 0x3, 0x4, 0x5, 0x6, 0x7:
					return BorBxorBandBld(EL, FH), size
				default:
					return ".word", 2
				}
			}
		case 0x8:
			switch EH {
			case 0x6:
				switch EL {
				case 0x0, 0x1, 0x2:
					return BSetBNotBClr(EL), 2
				case 0x7:
					if FH&0x8 == 0 {
						return "bst", size
					} else {
						return "bist", size
					}
				default:
					return ".word", 2
				}
			case 0x7:
				switch EL {
				case 0x0, 0x1, 0x2:
					return BSetBNotBClr(EL), 2
				default:
					return ".word", 2
				}
			}
		}
	case 0x3:
		switch BL {
		case 0x0:
			switch GH {
			case 0x6:
				switch EL {
				case 0x3:
					return "btst", size
				default:
					return ".word", 2
				}
			case 0x7:
				switch EL {
				case 0x3, 0x4, 0x5, 0x6, 0x7:
					return BorBxorBandBld(EL, HH), size
				default:
					return ".word", 2
				}
			}
		case 0x8:
			switch GH {
			case 0x6:
				switch EL {
				case 0x0, 0x1, 0x2:
					return BSetBNotBClr(EL), 2
				case 0x7:
					if HH&0x8 == 0 {
						return "bst", size
					} else {
						return "bist", size
					}
				default:
					return ".word", 2
				}
			case 0x7:
				switch EL {
				case 0x0, 0x1, 0x2:
					return BSetBNotBClr(EL), 2
				default:
					return ".word", 2
				}
			default:
				return ".word", 2
			}
		}
	default:
		return ".word", 2

	}
	return ".word", 2
}

func Branches(b byte) string {
	branchMap := map[byte]string{
		0x0: "bra",
		0x1: "brn",
		0x2: "bh",
		0x3: "bls",
		0x4: "bcc",
		0x5: "bcs",
		0x6: "bne",
		0x7: "beq",
		0x8: "bvc",
		0x9: "bvs",
		0xA: "bpl",
		0xB: "bmi",
		0xC: "bge",
		0xD: "blt",
		0xE: "bgt",
		0xF: "ble",
	}
	return branchMap[b]
}

func BSetBNotBClr(b byte) string {
	bsetbnotbclrMap := map[byte]string{
		0x0: "bset",
		0x1: "bnot",
		0x2: "bclr",
	}
	return bsetbnotbclrMap[b]
}

func OrXorAnd(b byte, shift bool) string {
	if shift {
		b = b - 0x8
	}
	orXorAndMap := map[byte]string{
		0x4: "or",
		0x5: "xor",
		0x6: "and",
	}
	return orXorAndMap[b]
}

func BorBxorBandBld(b byte, CB byte) string {

	// if DH&0x8 != 0 and it's not going to be btst, shift b up by 8, so we can lookup in one map
	if CB&0x8 != 0 && b != 0x3 {
		b = b + 0x8
	}
	m := map[byte]string{
		0x3: "btst",
		0x4: "bor",
		0x5: "bxor",
		0x6: "band",
		0x7: "bld",
		0xC: "bior",
		0xD: "bixor",
		0xE: "biand",
		0xF: "bild",
	}
	return m[b]
}

func PrintAssy(instructions []Inst) {
	unks := 0
	decoded := 0
	word := 0
	// Count the number of unks
	for _, inst := range instructions {
		if inst.Opcode == "unk" {
			unks++
		} else if inst.Opcode == ".word" {
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
