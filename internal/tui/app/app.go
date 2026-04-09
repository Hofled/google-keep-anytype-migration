package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
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
		case "ctrl+c":
			return a, tea.Quit
		}
	}

	// Update the current page
	model, cmd := a.pages[a.current].Update(msg)
	a.pages[a.current] = model.(models.Page)
	return a, cmd
}

// View renders the current page.
func (a *App) View() tea.View {
	if len(a.pages) == 0 {
		return tea.NewView("No pages available")
	}
	return a.pages[a.current].View()
}
