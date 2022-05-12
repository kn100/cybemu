// disassembler contains utilities to disassemble a binary for the Renesas
// (Previously Hitachi) h8s/2000 class CPUs. It is being developed as a part of
// a larger project to emulate the Cybiko Classic. It roughly follows the layout
// of the tables found in "H8S/2600 Series, H8S/2000 Series Software Manual
// Rev.4.00 2006.02", Section 2.5 Table 2.3 (p274-277).
package disassembler

import (
	"github.com/kn100/cybemu/addressingmode"
	"github.com/kn100/cybemu/instruction"
	"github.com/kn100/cybemu/opcode"
	"github.com/kn100/cybemu/size"
)

// Disassemble takes a sequence of bytes from a compiled binary and disassembles
// them. It will return a slice of instruction.Inst.
func Disassemble(bytes []byte) []instruction.Inst {
	instructions := []instruction.Inst{}

	i := 0
	for i < len(bytes) {
		inst := Decode(bytes[i:])
		inst.Pos = i

		for b := 0; b < inst.TotalBytes; b++ {
			inst.Bytes = append(inst.Bytes, bytes[i+b])
		}
		inst.DetermineOperandTypeAndSetData()
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
func Decode(bytes []byte) instruction.Inst {
	AH := bytes[0] >> 4
	AL := bytes[0] & 0x0F
	BH := bytes[1] >> 4
	inst := instruction.Inst{
		Opcode:     opcode.Invalid, // Everything is a word until proven otherwise
		TotalBytes: 2,
	}
	switch AH {
	case 0x0:
		switch AL {
		case 0x0:
			BL := bytes[1] & 0x0F
			if BH == 0x0 && BL == 0x0 {
				inst.Opcode = opcode.Nop
			}
		case 0x1:
			// 01
			inst = table232(inst, bytes[:8])
		case 0x2:
			BH := bytes[1] >> 4
			if BH == 0 || BH == 1 {
				inst.AddressingMode = addressingmode.RegisterDirect
				inst.BWL = size.Byte
				inst.Opcode = opcode.Stc
			}
		case 0x3:
			if BH == 0x0 {
				inst.AddressingMode = addressingmode.RegisterDirect
				inst.BWL = size.Byte
				inst.Opcode = opcode.Ldc
			} else if BH == 0x1 {
				inst.AddressingMode = addressingmode.RegisterDirect
				inst.BWL = size.Byte
				inst.Opcode = opcode.Ldc
			}
		case 0x4:
			inst.AddressingMode = addressingmode.Immediate
			inst.Opcode = opcode.Orc
		case 0x5:
			inst.AddressingMode = addressingmode.Immediate
			inst.Opcode = opcode.Xorc
		case 0x6:
			inst.AddressingMode = addressingmode.Immediate
			inst.Opcode = opcode.Andc
		case 0x7:
			inst.AddressingMode = addressingmode.Immediate
			inst.BWL = size.Byte
			inst.Opcode = opcode.Ldc
		case 0x8:
			inst.BWL = size.Byte
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.Opcode = opcode.Add
		case 0x9:
			inst.BWL = size.Word
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.Opcode = opcode.Add
		case 0xA, 0xB, 0xF:
			inst = table232(inst, bytes[:8])
		case 0xC:
			inst.Opcode = opcode.Mov
			inst.BWL = size.Byte
			inst.AddressingMode = addressingmode.RegisterDirect
		case 0xD:
			inst.Opcode = opcode.Mov
			inst.BWL = size.Word
			inst.AddressingMode = addressingmode.RegisterDirect
		case 0xE:
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.Opcode = opcode.Addx
		}
	case 0x1:
		switch AL {
		case 0x0, 0x1, 0x2, 0x3, 0x7, 0xA, 0xB, 0xF:
			inst = table232(inst, bytes[:8])
		case 0x4, 0x5, 0x6:
			inst.BWL = size.Byte
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.Opcode = orXorAnd(AL)
		case 0x8:
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.BWL = size.Byte
			inst.Opcode = opcode.Sub
		case 0x9:
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.BWL = size.Word
			inst.Opcode = opcode.Sub
		case 0xC:
			inst.BWL = size.Byte
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.Opcode = opcode.Cmp
		case 0xD:
			inst.BWL = size.Word
			inst.Opcode = opcode.Cmp
		case 0xE:
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.Opcode = opcode.Subx
		}
	case 0x2:
		inst.Opcode = opcode.Mov
		inst.AddressingMode = addressingmode.AbsoluteAddress
		inst.BWL = size.Byte
	case 0x3:
		inst.BWL = size.Byte
		inst.AddressingMode = addressingmode.AbsoluteAddress
		inst.Opcode = opcode.Mov
	case 0x4:
		// BL := bytes[1] & 0x0F
		// BL must be even HEX number according to spec...
		// MAMEs decompiler doesn't seem to validate this detail, so I won't either and assume I don't know what I am doing
		// BL == 0x0 || BL == 0x2 || BL == 0x4 || BL == 0x6 || BL == 0x8 || BL == 0xA || BL == 0xC || BL == 0xE
		// In execution, the H8S2000 will set last bit to 0, which results in going to the previous instruction to the requested one.
		inst.AddressingMode = addressingmode.ProgramCounterRelative
		inst.Opcode = branches(AL)
		inst.OperandSize = 8
	case 0x5:
		switch AL {
		case 0x0:
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.BWL = size.Byte
			inst.Opcode = opcode.Mulxu
		case 0x2:
			BL := bytes[1] & 0x0F
			if BL < 0x8 {
				inst.AddressingMode = addressingmode.RegisterDirect
				inst.BWL = size.Word
				inst.Opcode = opcode.Mulxu
			}
		case 0x1:
			inst.BWL = size.Byte
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.Opcode = opcode.Divxu
		case 0x3:
			BL := bytes[1] & 0x0F
			if BL < 0x8 {
				inst.BWL = size.Word
				inst.AddressingMode = addressingmode.RegisterDirect
				inst.Opcode = opcode.Divxu
			}

		case 0x4:
			BL := bytes[1] & 0x0F
			if BH == 0x7 && BL == 0x0 {
				inst.Opcode = opcode.Rts
			}

		case 0x5:
			inst.AddressingMode = addressingmode.ProgramCounterRelative
			inst.Opcode = opcode.Bsr
		case 0x6:
			BL := bytes[1] & 0x0F
			if BH == 0x7 && BL == 0x0 {
				inst.Opcode = opcode.Rte
			}
		case 0x7:
			BL := bytes[1] & 0x0F
			if BH < 0x4 && BL == 0 {
				inst.AddressingMode = addressingmode.RegisterDirect
				inst.Opcode = opcode.Trapa
			}
		case 0x8:
			inst = table232(inst, bytes[:8])
		case 0x9:
			BL := bytes[1] & 0x0F
			if BH < 0x8 && BL == 0 {
				inst.AddressingMode = addressingmode.RegisterIndirect
				inst.Opcode = opcode.Jmp
			}
		case 0xA:
			inst.AddressingMode = addressingmode.AbsoluteAddress
			inst.TotalBytes = 4
			inst.Opcode = opcode.Jmp
			inst.OperandSize = 24
		case 0xB:
			inst.AddressingMode = addressingmode.MemoryIndirect
			inst.Opcode = opcode.Jmp
			inst.OperandSize = 8
		case 0xC:
			BL := bytes[1] & 0x0F
			if BH == 0x0 && BL == 0x0 {
				inst.AddressingMode = addressingmode.ProgramCounterRelative
				inst.TotalBytes = 4
				inst.Opcode = opcode.Bsr
				inst.OperandSize = 16
			}
		case 0xD:
			BL := bytes[1] & 0x0F
			if BH < 0x8 && BL == 0 {
				inst.AddressingMode = addressingmode.RegisterIndirect
				inst.Opcode = opcode.Jsr
			}
		case 0xE:
			inst.AddressingMode = addressingmode.AbsoluteAddress
			inst.TotalBytes = 4
			inst.Opcode = opcode.Jsr
			inst.OperandSize = 24
		case 0xF:
			inst.AddressingMode = addressingmode.MemoryIndirect
			inst.Opcode = opcode.Jsr
			inst.OperandSize = 8
		}
	case 0x6:
		switch AL {
		case 0x0, 0x1, 0x2:
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.Opcode = bSetBNotBClr(AL)
		case 0x3:
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.Opcode = opcode.Btst
		case 0x4, 0x5, 0x6:
			inst.BWL = size.Word
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.Opcode = orXorAnd(AL)
		case 0x7:
			BH := bytes[1] >> 4
			inst.AddressingMode = addressingmode.RegisterDirect
			if BH&0x8 == 0 {
				inst.Opcode = opcode.Bst
			} else {
				inst.Opcode = opcode.Bist
			}
		case 0x8:
			inst.AddressingMode = addressingmode.RegisterIndirect
			inst.BWL = size.Byte
			inst.Opcode = opcode.Mov
		case 0xC:
			inst.BWL = size.Byte
			if BH > 0x7 {
				inst.AddressingMode = addressingmode.RegisterIndirectWithPreDecrement
				inst.Opcode = opcode.Mov
			} else {
				inst.AddressingMode = addressingmode.RegisterIndirectWithPostIncrement
				inst.Opcode = opcode.Mov
			}
		case 0xE:
			inst.Opcode = opcode.Mov
			inst.TotalBytes = 4
			inst.AddressingMode = addressingmode.RegisterIndirectWithDisplacement
			inst.BWL = size.Byte
		case 0x9:
			inst.BWL = size.Word
			inst.Opcode = opcode.Mov
			inst.AddressingMode = addressingmode.RegisterIndirect
		case 0xB:
			inst.BWL = size.Word
			inst.AddressingMode = addressingmode.AbsoluteAddress
			if BH == 0x0 || BH == 0x8 {
				inst.TotalBytes = 4
				inst.OperandSize = 16
				inst.Opcode = opcode.Mov
			} else if BH == 0x2 || BH == 0xA {
				inst.TotalBytes = 6
				inst.OperandSize = 32
				inst.Opcode = opcode.Mov
			}
		case 0xD:
			inst.BWL = size.Word
			inst.AddressingMode = addressingmode.RegisterIndirectWithPostIncrement
			switch BH {
			case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6:
				inst.Opcode = opcode.Mov
			case 0x7:
				inst.Opcode = opcode.Pop
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE:
				inst.Opcode = opcode.Mov
			case 0xF:
				inst.Opcode = opcode.Push
			}
		case 0xF:
			inst.OperandSize = 16
			inst.TotalBytes = 4
			inst.AddressingMode = addressingmode.RegisterIndirectWithDisplacement
			inst.BWL = size.Word
			inst.Opcode = opcode.Mov
		case 0xA:
			inst = table232(inst, bytes[:8])
		}
	case 0x7:
		switch AL {
		case 0x0, 0x1, 0x2:
			if BH < 0x8 {
				inst.AddressingMode = addressingmode.RegisterDirect
				inst.Opcode = bSetBNotBClr(AL)
			}

		case 0x3, 0x4, 0x5, 0x6, 0x7:
			inst.AddressingMode = addressingmode.RegisterDirect
			inst.Opcode = borBxorBandBld(AL, BH)
		case 0x8:
			BL := bytes[1] & 0x0F
			CH := bytes[2] >> 4
			CL := bytes[2] & 0x0F
			DH := bytes[3] >> 4
			if BH < 0x8 && BL == 0x0 && CH == 0x6 && (DH == 0x2 || DH == 0xA) {
				inst.AddressingMode = addressingmode.RegisterIndirectWithDisplacement
				inst.TotalBytes = 8
				if CL == 0xA {
					inst.Opcode = opcode.Mov
					inst.BWL = size.Byte
				} else if CL == 0xB {
					inst.Opcode = opcode.Mov
					inst.BWL = size.Word
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
			inst.TotalBytes = 4
			if CH == 0x5 && CL == 0x9 && DH == 0x8 && DL == 0xF {
				if BH == 5 && BL == 0xC {
					inst.BWL = size.Byte
					inst.Opcode = opcode.Eepmov
				} else if BH == 0xD && BL == 0x4 {
					inst.BWL = size.Word
					inst.Opcode = opcode.Eepmov
				}
			}
		}
	case 0x8:
		inst.BWL = size.Byte
		inst.AddressingMode = addressingmode.Immediate
		inst.Opcode = opcode.Add
	case 0x9:
		inst.AddressingMode = addressingmode.Immediate
		inst.Opcode = opcode.Addx
	case 0xA:
		inst.BWL = size.Byte
		inst.AddressingMode = addressingmode.Immediate
		inst.Opcode = opcode.Cmp
	case 0xB:
		inst.AddressingMode = addressingmode.Immediate
		inst.Opcode = opcode.Subx
	case 0xC, 0xD, 0xE:
		inst.BWL = size.Byte
		inst.AddressingMode = addressingmode.Immediate
		inst.Opcode = orXorAnd(AH - 8)
	case 0xF:
		inst.BWL = size.Byte
		inst.AddressingMode = addressingmode.Immediate
		inst.Opcode = opcode.Mov
	}

	return inst
}

func table232(inst instruction.Inst, bytes []byte) instruction.Inst {
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
		inst.BWL = size.Longword
		// 01 00 69
		if CH == 0x6 && CL == 0x9 && DL < 0x8 {
			inst.TotalBytes = 4
			inst.AddressingMode = addressingmode.RegisterIndirect
			inst.Opcode = opcode.Mov
			// 01 00 6F
		} else if CH == 0x6 && CL == 0xF {
			inst.TotalBytes = 6
			inst.AddressingMode = addressingmode.RegisterIndirectWithDisplacement
			inst.Opcode = opcode.Mov
		} else if CH == 0x7 && CL == 0x8 && DH < 0x8 && DL == 0x0 && EH == 0x6 && EL == 0xB && (FH == 0x2 || FH == 0xA) && FL < 0x8 {
			// 01 00 78 ?0 6B A? ?? ?? ?? ??
			// 01 00 78 ?0 6B 2? ?? ?? ?? ??
			// TODO: Potential bug in unidasm? should be checking `DH < 0x7` (p142/322)
			// Turns out the real hardware doesn't care what DH is so just set it to zero.
			inst.TotalBytes = 10
			inst.AddressingMode = addressingmode.RegisterIndirectWithDisplacement
			inst.Opcode = opcode.Mov
		} else if CH == 0x6 && CL == 0xB {
			inst.AddressingMode = addressingmode.AbsoluteAddress
			if DH == 0x0 || DH == 0x8 {
				inst.TotalBytes = 6
				inst.Opcode = opcode.Mov
			} else if DH == 0x2 || DH == 0xA {
				inst.TotalBytes = 8
				inst.Opcode = opcode.Mov
			}
		} else if CH == 0x6 && CL == 0xD && DH == 0x7 && DL < 0x8 {
			inst.TotalBytes = 4
			inst.AddressingMode = addressingmode.RegisterIndirectWithPreDecrement
			inst.Opcode = opcode.Pop
		} else if CH == 0x6 && CL == 0xD && DH < 0x7 && DL < 0x8 {
			inst.TotalBytes = 4
			inst.AddressingMode = addressingmode.RegisterIndirectWithPreDecrement
			inst.Opcode = opcode.Mov
		} else if CH == 0x6 && CL == 0xD && DH == 0xF && DL < 0x8 {
			inst.TotalBytes = 4
			inst.AddressingMode = addressingmode.RegisterIndirectWithPreDecrement
			inst.Opcode = opcode.Push
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
					inst.BWL = size.Longword
					if DH == 0x7 {
						inst.Opcode = opcode.Ldm
					} else if DH == 0xF {
						inst.Opcode = opcode.Stm
					}
				}
			case 0x4:
				switch BL {
				case 0x0:
					if CH == 0x6 && CL == 0x9 {
						// 01 40 69
						inst.TotalBytes = 4
						inst.AddressingMode = addressingmode.RegisterIndirect
						inst.BWL = size.Word
						if DH < 0x8 {
							inst.Opcode = opcode.Ldc
						} else {
							inst.Opcode = opcode.Stc
						}
					} else if CH == 0x6 && CL == 0xF {
						// 01 40 6F
						inst.TotalBytes = 6
						inst.AddressingMode = addressingmode.RegisterIndirectWithDisplacement
						inst.BWL = size.Word
						if DH < 0x8 {
							inst.Opcode = opcode.Ldc
						} else {
							inst.Opcode = opcode.Stc
						}
					} else if CH == 0x7 && CL == 0x8 && DH < 0x8 && EH == 0x6 && EL == 0xB && (FH == 0xA || FH == 0x2) && FL == 0x0 {
						// 01 40 78 ?? 6B A0
						inst.TotalBytes = 10
						inst.AddressingMode = addressingmode.RegisterIndirectWithDisplacement
						inst.BWL = size.Word
						if FH < 0xA {
							inst.Opcode = opcode.Ldc
						} else {
							inst.Opcode = opcode.Stc
						}
					} else if CH == 0x6 && CL == 0xD {
						// 01 40 6D
						inst.TotalBytes = 4
						inst.BWL = size.Word
						if DH > 0x7 {
							inst.AddressingMode = addressingmode.RegisterIndirectWithPreDecrement
							inst.Opcode = opcode.Stc
						} else {
							inst.AddressingMode = addressingmode.RegisterIndirectWithPostIncrement
							inst.Opcode = opcode.Ldc
						}
					} else if CH == 0x6 && CL == 0xB && DL == 0x0 {
						inst.AddressingMode = addressingmode.AbsoluteAddress
						inst.BWL = size.Word
						switch DH {
						case 0x0:
							// 01 40 6B 00
							inst.TotalBytes = 6
							inst.Opcode = opcode.Ldc
							inst.OperandSize = 16
						case 0x2:
							// 01 40 6B 20
							inst.TotalBytes = 8
							inst.Opcode = opcode.Ldc
							inst.OperandSize = 32
						case 0x8:
							// 01 40 6B 80
							inst.TotalBytes = 6
							inst.Opcode = opcode.Stc
							inst.OperandSize = 16
						case 0xA:
							// 01 40 6B A0
							inst.TotalBytes = 8
							inst.Opcode = opcode.Stc
							inst.OperandSize = 32
						}
					}
				case 0x1:
					if CH == 0x6 && CL == 0x9 {
						// 01 41 69
						inst.TotalBytes = 4
						inst.AddressingMode = addressingmode.RegisterIndirect
						inst.BWL = size.Word
						if DH < 0x8 {
							inst.Opcode = opcode.Ldc
						} else {
							inst.Opcode = opcode.Stc
						}
					} else if CH == 0x6 && CL == 0xF {
						// 01 41 6F
						inst.TotalBytes = 6
						inst.AddressingMode = addressingmode.RegisterIndirectWithDisplacement
						inst.BWL = size.Word
						if DH < 0x8 {
							inst.Opcode = opcode.Ldc
						} else {
							inst.Opcode = opcode.Stc
						}
					} else if CH == 0x7 && CL == 0x8 && DH < 0x8 && EH == 0x6 && EL == 0xB && (FH == 0xA || FH == 0x2) && FL == 0x0 {
						// 01 41 78 ?? 6B A0
						inst.TotalBytes = 10
						inst.AddressingMode = addressingmode.RegisterIndirectWithDisplacement
						inst.BWL = size.Word
						if FH < 0xA {
							inst.Opcode = opcode.Ldc
						} else {
							inst.Opcode = opcode.Stc
						}
					} else if CH == 0x6 && CL == 0xD {
						// 01 41 6D
						if DH > 0x7 {
							inst.TotalBytes = 4
							inst.BWL = size.Word
							inst.AddressingMode = addressingmode.RegisterIndirectWithPreDecrement
							inst.Opcode = opcode.Stc
						} else {
							inst.TotalBytes = 4
							inst.BWL = size.Word
							inst.AddressingMode = addressingmode.RegisterIndirectWithPostIncrement
							inst.Opcode = opcode.Ldc
						}
					} else if CH == 0x6 && CL == 0xB && DL == 0x0 {
						switch DH {
						case 0x0:
							// 01 41 6B 00
							inst.TotalBytes = 6
							inst.AddressingMode = addressingmode.AbsoluteAddress
							inst.BWL = size.Word
							inst.Opcode = opcode.Ldc
							inst.OperandSize = 16
						case 0x2:
							// 01 41 6B 20
							inst.TotalBytes = 8
							inst.AddressingMode = addressingmode.AbsoluteAddress
							inst.BWL = size.Word
							inst.Opcode = opcode.Ldc
							inst.OperandSize = 32
						case 0x8:
							// 01 41 6B 80
							inst.TotalBytes = 6
							inst.AddressingMode = addressingmode.AbsoluteAddress
							inst.BWL = size.Word
							inst.Opcode = opcode.Stc
							inst.OperandSize = 16
						case 0xA:
							// 01 41 6B A0
							inst.TotalBytes = 8
							inst.AddressingMode = addressingmode.AbsoluteAddress
							inst.BWL = size.Word
							inst.Opcode = opcode.Stc
							inst.OperandSize = 32
						}
					} else if CH == 0x0 && CL == 0x7 {
						// 01 41 07
						inst.TotalBytes = 4
						inst.AddressingMode = addressingmode.Immediate
						inst.BWL = size.Byte
						inst.Opcode = opcode.Ldc
					} else if CH == 0x0 && CL == 0x6 {
						// 01 41 06
						inst.TotalBytes = 4
						inst.AddressingMode = addressingmode.Immediate
						inst.Opcode = opcode.Andc
					} else if CH == 0x0 && CL == 0x5 {
						// 01 41 05
						inst.TotalBytes = 4
						inst.AddressingMode = addressingmode.Immediate
						inst.Opcode = opcode.Xorc

					} else if CH == 0x0 && CL == 0x4 {
						// 01 41 04
						inst.TotalBytes = 4
						inst.AddressingMode = addressingmode.Immediate
						inst.Opcode = opcode.Orc
					}
				}

			case 0x8:
				if BL == 0x0 {
					inst.Opcode = opcode.Sleep
				}

			case 0xC, 0xD, 0xF:
				inst = table233(bytes)
			case 0xE:
				if BL == 0x0 && CH == 0x7 && CL == 0xB && DH < 0x8 && DL == 0xC {
					inst.TotalBytes = 4
					inst.AddressingMode = addressingmode.RegisterIndirect
					inst.Opcode = opcode.Tas
				}
			}
		case 0xA:
			inst.AddressingMode = addressingmode.RegisterDirect
			switch BH {
			case 0x0:
				inst.BWL = size.Byte
				inst.Opcode = opcode.Inc
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				if BL < 0x8 {
					inst.Opcode = opcode.Add
					inst.BWL = size.Longword
				}

			}
		case 0xB:
			inst.AddressingMode = addressingmode.RegisterDirect
			switch BH {
			case 0x0, 0x8, 0x9:
				// 0B B0, 0B B8, 0B B9
				if BL < 0x8 {
					inst.Opcode = opcode.Adds
				}
			case 0x5, 0xD:
				inst.BWL = size.Word
				inst.Opcode = opcode.Inc
			case 0x7, 0xF:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Inc
				}
			}
		case 0xF:
			if BH > 0x7 && BL < 0x8 {
				inst.BWL = size.Longword
				inst.Opcode = opcode.Mov
			} else if BH == 0x0 {
				inst.Opcode = opcode.Daa
			}
		}
	case 0x1:
		inst.AddressingMode = addressingmode.RegisterDirect
		switch AL {
		case 0x0:
			switch BH {
			case 0x0, 0x4:

				inst.BWL = size.Byte
				inst.Opcode = opcode.Shll
			case 0x1, 0x5:
				inst.BWL = size.Word
				inst.Opcode = opcode.Shll
			case 0x3, 0x7:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Shll
				}
			case 0x8, 0xC:
				inst.BWL = size.Byte
				inst.Opcode = opcode.Shal
			case 0x9, 0xD:
				inst.BWL = size.Word
				inst.Opcode = opcode.Shal
			case 0xB, 0xF:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Shal
				}
			}
		case 0x1:
			switch BH {
			case 0x0, 0x4:
				inst.BWL = size.Byte
				inst.Opcode = opcode.Shlr
			case 0x1, 0x5:
				inst.BWL = size.Word
				inst.Opcode = opcode.Shlr
			case 0x3, 0x7:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Shlr
				}
			case 0x8, 0xC:
				inst.BWL = size.Byte
				inst.Opcode = opcode.Shar
			case 0x9, 0xD:
				inst.BWL = size.Word
				inst.Opcode = opcode.Shar
			case 0xB, 0xF:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Shar
				}
			}
		case 0x2:
			switch BH {
			case 0x0, 0x4:
				inst.BWL = size.Byte
				inst.Opcode = opcode.Rotxl
			case 0x1, 0x5:
				inst.BWL = size.Word
				inst.Opcode = opcode.Rotxl
			case 0x3, 0x7:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Rotxl
				}
			case 0x8, 0xC:
				inst.BWL = size.Byte
				inst.Opcode = opcode.Rotl
			case 0x9, 0xD:
				inst.BWL = size.Word
				inst.Opcode = opcode.Rotl
			case 0xB, 0xF:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Rotl
				}
			}
		case 0x3:
			switch BH {
			case 0x0, 0x4:
				inst.BWL = size.Byte
				inst.Opcode = opcode.Rotxr
			case 0x1, 0x5:
				inst.BWL = size.Word
				inst.Opcode = opcode.Rotxr
			case 0x3, 0x7:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Rotxr
				}
			case 0x8, 0xC:
				inst.BWL = size.Byte
				inst.Opcode = opcode.Rotr
			case 0x9, 0xD:
				inst.BWL = size.Word
				inst.Opcode = opcode.Rotr
			case 0xB, 0xF:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Rotr
				}
			}
		case 0x7:
			switch BH {
			case 0x0:
				inst.Opcode = opcode.Not
				inst.BWL = size.Byte
			case 0x1:
				inst.Opcode = opcode.Not
				inst.BWL = size.Word
			case 0x3:
				if BL < 0x8 {
					inst.Opcode = opcode.Not
					inst.BWL = size.Longword
				}
			case 0x5:
				inst.BWL = size.Word
				inst.Opcode = opcode.Extu
			case 0x7:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Extu
				}
			case 0x8:
				inst.Opcode = opcode.Neg
				inst.BWL = size.Byte
			case 0x9:
				inst.Opcode = opcode.Neg
				inst.BWL = size.Word
			case 0xB:
				if BL < 0x8 {
					inst.Opcode = opcode.Neg
					inst.BWL = size.Longword
				}
			case 0xD:
				inst.BWL = size.Word
				inst.Opcode = opcode.Exts
			case 0xF:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Exts
				}
			}
		case 0xA:
			switch BH {
			case 0x0:
				inst.BWL = size.Byte
				inst.Opcode = opcode.Dec
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Sub
				}
			}
		case 0xB:
			switch BH {
			case 0x0, 0x8, 0x9:
				if BL < 0x8 {
					inst.Opcode = opcode.Subs
				}
			case 0x5, 0xD:
				inst.BWL = size.Word
				inst.Opcode = opcode.Dec
			case 0x7, 0xF:
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Dec
				}
			}
		case 0xF:
			switch BH {
			case 0x0:
				// 1F 0?
				inst.Opcode = opcode.Das
			case 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF:
				// 1F 8?, 1F 9?, 1F A?, 1F B?, 1F C?, 1F D?, 1F E?, 1F F?
				if BL < 0x8 {
					inst.BWL = size.Longword
					inst.Opcode = opcode.Cmp
				}
			}
		}
	case 0x5:
		inst.AddressingMode = addressingmode.ProgramCounterRelative
		inst.TotalBytes = 4
		if AL == 0x8 && BL == 0x0 {
			// TODO: CPU does weird things with odd dest
			inst.Opcode = branches(BH)
		}
	case 0x6:
		inst.AddressingMode = addressingmode.AbsoluteAddress
		if AL == 0xA {
			switch BH {
			case 0x0:
				inst.Opcode = opcode.Mov
				inst.TotalBytes = 4
				inst.BWL = size.Byte
			case 0x2:
				inst.Opcode = opcode.Mov
				inst.TotalBytes = 6
				inst.BWL = size.Byte
			case 0x8:
				inst.TotalBytes = 4
				inst.BWL = size.Byte
				inst.Opcode = opcode.Mov
			case 0xA:
				inst.TotalBytes = 6
				inst.BWL = size.Byte
				inst.Opcode = opcode.Mov
			case 0x1, 0x3:
				inst = table234(bytes)
			case 0x4:
				inst.AddressingMode = addressingmode.AbsoluteAddress
				inst.TotalBytes = 4
				inst.Opcode = opcode.Movfpe
			case 0xC:
				inst.AddressingMode = addressingmode.AbsoluteAddress
				inst.TotalBytes = 4
				inst.Opcode = opcode.Movtpe
			}
		}
	case 0x7:
		inst.TotalBytes = 4
		inst.AddressingMode = addressingmode.Immediate
		switch AL {
		case 0x9:
			inst.BWL = size.Word
			switch BH {
			case 0x0:

				inst.Opcode = opcode.Mov
			case 0x1:
				inst.Opcode = opcode.Add
			case 0x2:
				inst.Opcode = opcode.Cmp
			case 0x3:
				inst.Opcode = opcode.Sub
			case 0x4, 0x5, 0x6:
				inst.Opcode = orXorAnd(BH)
			}
		case 0xA:
			inst.AddressingMode = addressingmode.Immediate
			inst.TotalBytes = 6
			inst.BWL = size.Longword
			if BL < 0x8 {
				switch BH {
				case 0x0:
					inst.Opcode = opcode.Mov
				case 0x1:
					inst.Opcode = opcode.Add
				case 0x2:
					inst.Opcode = opcode.Cmp
				case 0x3:
					inst.Opcode = opcode.Sub
				case 0x4, 0x5, 0x6:
					inst.Opcode = orXorAnd(BH)
				}
			}
		case 0xC:
			if BH < 0x8 && BL == 0x0 && CH == 0x7 && CL == 0x3 && DH < 0x8 && DL == 0x0 {
				// 7C ?0 73 ?0
				inst.AddressingMode = addressingmode.RegisterIndirect
				inst.TotalBytes = 4
				inst.Opcode = opcode.Btst
			} else {
				// 7C
				inst = table233(bytes)
			}

		case 0xD:
			if BH < 0x8 && BL == 0x0 && CH == 0x7 && CL == 0x0 && DH < 0x8 && DL == 0x0 {
				inst.AddressingMode = addressingmode.RegisterIndirect
				inst.TotalBytes = 4
				inst.Opcode = opcode.Bset
			} else {
				inst = table233(bytes)
			}

		case 0xE:
			inst.AddressingMode = addressingmode.AbsoluteAddress
			inst.TotalBytes = 4
			switch CH {
			case 0x6:
				// 7E ?? 63 ?0
				if CL == 0x3 && DL == 0x0 {
					inst.Opcode = opcode.Btst
				}

			case 0x7:
				// 7E ?? 73 ?0
				if CL == 0x3 && DH < 0x8 && DL == 0x0 {
					inst.Opcode = opcode.Btst
				} else {
					inst = table233(bytes)
				}

			}
		case 0xF:
			inst = table233(bytes)
		}
	}
	if inst.Opcode == opcode.Invalid {
		return instruction.Inst{Opcode: opcode.Invalid, TotalBytes: 2}
	}
	return inst
}

