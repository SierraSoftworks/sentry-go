package sentry

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"

	"testing"

	"github.com/smartystreets/goconvey/convey"
	. "github.com/smartystreets/goconvey/convey"
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
	Convey("Options", t, func() {
		oldOptionsProviders := defaultOptionProviders
		defaultOptionProviders = []func() Option{}
		defer func() {
			defaultOptionProviders = oldOptionsProviders
		}()

		Convey("AddDefaultOptionProvider", func() {
			So(defaultOptionProviders, ShouldBeEmpty)

			id, err := NewEventID()
			So(err, ShouldBeNil)
			provider := func() Option {
				return EventID(id)
			}

			AddDefaultOptionProvider(provider)
			So(defaultOptionProviders, ShouldHaveLength, 1)

			for _, provider := range defaultOptionProviders {
				So(provider(), ShouldResemble, EventID(id))
			}
		})

		Convey("AddDefaultOptions", func() {
			So(defaultOptionProviders, ShouldBeEmpty)

			id, err := NewEventID()
			So(err, ShouldBeNil)

			AddDefaultOptions(EventID(id), nil, EventID(id))
			So(defaultOptionProviders, ShouldHaveLength, 2)

			for _, provider := range defaultOptionProviders {
				So(provider(), ShouldResemble, EventID(id))
			}
		})
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

type testFinalizableOption struct {
	finalized bool
}

func (o *testFinalizableOption) Class() string {
	return "test"
}

func (o *testFinalizableOption) Finalize() {
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

func testGetOptionsProvider(sameType Option) Option {
	st := reflect.TypeOf(sameType)
	convey.So(st, convey.ShouldNotBeNil)

	for _, provider := range defaultOptionProviders {
		opt := provider()
		if reflect.TypeOf(opt) == st {
			return opt
		}
	}

	return nil
}

func testOptionsSerialize(opt Option) interface{} {
	if opt == nil {
		return nil
	}

	var data interface{}
	buf := bytes.NewBuffer([]byte{})
	convey.So(json.NewEncoder(buf).Encode(opt), convey.ShouldBeNil)
	convey.So(json.NewDecoder(buf).Decode(&data), convey.ShouldBeNil)
	return data
}
