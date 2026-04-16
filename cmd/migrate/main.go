package main

import (
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/app"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/pages/auth"
)

func main() {
	appState := &state.AppState{}
	pageState := state.NewAppPageState()

	apiKeyAuthPage, err := auth.NewApiKeyAuthPage(appState, pageState) // TODO refactor to lazy page construction
	if err != nil {
		log.Panicln(err)
	}

	authMethodPage, err := auth.NewMethodPage(pageState, apiKeyAuthPage.ID())
	if err != nil {
		log.Panicln(err)
	}

	apiKeyAuthPage.SetPrevPage(authMethodPage.ID())

	pageState.AddPages(authMethodPage, apiKeyAuthPage)

	pageState.ShowPage(authMethodPage.ID())

	tuiApp := app.NewApp(appState, pageState)

	p := tea.NewProgram(tuiApp)
	if _, err := p.Run(); err != nil {
		log.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
