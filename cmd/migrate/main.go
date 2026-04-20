package main

import (
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/app"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/pages/auth"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/pages/auth/challenge"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/pages/spaces"
)

func main() {
	appState := &state.AppState{}
	pageState := state.NewAppPageState()
	windowState := state.NewAppWindowState()

	apiKeyAuthPage, err := auth.NewApiKeyAuthPage(appState, pageState) // TODO refactor to lazy page construction
	if err != nil {
		log.Panicln(err)
	}

	challengeAuthPage, err := challenge.NewChallengeAuthPage(appState, pageState)
	if err != nil {
		log.Panicln(err)
	}

	authMethodPage, err := auth.NewMethodPage(pageState, apiKeyAuthPage.ID(), challengeAuthPage.ID())
	if err != nil {
		log.Panicln(err)
	}

	apiKeyAuthPage.SetPrevPage(authMethodPage.ID())
	challengeAuthPage.SetPrevPage(authMethodPage.ID())

	spacesListPage, err := spaces.NewSpacesModel(appState, pageState, windowState)
	if err != nil {
		log.Panicln(err)
	}
	spacesListPage.SetPrevPage(authMethodPage.ID())

	apiKeyAuthPage.SetNextPage(spacesListPage.ID())
	challengeAuthPage.SetNextPage(spacesListPage.ID())

	pageState.AddPages(authMethodPage, apiKeyAuthPage, challengeAuthPage, spacesListPage)

	pageState.SetCurrentPage(authMethodPage.ID())

	tuiApp := app.NewApp(appState, pageState, windowState)

	p := tea.NewProgram(tuiApp)
	if _, err := p.Run(); err != nil {
		log.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
