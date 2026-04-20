package list

import (
	"fmt"
	"io"
	"strings"

	bubblesList "charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/google/uuid"
)

type MultiSelectItemer interface {
	Id() uuid.UUID
	bubblesList.DefaultItem
}

type MultiSelectItem struct {
	id   uuid.UUID
	item bubblesList.DefaultItem
}

func NewMultiSelectItem(item bubblesList.DefaultItem) (MultiSelectItemer, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	return &MultiSelectItem{
		id:   id,
		item: item,
	}, nil
}

func (msi *MultiSelectItem) FilterValue() string {
	return msi.item.FilterValue()
}

func (msi *MultiSelectItem) Id() uuid.UUID {
	return msi.id
}

func (msi *MultiSelectItem) Title() string {
	return msi.item.Title()
}

func (msi *MultiSelectItem) Description() string {
	return msi.item.Description()
}

type MultiSelectDelegate struct {
	bubblesList.DefaultDelegate

	selectedIds map[uuid.UUID]bool
}

func NewMultiSelectDelegate() *MultiSelectDelegate {
	d := bubblesList.NewDefaultDelegate()
	return &MultiSelectDelegate{
		DefaultDelegate: d,
		selectedIds:     make(map[uuid.UUID]bool),
	}
}

func (d *MultiSelectDelegate) SetSelected(id uuid.UUID, selected bool) {
	d.selectedIds[id] = selected
}

func (d *MultiSelectDelegate) IsSelected(id uuid.UUID) bool {
	return d.selectedIds[id]
}

func (d *MultiSelectDelegate) Render(w io.Writer, m bubblesList.Model, index int, item bubblesList.Item) {
	if i, ok := item.(MultiSelectItemer); ok {
		isSelected := d.selectedIds[i.Id()]
		var prefix string
		if isSelected {
			prefix = fmt.Sprintf("%s ", selectedBox)
		} else {
			prefix = fmt.Sprintf("%s ", deselectedBox)
		}

		title := prefix + i.Title()
		desc := i.Description()

		var matchedRunes []int
		s := &d.Styles

		if m.Width() <= 0 {
			return
		}

		textWidth := m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight() - 4
		title = ansi.Truncate(title, textWidth, ellipsis)
		if d.ShowDescription {
			var lines []string
			for i, line := range strings.Split(desc, "\n") {
				if i >= d.Height()-1 {
					break
				}
				lines = append(lines, ansi.Truncate(line, textWidth, ellipsis))
			}
			desc = strings.Join(lines, "\n")
		}

		var (
			isCursor    = index == m.Index()
			emptyFilter = m.FilterState() == bubblesList.Filtering && m.FilterValue() == ""
			isFiltered  = m.FilterState() == bubblesList.Filtering || m.FilterState() == bubblesList.FilterApplied
		)

		if isFiltered && index < len(m.VisibleItems()) {
			matchedRunes = m.MatchesForItem(index)
		}

		if emptyFilter {
			title = s.DimmedTitle.Render(title)
			desc = s.DimmedDesc.Render(desc)
		} else if isCursor && m.FilterState() != bubblesList.Filtering {
			if isFiltered {
				unmatched := s.SelectedTitle.Inline(true)
				matched := unmatched.Inherit(s.FilterMatch)
				title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
			}
			title = s.SelectedTitle.Render(title)
			desc = s.SelectedDesc.Render(desc)
		} else {
			if isFiltered {
				unmatched := s.NormalTitle.Inline(true)
				matched := unmatched.Inherit(s.FilterMatch)
				title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
			}
			title = s.NormalTitle.Render(title)
			desc = s.NormalDesc.Render(desc)
		}

		if d.ShowDescription {
			fmt.Fprintf(w, "%s\n%s", title, desc)
			return
		}
		fmt.Fprintf(w, "%s", title)
	}
}

type MultiSelectModel struct {
	bubblesList.Model

	multiSelectDelegate *MultiSelectDelegate
}

func NewMultiSelect(items []bubblesList.DefaultItem, width, height int) (*MultiSelectModel, error) {
	d := NewMultiSelectDelegate()

	multiSelectItems := make([]bubblesList.Item, len(items))
	for i, v := range items {
		multiSelectItem, err := NewMultiSelectItem(v)
		if err != nil {
			return nil, err
		}

		multiSelectItems[i] = multiSelectItem
	}

	return &MultiSelectModel{
		Model:               bubblesList.New(multiSelectItems, d, width, height),
		multiSelectDelegate: d,
	}, nil
}

func (m *MultiSelectModel) ToggleSelection() {
	selectedItem := m.SelectedItem()
	if mi, ok := selectedItem.(MultiSelectItemer); ok {
		id := mi.Id()
		current := m.multiSelectDelegate.IsSelected(id)
		m.multiSelectDelegate.SetSelected(id, !current)
	}
}

func (m MultiSelectModel) SelectedItems() []bubblesList.Item {
	var selectedItems []bubblesList.Item
	allItems := m.Items()
	for _, i := range allItems {
		if mi, ok := i.(*MultiSelectItem); ok {
			if m.multiSelectDelegate.IsSelected(mi.Id()) {
				selectedItems = append(selectedItems, mi.item)
			}
		}
	}
	return selectedItems
}

func (m *MultiSelectModel) Update(msg tea.Msg) (*MultiSelectModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msgT := msg.(type) {
	case tea.KeyPressMsg:
		switch msgT.Code {
		case tea.KeySpace:
			if m.FilterState() != bubblesList.Filtering {
				m.ToggleSelection()
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		m.SetSize(msgT.Width, msgT.Height)
		return m, nil
	}

	m.Model, cmd = m.Model.Update(msg)

	return m, cmd
}