func table233(bytes []byte) instruction.Inst {
	inst := instruction.Inst{
		Opcode:     opcode.Invalid, // Everything is a word until proven otherwise
		TotalBytes: 2,
	}

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
		inst.AddressingMode = addressingmode.RegisterDirect
		inst.TotalBytes = 4
		if AL == 0x1 {
			switch BH {
			case 0xC:
				if BL == 0x0 && CH == 0x5 {
					if CL == 0x0 {
						inst.BWL = size.Byte
						inst.Opcode = opcode.Mulxs
					} else if CL == 0x2 && DL < 0x8 {
						inst.BWL = size.Word
						inst.Opcode = opcode.Mulxs
					}
				}

			case 0xD:
				if BL == 0x0 && CH == 0x5 {
					switch CL {
					case 0x1:
						inst.BWL = size.Byte
						inst.Opcode = opcode.Divxs
					case 0x3:
						if DL < 0x8 {
							inst.BWL = size.Word
							inst.Opcode = opcode.Divxs
						}
					}
				}

			case 0xF:
				inst.BWL = size.Longword
				if BL == 0x0 && CH == 0x6 && (CL == 0x4 || CL == 0x5 || CL == 0x6) && DH < 0x8 && DL < 0x8 {
					// 01 f0 64 ??, 01 f0 65 ??, 01 f0 66 ??
					inst.Opcode = orXorAnd(CL)
				}
			}
		}

	case 0x7:
		switch AL {
		case 0xC:
			inst.AddressingMode = addressingmode.RegisterIndirect
			inst.TotalBytes = 4
			if BH < 0x8 && BL == 0x0 {
				switch CH {
				case 0x6:
					if CL == 0x3 {
						inst.Opcode = opcode.Btst
					}

				case 0x7:
					switch CL {
					case 0x3, 0x4, 0x5, 0x6, 0x7:
						// 7C ?0 73, 7C ?0 74, 7C ?0 75, 7C ?0 76, 7C ?0 77
						inst.Opcode = borBxorBandBld(CL, DH)
					}
				}
			}

		case 0xD:
			inst.AddressingMode = addressingmode.RegisterIndirect
			inst.TotalBytes = 4
			if BH < 0x8 && BL == 0x0 {
				// 7D ?0
				switch CH {
				case 0x6:
					switch CL {
					case 0x0, 0x1, 0x2:
						// 7D ?0 60, 7D ?0 61, 7D ?0 62
						inst.Opcode = bSetBNotBClr(CL)
					case 0x7:
						// 7D ?0 67
						if DH&0x8 == 0 {
							inst.Opcode = opcode.Bst
						} else {
							inst.Opcode = opcode.Bist
						}
					}
				case 0x7:
					switch CL {
					case 0x0, 0x1, 0x2:
						if DH < 0x8 {
							// 7D ?0 70 ?0, 7D ?0 71 ?0, 7D ?0 72 ?0
							inst.Opcode = bSetBNotBClr(CL)
						}
					}
				}
			}
		case 0xE:
			inst.TotalBytes = 4
			inst.AddressingMode = addressingmode.AbsoluteAddress
			if CH == 0x7 && CL > 0x2 && CL < 0x8 {
				// 7E ?? 73, 7E ?? 74, 7E ?? 75, 7E ?? 76, 7E ?? 77
				inst.Opcode = borBxorBandBld(CL, DH)
			}
		case 0xF:
			inst.AddressingMode = addressingmode.AbsoluteAddress
			inst.TotalBytes = 4
			switch CH {
			case 0x6:
				switch CL {
				case 0x0, 0x1, 0x2:
					// 7F ?? 60, 7F ?? 61, 7F ?? 62
					inst.Opcode = bSetBNotBClr(CL)
				case 0x7:
					// 7F ?? 67
					if DH&0x8 == 0 {
						inst.Opcode = opcode.Bst
					} else {
						inst.Opcode = opcode.Bist
					}
				}
			case 0x7:
				switch CL {
				case 0x0, 0x1, 0x2:
					if DH < 0x8 {
						inst.Opcode = bSetBNotBClr(CL)
					}
				}
			}
		}
	}
	if inst.Opcode == opcode.Invalid {
		return instruction.Inst{Opcode: opcode.Invalid, TotalBytes: 2}
	}
	return inst
}

