# Cybemu
![GitHub CI](https://github.com/kn100/cybemu/actions/workflows/go.yml/badge.svg)
![Coveralls](https://coveralls.io/repos/github/kn100/cybemu/badge.svg?branch=master)
[![Go Documentation](https://godocs.io/github.com/kn100/cybemu/disassembler?status.svg)](https://godocs.io/github.com/kn100/cybemu/disassembler)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Discord](https://badgen.net/badge/icon/discord?icon=discord&label=The%20Cybiko%20Zone)](https://discord.gg/4E4Bjsjvyc)
[![DeepSource](https://deepsource.io/gh/kn100/cybemu.svg/?label=active+issues&show_trend=true&token=YSbPkpxvYCG4POBvgHCpL_5q)](https://deepsource.io/gh/kn100/cybemu/?ref=repository-badge)

Cybemu is an attempt to emulate the Cybiko Classic in Go. Currently, it can only do a very basic disassembly of a H8s/2000 binary. It'll only identify the length of the instruction, whether it's a Byte, Word, or Long, and the instruction opcode.

Special thanks to @_Tim_ on the Cybiko Zone Discord for providing constant guidance and useful information, including test cases.

## Playing with this
You can run this like this
```
go run main.go <file to disassemble>
```

A makefile is included which will help you to run the tests, build a binary, etc.

