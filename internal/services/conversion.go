package services

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

func StructToJSONBytes(structure any) ([]byte, error) {
	return json.Marshal(structure)
}

func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	gz := gzip.NewWriter(&buf)

	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
