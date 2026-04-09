package models

import tea "github.com/charmbracelet/bubbletea"

// Page represents a single page in the TUI application.
// It embeds tea.Model and adds methods for navigation.
type Page interface {
	tea.Model
	// CanProceed returns true if the page allows proceeding to the next page.
	CanProceed() bool
	// GetData returns the data collected by this page.
	GetData() interface{}
}
