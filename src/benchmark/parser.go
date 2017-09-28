package benchmark

import (
	"fmt"
	"strings"
)

// IOReader bridges gap between a Buffer and Reader types
type IOReader interface {
	ReadString(delim byte) (string, error)
	ReadByte() (byte, error)
}

func throwReadError(err error) {
	panic(fmt.Sprintf("Couldn't read response: %v", err))
}

// Parse parses responses from Redis server
func Parse(reader IOReader) string {
	char, err := reader.ReadByte()
	if err != nil {
		throwReadError(err)
	}

	var result string
	switch char {
	case 43:
		if result, err = reader.ReadString('\n'); err != nil {
			throwReadError(err)
		}
	case 45:
		if result, err = reader.ReadString('\n'); err != nil {
			throwReadError(err)
		}
	default:
		panic(fmt.Sprintf("Response symbole '%v' not supported", char))
	}
	return strings.TrimSpace(result)
}
