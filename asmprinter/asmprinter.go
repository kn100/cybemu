package asmprinter

import (
	"fmt"

	"github.com/kn100/cybemu/disassembler"
)

func PrintAssy(instructions []disassembler.Inst) {
	for _, inst := range instructions {
		if inst.TotalBytes == 2 {
			fmt.Printf("%04x: %02x%02x                 %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], getInstWithSize(inst))
		} else if inst.TotalBytes == 4 {
			fmt.Printf("%05x: %02x%02x %02x%02x            %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], getInstWithSize(inst))
		} else if inst.TotalBytes == 6 {
			fmt.Printf("%04x: %02x%02x %02x%02x %02x%02x       %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.Bytes[4], inst.Bytes[5], getInstWithSize(inst))
		} else if inst.TotalBytes == 8 {
			fmt.Printf("%04x: %02x%02x %02x%02x %02x%02x %02x%02x  %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.Bytes[4], inst.Bytes[5], inst.Bytes[6], inst.Bytes[7], getInstWithSize(inst))
		} else if inst.TotalBytes == 10 {
			fmt.Printf("%04x: %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x  %s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.Bytes[4], inst.Bytes[5], inst.Bytes[6], inst.Bytes[7], inst.Bytes[8], inst.Bytes[9], getInstWithSize(inst))
		} else {
			panic("Something really strange happened and we ended up with an instruction with more bytes than 10.")
		}
	}
}

func getInstWithSize(inst disassembler.Inst) string {
	sizeToString := map[disassembler.Size]string{
		disassembler.Byte:     ".b",
		disassembler.Word:     ".w",
		disassembler.Longword: ".l",
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
