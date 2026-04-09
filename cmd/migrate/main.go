package main

import (
	"fmt"
	"os"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/app"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/pages"
	tea "charm.land/bubbletea/v2"
)

func main() {
	// Initialize app state
	state := &models.AppState{}

	// Create pages
	authPage := pages.NewAuthPage()

	// Create app with pages
	tuiApp := app.NewApp([]models.Page{authPage}, state)

	// Run the TUI
	p := tea.NewProgram(tuiApp)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
