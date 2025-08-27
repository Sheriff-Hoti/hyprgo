package tui

import (
	"fmt"

	"github.com/Sheriff-Hoti/hyprgo/consts"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(21).
			Height(6).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(21).
				Height(6).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	selectedModelStyle = lipgloss.NewStyle().
				Width(21).
				Height(6).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("100"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	activeDot   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	inactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
)

type model struct {
	paginatedChoices          []string
	choices                   []string // items on the to-do list
	cursor                    int      // which to-do list item our cursor is pointing at
	selected                  int      // which to-do items are selected
	col_num                   int
	page                      int
	totalPage                 int
	size                      int
	onWallpaperSelectCallback func(t int)
	onPageChangeCallback      func(s []string)
}

func InitialModel(choices []string, selected int, page int, callback func(t int), onPageChangeCallback func(s []string)) model {
	// p := paginator.New()
	return model{
		choices:                   choices,
		paginatedChoices:          Paginate(choices, 0, consts.ITEMS_PER_PAGE),
		selected:                  selected,
		col_num:                   consts.CELL_COLS,
		page:                      page,
		size:                      consts.ITEMS_PER_PAGE,
		totalPage:                 totalPageCalculator(consts.ITEMS_PER_PAGE, len(choices)),
		onWallpaperSelectCallback: callback,
		onPageChangeCallback:      onPageChangeCallback,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if (m.cursor > 0) && (m.cursor-m.col_num+1 > 0) {
				m.cursor -= m.col_num
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if (m.cursor < len(m.paginatedChoices)-1) && (m.cursor+m.col_num < len(m.paginatedChoices)) {
				m.cursor += m.col_num
			}

			// The "up" and "k" keys move the cursor up
		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "right", "l":
			if m.cursor < len(m.paginatedChoices)-1 {
				m.cursor++
			}

		case "a":
			if m.page > 0 {
				m.page--
				m.paginatedChoices = Paginate(m.choices, m.page, int(consts.ITEMS_PER_PAGE))
				m.onPageChangeCallback(m.paginatedChoices)
			}
		case "d":
			if m.page < m.totalPage-1 {
				m.page++
				m.paginatedChoices = Paginate(m.choices, m.page, int(consts.ITEMS_PER_PAGE))
				m.onPageChangeCallback(m.paginatedChoices)
			}
		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.selected = m.cursor
			m.onWallpaperSelectCallback(m.selected)
			//here we need to make the change to update the wallpaper into ./local/.share or smth
			//config file prolly ./config/hypr/hyprgo.conf
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	// s := "What should we buy at the market?\n\n"
	s := ""

	accumulator := make([]string, 0, len(m.paginatedChoices))

	for idx, _ := range m.paginatedChoices {
		if idx == m.cursor {
			accumulator = append(accumulator, focusedModelStyle.Render())
		} else if idx == m.selected {
			accumulator = append(accumulator, selectedModelStyle.Render())
		} else {
			accumulator = append(accumulator, modelStyle.Render())

		}
		if (idx+1)%m.col_num == 0 || idx == len(m.paginatedChoices)-1 {
			s += lipgloss.JoinHorizontal(lipgloss.Top, accumulator...)
			s += "\n"
			accumulator = make([]string, 0, len(m.paginatedChoices))
		}
	}

	s += m.paginatorView()
	s += helpStyle.Render("\ntab: focus next • n: new %s • q: exit\n")
	s += fmt.Sprint(m.totalPage)
	s += fmt.Sprint(m.page)

	// The footer

	// Send the UI for rendering
	return s
}

func (m model) paginatorView() string {
	s := ""

	for i := range m.totalPage {
		// s += fmt.Sprint(i)
		if i == m.page {
			s += activeDot
		} else {
			s += inactiveDot

		}
	}

	// for idx := range m.totalPage - 1 {

	// }
	return s
}

func totalPageCalculator(size int, total int) int {
	return (total + size - 1) / size
}

func Paginate[T any](s []T, page int, size int) []T {
	length := len(s)

	if size <= 0 {
		return s
	}

	totalPages := (length + size - 1) / size // round up division

	if totalPages == 0 {
		totalPages = 1
	}

	if size >= length {
		return s
	}

	if page < 0 || page >= totalPages {
		return s
	}

	start := page * size
	end := min(start+size, len(s))

	return s[start:end]
}
