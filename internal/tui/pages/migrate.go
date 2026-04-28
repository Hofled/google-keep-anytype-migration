package pages

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/googlekeep"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/migrate"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/models/state"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/tui/styles"
)

var (
	spinnerStyle   = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
	doneTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Green)
	dotStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type migrationFinishedMsg struct {
}

type migrationFailedMsg struct {
	err error
}

type MigrationPageModel struct {
	*models.ModelInitOnce
	*models.PageIds

	authState   state.AppAuthStater
	spacesState state.ImportSpacesStater
	notesState  state.NotesStater

	spinner spinner.Model

	done bool

	currMigrationSpaceName string

	currErr error

	migrated []string

	totalNotesCount  int
	currentNoteIndex int
}

const displayedMigratedLines = 8

var dotsLine = dotStyle.Render(strings.Repeat(".", 30))

func NewMigrationPageModel(authState state.AppAuthStater, spacesState state.ImportSpacesStater, notesState state.NotesStater) (*MigrationPageModel, error) {
	ids, err := models.NewPageIds()
	if err != nil {
		return nil, err
	}

	model := &MigrationPageModel{
		PageIds:     ids,
		authState:   authState,
		spacesState: spacesState,
		notesState:  notesState,
		migrated:    make([]string, 0),
		spinner:     spinner.New(spinner.WithSpinner(spinner.MiniDot)),
	}

	model.ModelInitOnce = models.NewModelInitOnce(model)

	return model, nil
}

func (mpg *MigrationPageModel) Init() tea.Cmd {
	return mpg.spinner.Tick
}

func (mpg *MigrationPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		mpg.spinner, cmd = mpg.spinner.Update(msg)

		return mpg, cmd
	case migrateNotesMsg:
		mpg.totalNotesCount = len(msg.notes)
		return mpg, func() tea.Msg {
			if err := mpg.migrateNotes(msg.notes); err != nil {
				return migrationFailedMsg{err}
			} else {
				return migrationFinishedMsg{}
			}
		}
	case migrationFinishedMsg:
		mpg.done = true
		mpg.currErr = nil
		return mpg, tea.Quit
	case migrationFailedMsg:
		mpg.currErr = msg.err
	}

	return mpg, nil
}

func (mpg *MigrationPageModel) View() tea.View {
	var b strings.Builder

	var titleText string

	if mpg.done {
		titleText = doneTitleStyle.Render("Done Migrating!")
	} else {
		titleText = fmt.Sprintf("%s Migrating to %s [%d/%d]:", spinnerStyle.Render(mpg.spinner.View()), mpg.currMigrationSpaceName, mpg.currentNoteIndex+1, mpg.totalNotesCount)
	}

	fmt.Fprintf(&b, "%s\n", titleText)

	for i := displayedMigratedLines - 1; i >= 0; i-- {
		migratedNoteIndex := mpg.currentNoteIndex - 1 - i
		if migratedNoteIndex > len(mpg.migrated)-1 || migratedNoteIndex < 0 {
			b.WriteString(dotsLine)
		} else {
			fmt.Fprintf(&b, "✅ %s", mpg.migrated[migratedNoteIndex])
		}

		b.WriteString("\n")
	}

	if mpg.currErr != nil {
		fmt.Fprintf(&b, "%s\n", styles.ErrText.Render(mpg.currErr.Error()))
	}

	return tea.NewView(b.String())
}

func (mpg *MigrationPageModel) migrateNotes(notes []googlekeep.Note) error {
	for _, s := range mpg.spacesState.SelectedSpaces() {
		mpg.currMigrationSpaceName = s.Name

		collId, err := mpg.createMigrationCollection(s.Id)
		if err != nil {
			return fmt.Errorf("failed creating migration collection: %w", err)
		}

		client := mpg.authState.GetClient()
		if client == nil {
			return errors.New("client is nil")
		}

		ctx := context.Background()

		createdObjIds := make([]string, 0, len(notes))

		mpg.migrated = make([]string, 0, len(notes))

		for i, n := range notes {
			mpg.currentNoteIndex = i
			createObjReq := migrate.GoogleNoteToCreatePageRequest(n)
			resp, err := client.CreateObject(ctx, s.Id, createObjReq)
			if err != nil {
				return fmt.Errorf("failed creating object: %w", err)
			}

			createdObjIds = append(createdObjIds, resp.Object.Id)
			mpg.migrated = append(mpg.migrated, resp.Object.Name)
		}

		if err := client.AddObjectsToList(ctx, s.Id, collId, createdObjIds); err != nil {
			return fmt.Errorf("failed adding objects to list: %w", err)
		}
	}

	return nil
}

func (mpg *MigrationPageModel) createMigrationCollection(spaceId string) (string, error) {
	client := mpg.authState.GetClient()
	if client == nil {
		return "", errors.New("client is nil")
	}

	reqObj := rest.CreateObjectRequest{
		Name:    fmt.Sprintf("Google Keep Migration %s", time.Now().Format(time.RFC822)),
		TypeKey: "collection",
	}

	collectionResp, err := client.CreateObject(context.Background(), spaceId, reqObj)
	if err != nil {
		return "", err
	}

	return collectionResp.Object.Id, nil
}
