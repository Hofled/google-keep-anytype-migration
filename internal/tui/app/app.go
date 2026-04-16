package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
)

var appStyle = lipgloss.NewStyle().Margin(1, 1, 1, 1)

// App manages the overall TUI application state and page navigation.
type App struct {
	state     *state.AppState
	viewState *state.AppPageState
}

// NewApp creates a new TUI application with the given pages.
func NewApp(state *state.AppState, viewState *state.AppPageState) *App {
	return &App{
		state:     state,
		viewState: viewState,
	}
}

// Init initializes the application.
func (a *App) Init() tea.Cmd {
	if currView := a.viewState.CurrentPage(); currView != nil {
		return currView.InitOnce()
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
	_, cmd := a.viewState.CurrentPage().Update(msg)
	return a, cmd
}

// View renders the current page.
func (a *App) View() tea.View {
	if !a.viewState.HasPages() {
		return tea.NewView("⛔ No pages available ⛔")
	}

	pageView := a.viewState.CurrentPage().View()
	return tea.NewView(appStyle.Render(pageView.Content))
}
