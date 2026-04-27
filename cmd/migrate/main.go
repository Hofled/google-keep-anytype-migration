package main

import (
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/app"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/pages"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/pages/auth"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/pages/auth/challenge"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/pages/spaces"
)

func main() {
	authState := &state.AppAuthState{}
	pageState := state.NewAppPageState()
	windowState := state.NewAppWindowState()
	importSpacesState := &state.ImportSpacesState{}
	notesState := &state.NotesState{}

	apiKeyAuthPage, err := auth.NewApiKeyAuthPage(authState, pageState) // TODO refactor to lazy page construction
	if err != nil {
		log.Panicln(err)
	}

	challengeAuthPage, err := challenge.NewChallengeAuthPage(authState, pageState)
	if err != nil {
		log.Panicln(err)
	}

	authMethodPage, err := auth.NewMethodPage(pageState, apiKeyAuthPage.ID(), challengeAuthPage.ID())
	if err != nil {
		log.Panicln(err)
	}

	apiKeyAuthPage.SetPrevPage(authMethodPage.ID())
	challengeAuthPage.SetPrevPage(authMethodPage.ID())

	spacesListPage, err := spaces.NewSpacesModel(authState, pageState, windowState, importSpacesState)
	if err != nil {
		log.Panicln(err)
	}
	spacesListPage.SetPrevPage(authMethodPage.ID())

	apiKeyAuthPage.SetNextPage(spacesListPage.ID())
	challengeAuthPage.SetNextPage(spacesListPage.ID())

	notesSelectPage, err := pages.NewNoteSelectModel(pageState, windowState, notesState)
	if err != nil {
		log.Panicln(err)
	}
	notesSelectPage.SetPrevPage(spacesListPage.ID())

	spacesListPage.SetNextPage(notesSelectPage.ID())

	pageState.AddPages(authMethodPage, apiKeyAuthPage, challengeAuthPage, spacesListPage, notesSelectPage)

	pageState.SetCurrentPage(authMethodPage.ID())

	tuiApp := app.NewApp(authState, pageState, windowState)

	p := tea.NewProgram(tuiApp)
	if _, err := p.Run(); err != nil {
		log.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
