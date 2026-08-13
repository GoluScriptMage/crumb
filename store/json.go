package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Task struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"` // pending | done | canceled
}

type CrumbData struct {
	Next  string   `json:"next"`
	Tasks []Task   `json:"tasks"`
	Ideas []string `json:"ideas"`
	Notes []string `json:"notes"`
	Done  []string `json:"done"`
}

// dbPathOverride allows tests (or other callers) to redirect storage
// away from the default user config location.
var dbPathOverride string

// SetDbPathOverride redirects ReadData/WriteData to a custom path.
// Pass an empty string to revert to the default user config location.
func SetDbPathOverride(path string) {
	dbPathOverride = path
}

// GetDbPath resolves the path to ~/.config/crumb/data.json
func GetDbPath() (string, error) {
	if dbPathOverride != "" {
		return dbPathOverride, nil
	}
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

// Update reads data, applies the modification function, and writes back atomically.
// Eliminates duplicate read-modify-write boilerplate in commands.
func Update(fn func(*CrumbData) error) error {
	data, err := ReadData()
	if err != nil {
		return err
	}
	if err := fn(&data); err != nil {
		return err
	}
	return WriteData(data)
}
