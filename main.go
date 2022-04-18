package main

import (
	"fmt"
	"os"
)

type Inst struct {
	// Either 2, or 4
	Pos    int
	Size   int
	Opcode string
	Bytes  []byte
}

func main() {
	f, err := os.Open("emu_rom.bin")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	instructions := []Inst{}
	// Read entire file into byte array
	// TODO: Set array to correct size using os.Stat
	bytes := make([]byte, 10000)
	f.Read(bytes)
	i := 0
	for i < len(bytes) {
		// If we have less than 2 bytes left, we are done
		if i+2 > len(bytes) {
			break
		}
		AH := bytes[i] >> 4
		AL := bytes[i] & 0x0F
		inst := Inst{
			Pos:  i,
			Size: 2,
		}
		switch {
		// noop
		case AH == 0x0 && AL == 0x0:
			inst.Opcode = "nop"
		case AH == 0x0 && AL == 0x1:
			// Table 2.3
			inst.Opcode = "unk"
		case AH == 0x0 && AL == 0x2:
			inst.Opcode = "unk"
			// STC/STMAC
		case AH == 0x0 && AL == 0x3:
			inst.Opcode = "unk"
			// LDC/LDMAC
		case AH == 0x0 && AL == 0x4:
			inst.Opcode = "orc"
		case AH == 0x0 && AL == 0x5:
			inst.Opcode = "xorc"
		case AH == 0x0 && AL == 0x6:
			inst.Opcode = "andc"
		case AH == 0x0 && AL == 0x7:
			inst.Opcode = "ldc"
		case AH == 0x0 && AL == 0x8:
			inst.Opcode = "add"
		case AH == 0x0 && AL == 0x9:
			inst.Opcode = "add"
		case AH == 0x0 && AL == 0xA:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x0 && AL == 0xB:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x0 && AL == 0xC:
			inst.Opcode = "mov"
		case AH == 0x0 && AL == 0xD:
			inst.Opcode = "mov"
		case AH == 0x0 && AL == 0xE:
			inst.Opcode = "addx"
		case AH == 0x0 && AL == 0xF:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x1 && AL == 0x0:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x1 && AL == 0x1:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x1 && AL == 0x2:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x1 && AL == 0x3:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x1 && AL == 0x4:
			inst.Opcode = "or"
		case AH == 0x1 && AL == 0x5:
			inst.Opcode = "xor"
		case AH == 0x1 && AL == 0x6:
			inst.Opcode = "and"
		case AH == 0x1 && AL == 0x7:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x1 && AL == 0x8:
			inst.Opcode = "sub"
		case AH == 0x1 && AL == 0x9:
			inst.Opcode = "sub"
		case AH == 0x1 && AL == 0xA:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x1 && AL == 0xB:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x1 && AL == 0xC:
			inst.Opcode = "cmp"
		case AH == 0x1 && AL == 0xD:
			inst.Opcode = "cmp"
		case AH == 0x1 && AL == 0xE:
			inst.Opcode = "subx"
		case AH == 0x1 && AL == 0xF:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x2 || AH == 0x3:
			inst.Opcode = "mov"
		case AH == 0x4 && AL == 0x0:
			inst.Opcode = "bra"
		case AH == 0x4 && AL == 0x1:
			inst.Opcode = "brn"
		case AH == 0x4 && AL == 0x2:
			inst.Opcode = "bh"
		case AH == 0x4 && AL == 0x3:
			inst.Opcode = "bls"
		case AH == 0x4 && AL == 0x4:
			inst.Opcode = "bcc"
		case AH == 0x4 && AL == 0x5:
			inst.Opcode = "bcs"
		case AH == 0x4 && AL == 0x6:
			inst.Opcode = "bne"
		case AH == 0x4 && AL == 0x7:
			inst.Opcode = "beq"
		case AH == 0x4 && AL == 0x8:
			inst.Opcode = "bvc"
		case AH == 0x4 && AL == 0x9:
			inst.Opcode = "bvs"
		case AH == 0x4 && AL == 0xA:
			inst.Opcode = "bpl"
		case AH == 0x4 && AL == 0xB:
			inst.Opcode = "bmi"
		case AH == 0x4 && AL == 0xC:
			inst.Opcode = "bge"
		case AH == 0x4 && AL == 0xD:
			inst.Opcode = "blt"
		case AH == 0x4 && AL == 0xE:
			inst.Opcode = "bgt"
		case AH == 0x4 && AL == 0xF:
			inst.Opcode = "ble"
		case AH == 0x5 && AL == 0x0:
			inst.Opcode = "mulxu"
		case AH == 0x5 && AL == 0x1:
			inst.Opcode = "divxu"
		case AH == 0x5 && AL == 0x2:
			inst.Opcode = "mulxu"
		case AH == 0x5 && AL == 0x3:
			inst.Opcode = "divxu"
		case AH == 0x5 && AL == 0x4:
			inst.Opcode = "rts"
		case AH == 0x5 && AL == 0x5:
			inst.Opcode = "bsr"
		case AH == 0x5 && AL == 0x6:
			inst.Opcode = "rte"
		case AH == 0x5 && AL == 0x7:
			inst.Opcode = "trapa"
		case AH == 0x5 && AL == 0x8:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x5 && AL == 0x9:
			inst.Opcode = "jmp"
		case AH == 0x5 && AL == 0xA:
			inst.Opcode = "jmp"
		case AH == 0x5 && AL == 0xB:
			inst.Opcode = "jmp"
		case AH == 0x5 && AL == 0xC:
			inst.Opcode = "bsr"
		case AH == 0x5 && AL == 0xD:
			inst.Opcode = "jsr"
		case AH == 0x5 && AL == 0xE:
			inst.Opcode = "jsr"
		case AH == 0x5 && AL == 0xF:
			inst.Opcode = "jsr"
		case AH == 0x6 && AL == 0x0:
			inst.Opcode = "bset"
		case AH == 0x6 && AL == 0x1:
			inst.Opcode = "bnot"
		case AH == 0x6 && AL == 0x2:
			inst.Opcode = "bclr"
		case AH == 0x6 && AL == 0x3:
			inst.Opcode = "btst"
		case AH == 0x6 && AL == 0x4:
			inst.Opcode = "or"
		case AH == 0x6 && AL == 0x5:
			inst.Opcode = "xor"
		case AH == 0x6 && AL == 0x6:
			inst.Opcode = "and"
		case AH == 0x6 && AL == 0x7:
			inst.Opcode = "unk"
			// BST / BIAND
		case AH == 0x6 && AL == 0x8:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0x9:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0xA:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x6 && AL == 0xB:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0xC:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0xD:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0xE:
			inst.Opcode = "mov"
		case AH == 0x6 && AL == 0xF:
			inst.Opcode = "mov"
		case AH == 0x7 && AL == 0x0:
			inst.Opcode = "bset"
		case AH == 0x7 && AL == 0x1:
			inst.Opcode = "bnot"
		case AH == 0x7 && AL == 0x2:
			inst.Opcode = "bclr"
		case AH == 0x7 && AL == 0x3:
			inst.Opcode = "btst"
		case AH == 0x7 && AL == 0x4:
			inst.Opcode = "unk"
			//bor/bior
		case AH == 0x7 && AL == 0x5:
			inst.Opcode = "unk"
			//bxor/bixor
		case AH == 0x7 && AL == 0x6:
			inst.Opcode = "unk"
			//band/biand
		case AH == 0x7 && AL == 0x7:
			inst.Opcode = "unk"
			//bld/bild
		case AH == 0x7 && AL == 0x8:
			inst.Opcode = "mov"
		case AH == 0x7 && AL == 0x9:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x7 && AL == 0xA:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x7 && AL == 0xB:
			inst.Opcode = "eepmov"
		case AH == 0x7 && AL == 0xC:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x7 && AL == 0xD:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x7 && AL == 0xE:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x7 && AL == 0xF:
			inst.Opcode = "unk"
			// Table 2.3
		case AH == 0x8:
			inst.Opcode = "add"
		case AH == 0x9:
			inst.Opcode = "addx"
		case AH == 0xA:
			inst.Opcode = "cmp"
		case AH == 0xB:
			inst.Opcode = "subx"
		case AH == 0xC:
			inst.Opcode = "or"
		case AH == 0xD:
			inst.Opcode = "xor"
		case AH == 0xE:
			inst.Opcode = "and"
		case AH == 0xF:
			inst.Opcode = "mov"
		default:
			panic("wut")
		}
		for b := 0; b < inst.Size; b++ {
			inst.Bytes = append(inst.Bytes, bytes[i+b])
		}
		instructions = append(instructions, inst)
		i = i + inst.Size
	}
	PrintAssy(instructions)
}

func PrintAssy(instructions []Inst) {
	unks := 0
	decoded := 0
	// Count the number of unks
	for _, inst := range instructions {
		if inst.Opcode == "unk" {
			unks++
		} else {
			decoded++
		}
	}
	fmt.Printf("%d instructions decoded, %d unks\n", decoded, unks)
	// display percentage decoded
	fmt.Printf("%d%% decoded\n", (decoded*100)/(decoded+unks))
	for _, inst := range instructions {
		if inst.Size == 2 {
			fmt.Printf("%04X:\t%02X %02X\t%s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Opcode)
		} else if inst.Size == 4 {
			fmt.Printf("%04X:\t%02X %02X %02X %02X\t%s\n", inst.Pos, inst.Bytes[0], inst.Bytes[1], inst.Bytes[2], inst.Bytes[3], inst.Opcode)
		}

	}
}
