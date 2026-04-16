package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/app"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/pages/auth"
)

func main() {
	appState := &state.AppState{}
	viewState := &state.AppViewState{}

	authPage := auth.NewAuthPage(appState, viewState) // TODO refactor to lazy page construction

	viewState.AddPages(authPage)

	tuiApp := app.NewApp(appState, viewState)

	p := tea.NewProgram(tuiApp)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
