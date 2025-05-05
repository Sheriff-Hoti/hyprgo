package pkg

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Data struct {
	Current_wallpaper string `json:"current_wallpaper"`
	Init              bool   `json:"init"`
}

func GetOrCreateDataDir() (string, error) {
	home_dir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	data_dir := filepath.Join(home_dir, ".local", "share", "hyrgo")

	if err := os.MkdirAll(data_dir, 0755); err != nil {
		return "", err
	}

	return data_dir, nil
}

func GetDataContent() (*Data, error) {
	data_dir, err := GetOrCreateDataDir()

	if err != nil {
		return nil, err
	}

	data_path := filepath.Join(data_dir, "data.json")

	f, err := os.Open(data_path)

	if err != nil {
		var create_err error
		f, create_err = os.Create(data_path)

		if create_err != nil {
			return nil, create_err
		}

		default_data := Data{
			Current_wallpaper: "",
			Init:              true,
		}

		encoder_err := json.NewEncoder(f).Encode(&default_data)

		if encoder_err != nil {
			return nil, encoder_err
		}

		return &default_data, nil

	}

	var data Data
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		log.Println(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	return &data, nil

}
