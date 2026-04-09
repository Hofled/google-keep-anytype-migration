package models

import tea "charm.land/bubbletea/v2"

// Page represents a single page in the TUI application.
// It embeds tea.Model and adds methods for navigation.
type Page interface {
	tea.Model
	ModelOnceIniter
	// CanProceed returns true if the page allows proceeding to the next page.
	CanProceed() bool
	// GetData returns the data collected by this page.
	GetData() any
	NextPage()
	PrevPage()
}
