package main

import (
	"fmt"
	"os"
)

// AddressingMode describes the addressing modes the H8S2000 supports.
type AddressingMode int64

const (
	None AddressingMode = iota
	Immediate
	RegisterDirect
	RegisterIndirect
	AbsoluteAddress
	ProgramCounterRelative
	MemoryIndirect
	RegisterIndirectWithDisplacement
	RegisterIndirectWithPostIncrement
	RegisterIndirectWithPreDecrement
)

// Size ...
type Size int64

const (
	Unset Size = iota
	Byte
	Word
	Longword
)

// Inst is a single instruction
type Inst struct {
	// The position of the instruction in the file
	Pos int
	// The number of bytes the instruction is.
	TotalBytes int
	BWL        Size
	Opcode     string
	// The raw bytes that make up the instruction
	Bytes          []byte
	AddressingMode AddressingMode
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
		BH := bytes[i+1] >> 4

		inst := Inst{
			Opcode:     ".word", // Everything is a word until proven otherwise
			Pos:        i,
			TotalBytes: 2,
		}
		switch AH {
		case 0x0:
			switch AL {
			case 0x0:
				BL := bytes[i+1] & 0x0F
				if BH == 0x0 && BL == 0x0 {
					inst.Opcode = "nop"
				}
			case 0x1:
				inst = Table232(inst, bytes[i:i+8])
			case 0x2:
				BH := bytes[i+1] >> 4
				if BH == 0 {
					inst.AddressingMode = RegisterDirect
					inst.BWL = Byte
					inst.Opcode = "stc"
				}
			case 0x3:
				BH := bytes[i+1] >> 4
				if BH == 0x0 {
					inst.AddressingMode = RegisterDirect
					inst.BWL = Byte
					inst.Opcode = "ldc"
				} else {
					break
				}
			case 0x4:
				inst.AddressingMode = Immediate
				inst.Opcode = "orc"
			case 0x5:
				inst.AddressingMode = Immediate
				inst.Opcode = "xorc"
			case 0x6:
				inst.AddressingMode = Immediate
				inst.Opcode = "andc"
			case 0x7:
				inst.AddressingMode = Immediate
				inst.BWL = Byte
				inst.Opcode = "ldc"
			case 0x8:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "add"
			case 0x9:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "add"
			case 0xA, 0xB, 0xF:
				inst = Table232(inst, bytes[i:i+8])
			case 0xC:
				inst.Opcode = "mov"
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
			case 0xD:
				inst.Opcode = "mov"
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
			case 0xE:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "addx"
			}
		case 0x1:
			inst.BWL = Byte
			switch AL {
			case 0x0, 0x1, 0x2, 0x3, 0x7, 0xA, 0xB, 0xF:
				inst = Table232(inst, bytes[i:i+8])
			case 0x4, 0x5, 0x6:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = OrXorAnd(AL, false)
			case 0x8:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "sub"
			case 0x9:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "sub"
			case 0xC:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "cmp"
			case 0xD:
				inst.BWL = Word
				inst.Opcode = "cmp"
			case 0xE:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "subx"
			}
		case 0x2:
			inst.Opcode = "mov"
			inst.AddressingMode = AbsoluteAddress
			inst.BWL = Byte
		case 0x3:
			inst.BWL = Byte
			inst.AddressingMode = AbsoluteAddress
			inst.Opcode = "mov"
		case 0x4:
			// BL := bytes[i+1] & 0x0F
			// BL must be even HEX number according to spec...
			// MAMEs decompiler doesn't seem to validate this detail, so I won't either and assume I don't know what I am doing
			// BL == 0x0 || BL == 0x2 || BL == 0x4 || BL == 0x6 || BL == 0x8 || BL == 0xA || BL == 0xC || BL == 0xE
			// In execution, the H8S2000 will set last bit to 0, which results in going to the previous instruction to the requested one.
			inst.AddressingMode = ProgramCounterRelative
			inst.Opcode = Branches(AL)
		case 0x5:
			switch AL {
			case 0x0:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "mulxu"
			case 0x2:
				BL := bytes[i+1] & 0x0F
				if BL > 7 {
					break
				}
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "mulxu"
			case 0x1:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "divxu"
			case 0x3:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "divxu"
			case 0x4:
				BL := bytes[i+1] & 0x0F
				if BH != 0x7 || BL != 0x0 {
					break
				}
				inst.Opcode = "rts"
			case 0x5:
				inst.AddressingMode = ProgramCounterRelative
				inst.Opcode = "bsr"
			case 0x6:
				BL := bytes[i+1] & 0x0F
				if BH != 0x7 || BL != 0x0 {
					break
				}
				inst.Opcode = "rte"
			case 0x7:
				BL := bytes[i+1] & 0x0F
				if BH > 0x3 || BL != 0 {
					break
				}
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "trapa"
			case 0x8:
				inst = Table232(inst, bytes[i:i+8])
			case 0x9:
				BL := bytes[i+1] & 0x0F
				if BH > 0x7 || BL != 0 {
					break
				}
				inst.AddressingMode = RegisterIndirect
				inst.Opcode = "jmp"
			case 0xA:
				inst.AddressingMode = AbsoluteAddress
				inst.TotalBytes = 4
				inst.Opcode = "jmp"
			case 0xB:
				inst.AddressingMode = MemoryIndirect
				inst.Opcode = "jmp"
			case 0xC:
				BL := bytes[i+1] & 0x0F
				if BH != 0x0 || BL != 0x0 {
					break
				}
				inst.AddressingMode = ProgramCounterRelative
				inst.TotalBytes = 4
				inst.Opcode = "bsr"
			case 0xD:
				BL := bytes[i+1] & 0x0F
				if BH > 0x7 || BL != 0 {
					break
				}
				inst.AddressingMode = RegisterIndirect
				inst.Opcode = "jsr"
			case 0xE:
				inst.AddressingMode = AbsoluteAddress
				inst.TotalBytes = 4
				inst.Opcode = "jsr"
			case 0xF:
				inst.AddressingMode = MemoryIndirect
				inst.Opcode = "jsr"

			}
		case 0x6:
			switch AL {
			case 0x0, 0x1, 0x2:
				inst.AddressingMode = RegisterDirect
				inst.Opcode = BSetBNotBClr(AL)
			case 0x3:
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "btst"
			case 0x4, 0x5, 0x6:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = OrXorAnd(AL, false)
			case 0x7:
				BH := bytes[i+1] >> 4
				inst.AddressingMode = RegisterDirect
				if BH&0x8 == 0 {
					inst.Opcode = "bst"
				} else {
					inst.Opcode = "bist"
				}
			case 0x8:
				inst.AddressingMode = RegisterIndirect
				inst.BWL = Byte
				inst.Opcode = "mov"
			case 0xC:
				if BH > 0x7 {
					inst.AddressingMode = RegisterIndirectWithPreDecrement
				} else {
					inst.AddressingMode = RegisterIndirectWithPostIncrement
				}
				inst.BWL = Byte
				inst.Opcode = "mov"
			case 0xE:
				inst.Opcode = "mov"
				inst.TotalBytes = 4
				inst.AddressingMode = RegisterIndirectWithDisplacement
				inst.BWL = Byte
			case 0x9:
				inst.BWL = Word
				inst.Opcode = "mov"
				inst.AddressingMode = RegisterIndirect
			case 0xB:
				if BH == 0x0 || BH == 0x8 {
					inst.TotalBytes = 4
					inst.BWL = Word
				} else if BH == 0x2 || BH == 0xA {
					inst.TotalBytes = 6
					inst.BWL = Word
				} else {
					break
				}
				inst.AddressingMode = AbsoluteAddress
				inst.Opcode = "mov"
			case 0xD:
				inst.Opcode = "mov"
				inst.BWL = Word
				inst.AddressingMode = RegisterIndirectWithPostIncrement
			case 0xF:
				inst.TotalBytes = 4
				inst.AddressingMode = RegisterIndirectWithDisplacement
				inst.BWL = Word
				inst.Opcode = "mov"
			case 0xA:
				inst = Table232(inst, bytes[i:i+8])
			}
		case 0x7:
			switch AL {
			case 0x0, 0x1, 0x2:
				if BH > 0x7 {
					break
				}
				inst.AddressingMode = RegisterDirect
				inst.Opcode = BSetBNotBClr(AL)
			case 0x3, 0x4, 0x5, 0x6, 0x7:
				inst.Opcode = BorBxorBandBld(AL, BH)
				if inst.Opcode != ".word" {
					inst.AddressingMode = RegisterDirect
				}
			case 0x8:
				CL := bytes[i+2] & 0x0F
				if CL == 0xA {
					inst.BWL = Byte
				} else if CL == 0xB {
					inst.BWL = Word
				} else {
					break
				}
				inst.AddressingMode = RegisterIndirectWithDisplacement
				inst.TotalBytes = 8
				inst.Opcode = "mov"
			case 0x9, 0xA, 0xC, 0xD, 0xE, 0xF:
				inst = Table232(inst, bytes[i:i+8])
			case 0xB:
				BL := bytes[i+1] & 0x0F
				CH := bytes[i+2] >> 4
				CL := bytes[i+2] & 0x0F
				DH := bytes[i+3] >> 4
				DL := bytes[i+3] & 0x0F
				if !(CH == 0x5 && CL == 0x9 && DH == 0x8 && DL == 0xF) {
					break
				}
				if BH == 5 && BL == 0xC {
					inst.BWL = Byte
				} else if BH == 0xD && BL == 0x4 {
					inst.BWL = Word
				} else {
					break
				}
				inst.TotalBytes = 4
				inst.Opcode = "eepmov"
			}
		case 0x8:
			inst.BWL = Byte
			inst.AddressingMode = Immediate
			inst.Opcode = "add"
		case 0x9:
			inst.BWL = Byte
			inst.AddressingMode = Immediate
			inst.Opcode = "addx"
		case 0xA:
			inst.BWL = Byte
			inst.AddressingMode = Immediate
			inst.Opcode = "cmp"
		case 0xB:
			inst.BWL = Byte
			inst.AddressingMode = Immediate
			inst.Opcode = "subx"
		case 0xC, 0xD, 0xE:
			inst.BWL = Byte
			inst.AddressingMode = Immediate
			inst.Opcode = OrXorAnd(AH, true)
		case 0xF:
			inst.AddressingMode = Immediate
			inst.BWL = Byte
			inst.Opcode = "mov"
		default:
			panic("wut")
		}
		for b := 0; b < inst.TotalBytes; b++ {
			inst.Bytes = append(inst.Bytes, bytes[i+b])
		}
		if inst.Opcode == ".word" && (inst.TotalBytes != 2 || inst.BWL != Unset || inst.AddressingMode != None) {
			fmt.Printf("STATE LEAK! Pos: %04x\n", inst.Pos)
		}
		instructions = append(instructions, inst)
		i = i + inst.TotalBytes
	}

	PrintAssy(instructions)

}

