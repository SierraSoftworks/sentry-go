package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleTags() {
	cl := NewClient(
		// You can specify tags when creating your client
		Tags(map[string]string{
			"redis": "v1",
			"mgo":   "v2",
		}),
	)

	cl.Capture(
		// And override or expand on them when sending an event
		Tags(map[string]string{
			"redis":     "v2",
			"sentry-go": "v1",
		}),
	)
}

func TestTags(t *testing.T) {
	assert.Nil(t, Tags(nil), "it should return nil if the data provided is nil")

	data := map[string]string{
		"redis": "1.0.0",
	}

	o := Tags(data)
	require.NotNil(t, o, "it should not return nil if the data is non-nil")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "tags", o.Class(), "it should use the right option class")

	if assert.Implements(t, (*MergeableOption)(nil), o, "it should implement the MergeableOption interface") {
		t.Run("Merge()", func(t *testing.T) {
			om := o.(MergeableOption)

			assert.Equal(t, o, om.Merge(&testOption{}), "it should replace the old option if it is not recognized")

			t.Run("different entries", func(t *testing.T) {
				data2 := map[string]string{
					"pgsql": "5.4.0",
				}
				o2 := Tags(data2)
				require.NotNil(t, o2, "the second module option should not be nil")

				oo := om.Merge(o2)
				require.NotNil(t, oo, "it should not return nil when it merges")

				ooi, ok := oo.(*tagsOption)
				require.True(t, ok, "it should actually be a *tagsOption")

				if assert.Contains(t, ooi.tags, "redis", "it should contain the first key") {
					assert.Equal(t, data["redis"], ooi.tags["redis"], "it should have the right value for the first key")
				}

				if assert.Contains(t, ooi.tags, "pgsql", "it should contain the second key") {
					assert.Equal(t, data2["pgsql"], ooi.tags["pgsql"], "it should have the right value for the second key")
				}
			})

			t.Run("existing entries", func(t *testing.T) {
				data2 := map[string]string{
					"redis": "0.8.0",
				}
				o2 := Tags(data2)
				require.NotNil(t, o2, "the second module option should not be nil")

				oo := om.Merge(o2)
				require.NotNil(t, oo, "it should not return nil when it merges")

				ooi, ok := oo.(*tagsOption)
				require.True(t, ok, "it should actually be a *modulesOption")

				if assert.Contains(t, ooi.tags, "redis", "it should contain the first key") {
					assert.Equal(t, data["redis"], ooi.tags["redis"], "it should have the right value for the first key")
				}
			})
		})
	}

	t.Run("MarshalJSON()", func(t *testing.T) {
		assert.Equal(t, map[string]interface{}{
			"redis": "1.0.0",
		}, testOptionsSerialize(t, o))
	})
}
