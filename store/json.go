package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type CrumbData struct {
	Next  string   `json:"next"`
	Tasks []Task   `json:"tasks"`
	Ideas []string `json:"ideas"`
	Notes []string `json:"notes"`
	Done  []string `json:"done"`
}

// GetDbPath resolves the path to ~/.config/crumb/data.json
func GetDbPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "crumb", "data.json"), nil
}

// ReadData reads the JSON file and returns the CrumbData struct
func ReadData() (CrumbData, error) {
	dbPath, err := GetDbPath()
	if err != nil {
		return CrumbData{}, err
	}

	// If file doesn't exist, return empty stru0ct safely
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return CrumbData{Notes: []string{}}, nil
	}

	content, err := os.ReadFile(dbPath)
	if err != nil {
		return CrumbData{}, err
	}

	var data CrumbData
	err = json.Unmarshal(content, &data)
	if err != nil {
		return CrumbData{}, err
	}

	if data.Notes == nil {
		data.Notes = []string{}
	}

	return data, nil
}

// WriteData writes the CrumbData struct to the JSON file
func WriteData(data CrumbData) error {
	dbPath, err := GetDbPath()
	if err != nil {
		return err
	}

	// Create directory if missing
	dir := filepath.Dir(dbPath)
	err = os.MkdirAll(dir, 0755) // Idempotent operation; won't error if dir exists
	if err != nil {
		return err
	}

	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dbPath, payload, 0644)
}
