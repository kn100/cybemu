// contains a set of iota which define all the valid addressing modes found on a
// Renesas (Previously Hitachi) h8s/2000 cpu.
package addressingmode

//go:generate stringer -type=AddressingMode -output=addressingmode_string.go

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
