package sentry

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExamplePacket() {
	// Create a new packet object which can be sent to
	// Sentry by one of the transports or send queues.
	p := NewPacket().SetOptions(
		DSN(""),
		Message("Custom packet creation"),
	)

	// Create a clone of this packet if you want to use
	// it as a template
	p.Clone().SetOptions(
		Message("Overridden message which doesn't affect the original"),
	)
}

func TestPacket(t *testing.T) {
	p := NewPacket()
	assert.NotNil(t, p, "should return a non-nil packet")
	assert.Implements(t, (*Packet)(nil), p, "it should implement the Packet interface")

	t.Run("SetOptions()", func(t *testing.T) {
		assert.Equal(t, p, p.SetOptions(), "it should return the packet to support chaining")

		assert.Equal(t, p, p.SetOptions(nil), "it should ignore nil options")

		opt := &testOption{}
		assert.Equal(t, p.SetOptions(opt), p.Clone().SetOptions(nil, opt), "it should ignore nil options when other options are provided")

		pp, ok := p.(*packet)
		assert.True(t, ok, "it should actually be a *packet")

		p.SetOptions(&testOption{})
		assert.Contains(t, *pp, "test", "it should contain the option field")
		assert.Equal(t, &testOption{}, (*pp)["test"], "it should have the right value for the option field")

		t.Run("Option Replacement", func(t *testing.T) {
			opt1 := &testOption{}
			opt2 := &testOption{}

			p.SetOptions(opt1)
			assert.Same(t, opt1, (*pp)["test"], "the first option should be set in the packet")

			p.SetOptions(opt2)
			assert.Same(t, opt2, (*pp)["test"], "the first option should be replaced by the second")
		})

		t.Run("Omit()", func(t *testing.T) {
			p.SetOptions(&testOmitableOption{
				omit: true,
			})
			assert.NotEqual(t, &testOmitableOption{omit: true}, (*pp)["test"], "it should omit changes if Omit() returns true")

			p.SetOptions(&testOmitableOption{
				omit: false,
			})
			assert.Equal(t, &testOmitableOption{omit: false}, (*pp)["test"], "it should not omit changes if Omit() returns false")
		})

		t.Run("Finalize()", func(t *testing.T) {
			opt := &testFinalizeableOption{}
			assert.False(t, opt.finalized, "the option should initially not be finalized")

			p.SetOptions(opt)
			assert.True(t, opt.finalized, "the option should now be finalized")
			assert.Equal(t, opt, (*pp)["test"], "the option should be stored in the packet")
		})

		t.Run("Merge()", func(t *testing.T) {
			opt1 := &testMergeableOption{data: 1}
			opt2 := &testMergeableOption{data: 2}

			p.SetOptions(opt1)
			assert.Same(t, opt1, (*pp)["test"], "the packet should initially contain the first option")

			p.SetOptions(opt2)
			assert.Equal(t, &testMergeableOption{data: 3}, (*pp)["test"], "the packet should then contain the merged option")
			assert.Equal(t, 1, opt1.data, "the first option's data shouldn't be modified")
			assert.Equal(t, 2, opt2.data, "the second option's data shouldn't be modified")
		})

		t.Run("Apply()", func(t *testing.T) {
			opt := &testAdvancedOption{
				data: map[string]Option{
					"tested": Context("value", true),
				},
			}

			p.SetOptions(opt)
			assert.Contains(t, (*pp), "tested", "it should have run the Apply() method")
			assert.Equal(t, Context("value", true), (*pp)["tested"], "it should have stored the correct value")
		})
	})

	t.Run("Clone()", func(t *testing.T) {
		assert.False(t, p == p.Clone(), "it should clone to a new packet")
		assert.Equal(t, p, p.Clone(), "it should clone to an equivalent packet")

		p := NewPacket().SetOptions(DSN(""), Message("Test"))
		assert.Equal(t, p, p.Clone(), "the clone should copy any options across")
	})

	t.Run("MarshalJSON()", func(t *testing.T) {
		p := NewPacket()
		p.SetOptions(&testOption{})

		assert.Equal(t, map[string]interface{}{
			"test": map[string]interface{}{},
		}, testSerializePacket(t, p))

		p.SetOptions(&testSerializableOption{data: "testing"})
		assert.Equal(t, map[string]interface{}{
			"test": "testing",
		}, testSerializePacket(t, p))
	})
}

func testSerializePacket(t *testing.T, p Packet) interface{} {
	buf := bytes.NewBuffer([]byte{})
	assert.Nil(t, json.NewEncoder(buf).Encode(p), "it should not encounter any errors serializing the packet")

	var data interface{}
	assert.Nil(t, json.NewDecoder(buf).Decode(&data), "it should not encounter any errors deserializing the packet")

	return data
}
