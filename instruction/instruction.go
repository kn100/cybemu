package instruction

import (
	"fmt"
	"strings"

	"github.com/kn100/cybemu/addressingmode"
	"github.com/kn100/cybemu/opcode"
	"github.com/kn100/cybemu/operand"
	"github.com/kn100/cybemu/size"
)

// Inst is a single instruction
type Inst struct {
	// The position of the instruction in the file
	Pos int
	// The number of bytes the instruction is.
	TotalBytes int
	BWL        size.Size
	Opcode     opcode.Opcode
	// The raw bytes that make up the instruction
	Bytes          []byte
	AddressingMode addressingmode.AddressingMode
	OperandSize    int // Bits
	OperandType    operand.OperandType
	RegDst         []byte
	RegSrc         []byte
	Imm16          []byte
	Imm            []byte
	Imm32          []byte
	Reg            []byte
	RegCnt         []byte
	ImmR           []byte
	ImmL           []byte
}

func (i *Inst) DetermineOperandTypeAndSetData() {
	switch i.AddressingMode {
	// case addressingmode.Immediate:
	// 	panic("Immediate not implemented")
	case addressingmode.RegisterDirect:
		switch i.BWL {
		case size.Byte:
			switch i.Opcode {
			case opcode.Add, opcode.Sub, opcode.And, opcode.Cmp, opcode.Mov, opcode.Or, opcode.Xor:
				i.OperandType = operand.R8_R8
				i.RegSrc = []byte{i.Bytes[1] >> 4}
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Dec, opcode.Inc:
				i.OperandType = operand.Ix_R8_INC_DEC
				i.Imm = []byte{1}
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Divxs, opcode.Mulxs:
				i.OperandType = operand.R8_R16_MULXS_DIVXS
				i.RegSrc = []byte{i.Bytes[3] >> 4}
				i.RegDst = []byte{i.Bytes[3] & 0x0F}
			case opcode.Divxu, opcode.Mulxu:
				i.OperandType = operand.R8_R16
				i.RegSrc = []byte{i.Bytes[1] >> 4}
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Ldc:
				i.OperandType = operand.R8_LDC
				i.RegSrc = []byte{i.Bytes[1] & 0x0F}
			case opcode.Stc:
				i.OperandType = operand.R8_STC
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Neg, opcode.Not:
				i.OperandType = operand.R8
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Rotl, opcode.Rotr, opcode.Shal, opcode.Shar:
				i.OperandType = operand.Ix_R8_SH
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
				if (i.Bytes[1] >> 4) == 0x8 {
					i.Imm = []byte{1}
				} else {
					i.Imm = []byte{2}
				}
			case opcode.Rotxl, opcode.Rotxr, opcode.Shll, opcode.Shlr:
				i.OperandType = operand.Ix_R8_SH
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
				if (i.Bytes[1] >> 4) == 0x0 {
					i.Imm = []byte{1}
				} else {
					i.Imm = []byte{2}
				}
			}

		case size.Word:
			switch i.Opcode {
			case opcode.Add, opcode.Sub, opcode.And, opcode.Cmp, opcode.Mov, opcode.Or, opcode.Xor:
				i.OperandType = operand.R16_R16
				i.RegSrc = []byte{i.Bytes[1] >> 4}
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Dec, opcode.Inc:
				i.OperandType = operand.Ix_R16_INC_DEC
				BH := i.Bytes[1] >> 4
				if BH == 0x5 {
					i.Imm = []byte{1}
				} else if BH == 0xD {
					i.Imm = []byte{2}
				}
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Mulxs, opcode.Divxs:
				i.OperandType = operand.R16_R32_MULXS_DIVXS
				i.RegSrc = []byte{i.Bytes[3] >> 4}
				i.RegDst = []byte{i.Bytes[3] & 0x0F}
			case opcode.Divxu, opcode.Mulxu:
				i.OperandType = operand.R16_R32
				i.RegSrc = []byte{i.Bytes[1] >> 4}
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Exts:
				i.OperandType = operand.R16
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Extu:
				i.OperandType = operand.R16
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Neg, opcode.Not:
				i.OperandType = operand.R16
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Rotl, opcode.Rotr, opcode.Shal, opcode.Shar:
				i.OperandType = operand.Ix_R16_SH
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
				if (i.Bytes[1] >> 4) == 0x9 {
					i.Imm = []byte{1}
				} else {
					i.Imm = []byte{2}
				}
			case opcode.Rotxl, opcode.Rotxr, opcode.Shll, opcode.Shlr:
				i.OperandType = operand.Ix_R16_SH
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
				if (i.Bytes[1] >> 4) == 0x1 {
					i.Imm = []byte{1}
				} else {
					i.Imm = []byte{2}
				}
			}
		case size.Longword:
			switch i.Opcode {
			case opcode.Add, opcode.Sub, opcode.Cmp, opcode.Mov:
				i.OperandType = operand.R32_R32_S2
				i.RegSrc = []byte{i.Bytes[1] >> 4}
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.And, opcode.Or, opcode.Xor:
				i.OperandType = operand.R32_R32_S4
				i.RegSrc = []byte{i.Bytes[3] >> 4}
				i.RegDst = []byte{i.Bytes[3] & 0x0F}
			case opcode.Dec, opcode.Inc:
				i.OperandType = operand.Ix_R32_INC_DEC
				BH := i.Bytes[1] >> 4
				if BH == 0x7 {
					i.Imm = []byte{1}
				} else if BH == 0xF {
					i.Imm = []byte{2}
				}
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Exts, opcode.Extu, opcode.Neg, opcode.Not:
				i.OperandType = operand.R32
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Rotl, opcode.Rotr, opcode.Shal, opcode.Shar:
				i.OperandType = operand.Ix_R32_SH
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
				if (i.Bytes[1] >> 4) == 0xB {
					i.Imm = []byte{1}
				} else {
					i.Imm = []byte{2}
				}
			case opcode.Rotxl, opcode.Rotxr, opcode.Shll, opcode.Shlr:
				i.OperandType = operand.Ix_R32_SH
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
				if (i.Bytes[1] >> 4) == 0x3 {
					i.Imm = []byte{1}
				} else {
					i.Imm = []byte{2}
				}
			}
		case size.Unset:
			switch i.Opcode {
			case opcode.Adds, opcode.Subs:
				i.OperandType = operand.Ix_R32_ADDS_SUBS
				// 4 MSB of i.Bytes[1] is the immediate, 0x0 is 1, 0x8 is 2, 0x9 is 4.
				val := i.Bytes[1] >> 4
				if val == 0x0 {
					i.Imm = []byte{1}
				} else if val == 0x8 {
					i.Imm = []byte{2}
				} else if val == 0x9 {
					i.Imm = []byte{4}
				}
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Addx, opcode.Subx:
				i.OperandType = operand.R8_R8
				i.RegSrc = []byte{i.Bytes[1] >> 4}
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Band, opcode.Biand, opcode.Bild, opcode.Bior, opcode.Bist, opcode.Bixor, opcode.Bld, opcode.Bst:
				i.OperandType = operand.Ix_R8
				// 4 MSB of i.Bytes[1] is the immediate, 0x0 through 0x7 is 0 through 7
				i.Imm = []byte{i.Bytes[1] >> 4}
				if i.Imm[0] > 7 {
					i.Imm[0] = i.Imm[0] - 8
				}
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Bclr, opcode.Bnot, opcode.Bor, opcode.Bset, opcode.Btst, opcode.Bxor:
				// AH is the 4 MSB of i.Bytes[0]
				check := i.Bytes[0] >> 4
				if check == 0x7 {
					i.OperandType = operand.Ix_R8
					i.Imm = []byte{i.Bytes[1] >> 4}
					i.RegDst = []byte{i.Bytes[1] & 0x0F}
				} else {
					i.OperandType = operand.R8_R8
					i.RegSrc = []byte{i.Bytes[1] >> 4}
					i.RegDst = []byte{i.Bytes[1] & 0x0F}
				}
			case opcode.Daa, opcode.Das:
				i.OperandType = operand.R8
				i.RegDst = []byte{i.Bytes[1] & 0x0F}
			case opcode.Trapa:
				i.OperandType = operand.TRAPA_Ix
				i.Imm = []byte{i.Bytes[1] >> 4}
			}
		}
	case addressingmode.None:
		switch i.BWL {
		case size.Unset:
			switch i.Opcode {
			case opcode.Nop:
				i.OperandType = operand.None
			}
		}

		// case addressingmode.RegisterIndirect:
		// 	panic("RegisterIndirect not implemented")
		// case addressingmode.AbsoluteAddress:
		// 	panic("AbsoluteAddress not implemented")
		// case addressingmode.ProgramCounterRelative:
		// 	panic("ProgramCounterrelative not implemented")
		// case addressingmode.MemoryIndirect:
		// 	panic("MemoryIndirect not implemented")
		// case addressingmode.RegisterIndirectWithDisplacement:
		// 	panic("RegisterIndirectWithDisplacement not implemented")
		// case addressingmode.RegisterIndirectWithPostIncrement:
		// 	panic("RegisterIndirectWithPostIncrement not implemented")
		// case addressingmode.RegisterIndirectWithPreDecrement:
		// 	panic("RegisterIndirectWithPreDecrement not implemented")
	}
}