// Table 2.3 (2) - returns the instruction, and the size of the instruction (2 or 4 bytes)
func Table232(inst Inst, bytes []byte) Inst {
	AH := bytes[0] >> 4
	AL := bytes[0] & 0x0F
	BH := bytes[1] >> 4
	BL := bytes[1] & 0x0F
	CH := bytes[2] >> 4
	CL := bytes[2] & 0x0F
	DH := bytes[3] >> 4
	DL := bytes[3] & 0x0F
	EH := bytes[4] >> 4
	EL := bytes[4] & 0x0F
	FH := bytes[5] >> 4
	FL := bytes[5] & 0x0F
	switch AH {
	case 0x0:
		switch AL {
		case 0x1:
			switch BH {
			case 0x0:
				// Lord help me.. Also probably missing checks on size of registers (ie, extended registers are above 7, etc)
				if BL != 0x0 {
					break
				}
				// 01 00 69
				if CH == 0x6 && CL == 0x9 {
					inst.TotalBytes = 4
					inst.AddressingMode = RegisterIndirect
					inst.BWL = Longword
					// 01 00 6F
				} else if CH == 0x6 && CL == 0xF {
					inst.TotalBytes = 6
					inst.AddressingMode = RegisterIndirectWithDisplacement
					inst.BWL = Longword
					// 01 00 78 <80 6B A<8 NN NN NN NN
				} else if CH == 0x7 && CL == 0x8 && DL == 0x0 && EH == 0x6 && EL == 0xB && FH == 0x2 {
					inst.TotalBytes = 10
					inst.AddressingMode = RegisterIndirectWithDisplacement
					inst.BWL = Longword
					// 01 00 6B
				} else if CH == 0x6 && CL == 0xB {
					if DH == 0x0 {
						inst.TotalBytes = 6
					} else if DH == 0x2 {
						inst.TotalBytes = 8
					} else if DH == 0x8 {
						inst.TotalBytes = 6
					} else if DH == 0xA {
						inst.TotalBytes = 8
					} else {
						break
					}
					inst.AddressingMode = AbsoluteAddress
					inst.BWL = Longword
				} else if CH == 0x6 && CL == 0xE {
					inst.TotalBytes = 6
					inst.AddressingMode = RegisterIndirectWithPostIncrement
					inst.BWL = Longword
				} else if CH == 0x6 && CL == 0xD {
					inst.TotalBytes = 4
					inst.BWL = Longword
					inst.AddressingMode = RegisterIndirectWithPreDecrement
				} else {
					break
				}
				inst.Opcode = "mov"

			case 0x1, 0x2, 0x3:
				if !(BL == 0x0 &&

					CH == 0x6 &&
					CL == 0xD &&

					(DH == 0x7 || DH == 0xF) &&
					DL < 0x8) {
					break
				}
				inst.TotalBytes = 4
				inst.BWL = Longword
				if DH == 0x7 {
					inst.Opcode = "ldm"
				} else if DH == 0xF {
					inst.Opcode = "stm"
				}
			case 0x4:
				if DL != 0x0 {
					break
				}
				switch BL {
				case 0x0: // CCR
					if CH == 0x6 && CL == 0x9 {
						inst.TotalBytes = 4
						inst.AddressingMode = RegisterIndirect
						inst.BWL = Word
					} else if CH == 0x6 && CL == 0xF {
						inst.TotalBytes = 6
						inst.AddressingMode = RegisterIndirectWithDisplacement
						inst.BWL = Word
					} else if CH == 0x7 && CL == 0x8 && EH == 0x6 && EL == 0xB && FH == 0xA && FL == 0x0 {
						inst.TotalBytes = 10
						inst.AddressingMode = RegisterIndirectWithDisplacement
						inst.BWL = Word
					} else if CH == 0x6 && CL == 0xD {
						inst.TotalBytes = 4
						inst.AddressingMode = RegisterIndirectWithPreDecrement
						inst.BWL = Word
					} else if CH == 0x6 && CL == 0xB && DH == 0x8 {
						inst.TotalBytes = 6
						inst.AddressingMode = AbsoluteAddress
						inst.BWL = Word
					} else if CH == 0x6 && CL == 0xB && DH == 0xA {
						inst.TotalBytes = 8
						inst.AddressingMode = AbsoluteAddress
						inst.BWL = Word
					} else {
						break
					}
					if BH == 0 {
						inst.Opcode = "ldc"
					} else {
						inst.Opcode = "stc"
					}

				case 0x1: // EXR
					// TODO: Potentially missing an ORC here? (p166/322)
					if CH == 0x6 && CL == 0x9 {
						inst.TotalBytes = 4
						inst.AddressingMode = RegisterIndirect
						inst.BWL = Word
					} else if CH == 0x6 && CL == 0xF {
						inst.TotalBytes = 6
						inst.AddressingMode = RegisterIndirectWithDisplacement
						inst.BWL = Word
					} else if CH == 0x7 && CL == 0x8 && EH == 0x6 && EL == 0xB && FH == 0xA && FL == 0x0 {
						inst.TotalBytes = 10
						inst.AddressingMode = RegisterIndirectWithDisplacement
						inst.BWL = Word
					} else if CH == 0x6 && CL == 0xD {
						inst.TotalBytes = 4
						inst.AddressingMode = RegisterIndirectWithPreDecrement
						inst.BWL = Word
					} else if CH == 0x6 && CL == 0xB && DH == 0x8 {
						inst.TotalBytes = 6
						inst.AddressingMode = AbsoluteAddress
						inst.BWL = Word
					} else if CH == 0x6 && CL == 0xB && DH == 0xA {
						inst.TotalBytes = 8
						inst.AddressingMode = AbsoluteAddress
						inst.BWL = Word
					} else {
						break
					}
					if BH&0x8 == 0 {
						inst.Opcode = "ldc"
					} else {
						inst.Opcode = "stc"
					}
				default:
					break
				}

				inst.BWL = Word
				switch CH {
				case 0x6:
					switch CL {
					case 0x9:
						inst.TotalBytes = 4
						inst.AddressingMode = RegisterIndirect
					case 0xF:
						inst.TotalBytes = 6
						inst.AddressingMode = RegisterIndirectWithDisplacement
					case 0xD:
						inst.TotalBytes = 4
						inst.AddressingMode = RegisterIndirectWithPostIncrement
					case 0xB:
						switch DH {
						case 0x0:
							inst.TotalBytes = 6
							inst.AddressingMode = AbsoluteAddress
						case 0x2:
							// Address actually only is the fifth, sixth, and seventh byte, 8th byte is waste?
							inst.TotalBytes = 8
							inst.AddressingMode = AbsoluteAddress
						}
					}
				case 0x7:
					if CL == 0x8 && EH == 0x6 && EL == 0xB && FH == 0x2 && FL == 0x0 {
						inst.TotalBytes = 10
						inst.AddressingMode = RegisterIndirectWithDisplacement
					}
				}

			case 0x8:
				inst.Opcode = "sleep"
			case 0xC:
				inst = Table233(inst, bytes)
			case 0xD:
				inst = Table233(inst, bytes)
			case 0xE:
				inst.TotalBytes = 4
				inst.AddressingMode = RegisterIndirect
				inst.Opcode = "tas"
			case 0xF:
				inst = Table233(inst, bytes)
			}
		case 0xA:
			switch BH {
			case 0x0:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "inc"
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "add"
				inst.BWL = Longword
			}
		case 0xB:
			switch BH {
			case 0x0, 0x8, 0x9:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "adds"
			case 0x5, 0xD:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "inc"
			case 0x7, 0xF:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "inc"
			}
		case 0xF:
			if BH > 0x7 && BL < 0x8 {
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "mov"
			} else if BH == 0x0 {
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "daa"
			}
		}
	case 0x1:
		switch AL {
		case 0x0:
			switch BH {
			case 0x0:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shll"
			case 0x4:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shll"
			case 0x1:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shll"
			case 0x5:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shll"
			case 0x3:
				if BL > 0x7 {
					break
				}
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shll"
			case 0x7:
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shll"
			case 0x8:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shal"
			case 0xC:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shal"
			case 0x9:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shal"
			case 0xD:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shal"
			case 0xB:
				if BL > 0x7 {
					break
				}
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shal"
			case 0xF:
				if BL > 0x7 {
					break
				}
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shal"
			}
		case 0x1:
			switch BH {
			case 0x0:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shlr"
			case 0x4:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shlr"
			case 0x1:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shlr"
			case 0x5:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shlr"
			case 0x3:
				if BL > 0x7 {
					break
				}
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shlr"
			case 0x7:
				if BL > 0x7 {
					break
				}
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shlr"
			case 0x8:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shar"
			case 0xC:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shar"
			case 0x9:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shar"
			case 0xD:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shar"
			case 0xB:
				if BL > 0x7 {
					break
				}
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shar"
			case 0xF:
				if BL > 0x7 {
					break
				}
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "shar"
			}
		case 0x2:
			switch BH {
			case 0x0:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxl"
			case 0x4:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxl"
			case 0x1:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxl"
			case 0x5:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxl"
			case 0x3:
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxl"
			case 0x7:
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxl"
			case 0x8:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "rotl"
			case 0x9:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "rotl"
			case 0xC:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "rotl"
			case 0xD:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "rotl"
			case 0xB:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "rotl"
			case 0xF:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "rotl"
			}
		case 0x3:
			switch BH {
			case 0x0:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxr"
			case 0x4:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxr"
			case 0x1:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxr"
			case 0x5:
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxr"
			case 0x3:
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxr"
			case 0x7:
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "rotxr"
			case 0x8:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "rotr"
			case 0x9:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "rotr"
			case 0xC:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "rotr"
			case 0xD:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "rotr"
			case 0xB:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "rotr"
			case 0xF:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "rotr"
			}
		case 0x7:
			switch BH {
			case 0x0:
				inst.Opcode = "not"
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
			case 0x1:
				inst.Opcode = "not"
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
			case 0x3:
				inst.Opcode = "not"
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
			case 0x5:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "extu"
			case 0x7:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "extu"
			case 0x8:
				inst.Opcode = "neg"
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
			case 0x9:
				inst.Opcode = "neg"
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
			case 0xB:
				inst.Opcode = "neg"
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
			case 0xD:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "exts"
			case 0xF:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "exts"
			}
		case 0xA:
			switch BH {
			case 0x0:
				inst.BWL = Byte
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "dec"
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				if BL > 0x7 {
					break
				}
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "sub"
			}
		case 0xB:
			switch BH {
			case 0x0, 0x8, 0x9:
				if BL > 0x7 {
					break
				}
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "subs"
				inst.BWL = Longword
			case 0x5, 0xD:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "dec"
			case 0x7, 0xF:
				if BL > 0x7 {
					break
				}
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "dec"
			}

		case 0xF:
			switch BH {
			case 0x0:
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "das"
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				if BL > 0x7 {
					break
				}
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "cmp"
			}
		}
	case 0x5:
		switch AL {
		case 0x8:
			if BL != 0x0 {
				break
			}
			// TODO: CPU does weird things with odd dest
			inst.AddressingMode = ProgramCounterRelative
			inst.Opcode = Branches(BH)
			inst.TotalBytes = 4

		}
	case 0x6:
		switch AL {
		case 0xA:
			switch BH {
			case 0x0:
				inst.Opcode = "mov"
				inst.TotalBytes = 4
				inst.AddressingMode = AbsoluteAddress
				inst.BWL = Byte
			case 0x2:
				inst.Opcode = "mov"
				inst.TotalBytes = 6
				inst.AddressingMode = AbsoluteAddress
				inst.BWL = Byte
			case 0x8:
				inst.AddressingMode = AbsoluteAddress
				inst.TotalBytes = 4
				inst.BWL = Byte
				inst.Opcode = "mov"
			case 0xA:
				inst.AddressingMode = AbsoluteAddress
				inst.TotalBytes = 6
				inst.BWL = Byte
				inst.Opcode = "mov"
			case 0x1, 0x3:
				inst = Table234(inst, bytes)
			case 0x4:
				inst.AddressingMode = AbsoluteAddress
				inst.TotalBytes = 4
				inst.Opcode = "movfpe"
			case 0xC:
				inst.AddressingMode = AbsoluteAddress
				inst.TotalBytes = 4
				inst.Opcode = "movtpe"
			}
		}
	case 0x7:
		switch AL {
		case 0x9:
			inst.TotalBytes = 4
			inst.BWL = Word
			switch BH {
			case 0x0:
				inst.AddressingMode = Immediate
				inst.Opcode = "mov"
			case 0x1:
				inst.Opcode = "add"
			case 0x2:
				inst.Opcode = "cmp"
			case 0x3:
				inst.AddressingMode = Immediate
				inst.Opcode = "sub"
			case 0x4, 0x5, 0x6:
				if AL == 0x9 {
					inst.BWL = Word
				} else if AL == 0xA {
					inst.BWL = Longword
				}
				inst.AddressingMode = Immediate
				inst.Opcode = OrXorAnd(BH, false)
			}
		case 0xA:
			inst.TotalBytes = 6
			inst.BWL = Longword
			switch BH {
			case 0x0:
				inst.AddressingMode = Immediate
				inst.Opcode = "mov"
			case 0x1:
				inst.Opcode = "add"
			case 0x2:
				if BL > 0x7 {
					break
				}
				inst.Opcode = "cmp"
			case 0x3:
				if BL > 0x7 {
					break
				}
				inst.AddressingMode = Immediate
				inst.Opcode = "sub"
			case 0x4, 0x5, 0x6:
				if AL == 0x9 {
					inst.BWL = Word
				} else if AL == 0xA {
					inst.BWL = Longword
					if BL > 0x7 {
						break
					}
				}
				inst.AddressingMode = Immediate
				inst.Opcode = OrXorAnd(BH, false)
			}
		}
	}
	if inst.Opcode == ".word" {
		inst.TotalBytes = 2
		inst.BWL = Unset
	}
	return inst
}

