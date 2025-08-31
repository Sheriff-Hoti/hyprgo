package tui

import (
	"fmt"
	"os"

	"github.com/Sheriff-Hoti/hyprgo/kitty"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Grid struct {
	files         []string
	selected_file string
	windowWidth   int
	windowHeight  int
	keys          *listKeyMap
	columns       uint32
	rows          uint32
	cell_style    lipgloss.Style
	cells         []kitty.WriteFileResult
}

func (g *Grid) Init() tea.Cmd {

	grid_opts := kitty.KittyGridOpts{
		Cols:        3,
		Rows:        3,
		RowsSpacing: 1,
		ColsSpacing: 1,
		ImgWidth:    (g.windowWidth / 3) - 2,
		ImgHeight:   (g.windowHeight / 3) - 3,
		TopSpacing:  2,
		LeftSpacing: 3,
	}

	res, err := kitty.KittyWriteFiles(os.Stdout, []string{
		"./test_assets/img/test0.jpg",
		"./test_assets/img/test1.jpg",
		"./test_assets/img/test2.jpg",
		"./test_assets/img/test3.jpg",
		"./test_assets/img/test4.jpg",
		"./test_assets/img/test5.png",
		"./test_assets/img/test5.png",
		"./test_assets/img/test5.png",
		"./test_assets/img/test5.png",
		"./test_assets/img/test5.png",
		"./test_assets/img/test5.png",
		"./test_assets/img/test5.png",
		"./test_assets/img/test5.png",
		"./test_assets/img/test5.png",
		"./test_assets/img/test5.png",
		"./test_assets/img/test5.png",
	}, grid_opts)

	if err != nil {
		fmt.Println("Error writing files:", err)
	} else {
		g.cells = res
	}
	return nil
}

func NewGrid(files []string, selected_file string, init_term_width int, init_term_height int) *Grid {
	return &Grid{
		windowWidth:  init_term_width,
		windowHeight: init_term_height,
		files:        files,
		keys:         newListKeyMap(),
		columns:      3,
		rows:         3,
		cell_style: lipgloss.
			NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(10).Height(5).
			Align(lipgloss.Center, lipgloss.Center),
	}
}

func (g *Grid) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		g.windowWidth = msg.Width
		g.windowHeight = msg.Height
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch {
		case key.Matches(msg, g.keys.quit):
			return g, tea.Quit
		}
	}
	return g, nil
}

func (g *Grid) View() string {

	rowsCount := 3
	colsCount := 3

	selected_cell := g.cells[8]

	if rowsCount <= 0 {
		rowsCount = 1
	}
	if colsCount <= 0 {
		colsCount = 1
	}
	background := lipgloss.NewStyle().
		Width(g.windowWidth).
		Height(g.windowHeight).
		Background(lipgloss.Color("235"))

	square := lipgloss.NewStyle().
		Width(int(selected_cell.Width)).
		Height(int(selected_cell.Height)).
		// Background(lipgloss.Color("12")).        // blue square
		MarginTop(int(selected_cell.RowCell-2)). // y position
		MarginLeft(int(selected_cell.ColCell-2)).
		Border(lipgloss.RoundedBorder(), true).
		Render("")
		// x position
	// you can also set MarginRight/MarginBottom if needed
	return background.Render(square)

}
