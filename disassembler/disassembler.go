// disassembler contains utilities to disassemble a binary for the Renesas
// (Previously Hitachi) h8s/2000 class CPUs. It is being developed as a part of
// a larger project to emulate the Cybiko Classic. It roughly follows the layout
// of the tables found in "H8S/2600 Series, H8S/2000 Series Software Manual
// Rev.4.00 2006.02", Section 2.5 Table 2.3 (p274-277).
package disassembler

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

// Size indicates the size of an instruction (Byte, Word, Longword)
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
	OperandSize    int // Bits
}

// Disassemble takes a sequence of bytes from a compiled binary and disassembles
// them. It will return a slice of Inst.
func Disassemble(bytes []byte) []Inst {
	instructions := []Inst{}

	i := 0
	for i < len(bytes) {
		inst := Decode(bytes[i:])
		inst.Pos = i

		for b := 0; b < inst.TotalBytes; b++ {
			inst.Bytes = append(inst.Bytes, bytes[i+b])
		}

		instructions = append(instructions, inst)
		i = i + inst.TotalBytes
	}
	return instructions
}

// Decode takes a sequence of bytes, and returns the very first valid
// instruction found in the sequence. TotalBytes can be used to determine the
// end of the instruction, and therefore how many bytes passed were read. You
// can then call Decode again with those bytes to get the next instruction. The
// instructions specifically implemented here roughly follow p274 (see package
// comment), but will call the functions required to decode instructions found
// in the other tables too.
func Decode(bytes []byte) Inst {
	AH := bytes[0] >> 4
	AL := bytes[0] & 0x0F
	BH := bytes[1] >> 4
	inst := Inst{
		Opcode:     ".word", // Everything is a word until proven otherwise
		TotalBytes: 2,
	}
	switch AH {
	case 0x0:
		switch AL {
		case 0x0:
			BL := bytes[1] & 0x0F
			if BH == 0x0 && BL == 0x0 {
				inst.Opcode = "nop"
			}
		case 0x1:
			// 01
			inst = table232(inst, bytes[:8])
		case 0x2:
			BH := bytes[1] >> 4
			if BH == 0 || BH == 1 {
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "stc"
			}
		case 0x3:
			if BH == 0x0 {
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "ldc"
			} else if BH == 0x1 {
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "ldc"
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
			inst = table232(inst, bytes[:8])
		case 0xC:
			inst.Opcode = "mov"
			inst.BWL = Byte
			inst.AddressingMode = RegisterDirect
		case 0xD:
			inst.Opcode = "mov"
			inst.BWL = Word
			inst.AddressingMode = RegisterDirect
		case 0xE:
			inst.AddressingMode = RegisterDirect
			inst.Opcode = "addx"
		}
	case 0x1:
		switch AL {
		case 0x0, 0x1, 0x2, 0x3, 0x7, 0xA, 0xB, 0xF:
			inst = table232(inst, bytes[:8])
		case 0x4, 0x5, 0x6:
			inst.BWL = Byte
			inst.AddressingMode = RegisterDirect
			inst.Opcode = orXorAnd(AL, false)
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
			// TODO
			inst.BWL = Unset
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
		// BL := bytes[1] & 0x0F
		// BL must be even HEX number according to spec...
		// MAMEs decompiler doesn't seem to validate this detail, so I won't either and assume I don't know what I am doing
		// BL == 0x0 || BL == 0x2 || BL == 0x4 || BL == 0x6 || BL == 0x8 || BL == 0xA || BL == 0xC || BL == 0xE
		// In execution, the H8S2000 will set last bit to 0, which results in going to the previous instruction to the requested one.
		inst.AddressingMode = ProgramCounterRelative
		inst.Opcode = branches(AL)
		inst.OperandSize = 8
	case 0x5:
		switch AL {
		case 0x0:
			inst.AddressingMode = RegisterDirect
			inst.BWL = Byte
			inst.Opcode = "mulxu"
		case 0x2:
			BL := bytes[1] & 0x0F
			if BL < 0x8 {
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "mulxu"
			}
		case 0x1:
			inst.BWL = Byte
			inst.AddressingMode = RegisterDirect
			inst.Opcode = "divxu"
		case 0x3:
			BL := bytes[1] & 0x0F
			if BL < 0x8 {
				inst.BWL = Word
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "divxu"
			}

		case 0x4:
			BL := bytes[1] & 0x0F
			if BH == 0x7 && BL == 0x0 {
				inst.Opcode = "rts"
			}

		case 0x5:
			inst.AddressingMode = ProgramCounterRelative
			inst.Opcode = "bsr"
		case 0x6:
			BL := bytes[1] & 0x0F
			if BH == 0x7 && BL == 0x0 {
				inst.Opcode = "rte"
			}
		case 0x7:
			BL := bytes[1] & 0x0F
			if BH < 0x4 && BL == 0 {
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "trapa"
			}
		case 0x8:
			inst = table232(inst, bytes[:8])
		case 0x9:
			BL := bytes[1] & 0x0F
			if BH < 0x8 && BL == 0 {
				inst.AddressingMode = RegisterIndirect
				inst.Opcode = "jmp"
			}
		case 0xA:
			inst.AddressingMode = AbsoluteAddress
			inst.TotalBytes = 4
			inst.Opcode = "jmp"
			inst.OperandSize = 24
		case 0xB:
			inst.AddressingMode = MemoryIndirect
			inst.Opcode = "jmp"
			inst.OperandSize = 8
		case 0xC:
			BL := bytes[1] & 0x0F
			if BH == 0x0 && BL == 0x0 {
				inst.AddressingMode = ProgramCounterRelative
				inst.TotalBytes = 4
				inst.Opcode = "bsr"
				inst.OperandSize = 16
			}
		case 0xD:
			BL := bytes[1] & 0x0F
			if BH < 0x8 && BL == 0 {
				inst.AddressingMode = RegisterIndirect
				inst.Opcode = "jsr"
			}
		case 0xE:
			inst.AddressingMode = AbsoluteAddress
			inst.TotalBytes = 4
			inst.Opcode = "jsr"
			inst.OperandSize = 24
		case 0xF:
			inst.AddressingMode = MemoryIndirect
			inst.Opcode = "jsr"
			inst.OperandSize = 8
		}
	case 0x6:
		switch AL {
		case 0x0, 0x1, 0x2:
			inst.AddressingMode = RegisterDirect
			inst.Opcode = bSetBNotBClr(AL)
		case 0x3:
			inst.AddressingMode = RegisterDirect
			inst.Opcode = "btst"
		case 0x4, 0x5, 0x6:
			inst.BWL = Word
			inst.AddressingMode = RegisterDirect
			inst.Opcode = orXorAnd(AL, false)
		case 0x7:
			BH := bytes[1] >> 4
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
				inst.BWL = Byte
				inst.Opcode = "mov"
			} else {
				inst.AddressingMode = RegisterIndirectWithPostIncrement
				inst.BWL = Byte
				inst.Opcode = "mov"
			}
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
				inst.OperandSize = 16
				inst.BWL = Word
				inst.AddressingMode = AbsoluteAddress
				inst.Opcode = "mov"
			} else if BH == 0x2 || BH == 0xA {
				inst.TotalBytes = 6
				inst.OperandSize = 32
				inst.BWL = Word
				inst.AddressingMode = AbsoluteAddress
				inst.Opcode = "mov"
			}
		case 0xD:
			switch BH {
			case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6:
				inst.Opcode = "mov"
				inst.BWL = Word
				inst.AddressingMode = RegisterIndirectWithPostIncrement
			case 0x7:
				inst.Opcode = "pop"
				inst.BWL = Word
				inst.AddressingMode = RegisterIndirectWithPostIncrement
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE:
				inst.Opcode = "mov"
				inst.BWL = Word
				inst.AddressingMode = RegisterIndirectWithPostIncrement
			case 0xF:
				inst.Opcode = "push"
				inst.BWL = Word
				inst.AddressingMode = RegisterIndirectWithPostIncrement
			}
		case 0xF:
			inst.OperandSize = 16
			inst.TotalBytes = 4
			inst.AddressingMode = RegisterIndirectWithDisplacement
			inst.BWL = Word
			inst.Opcode = "mov"
		case 0xA:
			inst = table232(inst, bytes[:8])
		}
	case 0x7:
		switch AL {
		case 0x0, 0x1, 0x2:
			if BH < 0x8 {
				inst.AddressingMode = RegisterDirect
				inst.Opcode = bSetBNotBClr(AL)
			}

		case 0x3, 0x4, 0x5, 0x6, 0x7:
			inst.Opcode = borBxorBandBld(AL, BH)
			if inst.Opcode != ".word" {
				inst.AddressingMode = RegisterDirect
			}
		case 0x8:
			BL := bytes[1] & 0x0F
			CH := bytes[2] >> 4
			CL := bytes[2] & 0x0F
			DH := bytes[3] >> 4
			if BH < 0x8 && BL == 0x0 && CH == 0x6 && (DH == 0x2 || DH == 0xA) {
				if CL == 0xA {
					inst.AddressingMode = RegisterIndirectWithDisplacement
					inst.TotalBytes = 8
					inst.Opcode = "mov"
					inst.BWL = Byte
				} else if CL == 0xB {
					inst.AddressingMode = RegisterIndirectWithDisplacement
					inst.TotalBytes = 8
					inst.Opcode = "mov"
					inst.BWL = Word
				}
			}
		case 0x9, 0xA, 0xC, 0xD, 0xE, 0xF:
			// 79, 7A, 7B, 7C, 7D, 7E, 7F
			inst = table232(inst, bytes[:8])
		case 0xB:
			BL := bytes[1] & 0x0F
			CH := bytes[2] >> 4
			CL := bytes[2] & 0x0F
			DH := bytes[3] >> 4
			DL := bytes[3] & 0x0F
			if CH == 0x5 && CL == 0x9 && DH == 0x8 && DL == 0xF {
				if BH == 5 && BL == 0xC {
					inst.BWL = Byte
					inst.TotalBytes = 4
					inst.Opcode = "eepmov"
				} else if BH == 0xD && BL == 0x4 {
					inst.BWL = Word
					inst.TotalBytes = 4
					inst.Opcode = "eepmov"
				}
			}
		}
	case 0x8:
		inst.BWL = Byte
		inst.AddressingMode = Immediate
		inst.Opcode = "add"
	case 0x9:
		inst.AddressingMode = Immediate
		inst.Opcode = "addx"
	case 0xA:
		inst.BWL = Byte
		inst.AddressingMode = Immediate
		inst.Opcode = "cmp"
	case 0xB:
		inst.AddressingMode = Immediate
		inst.Opcode = "subx"
	case 0xC, 0xD, 0xE:
		inst.BWL = Byte
		inst.AddressingMode = Immediate
		inst.Opcode = orXorAnd(AH, true)
	case 0xF:
		inst.BWL = Byte
		inst.AddressingMode = Immediate
		inst.Opcode = "mov"
	}

	return inst
}