// Returns the mnemonic, followed by the size, and the operand data.
func (i *Inst) String() string {
	build := ""
	switch i.OperandType {
	case operand.R8_R8, operand.R16_R16:
		bwl := i.BWL
		if bwl == size.Unset {
			bwl = size.Byte
		}
		build = fmt.Sprintf("%s, %s", toRegister(i.RegSrc[0], bwl), toRegister(i.RegDst[0], bwl))
	case operand.R32_R32_S2:
		build = fmt.Sprintf("%s, %s", toShiftedRegister(i.RegSrc[0], i.BWL), toRegister(i.RegDst[0], i.BWL))
	case operand.R32_R32_S4:
		build = fmt.Sprintf("%s, %s", toRegister(i.RegSrc[0], i.BWL), toRegister(i.RegDst[0], i.BWL))
	case operand.Ix_R32_ADDS_SUBS:
		build = fmt.Sprintf("#%d, %s", int(i.Imm[0]), toRegister(i.RegDst[0], size.Longword))
	case operand.Ix_R8:
		build = fmt.Sprintf("#%d, %s", int(i.Imm[0]), toRegister(i.RegDst[0], size.Byte))
	case operand.R8:
		build = toRegister(i.RegDst[0], size.Byte)
	case operand.Ix_R8_INC_DEC, operand.Ix_R8_SH:
		build = fmt.Sprintf("#%d, %s", int(i.Imm[0]), toRegister(i.RegDst[0], size.Byte))
	case operand.Ix_R16_INC_DEC, operand.Ix_R16_SH:
		build = fmt.Sprintf("#%d, %s", int(i.Imm[0]), toRegister(i.RegDst[0], size.Word))
	case operand.Ix_R32_INC_DEC, operand.Ix_R32_SH:
		build = fmt.Sprintf("#%d, %s", int(i.Imm[0]), toRegister(i.RegDst[0], size.Longword))
	case operand.R8_R16_MULXS_DIVXS, operand.R8_R16:
		build = fmt.Sprintf("%s, %s", toRegister(i.RegSrc[0], i.BWL), toRegister(i.RegDst[0], size.Word))
	case operand.R16_R32_MULXS_DIVXS, operand.R16_R32:
		build = fmt.Sprintf("%s, %s", toRegister(i.RegSrc[0], i.BWL), toRegister(i.RegDst[0], size.Longword))
	case operand.R16, operand.R32:
		build = toRegister(i.RegDst[0], i.BWL)
	case operand.R8_LDC:
		ccrexr := "ccr"
		if i.Bytes[1]>>4 == 0x1 {
			ccrexr = "exr"
		}
		build = fmt.Sprintf("%s, %s", toRegister(i.RegSrc[0], i.BWL), ccrexr)
	case operand.R8_STC:
		ccrexr := "ccr"
		if i.Bytes[1]>>4 == 0x1 {
			ccrexr = "exr"
		}
		build = fmt.Sprintf("%s, %s", ccrexr, toRegister(i.RegDst[0], i.BWL))
	case operand.TRAPA_Ix:
		build = fmt.Sprintf("#%d", int(i.Imm[0]))
	case operand.None:
		build = ""
	}

	mnemonic := strings.ToLower(i.Opcode.String())
	sizeSuffix := size.GetSizeAsSuffix(i.BWL)
	// TODO: handling around none operand is horrible
	if build != "" {
		return fmt.Sprintf("%s%s %s", mnemonic, sizeSuffix, build)
	} else if i.OperandType == operand.None {
		return fmt.Sprintf("%s%s", mnemonic, sizeSuffix)
	}

	return fmt.Sprintf("%s%s :(", mnemonic, sizeSuffix)
}

func toRegister(b byte, s size.Size) string {
	register := ""
	intb := int(b)

	switch s {
	case size.Byte:
		hl := "h"

		if int(b) > 7 {
			hl = "l"
			intb -= 8
		}
		register = fmt.Sprintf("r%d%s", intb, hl)
	case size.Word:
		re := "r"
		if int(b) > 7 {
			re = "e"
			intb -= 8
		}
		register = fmt.Sprintf("%s%d", re, intb)
	case size.Longword:
		if intb == 15 || intb == 7 {
			register = "sp"
		} else {
			register = fmt.Sprintf("er%d", intb)
		}
	}
	return register
}

func toShiftedRegister(b byte, s size.Size) string {
	register := ""
	intb := int(b)
	intb = intb - 8

	if s == size.Longword {
		if intb == 15 || intb == 7 {
			register = "sp"
		} else {
			register = fmt.Sprintf("er%d", intb)
		}
	}
	return register
}
