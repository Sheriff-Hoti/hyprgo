package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sheriff-Hoti/hyprgo/consts"
)

func GetWallpapers(dir string) (filenames []string, erro error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	file_names := make([]string, 0, 10)

	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				fullPath := filepath.Join(dir, entry.Name())
				file_names = append(file_names, fullPath)
			}
		}
	}

	return file_names, nil
}

func ReadConfigFile(config *string) (map[string]string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	config_path := fmt.Sprintf("%v/%v", home_dir, consts.CONFIG_PATH)
	log.Println(config_path)
	if _, err := os.Stat(config_path); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist and if it does not exists just return the defaults
		log.Print("congrats 404")

		return GetDefaultConfigVals(), nil
	}

	file, err := os.Open(config_path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	kvpairmap := make(map[string]string, 0)

	for scanner.Scan() {
		key, val, err := ExtractKVPair(scanner.Text())
		if err != nil {
			log.Println(err)
		} else {
			kvpairmap[key] = val
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return kvpairmap, nil

	// this method vill return a valid kv pair
}

func ExtractKVPair(line string) (string, string, error) {
	split := strings.Split(line, "=")
	if len(split) != 2 {
		return "", "", errors.New("it should be 2")
	}
	key := strings.Trim(split[0], " ")
	value := strings.Trim(split[1], " ")

	if key == "" {
		return "", "", errors.New("key must not be empty")
	}

	if value == "" {
		return "", "", errors.New("value must not be empty")
	}

	return key, value, nil
}

func GetDefaultConfigPath() string {
	xdg_config_home := "XDG_CONFIG_HOME"
	home := "HOME"

	if _, ok := os.LookupEnv(xdg_config_home); ok {
		return os.ExpandEnv(filepath.Join(fmt.Sprintf("$%v", xdg_config_home), "hyprgo.conf"))
	}
	return os.ExpandEnv(filepath.Join(fmt.Sprintf("$%v", home), ".config", "hyprgo.conf"))

}

func GetDefaultConfigVals() map[string]string {
	return map[string]string{
		"backend":       "swaync",
		"wallpaper_dir": os.ExpandEnv("$HOME"),
		"data_dir":      GetDefaultDataPath(),
	}
}

// after this it will be ging to other validator metho
// first the default config and then each config function will get
// the map value and pull their own key
// also add a map func i guess ??
