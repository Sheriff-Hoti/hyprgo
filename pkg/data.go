package pkg

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type DataActionMode int

const (
	Read DataActionMode = iota
	Write
)

type Data struct {
	Current_wallpaper string `json:"current_wallpaper"`
	Init              bool   `json:"init"`
}

type DataAction struct {
	Mode DataActionMode
	Data *Data
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

func DataContent(action DataAction) (*Data, error) {
	data_dir, err := GetOrCreateDataDir()

	if err != nil {
		return nil, err
	}

	data_path := filepath.Join(data_dir, "data.json")

	f, err := os.OpenFile(data_path, os.O_RDWR|os.O_CREATE, 0644)

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

	defer f.Close()

	decoder := json.NewDecoder(f)
	encoder := json.NewEncoder(f)

	var data Data

	switch action.Mode {
	case Read:
		if err := decoder.Decode(&data); err != nil {
			log.Println(err)
			return nil, err
		}
		return &data, nil
	case Write:
		new_data := action.Data
		if err := encoder.Encode(new_data); err != nil {
			log.Println(err)
			return nil, err

		}
		return new_data, nil
	}

	// if err := f.Close(); err != nil {
	// 	log.Fatal(err)
	// }

	return &data, nil

}
