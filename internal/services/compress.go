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

func DecompressJSONList(compressedData []byte) ([]interface{}, error) {
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

	// Konwersja do listy obiektów JSON
	var result []interface{}
	if err := json.Unmarshal(decompressedData, &result); err != nil {
		return nil, err
	}

	return result, nil
}
