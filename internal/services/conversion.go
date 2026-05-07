package services

import "encoding/json"

func StructToJSONBytes(structure any) ([]byte, error) {
	return json.Marshal(structure)
}
