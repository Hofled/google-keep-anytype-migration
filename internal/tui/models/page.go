package models

import (
	tea "charm.land/bubbletea/v2"
	"github.com/google/uuid"
)

// Page represents a single page in the TUI application.
// It embeds tea.Model and adds methods for navigation.
type Page interface {
	tea.Model
	ModelInitOncer
	// CanProceed returns true if the page allows proceeding to the next page.
	CanProceed() bool
	// GetData returns the data collected by this page.
	GetData() any
	NextPageId() PageId
	PrevPageId() PageId
	ID() PageId
}

type PageId uuid.UUID

func NewPageID() (PageId, error) {
	pageId, err := uuid.NewUUID()
	return PageId(pageId), err
}
