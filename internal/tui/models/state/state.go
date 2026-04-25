package state

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/google/uuid"
)

// AppState holds the global state of the TUI application.
type AppState struct {
	AppAuthState
}

type AppPageState struct {
	pages         map[models.PageId]models.Page
	currentPageId models.PageId
}

type AppPageStater interface {
	HasPages() bool
	NextPage() (tea.Cmd, error)
	PrevPage() (tea.Cmd, error)
	ShowPage(pageId models.PageId) (tea.Cmd, error)
	CurrentPage() models.Page
	AddPages(pages ...models.Page)
}

func NewAppPageState() *AppPageState {
	return &AppPageState{
		pages:         make(map[models.PageId]models.Page),
		currentPageId: models.PageId(uuid.Nil),
	}
}

func (aps *AppPageState) SetCurrentPage(pageId models.PageId) {
	aps.currentPageId = pageId
}

func (avs *AppPageState) HasPages() bool {
	return len(avs.pages) > 0
}

func (avs *AppPageState) AddPages(pages ...models.Page) {
	for _, page := range pages {
		avs.pages[page.ID()] = page
	}
}

func (avs *AppPageState) NextPage() (tea.Cmd, error) {
	if currentPage := avs.CurrentPage(); currentPage != nil {
		if nextPageId := currentPage.NextPageId(); nextPageId != models.PageId(uuid.Nil) {
			return avs.ShowPage(nextPageId)
		}
	}

	return nil, nil
}

func (avs *AppPageState) PrevPage() (tea.Cmd, error) {
	if currentPage := avs.CurrentPage(); currentPage != nil {
		if prevPageId := currentPage.PrevPageId(); prevPageId != models.PageId(uuid.Nil) {
			return avs.ShowPage(prevPageId)
		}
	}

	return nil, nil
}

func (avs *AppPageState) CurrentPage() models.Page {
	currentPage, exists := avs.pages[avs.currentPageId]
	if !exists || currentPage == nil {
		return nil
	}

	return currentPage
}

func (avs *AppPageState) ShowPage(pageId models.PageId) (tea.Cmd, error) {
	page, exists := avs.pages[pageId]
	if !exists || page == nil {
		return nil, fmt.Errorf("page with id %s not found", pageId)
	}

	avs.currentPageId = pageId

	return page.InitOnce(), nil
}