func table232(inst Inst, bytes []byte) Inst {
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
	if AH == 0x0 && AL == 0x1 && BH == 0x0 && BL == 0x0 {
		// 01 00 69
		if CH == 0x6 && CL == 0x9 && DL < 0x8 {
			inst.TotalBytes = 4
			inst.AddressingMode = RegisterIndirect
			inst.BWL = Longword
			inst.Opcode = "mov"
			// 01 00 6F
		} else if CH == 0x6 && CL == 0xF {
			inst.TotalBytes = 6
			inst.AddressingMode = RegisterIndirectWithDisplacement
			inst.BWL = Longword
			inst.Opcode = "mov"
		} else if CH == 0x7 && CL == 0x8 && DH < 0x8 && DL == 0x0 && EH == 0x6 && EL == 0xB && (FH == 0x2 || FH == 0xA) && FL < 0x8 {
			// 01 00 78 ?0 6B A? ?? ?? ?? ??
			// 01 00 78 ?0 6B 2? ?? ?? ?? ??
			// TODO: Potential bug in unidasm? should be checking `DH < 0x7` (p142/322)
			// Turns out the real hardware doesn't care what DH is so just set it to zero.
			inst.TotalBytes = 10
			inst.AddressingMode = RegisterIndirectWithDisplacement
			inst.BWL = Longword
			inst.Opcode = "mov"
		} else if CH == 0x6 && CL == 0xB {
			if DH == 0x0 || DH == 0x8 {
				inst.TotalBytes = 6
				inst.AddressingMode = AbsoluteAddress
				inst.BWL = Longword
				inst.Opcode = "mov"
			} else if DH == 0x2 || DH == 0xA {
				inst.TotalBytes = 8
				inst.AddressingMode = AbsoluteAddress
				inst.BWL = Longword
				inst.Opcode = "mov"
			}
		} else if CH == 0x6 && CL == 0xD && DH == 0x7 && DL < 0x8 {
			inst.TotalBytes = 4
			inst.BWL = Longword
			inst.AddressingMode = RegisterIndirectWithPreDecrement
			inst.Opcode = "pop"
		} else if CH == 0x6 && CL == 0xD && DH < 0x7 && DL < 0x8 {
			inst.TotalBytes = 4
			inst.BWL = Longword
			inst.AddressingMode = RegisterIndirectWithPreDecrement
			inst.Opcode = "mov"
		} else if CH == 0x6 && CL == 0xD && DH == 0xF && DL < 0x8 {
			inst.TotalBytes = 4
			inst.BWL = Longword
			inst.AddressingMode = RegisterIndirectWithPreDecrement
			inst.Opcode = "push"
		}
	}
	if AH == 0x0 && AL == 0x1 && BH == 0x4 {
		if BL == 0x0 {
			if CH == 0x6 && CL == 0x9 {
				// 01 40 69
				inst.TotalBytes = 4
				inst.AddressingMode = RegisterIndirect
				inst.BWL = Word
				if DH < 0x8 {
					inst.Opcode = "ldc"
				} else {
					inst.Opcode = "stc"
				}
			} else if CH == 0x6 && CL == 0xF {
				// 01 40 6F
				inst.TotalBytes = 6
				inst.AddressingMode = RegisterIndirectWithDisplacement
				inst.BWL = Word
				if DH < 0x8 {
					inst.Opcode = "ldc"
				} else {
					inst.Opcode = "stc"
				}
			} else if CH == 0x7 && CL == 0x8 && DH < 0x8 && EH == 0x6 && EL == 0xB && (FH == 0xA || FH == 0x2) && FL == 0x0 {
				// 01 40 78 ?? 6B A0
				inst.TotalBytes = 10
				inst.AddressingMode = RegisterIndirectWithDisplacement
				inst.BWL = Word
				if FH < 0xA {
					inst.Opcode = "ldc"
				} else {
					inst.Opcode = "stc"
				}
			} else if CH == 0x6 && CL == 0xD {
				// 01 40 6D
				inst.TotalBytes = 4
				inst.BWL = Word
				if DH > 0x7 {
					inst.AddressingMode = RegisterIndirectWithPreDecrement
					inst.Opcode = "stc"

				} else {
					inst.AddressingMode = RegisterIndirectWithPostIncrement
					inst.Opcode = "ldc"
					inst.BWL = Word
				}
			} else if CH == 0x6 && CL == 0xB && DL == 0x0 {
				switch DH {
				case 0x0:
					// 01 40 6B 00
					inst.TotalBytes = 6
					inst.AddressingMode = AbsoluteAddress
					inst.BWL = Word
					inst.Opcode = "ldc"
					inst.OperandSize = 16
				case 0x2:
					// 01 40 6B 20
					inst.TotalBytes = 8
					inst.AddressingMode = AbsoluteAddress
					inst.BWL = Word
					inst.Opcode = "ldc"
					inst.OperandSize = 32
				case 0x8:
					// 01 40 6B 80
					inst.TotalBytes = 6
					inst.AddressingMode = AbsoluteAddress
					inst.BWL = Word
					inst.Opcode = "stc"
					inst.OperandSize = 16
				case 0xA:
					// 01 40 6B A0
					inst.TotalBytes = 8
					inst.AddressingMode = AbsoluteAddress
					inst.BWL = Word
					inst.Opcode = "stc"
					inst.OperandSize = 32
				}

			}
		} else if BL == 0x1 {
			if CH == 0x6 && CL == 0x9 {
				// 01 41 69
				inst.TotalBytes = 4
				inst.AddressingMode = RegisterIndirect
				inst.BWL = Word
				if DH < 0x8 {
					inst.Opcode = "ldc"
				} else {
					inst.Opcode = "stc"
				}
			} else if CH == 0x6 && CL == 0xF {
				// 01 41 6F
				inst.TotalBytes = 6
				inst.AddressingMode = RegisterIndirectWithDisplacement
				inst.BWL = Word
				if DH < 0x8 {
					inst.Opcode = "ldc"
				} else {
					inst.Opcode = "stc"
				}
			} else if CH == 0x7 && CL == 0x8 && DH < 0x8 && EH == 0x6 && EL == 0xB && (FH == 0xA || FH == 0x2) && FL == 0x0 {
				// 01 41 78 ?? 6B A0
				inst.TotalBytes = 10
				inst.AddressingMode = RegisterIndirectWithDisplacement
				inst.BWL = Word
				if FH < 0xA {
					inst.Opcode = "ldc"
				} else {
					inst.Opcode = "stc"
				}
			} else if CH == 0x6 && CL == 0xD {
				// 01 41 6D
				if DH > 0x7 {
					inst.TotalBytes = 4
					inst.BWL = Word
					inst.AddressingMode = RegisterIndirectWithPreDecrement
					inst.Opcode = "stc"
				} else {
					inst.TotalBytes = 4
					inst.BWL = Word
					inst.AddressingMode = RegisterIndirectWithPostIncrement
					inst.Opcode = "ldc"
				}
			} else if CH == 0x6 && CL == 0xB && DL == 0x0 {
				switch DH {
				case 0x0:
					// 01 41 6B 00
					inst.TotalBytes = 6
					inst.AddressingMode = AbsoluteAddress
					inst.BWL = Word
					inst.Opcode = "ldc"
					inst.OperandSize = 16
				case 0x2:
					// 01 41 6B 20
					inst.TotalBytes = 8
					inst.AddressingMode = AbsoluteAddress
					inst.BWL = Word
					inst.Opcode = "ldc"
					inst.OperandSize = 32
				case 0x8:
					// 01 41 6B 80
					inst.TotalBytes = 6
					inst.AddressingMode = AbsoluteAddress
					inst.BWL = Word
					inst.Opcode = "stc"
					inst.OperandSize = 16
				case 0xA:
					// 01 41 6B A0
					inst.TotalBytes = 8
					inst.AddressingMode = AbsoluteAddress
					inst.BWL = Word
					inst.Opcode = "stc"
					inst.OperandSize = 32
				}
			} else if CH == 0x0 && CL == 0x7 {
				// 01 41 07
				inst.TotalBytes = 4
				inst.AddressingMode = Immediate
				inst.BWL = Byte
				inst.Opcode = "ldc"
			} else if CH == 0x0 && CL == 0x6 {
				// 01 41 06
				inst.TotalBytes = 4
				inst.AddressingMode = Immediate
				inst.Opcode = "andc"
			} else if CH == 0x0 && CL == 0x5 {
				// 01 41 05
				inst.TotalBytes = 4
				inst.AddressingMode = Immediate
				inst.Opcode = "xorc"

			} else if CH == 0x0 && CL == 0x4 {
				// 01 41 04
				inst.TotalBytes = 4
				inst.AddressingMode = Immediate
				inst.Opcode = "orc"
			}
		}
	}
	switch AH {
	case 0x0:
		switch AL {
		case 0x1:
			switch BH {
			case 0x1, 0x2, 0x3:
				if BL == 0x0 && CH == 0x6 && CL == 0xD && (DH == 0x7 || DH == 0xF) && DL < 0x8 {
					inst.TotalBytes = 4
					inst.BWL = Longword
					if DH == 0x7 {
						inst.Opcode = "ldm"
					} else if DH == 0xF {
						inst.Opcode = "stm"
					}
				}
			case 0x8:
				if BL == 0x0 {
					inst.Opcode = "sleep"
				}

			case 0xC:
				inst = table233(inst, bytes)
			case 0xD:
				inst = table233(inst, bytes)
			case 0xE:
				if BL == 0x0 && CH == 0x7 && CL == 0xB && DH < 0x8 && DL == 0xC {
					inst.TotalBytes = 4
					inst.AddressingMode = RegisterIndirect
					inst.Opcode = "tas"
				}
			case 0xF:
				inst = table233(inst, bytes)
			}
		case 0xA:
			switch BH {
			case 0x0:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "inc"
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				if BL < 0x8 {
					inst.AddressingMode = RegisterDirect
					inst.Opcode = "add"
					inst.BWL = Longword
				}

			}
		case 0xB:
			switch BH {
			case 0x0, 0x8, 0x9:
				// 0B B0, 0B B8, 0B B9
				if BL < 0x8 {
					inst.AddressingMode = RegisterDirect
					inst.Opcode = "adds"
				}
			case 0x5, 0xD:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "inc"
			case 0x7, 0xF:
				if BL < 0x8 {
					inst.AddressingMode = RegisterDirect
					inst.BWL = Longword
					inst.Opcode = "inc"
				}
			}
		case 0xF:
			if BH > 0x7 && BL < 0x8 {
				inst.AddressingMode = RegisterDirect
				inst.BWL = Longword
				inst.Opcode = "mov"
			} else if BH == 0x0 {
				inst.AddressingMode = RegisterDirect
				inst.Opcode = "daa"
			}
		}
	case 0x1:
		switch AL {
		case 0x0:
			switch BH {
			case 0x0:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "shll"
			case 0x4:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "shll"
			case 0x1:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "shll"
			case 0x5:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "shll"
			case 0x3:
				if BL < 0x8 {
					inst.AddressingMode = RegisterDirect
					inst.BWL = Longword
					inst.Opcode = "shll"
				}

			case 0x7:
				if BL < 0x8 {
					inst.AddressingMode = RegisterDirect
					inst.BWL = Longword
					inst.Opcode = "shll"
				}
			case 0x8:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "shal"
			case 0xC:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Byte
				inst.Opcode = "shal"
			case 0x9:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "shal"
			case 0xD:
				inst.AddressingMode = RegisterDirect
				inst.BWL = Word
				inst.Opcode = "shal"
			case 0xB:
				if BL < 0x8 {
					inst.AddressingMode = RegisterDirect
					inst.BWL = Longword
					inst.Opcode = "shal"
				}

			case 0xF:
				if BL < 0x8 {
					inst.AddressingMode = RegisterDirect
					inst.BWL = Longword
					inst.Opcode = "shal"
				}
			}
		case 0x1:
			switch BH {
			case 0x0:
				inst.BWL = Byte
				inst.Opcode = "shlr"
			case 0x4:
				inst.BWL = Byte
				inst.Opcode = "shlr"
			case 0x1:
				inst.BWL = Word
				inst.Opcode = "shlr"
			case 0x5:
				inst.BWL = Word
				inst.Opcode = "shlr"
			case 0x3:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "shlr"
				}
			case 0x7:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "shlr"
				}
			case 0x8:
				inst.BWL = Byte
				inst.Opcode = "shar"
			case 0xC:
				inst.BWL = Byte
				inst.Opcode = "shar"
			case 0x9:
				inst.BWL = Word
				inst.Opcode = "shar"
			case 0xD:
				inst.BWL = Word
				inst.Opcode = "shar"
			case 0xB:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "shar"
				}
			case 0xF:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "shar"
				}
			}
		case 0x2:
			switch BH {
			case 0x0:
				inst.BWL = Byte
				inst.Opcode = "rotxl"
			case 0x4:
				inst.BWL = Byte
				inst.Opcode = "rotxl"
			case 0x1:
				inst.BWL = Word
				inst.Opcode = "rotxl"
			case 0x5:
				inst.BWL = Word
				inst.Opcode = "rotxl"
			case 0x3:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "rotxl"
				}

			case 0x7:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "rotxl"
				}

			case 0x8:
				inst.BWL = Byte
				inst.Opcode = "rotl"
			case 0x9:
				inst.BWL = Word
				inst.Opcode = "rotl"
			case 0xC:
				inst.BWL = Byte
				inst.Opcode = "rotl"
			case 0xD:
				inst.BWL = Word
				inst.Opcode = "rotl"
			case 0xB:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "rotl"
				}
			case 0xF:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "rotl"
				}
			}
		case 0x3:
			switch BH {
			case 0x0:
				inst.BWL = Byte
				inst.Opcode = "rotxr"
			case 0x4:
				inst.BWL = Byte
				inst.Opcode = "rotxr"
			case 0x1:
				inst.BWL = Word
				inst.Opcode = "rotxr"
			case 0x5:
				inst.BWL = Word
				inst.Opcode = "rotxr"
			case 0x3:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "rotxr"
				}
			case 0x7:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "rotxr"
				}
			case 0x8:
				inst.BWL = Byte
				inst.Opcode = "rotr"
			case 0x9:
				inst.BWL = Word
				inst.Opcode = "rotr"
			case 0xC:
				inst.BWL = Byte
				inst.Opcode = "rotr"
			case 0xD:
				inst.BWL = Word
				inst.Opcode = "rotr"
			case 0xB:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "rotr"
				}
			case 0xF:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "rotr"
				}
			}
		case 0x7:
			switch BH {
			case 0x0:
				inst.Opcode = "not"
				inst.BWL = Byte
			case 0x1:
				inst.Opcode = "not"
				inst.BWL = Word
			case 0x3:
				if BL < 0x8 {
					inst.Opcode = "not"
					inst.BWL = Longword
				}
			case 0x5:
				inst.BWL = Word
				inst.Opcode = "extu"
			case 0x7:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "extu"
				}
			case 0x8:
				inst.Opcode = "neg"
				inst.BWL = Byte
			case 0x9:
				inst.Opcode = "neg"
				inst.BWL = Word
			case 0xB:
				if BL < 0x8 {
					inst.Opcode = "neg"
					inst.BWL = Longword
				}
			case 0xD:
				inst.BWL = Word
				inst.Opcode = "exts"
			case 0xF:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "exts"
				}
			}
		case 0xA:
			switch BH {
			case 0x0:
				inst.BWL = Byte
				inst.Opcode = "dec"
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "sub"
				}
			}
		case 0xB:
			switch BH {
			case 0x0, 0x8, 0x9:
				if BL < 0x8 {
					inst.Opcode = "subs"
				}
			case 0x5, 0xD:
				inst.BWL = Word
				inst.Opcode = "dec"
			case 0x7, 0xF:
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "dec"
				}
			}
		case 0xF:
			switch BH {
			case 0x0:
				// 1F 0?
				inst.Opcode = "das"
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				// 1F 8?, 1F 9?, 1F A?, 1F B?, 1F C?, 1F D?, 1F E?, 1F F?
				if BL < 0x8 {
					inst.BWL = Longword
					inst.Opcode = "cmp"
				}
			}
		}
	case 0x5:
		inst.AddressingMode = ProgramCounterRelative
		switch AL {
		case 0x8:
			if BL == 0x0 {
				// TODO: CPU does weird things with odd dest
				inst.Opcode = branches(BH)
				inst.TotalBytes = 4
			}

		}
	case 0x6:
		switch AL {
		case 0xA:
			switch BH {
			case 0x0:
				inst.AddressingMode = AbsoluteAddress
				inst.Opcode = "mov"
				inst.TotalBytes = 4
				inst.BWL = Byte
			case 0x2:
				inst.AddressingMode = AbsoluteAddress
				inst.Opcode = "mov"
				inst.TotalBytes = 6
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
				inst = table234(inst, bytes)
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
			switch BH {
			case 0x0:
				inst.TotalBytes = 4
				inst.AddressingMode = Immediate
				inst.BWL = Word
				inst.Opcode = "mov"
			case 0x1:
				inst.TotalBytes = 4
				inst.AddressingMode = Immediate
				inst.BWL = Word
				inst.Opcode = "add"
			case 0x2:
				inst.TotalBytes = 4
				inst.AddressingMode = Immediate
				inst.BWL = Word
				inst.Opcode = "cmp"
			case 0x3:
				inst.TotalBytes = 4
				inst.AddressingMode = Immediate
				inst.BWL = Word
				inst.Opcode = "sub"
			case 0x4, 0x5, 0x6:
				inst.TotalBytes = 4
				inst.AddressingMode = Immediate
				inst.BWL = Word
				inst.Opcode = orXorAnd(BH, false)
			}
		case 0xA:
			// TODO: Possible too broad
			if BL < 0x8 {
				switch BH {
				case 0x0:
					inst.AddressingMode = Immediate
					inst.TotalBytes = 6
					inst.BWL = Longword
					inst.Opcode = "mov"
				case 0x1:
					inst.AddressingMode = Immediate
					inst.TotalBytes = 6
					inst.BWL = Longword
					inst.Opcode = "add"
				case 0x2:
					inst.AddressingMode = Immediate
					inst.TotalBytes = 6
					inst.BWL = Longword
					inst.Opcode = "cmp"
				case 0x3:
					inst.AddressingMode = Immediate
					inst.TotalBytes = 6
					inst.BWL = Longword
					inst.Opcode = "sub"
				case 0x4, 0x5, 0x6:
					inst.Opcode = orXorAnd(BH, false)
					if inst.Opcode != ".word" {
						inst.AddressingMode = Immediate
						inst.TotalBytes = 6
						inst.BWL = Longword
					}
				}
			}
		case 0xC:
			if BH < 0x8 && BL == 0x0 && CH == 0x7 && CL == 0x3 && DH < 0x8 && DL == 0x0 {
				// 7C ?0 73 ?0
				inst.AddressingMode = RegisterIndirect
				inst.TotalBytes = 4
				inst.Opcode = "btst"
			} else {
				// 7C
				inst = table233(inst, bytes)
			}

		case 0xD:
			if BH < 0x8 && BL == 0x0 && CH == 0x7 && CL == 0x0 && DH < 0x8 && DL == 0x0 {
				inst.AddressingMode = RegisterIndirect
				inst.TotalBytes = 4
				inst.Opcode = "bset"
			} else {
				inst = table233(inst, bytes)
			}

		case 0xE:
			switch CH {
			case 0x6:
				// 7E ?? 63 ?0
				if CL == 0x3 && DL == 0x0 {
					inst.AddressingMode = AbsoluteAddress
					inst.TotalBytes = 4
					inst.Opcode = "btst"
				}

			case 0x7:
				// 7E ?? 73 ?0
				if CL == 0x3 && DH < 0x8 && DL == 0x0 {
					inst.AddressingMode = AbsoluteAddress
					inst.TotalBytes = 4
					inst.Opcode = "btst"
				} else {
					inst = table233(inst, bytes)
				}

			}
		case 0xF:
			inst = table233(inst, bytes)
		}
	}

	return inst
}

