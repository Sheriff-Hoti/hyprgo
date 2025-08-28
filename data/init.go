package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Data struct {
	Current_wallpaper string `json:"current_wallpaper"`
	Init              bool   `json:"init"`
}

func GetDefaultDataPath() string {
	const (
		xdgConfigHome = "XDG_DATA_HOME"
	)

	if val, ok := os.LookupEnv(xdgConfigHome); ok {
		return filepath.Join(val, "hyprgo", "data.json")
	}

	// fallback to $HOME/.config/hyprgo/config.json
	home, err := os.UserHomeDir()
	if err != nil {
		// if home can't be resolved, fallback to current working directory
		return filepath.Join(".", "data.json")
	}

	return filepath.Join(home, ".local", "share", "hyprgo", "data.json")
}

func ReadDataFile(data_path string) (*Data, error) {
	data_path = os.ExpandEnv(data_path)
	if err := os.MkdirAll(filepath.Dir(data_path), 0755); err != nil {
		return nil, fmt.Errorf("creating parent dirs: %w", err)
	}

	// Open or create the file
	file, err := os.OpenFile(data_path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	// If the file is empty, return nil
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}
	if stat.Size() == 0 {
		return &Data{}, nil
	}
	// Decode JSON into Data
	var data Data
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding json: %w", err)
	}

	return &data, nil
}
