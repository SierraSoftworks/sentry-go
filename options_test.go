package sentry

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/smartystreets/goconvey/convey"
)

type testOption struct {
}

func (o *testOption) Class() string {
	return "test"
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
