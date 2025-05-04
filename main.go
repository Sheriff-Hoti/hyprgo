package main

import (
	"fmt"
	"os"

	"github.com/Sheriff-Hoti/hyprgo/consts"
	"github.com/Sheriff-Hoti/hyprgo/pkg"
	"github.com/Sheriff-Hoti/hyprgo/tui"
	tea "github.com/charmbracelet/bubbletea"
)

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.

func RenderImages(filenames []string) {
	for idx, filename := range filenames {
		pkg.IcatCmdHalder(pkg.ICatOptions{
			Place: pkg.Place{
				Width:  consts.ICAT_IMAGE_WIDTH,
				Height: consts.ICAT_IMAGE_HEIGHT,
				Top:    consts.ICAT_IMAGE_TOP_OFFSET + ((idx / consts.CELL_COLS) * 8),
				Left:   consts.ICAT_IMAGE_LEFT_OFFSET + ((idx % consts.CELL_COLS) * (consts.ICAT_IMAGE_WIDTH + 3)),
			},
			Extra_args:     []string{"--z-index=--1"},
			Scale_up:       true,
			Wallpaper_path: filename,
		})
	}
}

func main() {

	kvpairmap, err := pkg.ReadConfigFile()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
		return
	}

	backend := pkg.InitBackend(kvpairmap)
	fmt.Println(kvpairmap)
	fmt.Println(backend)

	//first read the config then start rendering images
	//  then start the tea program

	filenames, filenames_error := pkg.GetWallpapers("./img")

	if filenames_error != nil {
		fmt.Println("Error:", filenames_error)
		os.Exit(1)
	}

	// RenderImages(filenames)

	// fmt.Print("\033[H")

	p := tea.NewProgram(tui.InitialModel(filenames, 0, func(t int) {
		backend.SetImage(filenames[t])

	}))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

//use icat
//kitten icat --stdin=no --align=left --place=20x20@1x3 ./img/night-sky.jpg
//&& kitten icat --stdin=no --align=left --place=20x20@1x10 ./img/test.jpg