func table233(inst Inst, bytes []byte) Inst {
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
		if AL == 0x1 {
			switch BH {
			case 0xC:
				if BL == 0x0 && CH == 0x5 {
					if CL == 0x0 {
						inst.AddressingMode = RegisterDirect
						inst.TotalBytes = 4
						inst.BWL = Byte
						inst.Opcode = "mulxs"
					} else if CL == 0x2 && DL < 0x8 {
						inst.AddressingMode = RegisterDirect
						inst.TotalBytes = 4
						inst.BWL = Word
						inst.Opcode = "mulxs"
					}
				}

			case 0xD:
				if BL == 0x0 && CH == 0x5 {
					switch CL {
					case 0x1:
						inst.AddressingMode = RegisterDirect
						inst.TotalBytes = 4
						inst.BWL = Byte
						inst.Opcode = "divxs"
					case 0x3:
						if DL < 0x8 {
							inst.AddressingMode = RegisterDirect
							inst.TotalBytes = 4
							inst.BWL = Word
							inst.Opcode = "divxs"
						}
					}
				}

			case 0xF:
				if BL == 0x0 && CH == 0x6 && CL == 0x4 && DH < 0x8 && DL < 0x8 {
					inst.AddressingMode = RegisterDirect
					inst.TotalBytes = 4
					inst.BWL = Longword
					inst.Opcode = "or"
				}
				//01f0 6524
				if BL == 0x0 && CH == 0x6 && CL == 0x5 && DH < 0x8 && DL < 0x8 {
					inst.AddressingMode = RegisterDirect
					inst.TotalBytes = 4
					inst.BWL = Longword
					inst.Opcode = "xor"
				}
				if BL == 0x0 && CH == 0x6 && CL == 0x6 && DH < 0x8 && DL < 0x8 {
					inst.AddressingMode = RegisterDirect
					inst.TotalBytes = 4
					inst.BWL = Longword
					inst.Opcode = "and"
				}
			}
		}

	case 0x7:
		switch AL {
		case 0xC:
			if BH < 0x8 && BL == 0x0 {
				switch CH {
				case 0x6:
					if CL == 0x3 {
						inst.AddressingMode = RegisterIndirect
						inst.TotalBytes = 4
						inst.Opcode = "btst"
					}

				case 0x7:
					switch CL {
					case 0x3, 0x4, 0x5, 0x6, 0x7:
						// 7C ?0 73, 7C ?0 74, 7C ?0 75, 7C ?0 76, 7C ?0 77
						inst.Opcode = borBxorBandBld(CL, DH)
						if inst.Opcode != ".word" {
							inst.AddressingMode = RegisterIndirect
							inst.TotalBytes = 4
						}
					}
				}
			}

		case 0xD:
			if BH < 0x8 && BL == 0x0 {
				// 7D ?0
				switch CH {
				case 0x6:
					switch CL {
					case 0x0, 0x1, 0x2:
						// 7D ?0 60, 7D ?0 61, 7D ?0 62
						inst.Opcode = bSetBNotBClr(CL)
						if inst.Opcode != ".word" {
							inst.AddressingMode = RegisterIndirect
							inst.TotalBytes = 4
						}
					case 0x7:
						// 7D ?0 67
						inst.TotalBytes = 4
						if DH&0x8 == 0 {
							inst.AddressingMode = RegisterIndirect
							inst.Opcode = "bst"
						} else {
							inst.Opcode = "bist"
						}
					}
				case 0x7:
					switch CL {
					case 0x0, 0x1, 0x2:
						if DH < 0x8 {
							// 7D ?0 70 ?0, 7D ?0 71 ?0, 7D ?0 72 ?0
							inst.Opcode = bSetBNotBClr(CL)
							if inst.Opcode != ".word" {
								inst.AddressingMode = RegisterIndirect
								inst.TotalBytes = 4
							}
						}
					}
				}
			}
		case 0xE:
			switch CH {
			case 0x7:
				switch CL {
				case 0x3, 0x4, 0x5, 0x6, 0x7:
					// 7E ?? 73, 7E ?? 74, 7E ?? 75, 7E ?? 76, 7E ?? 77
					inst.Opcode = borBxorBandBld(CL, DH)
					if inst.Opcode != ".word" {
						inst.TotalBytes = 4
						inst.AddressingMode = AbsoluteAddress
					}
				}
			}
		case 0xF:
			switch CH {
			case 0x6:
				switch CL {
				case 0x0, 0x1, 0x2:
					// 7F ?? 60, 7F ?? 61, 7F ?? 62
					inst.Opcode = bSetBNotBClr(CL)
					if inst.Opcode != ".word" {
						inst.AddressingMode = AbsoluteAddress
						inst.TotalBytes = 4
					}
				case 0x7:
					// 7F ?? 67
					inst.TotalBytes = 4
					if DH&0x8 == 0 {
						inst.AddressingMode = AbsoluteAddress
						inst.Opcode = "bst"
					} else {
						inst.AddressingMode = AbsoluteAddress
						inst.Opcode = "bist"
					}
				}
			case 0x7:
				switch CL {
				case 0x0, 0x1, 0x2:
					if DH < 0x8 {
						inst.Opcode = bSetBNotBClr(CL)
						if inst.Opcode != ".word" {
							inst.AddressingMode = AbsoluteAddress
							inst.TotalBytes = 4
						}
					}
				}
			}
		}
	}

	return inst
}

