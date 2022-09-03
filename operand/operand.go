// contains an set of iota which define all the valid operand types on a Renesas
// (Previously Hitachi) h8s/2000. They're a direct copy of what is found in C4PC
// - Thanks to @_tim_ for providing these.
package operand

//go:generate stringer -type=OperandType -output=operandtype_string.go

// OperandType is a type of Operand.
type OperandType int64

const (
	Unknown OperandType = iota
	None
	I8_R8
	I16_R16
	I32_R32
	Ix_R8
	Ix_AR32
	Ix_AI8
	Ix_AI16
	Ix_AI32
	R32_R32_S2
	R32_R32_S4
	R16_R16
	R8_R16
	R16_R32
	R8_R16_MULXS_DIVXS
	R16_R32_MULXS_DIVXS
	R32_AI16
	AI16_R32
	R16_AI16
	AI16_R16
	R16_AI32
	AI32_R16
	AR32_R32
	R32_AR32
	AR32_R16
	R16_AR32
	AR32_R8
	R8_AR32
	O16
	I24
	S2_IMM
	I8_EXR
	Ix_R8_SH
	Ix_R16_SH
	Ix_R32_SH
	R8
	R8_STC
	R8_LDC
	R16
	R32
	R8_R8
	AR32_R32R32
	R32R32_AR32
	R32_AI32
	AI32_R32
	AR32_S2
	AI16R32_R16
	R16_AI16R32
	R8_AI32_BCLR
	R8_AI16
	R8_AI16_S6
	AI16_R8
	AI32R32_R8
	AI32R32_R16
	R16_AI32R32
	R8_AI32R32
	R8_AR32_BCLR
	AI32R32_R32
	R32_AI32R32
	AI16R32_R32
	R32_AI16R32
	AI16R32_CCR
	CCR_AI16R32
	AI32R32_CCR
	CCR_AI32R32
	AI16R32_R8
	R8_AI16R32
	Ix_R32_ADDS_SUBS
	Ix_R8_INC_DEC
	Ix_R16_INC_DEC
	Ix_R32_INC_DEC
	R8_AI8
	AI8_R8
	R8_AI8_BCLR
	S4_R32
	AI16_CCR
	CCR_AI16
	AI32_CCR
	CCR_AI32
	CCR_AR32
	TRAPA_Ix
	// TODO: Uhhhh
	AI32_R8
)
