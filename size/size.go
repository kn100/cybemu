// contains a set of iota which define all the valid sizes for a instruction
// found on a Renesas (Previously Hitachi) h8s/2000 cpu.
package size

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=Size -output=size_string.go

// Size indicates the size of an instruction (Byte, Word, Longword)
type Size int64

const (
	Unset Size = iota
	Byte
	Word
	Longword
)

// GetSizeAsString, given a Size, returns the string representation of a size.
// For example, if the size is Byte, then the string returned is ".b". It is
// intended to be used as part of the generation of a line in a disassembler.
func GetSizeAsSuffix(size Size) string {
	if size == Unset {
		return ""
	}
	return fmt.Sprintf(".%s", strings.ToLower(size.String()[0:1]))
}