func table234(inst Inst, bytes []byte) Inst {
	AH := bytes[0] >> 4
	AL := bytes[0] & 0x0F
	if AH == 0x6 && AL == 0xA {
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
						inst.TotalBytes = 6
						inst.AddressingMode = AbsoluteAddress
						inst.Opcode = "btst"
					}
				case 0x7:
					switch EL {
					case 0x3, 0x4, 0x5, 0x6, 0x7:
						inst.Opcode = borBxorBandBld(EL, FH)
						if inst.Opcode != ".word" {
							inst.TotalBytes = 6
							inst.AddressingMode = AbsoluteAddress
						}
					}
				}
			case 0x8:
				switch EH {
				case 0x6:
					switch EL {
					case 0x0, 0x1, 0x2:
						inst.Opcode = bSetBNotBClr(EL)
						if inst.Opcode != ".word" {
							inst.TotalBytes = 6
							inst.AddressingMode = AbsoluteAddress
						}
					case 0x7:
						if FH&0x8 == 0 {
							inst.TotalBytes = 6
							inst.AddressingMode = AbsoluteAddress
							inst.Opcode = "bst"
						} else {
							inst.TotalBytes = 6
							inst.AddressingMode = AbsoluteAddress
							inst.Opcode = "bist"
						}
					}
				case 0x7:
					switch EL {
					case 0x0, 0x1, 0x2:
						// 6A 18 ?? ?? 72 ?0
						if FH < 0x8 {
							inst.Opcode = bSetBNotBClr(EL)
							if inst.Opcode != ".word" {
								inst.AddressingMode = AbsoluteAddress
								inst.TotalBytes = 6
							}
						}
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
					}
				case 0x7:
					switch GL {
					case 0x3, 0x4, 0x5, 0x6, 0x7:
						// if HH < 0x8 {
						inst.Opcode = borBxorBandBld(GL, HH)
						if inst.Opcode != ".word" {
							inst.AddressingMode = AbsoluteAddress
							inst.TotalBytes = 8
						}
						// }
					}
				}
			case 0x8:
				switch GH {
				case 0x6:
					switch GL {
					case 0x0, 0x1, 0x2:
						inst.Opcode = bSetBNotBClr(GL)
						if inst.Opcode != ".word" {
							inst.AddressingMode = AbsoluteAddress
							inst.TotalBytes = 8
						}
					case 0x7:
						inst.TotalBytes = 8
						inst.AddressingMode = AbsoluteAddress
						if HH&0x8 == 0 {
							inst.Opcode = "bst"
						} else {
							inst.Opcode = "bist"
						}
					}
				case 0x7:
					switch GL {
					case 0x0, 0x1, 0x2:
						// 6A 38 ?? ?? ?? ?? 72 ?0
						if HH < 0x8 {
							inst.Opcode = bSetBNotBClr(GL)
							if inst.Opcode != ".word" {
								inst.AddressingMode = AbsoluteAddress
								inst.TotalBytes = 8
							}
						}
					}
				}
			}
		}
	}
	return inst
}

func branches(b byte) string {
	branchMap := map[byte]string{
		0x0: "bra", // bra in the manual
		0x1: "brn", // brn in the manual
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

func bSetBNotBClr(b byte) string {
	bsetbnotbclrMap := map[byte]string{
		0x0: "bset",
		0x1: "bnot",
		0x2: "bclr",
	}
	return bsetbnotbclrMap[b]
}

func orXorAnd(b byte, shift bool) string {
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

func borBxorBandBld(b byte, CB byte) string {
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
	} else {
		if val, ok := m[b]; ok {
			return val
		}
	}
	return ".word"
}
