package opcode

//go:generate stringer -type=Opcode -output=opcode_string.go

// Opcode describes the opcode that is used to identify an instruction.
type Opcode int64

const (
	Invalid Opcode = iota
	Add
	Adds
	Addx
	And
	Andc
	Band
	Bcc
	Bclr
	Bcs
	Beq
	Bge
	Bgt
	Bhi
	Biand
	Bild
	Bior
	Bist
	Bixor
	Bld
	Ble
	Bls
	Blt
	Bmi
	Bne
	Bnot
	Bor
	Bpl
	Bra
	Brn
	Bset
	Bsr
	Bst
	Btst
	Bvc
	Bvs
	Bxor
	Clrmac
	Cmp
	Daa
	Das
	Dec
	Divxs
	Divxu
	Eepmov
	Exts
	Extu
	Inc
	Jmp
	Jsr
	Ldc
	Ldm
	Ldmac
	Mac
	Mov
	Movfpe
	Movtpe
	Mulxs
	Mulxu
	Neg
	Nop
	Not
	Or
	Orc
	Pop
	Push
	Rotl
	Rotr
	Rotxl
	Rotxr
	Rte
	Rts
	Shal
	Shar
	Shll
	Shlr
	Sleep
	Stc
	Stm
	Stmac
	Sub
	Subs
	Subx
	Tas
	Trapa
	Xor
	Xorc
)
