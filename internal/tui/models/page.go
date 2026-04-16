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
	PageIder
	// CanProceed returns true if the page allows proceeding to the next page.
	CanProceed() bool
	// GetData returns the data collected by this page.
	GetData() any
}

type PageId uuid.UUID

func newPageID() (PageId, error) {
	pageId, err := uuid.NewUUID()
	return PageId(pageId), err
}

type PageIder interface {
	NextPageId() PageId
	PrevPageId() PageId
	SetNextPage(pageId PageId)
	SetPrevPage(pageId PageId)
	ID() PageId
}

type PageIds struct {
	id         PageId
	nextPageId *PageId
	prevPageId *PageId
}

func NewPageIds() (*PageIds, error) {
	id, err := newPageID()
	if err != nil {
		return nil, err
	}

	return &PageIds{
		id:         id,
		nextPageId: new(PageId),
		prevPageId: new(PageId),
	}, nil
}

func (pi *PageIds) NextPageId() PageId {
	if pi.nextPageId != nil {
		return *pi.nextPageId
	} else {
		return PageId(uuid.Nil)
	}
}

func (pi *PageIds) PrevPageId() PageId {
	if pi.prevPageId != nil {
		return *pi.prevPageId
	} else {
		return PageId(uuid.Nil)
	}
}

func (pi *PageIds) SetNextPage(pageId PageId) {
	*pi.nextPageId = pageId
}

func (pi *PageIds) SetPrevPage(pageId PageId) {
	*pi.prevPageId = pageId
}

func (pi *PageIds) ID() PageId {
	return pi.id
}
