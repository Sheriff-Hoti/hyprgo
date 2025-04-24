package main

import (
	"fmt"
	"os"

	"github.com/Sheriff-Hoti/hyprgo/tui"
	tea "github.com/charmbracelet/bubbletea"
)

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.

func main() {
	p := tea.NewProgram(tui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

//use catimg
