package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Task struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"` // pending | done | failed
}

type TimerState struct {
	Task      string `json:"task"`
	Minutes   int    `json:"minutes"`
	StartedAt int64  `json:"started_at"`
	Duration  int    `json:"duration"` // in seconds
}

type CrumbData struct {
	Tasks []Task      `json:"tasks"`
	Ideas []string    `json:"ideas"`
	Notes []string    `json:"notes"`
	Timer *TimerState `json:"timer,omitempty"`
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

	// If file doesn't exist, return empty struct safely
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return CrumbData{
			Tasks: []Task{},
			Notes: []string{},
			Ideas: []string{},
		}, nil
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

	if data.Tasks == nil {
		data.Tasks = []Task{}
	}
	if data.Notes == nil {
		data.Notes = []string{}
	}
	if data.Ideas == nil {
		data.Ideas = []string{}
	}

	return data, nil
}

// WriteData writes the CrumbData struct to the JSON file atomically.
// It writes to a temp file in the same directory and then renames it over the target.
func WriteData(data CrumbData) error {
	dbPath, err := GetDbPath()
	if err != nil {
		return err
	}

	// Create directory if missing
	dir := filepath.Dir(dbPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(dir, "crumb-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.Write(payload); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err := os.Chmod(tmpPath, 0644); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, dbPath)
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
