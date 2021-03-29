package util

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type BodyReader struct {
	io.LimitedReader
}

func (bodyReader BodyReader) Length() int64 {
	return bodyReader.N
}

func BodyFromMap(data map[string]string) BodyReader {
	body, err := json.Marshal(data)
	if err != nil {
		return BodyReader{}
	}
	return BodyReader{io.LimitedReader{
		R: bytes.NewBuffer(body),
		N: int64(len(body))}}
}

func BodyFromFile(path string) BodyReader {
	info, err := os.Stat(path)

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
