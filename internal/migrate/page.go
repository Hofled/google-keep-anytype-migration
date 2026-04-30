package migrate

import (
	"fmt"
	"strings"
	"time"

	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
	"github.com/Hofled/go-google-keep-anytype-migration/pkg/googlekeep"
)

func GoogleNoteToCreatePageRequest(note googlekeep.Note) rest.CreateObjectRequest {
	title := note.Title
	if len(title) == 0 {
		title = time.UnixMicro(int64(note.CreatedTimestampUsec)).Format(time.RFC822)
	}

	var bodyBuilder strings.Builder

	if len(note.ListContent) > 0 {
		bodyBuilder.WriteString(listContentToPageBody(note.ListContent))
		bodyBuilder.WriteString("\n---\n")
		for _, a := range note.Annotations {
			fmt.Fprintf(&bodyBuilder, "[%s](%s)\n", a.Title, a.Url)
		}
	} else if len(note.TextContent) > 0 {
		bodyBuilder.WriteString(note.TextContent)
	}

	return rest.CreateObjectRequest{
		TypeKey: "page",
		Name:    title,
		Body:    bodyBuilder.String(),
	}
}

func listContentToPageBody(listContent []googlekeep.ListContent) string {
	var b strings.Builder

	for _, item := range listContent {
		boxMarking := " "
		if item.IsChecked {
			boxMarking = "x"
		}

		fmt.Fprintf(&b, "- [%s] %s\n", boxMarking, item.Text)
	}

	return b.String()
}
