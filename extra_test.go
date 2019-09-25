package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleExtra() {
	cl := NewClient(
		// You can define extra fields when you create your client
		Extra(map[string]interface{}{
			"redis": map[string]interface{}{
				"host": "redis",
				"port": 6379,
			},
		}),
	)

	cl.Capture(
		// You can also define extra info when you send the event
		// The extra object will be shallowly merged automatically,
		// so this would send both `redis` and `cache`.
		Extra(map[string]interface{}{
			"cache": map[string]interface{}{
				"key": "user.127.profile",
				"hit": false,
			},
		}),
	)
}

func TestExtra(t *testing.T) {
	data := map[string]interface{}{
		"redis": map[string]interface{}{
			"host": "redis",
			"port": 6379,
		},
	}

	assert.Nil(t, Extra(nil), "it should return nil if the data is nil")

	o := Extra(data)
	assert.NotNil(t, o, "it should return a non-nil option when the data is not nil")
	assert.Implements(t, (*Option)(nil), o, "it should implement the Option interface")
	assert.Equal(t, "extra", o.Class(), "it should use the right option class")
	
	if assert.Implements(t, (*MergeableOption)(nil), o, "it should implement the MergeableOption interface") {
		om := o.(MergeableOption)

		t.Run("Merge()", func(t *testing.T) {
			data2 := map[string]interface{}{
				"cache": map[string]interface{}{
					"key": "user.127.profile",
					"hit": false,
				},
			}

			o2 := Extra(data2)

			assert.Equal(t, o, om.Merge(&testOption{}), "it should replace the old option if it is not recognized")

			oo := om.Merge(o2)
			assert.NotNil(t, oo, "it should return a non-nil merge result")
			assert.NotEqual(t, o, oo, "it should not re-purpose the first option")
			assert.NotEqual(t, o2, oo, "it should not re-purpose the second option")

			eo, ok := oo.(*extraOption)
			assert.True(t, ok, "it should actually be an *extraOption")

			assert.Contains(t, eo.extra, "redis", "it should contain the key from the first option")
			assert.Contains(t, eo.extra, "cache", "it should contain the key from the second option")

			data2 = map[string]interface{}{
				"redis": map[string]interface{}{
					"host": "redis-dev",
					"port": 6379,
				},
			}

			o2 = Extra(data2)
			assert.NotNil(t, o2, "it should not be nil")
			oo = om.Merge(o2)
			assert.NotNil(t, oo, "it should return a non-nil merge result")

			eo, ok = oo.(*extraOption)
			assert.True(t, ok, "it should actually be an *extraOption")

			assert.Contains(t, eo.extra, "redis", "it should contain the key")
			assert.Equal(t, data, eo.extra, "it should use the new option's data")
		})
	}

	t.Run("MarshalJSON()", func(t *testing.T) {
		data := map[string]interface{}{
			"redis": map[string]interface{}{
				"host": "redis",
				// Float mode required since we aren't deserializing into an int
				"port": 6379.,
			},
		}

		serialized := testOptionsSerialize(t, Extra(data))
		assert.Equal(t, data, serialized, "it should serialize to the data")
	})
}
