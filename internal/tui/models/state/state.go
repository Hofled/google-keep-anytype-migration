package state

import (
	"fmt"

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
	NextPage() error
	PrevPage() error
	ShowPage(pageId models.PageId) error
	CurrentPage() models.Page
	AddPages(pages ...models.Page)
}

func NewAppPageState() *AppPageState {
	return &AppPageState{
		pages:         make(map[models.PageId]models.Page),
		currentPageId: models.PageId(uuid.Nil),
	}
}

func (avs *AppPageState) HasPages() bool {
	return len(avs.pages) > 0
}

func (avs *AppPageState) AddPages(pages ...models.Page) {
	for _, page := range pages {
		avs.pages[page.ID()] = page
	}
}

func (avs *AppPageState) NextPage() error {
	if currentPage := avs.CurrentPage(); currentPage != nil {
		if nextPageId := currentPage.NextPageId(); nextPageId != models.PageId(uuid.Nil) {
			return avs.ShowPage(nextPageId)
		}
	}

	return nil
}

func (avs *AppPageState) PrevPage() error {
	if currentPage := avs.CurrentPage(); currentPage != nil {
		if prevPageId := currentPage.PrevPageId(); prevPageId != models.PageId(uuid.Nil) {
			return avs.ShowPage(prevPageId)
		}
	}

	return nil
}

func (avs *AppPageState) CurrentPage() models.Page {
	currentPage, exists := avs.pages[avs.currentPageId]
	if !exists || currentPage == nil {
		return nil
	}

	return currentPage
}

func (avs *AppPageState) ShowPage(pageId models.PageId) error {
	page, exists := avs.pages[pageId]
	if !exists || page == nil {
		return fmt.Errorf("page with id %s not found", pageId)
	}

	avs.currentPageId = pageId

	page.InitOnce()

	return nil
}