// TODO: Potentially missing register direct sub.w (p235/322)
// TODO: Potentially missing XORC (p246/322)

// Table2.3 (3) - returns the instruction and the size in bytes (4)
func Table233(inst Inst, bytes []byte) Inst {
	AH := bytes[0] >> 4
	AL := bytes[0] & 0x0F
	BH := bytes[1] >> 4
	BL := bytes[1] & 0x0F
	CH := bytes[2] >> 4
	CL := bytes[2] & 0x0F
	DH := bytes[3] >> 4
	DL := bytes[3] & 0x0F
	switch AH {
	case 0x0:
		if AL != 0x1 {
			break
		}
		switch BH {
		case 0xC:
			if !(BL == 0x0 && CH == 0x5) {
				break
			}
			if CL == 0x0 {
				inst.BWL = Byte
			} else if CL == 0x2 {
				inst.BWL = Word
			} else {
				break
			}
			inst.TotalBytes = 4
			inst.Opcode = "mulxs"

		case 0xD:
			if !(BL == 0x0 && CH == 0x5) {
				break
			}
			switch CL {
			case 0x1:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "divxs"
				inst.TotalBytes = 4
			case 0x3:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "divxs"
				inst.TotalBytes = 4
			}
		case 0xF:
			if BL == 0x0 && CH == 0x6 && CL == 0x4 && DH < 0x8 && DL < 0x8 {
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "or"
				inst.TotalBytes = 4
			}
			//01f0 6524
			if BL == 0x0 && CH == 0x6 && CL == 0x5 && DH < 0x8 && DL < 0x8 {
				inst.BWL = Longword
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "xor"
				inst.TotalBytes = 4
			}
		}
	case 0xF:
		if !(BL == 0x0 && CH == 0x6) ||
			(CL != 0x4 && CL != 0x5 && CL != 0x6) {
			break
		}
		inst.Opcode = OrXorAnd(CL, false)

	case 0x7:
		switch AL {
		case 0xC:
			if BL != 0x0 {
				break
			}
			switch CH {
			case 0x6:
				if CL != 0x3 {
					break
				}
				inst.AddressingMode = RegisterIndirect
				inst.Opcode = "btst"
			case 0x7:
				switch CL {
				case 0x3, 0x4, 0x5, 0x6, 0x7:
					inst.AddressingMode = RegisterIndirect
					inst.Opcode = BorBxorBandBld(CL, DH)
				default:
					break
				}
			}
		case 0xD:
			if BL != 0x0 {
				break
			}
			switch CH {
			case 0x6:
				switch CL {
				case 0x0, 0x1, 0x2:
					inst.TotalBytes = 4
					inst.AddressingMode = RegisterIndirect
					inst.Opcode = BSetBNotBClr(CL)
				case 0x7:
					inst.AddressingMode = RegisterIndirect
					if DH&0x8 == 0 {
						inst.Opcode = "bst"
					} else {
						inst.Opcode = "bist"
					}
				default:
					break
				}
			case 0x7:
				switch CL {
				case 0x0, 0x1, 0x2:
					inst.TotalBytes = 4
					inst.AddressingMode = RegisterIndirect
					inst.Opcode = BSetBNotBClr(CL)
				default:
					break
				}

			}
		case 0xE:
			switch CH {
			case 0x6:
				if CL != 0x3 {
					break
				}
				inst.AddressingMode = AbsoluteAddress
				inst.Opcode = "btst"
			case 0x7:
				switch CL {
				case 0x3, 0x4, 0x5, 0x6, 0x7:
					inst.AddressingMode = AbsoluteAddress
					inst.Opcode = BorBxorBandBld(CL, DH)
				default:
					break
				}
			}
		case 0xF:
			switch CH {
			case 0x6:
				switch CL {
				case 0x0, 0x1, 0x2:
					inst.TotalBytes = 4
					inst.AddressingMode = AbsoluteAddress
					inst.Opcode = BSetBNotBClr(CL)
				case 0x7:
					inst.AddressingMode = AbsoluteAddress
					if DH&0x8 == 0 {
						inst.Opcode = "bst"
					} else {
						inst.Opcode = "bist"
					}
				default:
					break
				}
			case 0x7:
				switch CL {
				case 0x0, 0x1, 0x2:
					inst.TotalBytes = 4
					inst.AddressingMode = AbsoluteAddress
					inst.Opcode = BSetBNotBClr(CL)
				default:
					break
				}
			}
		}
	}
	if inst.Opcode == ".word" {
		inst.TotalBytes = 2
	}
	return inst
}

