package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// BodyReader ... TODO
type BodyReader struct {
	io.LimitedReader
}

// Length ... TODO
func (bodyReader BodyReader) Length() int64 {
	return bodyReader.N
}

// BodyFromMap ... TODO
func BodyFromMap(data map[string]string) BodyReader {
	body, err := json.Marshal(data)
	if err != nil {
		return BodyReader{}
	}
	return BodyReader{io.LimitedReader{
		R: bytes.NewBuffer(body),
		N: int64(len(body))}}
}

// BodyFromFile ... TODO
func BodyFromFile(path string) BodyReader {
	info, err := os.Stat(path)

	fmt.Printf("Error: %v\nInfo: %v\n", err, info) // TODO DEBUG

	if err != nil {
		return BodyReader{}
	}

	input, err := os.Open(path)
	if err != nil {
		return BodyReader{}
	}

	return BodyReader{io.LimitedReader{
		R: input,
		N: info.Size()}}
}
