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

func GetSizeAsSuffix(size Size) string {
	if size == Unset {
		return ""
	}
	return fmt.Sprintf(".%s", strings.ToLower(size.String()[0:1]))
}
