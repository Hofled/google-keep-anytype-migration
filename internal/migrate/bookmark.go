package migrate

import (
	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype"
	"github.com/Hofled/go-google-keep-anytype-migration/internal/anytype/rest"
	"github.com/Hofled/go-google-keep-anytype-migration/pkg/googlekeep"
)

func AnnotationToBookmark(annotation googlekeep.Annotations) rest.CreateObjectRequest {
	return rest.CreateObjectRequest{
		TypeKey: "bookmark",
		Name:    annotation.Title,
		Body:    annotation.Description,
		Properties: []anytype.PropertyLinkWithValue{
			anytype.NewURLProperty("source", annotation.Url),
		},
	}
}
