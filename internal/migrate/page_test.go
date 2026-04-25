package migrate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/googlekeep"
	"github.com/stretchr/testify/require"
)

func TestGoogleNoteToCreatePageRequest(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "test", "testdata", "googlekeep")
	testFile := filepath.Join(testDataDir, "note_listcontent.json")
	require.FileExists(t, testFile)

	fileContent, err := os.ReadFile(testFile)
	require.NoError(t, err)

	var note googlekeep.Note
	require.NoError(t, json.Unmarshal(fileContent, &note))

	createPageReq := GoogleNoteToCreatePageRequest(note)
	require.Equal(t, "page", createPageReq.TypeKey)
	require.Equal(t, note.Title, createPageReq.Name)
	require.Equal(t, "- [ ] Content 1\n- [x] Content 2", createPageReq.Body)
}

func TestGoogleNoteToCreatePageRequest_NoTitle(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "test", "testdata", "googlekeep")
	testFile := filepath.Join(testDataDir, "note_no_title.json")
	require.FileExists(t, testFile)

	fileContent, err := os.ReadFile(testFile)
	require.NoError(t, err)

	var note googlekeep.Note
	require.NoError(t, json.Unmarshal(fileContent, &note))

	createPageReq := GoogleNoteToCreatePageRequest(note)
	require.Equal(t, "page", createPageReq.TypeKey)
	createdDateStr := time.UnixMicro(int64(note.CreatedTimestampUsec)).UTC().Format(time.RFC822)
	require.Equal(t, createdDateStr, createPageReq.Name)
}
