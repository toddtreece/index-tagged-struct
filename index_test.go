package index_test

import (
	"testing"

	segment "github.com/blugelabs/bluge_segment_api"
	"github.com/stretchr/testify/require"
	index "github.com/toddtreece/index-tagged-struct"
)

func TestStructTags(t *testing.T) {
	d := index.DashboardObjectSummary{
		ObjectSummary: index.ObjectSummary{
			Name:        "Name",
			Description: "test",
		},
		Tags: []string{"tag1"},
	}

	index := index.NewIndexDocument("1234")
	index.Parse(d)
	doc := index.ToBlugeDocument()

	names := []string{}
	doc.EachField(func(field segment.Field) {
		names = append(names, field.Name())
	})
	require.Equal(t, []string{"_id", "name", "tags"}, names)
}