func Table234(inst Inst, bytes []byte) Inst {
	inst.TotalBytes = 6
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
	GL := bytes[6] & 0x0F
	// GL appears unused?
	HH := bytes[7] >> 4
	//HL := bytes[7] & 0x0F

	switch BH {
	case 0x1:
		switch BL {
		case 0x0:
			switch EH {
			case 0x6:
				switch EL {
				case 0x3:
					inst.AddressingMode = AbsoluteAddress
					inst.Opcode = "btst"
				default:
					break
				}
			case 0x7:
				switch EL {
				case 0x3, 0x4, 0x5, 0x6, 0x7:
					inst.AddressingMode = AbsoluteAddress
					inst.Opcode = BorBxorBandBld(EL, FH)
				default:
					break
				}
			}
		case 0x8:
			switch EH {
			case 0x6:
				switch EL {
				case 0x0, 0x1, 0x2:
					inst.TotalBytes = 6
					inst.AddressingMode = AbsoluteAddress
					inst.Opcode = BSetBNotBClr(EL)
				case 0x7:
					inst.AddressingMode = AbsoluteAddress
					if FH&0x8 == 0 {
						inst.Opcode = "bst"
					} else {
						inst.Opcode = "bist"
					}
				default:
					break
				}
			case 0x7:
				switch EL {
				case 0x0, 0x1, 0x2:
					inst.TotalBytes = 6
					inst.AddressingMode = AbsoluteAddress
					inst.Opcode = BSetBNotBClr(EL)
				default:
					break
				}
			}
		}
	case 0x3:
		switch BL {
		case 0x0:
			switch GH {
			case 0x6:
				switch GL {
				case 0x3:
					inst.AddressingMode = AbsoluteAddress
					inst.TotalBytes = 8
					inst.Opcode = "btst"
				default:
					break
				}
			case 0x7:
				switch GL {
				case 0x3, 0x4, 0x5, 0x6, 0x7:
					inst.AddressingMode = AbsoluteAddress
					inst.TotalBytes = 8
					inst.Opcode = BorBxorBandBld(EL, HH)
				default:
					break
				}
			}
		case 0x8:
			switch GH {
			case 0x6:
				switch EL {
				case 0x0, 0x1, 0x2:
					inst.TotalBytes = 8
					inst.AddressingMode = AbsoluteAddress
					inst.Opcode = BSetBNotBClr(EL)
				case 0x7:
					inst.TotalBytes = 8
					inst.AddressingMode = AbsoluteAddress
					if HH&0x8 == 0 {
						inst.Opcode = "bst"
					} else {
						inst.Opcode = "bist"
					}
				default:
					break
				}
			case 0x7:
				switch GL {
				case 0x0, 0x1, 0x2:
					inst.TotalBytes = 8
					inst.AddressingMode = AbsoluteAddress
					inst.Opcode = BSetBNotBClr(GL)
				}
			}
		}
	default:
		inst.TotalBytes = 2
	}
	if inst.Opcode == ".word" {
		inst.TotalBytes = 2
	}
	return inst
}

