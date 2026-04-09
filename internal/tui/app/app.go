package app

import (
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	tea "github.com/charmbracelet/bubbletea"
)

// App manages the overall TUI application state and page navigation.
type App struct {
	pages   []models.Page
	current int
	state   *models.AppState
}

// NewApp creates a new TUI application with the given pages.
func NewApp(pages []models.Page, state *models.AppState) *App {
	return &App{
		pages:   pages,
		current: 0,
		state:   state,
	}
}

// Init initializes the application.
func (a *App) Init() tea.Cmd {
	if len(a.pages) > 0 {
		return a.pages[0].Init()
	}
	return nil
}

// Update handles messages and updates the current page.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		case "right", "l":
			if a.current < len(a.pages)-1 && a.pages[a.current].CanProceed() {
				a.current++
				return a, a.pages[a.current].Init()
			}
		case "left", "h":
			if a.current > 0 {
				a.current--
				return a, a.pages[a.current].Init()
			}
		}
	}

	// Update the current page
	model, cmd := a.pages[a.current].Update(msg)
	a.pages[a.current] = model.(models.Page)
	return a, cmd
}

// View renders the current page.
func (a *App) View() string {
	if len(a.pages) == 0 {
		return "No pages available"
	}
	return a.pages[a.current].View()
}
