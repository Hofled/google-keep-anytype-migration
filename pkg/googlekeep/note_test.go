package googlekeep

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestUnmarshalNote_JSON(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "test", "testdata", "googlekeep")
	testFiles, err := filepath.Glob(filepath.Join(testDataDir, "note*.json"))
	require.NoError(t, err)

	for _, filePath := range testFiles {
		t.Logf("Unmarshaling %s", filePath)
		b, err := os.ReadFile(filePath)
		require.NoError(t, err)

		var note Note
		require.NoErrorf(t, json.Unmarshal(b, &note), "failed unmarshaling %s", filePath)
		require.NotEmpty(t, note)
	}
}
