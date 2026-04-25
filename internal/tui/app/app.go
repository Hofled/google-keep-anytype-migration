package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
)

// App manages the overall TUI application state and page navigation.
type App struct {
	authState   state.AppAuthStater
	pageState   state.AppPageStater
	windowState state.AppWindowStater
}

// NewApp creates a new TUI application with the given pages.
func NewApp(authState state.AppAuthStater, viewState state.AppPageStater, windowState state.AppWindowStater) *App {
	return &App{
		authState:   authState,
		pageState:   viewState,
		windowState: windowState,
	}
}

// Init initializes the application.
func (a *App) Init() tea.Cmd {
	if currView := a.pageState.CurrentPage(); currView != nil {
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

	case tea.WindowSizeMsg:
		a.windowState.SetWindowWidth(msg.Width)
		a.windowState.SetWindowHeight(msg.Height)
	}

	// Update the current page
	_, cmd := a.pageState.CurrentPage().Update(msg)
	return a, cmd
}

// View renders the current page.
func (a *App) View() tea.View {
	if !a.pageState.HasPages() {
		return tea.NewView("⛔ No pages available ⛔")
	}

	if currentPage := a.pageState.CurrentPage(); currentPage != nil {
		pageView := currentPage.View()
		return tea.NewView(pageView.Content)
	}

	return tea.NewView("Current page is empty")
}
