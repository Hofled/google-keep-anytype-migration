package widgets

type Focusable interface {
	SetFocus(f bool)
	Focused() bool
}

type FocusableWidget struct {
	focused bool
}

func (fw *FocusableWidget) SetFocus(f bool) {
	fw.focused = f
}

func (fw *FocusableWidget) Focused() bool {
	return fw.focused
}
