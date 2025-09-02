package tui

import (
	"os"

	"github.com/Sheriff-Hoti/hyprgo/kitty"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Grid struct {
	files         [][]string
	selected_file string
	windowWidth   int
	windowHeight  int
	keys          *listKeyMap
	columns       uint32
	rows          uint32
	cursor_cell   kitty.Cell
	hidden        bool
	pages         []Page
	page_cursor   int
	page_total    int
	grid_opts     kitty.KittyGridOpts
	paginator     paginator.Model
}

func (g *Grid) Init() tea.Cmd {

	return func() tea.Msg {
		firstPage := g.GetFirstPage()
		g.paginator.SetTotalPages(len(g.pages))
		_ = firstPage.RenderImages(os.Stdout)
		return nil
	}
}

func NewGrid(absfiles []string, selected_file string, init_term_width int, init_term_height int) *Grid {

	cols := 3
	rows := 3
	pageSize := cols * rows

	// chunk absfiles into pages where each inner slice has up to pageSize entries
	chunked := make([][]string, 0, (len(absfiles)+pageSize-1)/pageSize)
	for i := 0; i < len(absfiles); i += pageSize {
		end := min(i+pageSize, len(absfiles))
		chunked = append(chunked, absfiles[i:end])
	}

	p := paginator.New()
	p.Type = paginator.Dots
	// p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	return &Grid{
		windowWidth:  init_term_width,
		windowHeight: init_term_height,
		files:        chunked,
		keys:         newListKeyMap(),
		columns:      3,
		rows:         3,
		paginator:    p,
		grid_opts: kitty.KittyGridOpts{
			Cols:        3,
			Rows:        3,
			RowsSpacing: 1,
			ColsSpacing: 1,
			ImgWidth:    uint32((init_term_width / 3) - 2),
			ImgHeight:   uint32((init_term_height / 3) - 2),
			TopSpacing:  2,
			LeftSpacing: 3,
		},
	}
}

func (g *Grid) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd
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

		case key.Matches(msg, g.keys.hide):
			// selected_cell := g.cells[5]
			if g.hidden {
				// g.pages[1].Hide(os.Stdout)
				g.pages[0].Show(os.Stdout)
			} else {
				g.pages[0].Hide(os.Stdout)
				// g.pages[1].Show(os.Stdout)

			}
			g.hidden = !g.hidden

			return g, nil
		}
	}

	paginator, cmd := g.paginator.Update(msg)
	cmds = append(cmds, cmd)
	g.paginator = paginator
	return g, tea.Batch(cmds...)
}

func (g *Grid) View() string {

	background := lipgloss.NewStyle().
		Width(g.windowWidth).
		Height(g.windowHeight).
		Background(lipgloss.Color("235")).Render("")

	// square := lipgloss.NewStyle().
	// 	Width(int(selected_cell.Width)).
	// 	Height(int(selected_cell.Height)).
	// 	// Background(lipgloss.Color("12")).        // blue square
	// 	MarginTop(int(selected_cell.RowCell-2)). // y position
	// 	MarginLeft(int(selected_cell.ColCell-2)).
	// 	Border(lipgloss.RoundedBorder(), true).
	// 	Render("")
	// x position
	// you can also set MarginRight/MarginBottom if needed
	// s := ""

	// if len(g.pages) > 0 {
	// 	for _, cells := range g.pages[0].cells {
	// 		s += fmt.Sprintf(" %s", cells.Filename)
	// 	}
	// }
	view := lipgloss.JoinVertical(lipgloss.Center, background, g.paginator.View())
	return view

}

func (g *Grid) ComputePages() []Page {
	pageSize := int(g.columns) * int(g.rows)
	fileLen := len(g.files)

	if fileLen == 0 || pageSize == 0 {
		return nil
	}

	pageCount := (fileLen + pageSize - 1) / pageSize // ceiling division

	g.pages = make([]Page, 0, pageCount)

	for idx, files := range g.files {

		page := NewPage(idx, files, g.grid_opts)
		g.pages = append(g.pages, *page)
	}

	return g.pages
}

func (g *Grid) GetFirstPage() Page {
	// compute all pages
	pages := g.ComputePages()

	if len(pages) > 0 {
		return pages[0]
	}

	// fallback: if somehow no pages were created, make one manually
	firstPage := NewPage(0, []string{}, g.grid_opts)
	g.pages = []Page{*firstPage}
	return *firstPage
}
