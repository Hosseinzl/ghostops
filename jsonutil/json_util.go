package jsonutil

import (
	"encoding/json"
	"fmt"
	"os"
)

// generic function to save any object to a JSON file
func SaveToJSONFile[T any](filePath string, data T) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// generic function to load any object from a JSON file
func LoadFromJSONFile[T any](filePath string) (T, error) {

	var data T

	file, err := os.Open(filePath)
	if err != nil {
		return *new(T), fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		return *new(T), fmt.Errorf("failed to decode JSON: %w", err)
	}

	return data, nil
}
