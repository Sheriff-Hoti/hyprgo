package main

import (
	"flag"
	"log"
	"os"

	"github.com/Sheriff-Hoti/hyprgo/config"
	"github.com/Sheriff-Hoti/hyprgo/data"
)

func main() {
	config_file := flag.String("config", config.GetDefaultConfigPath(), "config file")
	init := false
	if len(os.Args) > 2 {
		log.Fatal("Too many arguments")
	}
	if len(os.Args) == 2 {
		init = os.Args[1] == "init"
	}

	flag.Parse()

	config, err := config.ReadConfigFile(*config_file)
	if err != nil {
		log.Fatal(err)
	}

	data, err := data.ReadDataFile(config.Data_dir)
	if err != nil {
		log.Fatal(err)
	}

	//check init here
	if init {
		//and check if data is already initialized
	}

	log.Printf("Using config: %+v\ninit: %v", config, init)
	log.Printf("Using data: %+v", data)

}
