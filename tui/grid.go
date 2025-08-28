package tui

import (
	"fmt"

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
	row_spacing   int
	col_spacing   int
	padding       int
	cell_style    lipgloss.Style
}

func (g *Grid) Init() tea.Cmd {
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
		row_spacing:  1,
		col_spacing:  1,
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

	// return fmt.Sprint("Grid dims:", g.windowWidth, "x", g.windowHeight)

	rowsCount := 4
	colsCount := 4

	if rowsCount <= 0 {
		rowsCount = 1
	}
	if colsCount <= 0 {
		colsCount = 1
	}

	cellWidth := (g.windowWidth - (colsCount*2)*g.col_spacing) / colsCount
	cellHeight := (g.windowHeight - (rowsCount*2)*g.row_spacing) / rowsCount

	cellStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(cellWidth-2).
		Height(cellHeight-2).
		Align(lipgloss.Center, lipgloss.Center).
		Margin(g.row_spacing, g.col_spacing)

	// Create cells
	cells := make([]string, rowsCount*colsCount)
	for i := range rowsCount * colsCount {
		cells[i] = cellStyle.Render(fmt.Sprintf("%d", i+1))
	}

	// Build rows
	rows := make([]string, rowsCount)
	for r := range rowsCount {
		start := r * colsCount
		end := start + colsCount
		rows[r] = lipgloss.JoinHorizontal(lipgloss.Center, cells[start:end]...)
	}

	// Join rows vertically
	grid := lipgloss.JoinVertical(lipgloss.Center, rows...)

	// Grid-level style
	gridStyle := lipgloss.NewStyle().
		// Padding(4, 8). // top/bottom, left/right
		// Margin(10, 10). // space outside the grid
		Align(lipgloss.Center, lipgloss.Center).Padding(2)

	return gridStyle.Render(grid)
}
