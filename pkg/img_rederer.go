package pkg

// https://www.sohamkamani.com/golang/options-pattern/

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

type ICatCmd struct {
	Agregated_args []string
}

type ICatOption func(*ICatCmd)

func WithStdIn(yes bool) ICatOption {
	return func(ic *ICatCmd) {
		if yes {
			ic.Agregated_args = append(ic.Agregated_args, "--stdin=yes")
		} else {
			ic.Agregated_args = append(ic.Agregated_args, "--stdin=no")
		}
	}
}

func WithScaleUp() ICatOption {
	return func(ic *ICatCmd) {
		ic.Agregated_args = append(ic.Agregated_args, "--scale-up")
	}
}

func WithPlace(place Place) ICatOption {
	return func(ic *ICatCmd) {
		ic.Agregated_args = append(ic.Agregated_args, fmt.Sprintf("--place=%vx%v@%vx%v", place.Width, place.Height, place.Left, place.Top))
	}
}

func WithExtraArgs(args ...string) ICatOption {
	return func(ic *ICatCmd) {
		ic.Agregated_args = append(ic.Agregated_args, args...)
	}
}

func withWallpaperPath(path string) ICatOption {
	return func(ic *ICatCmd) {
		ic.Agregated_args = append(ic.Agregated_args, path)
	}
}

func ICatCmdBuilder(path string, opts ...ICatOption) {

	cmdArgs := &ICatCmd{
		Agregated_args: []string{icat_cmd},
	}

	for _, opt := range opts {
		opt(cmdArgs)
	}

	withWallpaperPath(path)(cmdArgs)

	cmd := exec.Command(base_cmd, cmdArgs.Agregated_args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err, cmd.Args)
		return
	}

}
