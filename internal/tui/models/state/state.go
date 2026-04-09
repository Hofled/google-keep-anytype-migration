package state

import "github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"

// AppState holds the global state of the TUI application.
type AppState struct {
	AppAuthState
}

type AppViewState struct {
	pages            []models.Page
	currentViewIndex int
}

type AppViewStater interface {
	HasViews() bool
	NextView()
	PrevView()
	CurrentView() models.Page
	CurrentViewIndex() int
	AddPages(pages ...models.Page)
}

func (avs *AppViewState) HasViews() bool {
	return len(avs.pages) > 0
}

func (avs *AppViewState) AddPages(pages ...models.Page) {
	avs.pages = append(avs.pages, pages...)
}

func (avs *AppViewState) NextView() {
	avs.currentViewIndex = (avs.currentViewIndex + 1) % len(avs.pages)
	avs.pages[avs.currentViewIndex].InitOnce()
}

func (avs *AppViewState) PrevView() {
	avs.currentViewIndex = (avs.currentViewIndex - 1) % len(avs.pages)
	avs.pages[avs.currentViewIndex].InitOnce()
}

func (avs *AppViewState) CurrentView() models.Page {
	if avs.currentViewIndex < len(avs.pages) {
		return avs.pages[avs.CurrentViewIndex()]
	}

	return nil
}

func (avs *AppViewState) CurrentViewIndex() int {
	return avs.currentViewIndex
}