func table234(bytes []byte) instruction.Inst {
	AH := bytes[0] >> 4
	AL := bytes[0] & 0x0F
	inst := instruction.Inst{
		Opcode:     opcode.Invalid,
		TotalBytes: 2,
	}

	if AH != 0x6 || AL != 0xA {
		return inst
	}

	BH := bytes[1] >> 4
	BL := bytes[1] & 0x0F
	EH := bytes[4] >> 4
	EL := bytes[4] & 0x0F
	FH := bytes[5] >> 4
	GH := bytes[6] >> 4
	GL := bytes[6] & 0x0F
	HH := bytes[7] >> 4
	inst.AddressingMode = addressingmode.AbsoluteAddress
	switch BH {
	case 0x1:
		inst.TotalBytes = 6
		switch BL {
		case 0x0:
			switch EH {
			case 0x6:
				if EL == 0x3 {
					inst.Opcode = opcode.Btst
				}
			case 0x7:
				if EL > 0x2 && EL < 0x8 {
					inst.Opcode = borBxorBandBld(EL, FH)
				}
			}
		case 0x8:
			switch EH {
			case 0x6:
				switch EL {
				case 0x0, 0x1, 0x2:
					inst.Opcode = bSetBNotBClr(EL)
				case 0x7:
					if FH&0x8 == 0 {
						inst.Opcode = opcode.Bst
					} else {
						inst.Opcode = opcode.Bist
					}
				}
			case 0x7:
				// 6A 18 ?? ?? 72 ?0
				if EL < 0x3 && FH < 0x8 {
					inst.Opcode = bSetBNotBClr(EL)
				}
			}
		}
	case 0x3:
		inst.TotalBytes = 8
		switch BL {
		case 0x0:
			switch GH {
			case 0x6:
				if GL == 0x3 {
					inst.Opcode = opcode.Btst
				}
			case 0x7:
				if GL > 0x2 && GL < 0x8 {
					inst.Opcode = borBxorBandBld(GL, HH)
				}
			}
		case 0x8:
			switch GH {
			case 0x6:
				switch GL {
				case 0x0, 0x1, 0x2:
					inst.Opcode = bSetBNotBClr(GL)
				case 0x7:
					if HH&0x8 == 0 {
						inst.Opcode = opcode.Bst
					} else {
						inst.Opcode = opcode.Bist
					}
				}
			case 0x7:
				if GL < 0x3 && HH < 0x8 {
					// 6A 38 ?? ?? ?? ?? 72 ?0
					inst.Opcode = bSetBNotBClr(GL)
				}
			}
		}
	}

	if inst.Opcode == opcode.Invalid {
		return instruction.Inst{Opcode: opcode.Invalid, TotalBytes: 2}
	}
	return inst
}

