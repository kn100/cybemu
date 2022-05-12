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
