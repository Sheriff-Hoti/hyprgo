package pkg

import (
	"fmt"
	"os"
	"os/exec"
)

var (
	base_cmd = "kitten"
	icat_cmd = "icat"
)

type Place struct {
	Width  int
	Height int
	Left   int
	Top    int
}

type ICatOptions struct {
	Stdin          bool
	Scale_up       bool
	Place          Place
	Extra_args     []string
	Wallpaper_path string
}

func IcatCmdHalder(options ICatOptions) {
	stdin := "--stdin=no"
	scale_up := ""
	place := fmt.Sprintf("--place=%vx%v@%vx%v", options.Place.Width, options.Place.Height, options.Place.Left, options.Place.Top)
	if options.Stdin {
		stdin = "--stdin=yes"
	}
	if options.Scale_up {
		scale_up = "--scale-up"
	}

	fullArgs := append([]string{icat_cmd, stdin, scale_up, place}, options.Extra_args...)

	cmd := exec.Command(base_cmd, append(fullArgs, options.Wallpaper_path)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err, cmd.Args)
		return
	}

}
