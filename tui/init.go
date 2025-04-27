package tui

import (
	"github.com/Sheriff-Hoti/hyprgo/pkg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	choices   = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
)

type model struct {
	choices   []string // items on the to-do list
	cursor    int      // which to-do list item our cursor is pointing at
	selected  int      // which to-do items are selected
	term_info *pkg.TerminalInfo
}

func InitialModel() model {
	return model{
		choices:   choices,
		selected:  0,
		term_info: nil,
	}
}

// here we ll need to read the images and parse them to kitty image
// and also to read configs from config file
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
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.selected = m.cursor
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

	accumulator := make([]string, 0, len(choices))

	for idx, val := range m.choices {
		if idx == m.cursor {
			accumulator = append(accumulator, focusedModelStyle.Render(val))
		} else {
			accumulator = append(accumulator, modelStyle.Render(val))

		}
		if (idx+1)%3 == 0 || idx == len(m.choices)-1 {
			s += lipgloss.JoinHorizontal(lipgloss.Top, accumulator...)
			s += "\n"
			accumulator = make([]string, 0, len(choices))
		}
	}

	s += helpStyle.Render("\ntab: focus next • n: new %s • q: exit\n")

	// The footer

	// Send the UI for rendering
	return s
}
