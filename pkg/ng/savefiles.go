package ng

import (
	"encoding/json"
	"os"
)

// SaveGameToFile writes the provided GameData to the specified file path in JSON format.
// The file will be created or truncated if it already exists.
func SaveGameToFile(path string, data GameData) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// LoadGameFromFile reads the specified file path and decodes JSON into dst.
func LoadGameFromFile(path string, dst *GameData) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	return nil
}