func Branches(b byte) string {
	branchMap := map[byte]string{
		0x0: "bt", // bra in the manual
		0x1: "bf", // brn in the manual
		0x2: "bhi",
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
	m := map[byte]string{
		0x3: "btst",
		0x4: "bor",
		0x5: "bxor",
		0x6: "band",
		0x7: "bld",
	}
	nm := map[byte]string{
		0x4: "bior",
		0x5: "bixor",
		0x6: "biand",
		0x7: "bild",
	}
	if CB > 0x7 {
		if val, ok := nm[b]; ok {
			return val
		}
		return ".word"

	} else {
		if val, ok := m[b]; ok {
			return val
		}
		return ".word"
	}
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
	// fmt.Printf("%d instructions decoded, %d unks, %d words\n", decoded, unks, word)
	// // display percentage decoded
	// fmt.Printf("%d%% decoded\n", (decoded*100)/(decoded+unks))
	for _, inst := range instructions {
		if inst.TotalBytes == 2 {
			fmt.Printf("%04x: %02x%02x                      %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], getInstWithSize(inst))
		} else if inst.TotalBytes == 4 {
			fmt.Printf("%04x: %02x%02x %02x%02x                 %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], getInstWithSize(inst))
		} else if inst.TotalBytes == 6 {
			fmt.Printf("%04x: %02x%02x %02x%02x %02x%02x            %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.Bytes[4], inst.Bytes[5], getInstWithSize(inst))
		} else if inst.TotalBytes == 8 {
			fmt.Printf("%04x: %02x%02x %02x%02x %02x%02x %02x%02x       %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.Bytes[4], inst.Bytes[5], inst.Bytes[6], inst.Bytes[7], getInstWithSize(inst))
		} else if inst.TotalBytes == 10 {
			fmt.Printf("%04x: %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x  %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.Bytes[4], inst.Bytes[5], inst.Bytes[6], inst.Bytes[7], inst.Bytes[8], inst.Bytes[9], getInstWithSize(inst))
		} else {
			panic("wat")
		}

	}
}

func getInstWithSize(inst Inst) string {
	sizeToString := map[Size]string{
		Byte:     ".b",
		Word:     ".w",
		Longword: ".l",
	}

	instsInexplicablyMissingSize := map[string]bool{
		"stc": true,
		"ldc": true,
	}

	if _, ok := instsInexplicablyMissingSize[inst.Opcode]; ok {
		return inst.Opcode
	} else {
		return inst.Opcode + sizeToString[inst.BWL]
	}
}

// Below is the worlds worst test.

// func Verifier(instructions []Inst) {
// 	sizeToString := map[Size]string{
// 		Byte:     "b",
// 		Word:     "w",
// 		Longword: "l",
// 	}

// 	debug := true
// 	// read in `output-from-h8disasm`
// 	f, err := os.Open("output-from-h8disasm")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	lines := []string{}
// 	scanner := bufio.NewScanner(f)
// 	for scanner.Scan() {
// 		lines = append(lines, scanner.Text())
// 	}
// 	if err := scanner.Err(); err != nil {
// 		panic(err)
// 	}
// 	missingPC, wrongInstruction, correct := 0, 0, 0
// 	// loop through instructions, and find the line which matches Pos
// 	for _, inst := range instructions {
// 		toFind := fmt.Sprintf("%04x:", inst.Pos)
// 		index := awfulFind(toFind, lines)
// 		if index == -1 {
// 			if inst.Opcode == "nop" {
// 				continue
// 			}
// 			missingPC++
// 			if debug {
// 				fmt.Printf("Couldn't find pc %s\n", toFind)
// 			}
// 			continue
// 		}
// 		opcodesize := inst.Opcode
// 		if inst.BWL != Unset && inst.Opcode != "ldc" && inst.Opcode != "stc" && inst.Opcode != "inc" && inst.Opcode != "dec" {
// 			opcodesize = fmt.Sprintf("%s.%s", inst.Opcode, sizeToString[inst.BWL])
// 		}

// 		if !strings.Contains(lines[index], opcodesize) {
// 			wrongInstruction++
// 			if debug {
// 				fmt.Printf("At position %x, we expected opcode %s, but wasn't present. Line was: %s\n", inst.Pos, opcodesize, lines[index])
// 			}
// 			continue
// 		}
// 		correct++
// 	}
// 	fmt.Printf("%d correct instructions, %d wrong program counter, %d wrong instructions\n", correct, missingPC, wrongInstruction)
// 	// percentages
// 	fmt.Printf("%d%% correct\n", (correct*100)/(correct+missingPC+wrongInstruction))
// }

// func awfulFind(substring string, in []string) int {
// 	for i, s := range in {
// 		if strings.Contains(s, substring) {
// 			return i
// 		}
// 	}
// 	return -1
// }

// func Hex2Bin(in byte) string {
// 	var out []byte
// 	for i := 7; i >= 0; i-- {
// 		b := (in >> uint(i))
// 		out = append(out, (b%2)+48)
// 	}
// 	return string(out)
// }
