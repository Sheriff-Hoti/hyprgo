package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Sheriff-Hoti/hyprgo/config"
	"github.com/Sheriff-Hoti/hyprgo/data"
	"github.com/Sheriff-Hoti/hyprgo/tui"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

func main() {
	config_file := flag.String("config", config.GetDefaultConfigPath(), "config file")
	init := false
	// if len(os.Args) > 2 {
	// 	log.Fatal("Too many arguments")
	// }
	if len(os.Args) == 2 {
		init = os.Args[1] == "init"
	}

	flag.Parse()

	config_struct, err := config.ReadConfigFile(*config_file)
	if err != nil {
		log.Fatal(err)
	}

	data, err := data.ReadDataFile(config_struct.Data_dir)
	if err != nil {
		log.Fatal(err)
	}

	//check init here
	if init {
		log.Print(init)
		//and check if data is already initialized
	}

	files, err := config.GetWallpapers(config_struct.Wallpaper_dir)
	if err != nil {
		log.Fatal("Error trying to get wallpapers:", err)
	}

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	fmt.Print("Width:", width, "Height:", height, "\n")

	p := tea.NewProgram(
		tui.NewGrid(files, data.Current_wallpaper, width, height), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
