package sentry

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"

	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleAddDefaultOptions() {
	// You can add default options to all of your root Sentry
	// clients like this.
	AddDefaultOptions(
		Release("v1.0.0"),
		DSN("..."),
	)
}

func ExampleAddDefaultOptionProvider() {
	// You can also register options providers which will dynamically
	// generate options for each new event that is sent
	AddDefaultOptionProvider(func() Option {
		if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
			return DSN(dsn)
		}

		return nil
	})
}

func TestOptions(t *testing.T) {
	oldOptionsProviders := defaultOptionProviders
	defer func() {
		defaultOptionProviders = oldOptionsProviders
	}()

	id, err := NewEventID()
	assert.Nil(t, err, "there should be no errors creating the ID")

	t.Run("AddDefaultOptionProvider()", func(t *testing.T) {
		defaultOptionProviders = []func() Option{}

		provider := func() Option {
			return EventID(id)
		}

		AddDefaultOptionProvider(provider)
		assert.Len(t, defaultOptionProviders, 1, "the provider should now be present in the default options providers list")

		for _, provider := range defaultOptionProviders {
			assert.Equal(t, EventID(id), provider(), "the provider should return the right option")
		}
	})

	t.Run("AddDefaultOptions()", func(t *testing.T) {
		defaultOptionProviders = []func() Option{}

		AddDefaultOptions(EventID(id), nil, EventID(id))
		assert.Len(t, defaultOptionProviders, 2, "the provider should now be present in the default options providers list")

		for _, provider := range defaultOptionProviders {
			assert.Equal(t, EventID(id), provider(), "the provider should return the right option")
		}
	})
}

type testOption struct {
}

func (o *testOption) Class() string {
	return "test"
}

type testCustomClassOption struct {
	class string
}

func (o *testCustomClassOption) Class() string {
	return o.class
}

type testOmitableOption struct {
	omit bool
}

func (o *testOmitableOption) Class() string {
	return "test"
}

func (o *testOmitableOption) Omit() bool {
	return o.omit
}

type testFinalizeableOption struct {
	finalized bool
}

func (o *testFinalizeableOption) Class() string {
	return "test"
}

func (o *testFinalizeableOption) Finalize() {
	o.finalized = true
}

type testMergeableOption struct {
	data int
}

func (o *testMergeableOption) Class() string {
	return "test"
}

func (o *testMergeableOption) Merge(other Option) Option {
	if oo, ok := other.(*testMergeableOption); ok {
		return &testMergeableOption{
			data: o.data + oo.data,
		}
	}

	return o
}

type testAdvancedOption struct {
	data map[string]Option
}

func (o *testAdvancedOption) Class() string {
	return "test"
}

func (o *testAdvancedOption) Apply(packet map[string]Option) {
	for k, v := range o.data {
		packet[k] = v
	}
}

type testSerializableOption struct {
	data string
}

func (o *testSerializableOption) Class() string {
	return "test"
}

func (o *testSerializableOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.data)
}

func testGetOptionsProvider(t *testing.T, sameType Option) Option {
	st := reflect.TypeOf(sameType)
	assert.NotNil(t, st, "getting the reflection type should not fail")

	for _, provider := range defaultOptionProviders {
		opt := provider()
		if reflect.TypeOf(opt) == st {
			return opt
		}
	}

	return nil
}

func testOptionsSerialize(t *testing.T, opt Option) interface{} {
	if opt == nil {
		return nil
	}

	var data interface{}
	buf := bytes.NewBuffer([]byte{})
	assert.Nil(t, json.NewEncoder(buf).Encode(opt), "no error should occur when serializing to JSON")
	assert.Nil(t, json.NewDecoder(buf).Decode(&data), "no error should occur when deserializing from JSON")
	return data
}
