package main

import (
	"fmt"
	"os"

	"github.com/kn100/cybemu/asmprinter"
	"github.com/kn100/cybemu/disassembler"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: disassembler <file>")
		os.Exit(1)
	}

	args := os.Args[1:]
	f, err := os.Open(args[0])
	if err != nil {
		fmt.Printf("Couldn't open file for some reason. Error was: %s\n", err)
		return
	}
	defer f.Close()

	finfo, err := os.Stat(args[0])
	if err != nil {
		fmt.Printf("Couldn't get info about file for some reason. Error was: %s\n", err)
		return
	}

	bytes := make([]byte, finfo.Size())
	f.Read(bytes)
	instructions := disassembler.Disassemble(bytes)
	asmprinter.PrintAssy(instructions)
}
