package pages

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/styles"
	"github.com/Hofled/go-google-keep-anytype-migration/pkg/googlekeep"
)

type parsedNotesMsg struct {
	notes []googlekeep.Note
}

type migrateNotesMsg struct {
	notes []googlekeep.Note
}

type selectedDirErrMsg struct {
	err error
}

type NoteSelectModel struct {
	*models.ModelInitOnce
	*models.PageIds

	pageState   state.AppPageStater
	windowState state.AppWindowStater
	notesState  state.NotesStater

	picker filepicker.Model

	isParsingNotes bool
	spinner        spinner.Model

	selectedDir    string
	selectedDirErr error
}

func NewNoteSelectModel(pageState state.AppPageStater, windowState state.AppWindowStater, notesState state.NotesStater) (*NoteSelectModel, error) {
	ids, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	picker := filepicker.New()
	picker.DirAllowed = true
	picker.FileAllowed = false
	picker.CurrentDirectory, _ = os.UserHomeDir()

	m := &NoteSelectModel{
		PageIds:     ids,
		pageState:   pageState,
		windowState: windowState,
		notesState:  notesState,
		picker:      picker,
		spinner:     spinner.New(spinner.WithSpinner(spinner.Dot)),
	}

	m.ModelInitOnce = models.NewModelInitOnce(m)

	return m, nil
}

func (nsm *NoteSelectModel) Init() tea.Cmd {
	nsm.picker.SetHeight(nsm.windowState.GetWindowHeight() - 5)

	return tea.Batch(nsm.picker.Init(), nsm.spinner.Tick)
}

func (nsm *NoteSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case parsedNotesMsg:
		nsm.notesState.SetParsedNotes(msg.notes)
		var err error
		if cmd, err = nsm.pageState.NextPage(); err == nil {
			return nsm, tea.Sequence(cmd, func() tea.Msg {
				return migrateNotesMsg{msg.notes}
			})
		}
	case selectedDirErrMsg:
		nsm.selectedDirErr = msg.err
	case spinner.TickMsg:
		nsm.spinner, cmd = nsm.spinner.Update(msg)
		return nsm, cmd
	}

	nsm.picker, cmd = nsm.picker.Update(msg)

	if didSelect, path := nsm.picker.DidSelectFile(msg); didSelect {
		nsm.selectedDir = path
		return nsm, tea.Batch(cmd, func() tea.Msg {
			nsm.isParsingNotes = true
			notes, err := parseNotes(path)
			nsm.isParsingNotes = false
			if err != nil {
				return selectedDirErrMsg{err}
			} else {
				nsm.selectedDirErr = nil
				return parsedNotesMsg{notes}
			}
		})
	}

	return nsm, cmd
}

func (nsm *NoteSelectModel) View() tea.View {
	var b strings.Builder

	b.WriteString("Select Directory With Google Takeout 'Keep' Export:")

	if len(nsm.selectedDir) != 0 {
		fmt.Fprintf(&b, " %s", nsm.picker.Styles.Selected.Render(nsm.selectedDir))
	}

	b.WriteString("\n")

	if nsm.isParsingNotes {
		fmt.Fprintf(&b, "%s Parsing notes...", nsm.spinner.View())
	} else {
		fmt.Fprintf(&b, "Currently viewing: %s\n", nsm.picker.Styles.Directory.Render(nsm.picker.CurrentDirectory))

		fmt.Fprintf(&b, "\n%s\n", nsm.picker.View())
	}

	if nsm.selectedDirErr != nil {
		b.WriteString(styles.ErrText.Render(nsm.selectedDirErr.Error()))
	}

	v := tea.NewView(b.String())
	v.AltScreen = true

	return v
}

func parseNotes(dirPath string) ([]googlekeep.Note, error) {
	notes := make([]googlekeep.Note, 0)
	var encounteredJson bool

	if err := filepath.WalkDir(dirPath, func(entryPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(entryPath) != ".json" {
			return nil
		}

		encounteredJson = true

		jsonB, err := os.ReadFile(entryPath)
		if err != nil {
			return err
		}

		var note googlekeep.Note
		if unmarshalErr := json.Unmarshal(jsonB, &note); unmarshalErr != nil {
			return unmarshalErr
		}

		notes = append(notes, note)

		return nil
	}); err != nil {
		return nil, err
	}

	if !encounteredJson {
		return nil, errors.New("directory does not contain any valid .json files")
	}

	return notes, nil
}