func branches(b byte) opcode.Opcode {
	branchMap := map[byte]opcode.Opcode{
		0x0: opcode.Bra, // bra in the manual
		0x1: opcode.Brn, // brn in the manual
		0x2: opcode.Bhi,
		0x3: opcode.Bls,
		0x4: opcode.Bcc,
		0x5: opcode.Bcs,
		0x6: opcode.Bne,
		0x7: opcode.Beq,
		0x8: opcode.Bvc,
		0x9: opcode.Bvs,
		0xA: opcode.Bpl,
		0xB: opcode.Bmi,
		0xC: opcode.Bge,
		0xD: opcode.Blt,
		0xE: opcode.Bgt,
		0xF: opcode.Ble,
	}
	return branchMap[b]
}

func bSetBNotBClr(b byte) opcode.Opcode {
	bsetbnotbclrMap := map[byte]opcode.Opcode{
		0x0: opcode.Bset,
		0x1: opcode.Bnot,
		0x2: opcode.Bclr,
	}
	return bsetbnotbclrMap[b]
}

func orXorAnd(b byte) opcode.Opcode {
	orXorAndMap := map[byte]opcode.Opcode{
		0x4: opcode.Or,
		0x5: opcode.Xor,
		0x6: opcode.And,
	}
	return orXorAndMap[b]
}

func borBxorBandBld(b byte, CB byte) opcode.Opcode {
	m := map[byte]opcode.Opcode{
		0x3: opcode.Btst,
		0x4: opcode.Bor,
		0x5: opcode.Bxor,
		0x6: opcode.Band,
		0x7: opcode.Bld,
	}
	nm := map[byte]opcode.Opcode{
		0x4: opcode.Bior,
		0x5: opcode.Bixor,
		0x6: opcode.Biand,
		0x7: opcode.Bild,
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
	return opcode.Invalid
}
