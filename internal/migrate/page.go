package migrate

import (
	"fmt"
	"strings"
	"time"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/googlekeep"
)

func GoogleNoteToCreatePageRequest(note googlekeep.Note) rest.CreateObjectRequest {
	title := note.Title
	if len(title) == 0 {
		title = time.UnixMicro(int64(note.CreatedTimestampUsec)).UTC().Format(time.RFC822)
	}

	body := note.TextContent
	if len(note.ListContent) > 0 {
		body = listContentToPageBody(note.ListContent)
	}

	return rest.CreateObjectRequest{
		TypeKey: "page",
		Name:    title,
		Body:    body,
	}
}

func listContentToPageBody(listContent []googlekeep.ListContent) string {
	var b strings.Builder

	for i, item := range listContent {
		boxMarking := " "
		if item.IsChecked {
			boxMarking = "x"
		}

		b.WriteString(fmt.Sprintf("- [%s] %s", boxMarking, item.Text))

		if i < len(listContent)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}
