// contains code intended to print a slice of Inst as standard assembly
// notation (ie, a textual output of the instructions).
package asmprinter

import (
	"fmt"

	"github.com/kn100/cybemu/instruction"
)

// PrintAssy prints a slice of instructions as standard assembly notation.
func PrintAssy(instructions []instruction.Inst) {
	for _, inst := range instructions {
		if inst.TotalBytes == 2 {
			fmt.Printf("%05x: %02x%02x                      %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.String())
		} else if inst.TotalBytes == 4 {
			fmt.Printf("%05x: %02x%02x %02x%02x                 %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.String())
		} else if inst.TotalBytes == 6 {
			fmt.Printf("%05x: %02x%02x %02x%02x %02x%02x            %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.Bytes[4], inst.Bytes[5], inst.String())
		} else if inst.TotalBytes == 8 {
			fmt.Printf("%05x: %02x%02x %02x%02x %02x%02x %02x%02x       %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.Bytes[4], inst.Bytes[5], inst.Bytes[6], inst.Bytes[7], inst.String())
		} else if inst.TotalBytes == 10 {
			fmt.Printf("%05x: %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x  %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.Bytes[4], inst.Bytes[5], inst.Bytes[6], inst.Bytes[7], inst.Bytes[8], inst.Bytes[9], inst.String())
		} else {
			panic("Something really strange happened and we ended up with an instruction with more bytes than 10.")
		}
	}
}
