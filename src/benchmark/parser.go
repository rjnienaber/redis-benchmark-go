package benchmark

import (
	"fmt"
	"io"
	"strings"
)

// IOReader bridges gap between Buffer and Reader types
type IOReader interface {
	io.Reader
	io.ByteReader
	ReadString(delim byte) (string, error)
}

func throwReadError(err error) {
	panic(fmt.Sprintf("Couldn't read response: %v", err))
}

func readByte(reader io.ByteReader) byte {
	value, err := reader.ReadByte()
	if err != nil {
		throwReadError(err)
	}
	return value
}

func readBytes(n int, reader io.Reader) (bytesRead int, bytes []byte) {
	var err error
	bytes = make([]byte, n)
	bytesRead, err = reader.Read(bytes)
	if err != nil {
		throwReadError(err)
	}
	return
}

func parseInteger(reader io.ByteReader) int {
	number := 0
	negativeMultipler := 1
	for {
		value := readByte(reader)
		if value == 45 {
			negativeMultipler = -1
			continue
		}

		if value == 13 {
			break
		}

		number = (number * 10) + int(value-48)
	}
	readByte(reader) // discard '\n'
	return number * negativeMultipler
}

func parseBulkString(reader IOReader) interface{} {
	length := parseInteger(reader)
	if length == 0 {
		return ""
	}
	if length == -1 {
		return nil
	}

	bytesRead, bytes := readBytes(length, reader)
	if bytesRead != length {
		panic(fmt.Sprintf("Read less bytes than required (should have been %d, was %d)", length, bytesRead))
	}
	readBytes(2, reader) // discard "\r\n"
	return string(bytes)
}

func parseArray(reader IOReader) interface{} {
	length := parseInteger(reader)
	if length == 0 {
		var empty []interface{}
		return empty
	}
	if length == -1 {
		return nil
	}

	values := make([]interface{}, length)
	for index := range values {
		values[index] = Parse(reader)
	}

	return values
}

// Parse parses responses from Redis server
func Parse(reader IOReader) interface{} {
	char := readByte(reader)

	switch char {
	case 36:
		return parseBulkString(reader)
	case 42:
		return parseArray(reader)
	case 43, 45:
		result, err := reader.ReadString('\n')
		if err != nil {
			throwReadError(err)
		}
		return strings.TrimSpace(result)
	case 58:
		return parseInteger(reader)
	default:
		panic(fmt.Sprintf("Response symbole '%v' not supported", char))
	}
}
