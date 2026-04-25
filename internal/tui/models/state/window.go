package state

type AppWindowStater interface {
	GetWindowHeight() int
	GetWindowWidth() int
	SetWindowHeight(int)
	SetWindowWidth(int)
}

type AppWindowState struct {
	width, height int
}

func NewAppWindowState() *AppWindowState {
	return &AppWindowState{}
}

func (aws *AppWindowState) GetWindowHeight() int {
	return aws.height
}

func (aws *AppWindowState) GetWindowWidth() int {
	return aws.width
}

func (aws *AppWindowState) SetWindowHeight(h int) {
	aws.height = h
}

func (aws *AppWindowState) SetWindowWidth(w int) {
	aws.width = w
}
