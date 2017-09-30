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

func formatError(err error) error {
	return fmt.Errorf("Parser: Couldn't read response: %v", err)
}

func readByte(reader io.ByteReader) (value byte, err error) {
	value, err = reader.ReadByte()
	if err != nil {
		err = formatError(err)
	}
	return
}

func readBytes(n int, reader io.Reader) (bytesRead int, bytes []byte, err error) {
	bytes = make([]byte, n)
	bytesRead, err = reader.Read(bytes)
	if err != nil {
		err = formatError(err)
	}
	return
}

func parseInteger(reader io.ByteReader) (int, error) {
	number := 0
	negativeMultipler := 1
	for {
		value, err := readByte(reader)
		if err != nil {
			return -1, err
		}

		if value == 45 {
			negativeMultipler = -1
			continue
		}

		if value == 13 {
			_, err := readByte(reader) // discard '\n'
			if err != nil {
				return -1, err
			}
			break
		}

		number = (number * 10) + int(value-48)
	}
	return number * negativeMultipler, nil
}

func parseBulkString(reader IOReader) (interface{}, error) {
	length, err := parseInteger(reader)
	if err != nil {
		return nil, err
	}
	if length == 0 {
		return "", nil
	}
	if length == -1 {
		return nil, nil
	}

	bytesRead, bytes, err := readBytes(length, reader)
	if err != nil {
		return nil, err
	}
	if bytesRead != length {
		return nil, fmt.Errorf("Read less bytes than required (should have been %d, was %d)", length, bytesRead)
	}
	_, _, err = readBytes(2, reader) // discard "\r\n"
	if err != nil {
		return nil, err
	}
	return string(bytes), nil
}

func parseArray(reader IOReader) (interface{}, error) {
	length, err := parseInteger(reader)
	if err != nil {
		return nil, err
	}
	if length == 0 {
		var empty []interface{}
		return empty, nil
	}
	if length == -1 {
		return nil, nil
	}

	values := make([]interface{}, length)
	for index := range values {
		values[index], err = Parse(reader)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

// Parse parses responses from Redis server
func Parse(reader IOReader) (interface{}, error) {
	char, err := readByte(reader)
	if err != nil {
		return nil, err
	}
	switch char {
	case 36:
		return parseBulkString(reader)
	case 42:
		return parseArray(reader)
	case 43:
		result, err := reader.ReadString('\n')
		if err != nil {
			return nil, formatError(err)
		}
		return strings.TrimSpace(result), nil
	case 45:
		result, err := reader.ReadString('\n')
		if err != nil {
			return nil, formatError(err)
		}
		return nil, fmt.Errorf("Redis: %v", strings.TrimSpace(result))
	case 58:
		return parseInteger(reader)
	default:
		panic(fmt.Sprintf("Response symbole '%v' not supported", char))
	}
}
