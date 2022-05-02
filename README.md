# Cybemu

Cybemu is an attempt to emulate the Cybiko Classic in Go. Currently, it can only do a very basic disassembly of a H8s/2000 binary. It'll only identify the length of the instruction, whether it's a Byte, Word, or Long, and the instruction opcode.

Special thanks to @_Tim_ on the Cybiko Zone Discord for providing constant guidance and useful information, including test cases.

## Playing with this
You can run this like this
```
go run main.go <file to disassemble>
```

A makefile is included which will help you to run the tests, build a binary, etc.

