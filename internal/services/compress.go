package services

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
)

func CompressJSONList(data []interface{}) ([]byte, error) {
	// Konwersja listy obiektów do JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Kompresja do GZ
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	if _, err := writer.Write(jsonData); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func gzipDecompress(compressedData []byte) ([]byte, error) {
	// Dekompresja z GZ
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Odczyt zdekompresowanych danych
	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return decompressedData, nil
}
