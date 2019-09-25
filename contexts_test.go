package sentry

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleRuntimeContext() {
	cl := NewClient(
		// You can configure this when creating your client
		RuntimeContext("go", runtime.Version()),
	)

	cl.Capture(
		// Or when sending an event
		RuntimeContext("go", runtime.Version()),
	)
}

func ExampleOSContext() {
	osInfo := OSContextInfo{
		Version:       "CentOS 7.3",
		Build:         "centos7.3.1611",
		KernelVersion: "3.10.0-514",
		Rooted:        false,
	}

	cl := NewClient(
		// You can provide this when creating your client
		OSContext(&osInfo),
	)

	cl.Capture(
		// Or when you send an event
		OSContext(&osInfo),
	)
}

func ExampleDeviceContext() {
	deviceInfo := DeviceContextInfo{
		Architecture: "arm",
		BatteryLevel: 100,
		Family:       "Samsung Galaxy",
		Model:        "Samsung Galaxy S8",
		ModelID:      "SM-G95550",
		Name:         "Samsung Galaxy S8",
		Orientation:  "portrait",
	}

	cl := NewClient(
		// You can provide this when creating your client
		DeviceContext(&deviceInfo),
	)

	cl.Capture(
		// Or when you send an event
		DeviceContext(&deviceInfo),
	)
}

func TestRuntimeContext(t *testing.T) {
	c := RuntimeContext("go", runtime.Version())

	assert.NotNil(t, c, "it should not return a nil option")
	assert.IsType(t, Context("runtime", nil), c, "it should return the same thing as a Context()")
	
	cc, ok := c.(*contextOption)
	assert.True(t, ok, "it should actually return a *contextOption")
	if assert.Contains(t, cc.contexts, "runtime", "it should specify a runtime context") {
		assert.Equal(t, map[string]string{
			"name": "go",
			"version": runtime.Version(),
		}, cc.contexts["runtime"], "it should specify the correct context values")
	}
}

func TestOSContext(t *testing.T) {
	osInfo := OSContextInfo{
		Version:       "CentOS 7.3",
		Build:         "centos7.3.1611",
		KernelVersion: "3.10.0-514",
		Rooted:        false,
	}

	c := OSContext(&osInfo)

	assert.NotNil(t, c, "it should not return a nil option")
	assert.IsType(t, Context("os", nil), c, "it should return the same thing as a Context()")
	
	cc, ok := c.(*contextOption)
	assert.True(t, ok, "it should actually return a *contextOption")
	if assert.Contains(t, cc.contexts, "os", "it should specify an os context") {
		assert.Equal(t, &osInfo, cc.contexts["os"], "it should specify the correct context values")
	}
}

func TestDeviceContext(t *testing.T) {
	deviceInfo := DeviceContextInfo{
		Architecture: "arm",
		BatteryLevel: 100,
		Family:       "Samsung Galaxy",
		Model:        "Samsung Galaxy S8",
		ModelID:      "SM-G95550",
		Name:         "Samsung Galaxy S8",
		Orientation:  "portrait",
	}

	c := DeviceContext(&deviceInfo)

	assert.NotNil(t, c, "it should not return a nil option")
	assert.IsType(t, Context("device", nil), c, "it should return the same thing as a Context()")
	
	cc, ok := c.(*contextOption)
	assert.True(t, ok, "it should actually return a *contextOption")
	if assert.Contains(t, cc.contexts, "device", "it should specify an os context") {
		assert.Equal(t, &deviceInfo, cc.contexts["device"], "it should specify the correct context values")
	}
}

func TestContext(t *testing.T) {
	c := Context("test", "data")
	assert.NotNil(t, c, "it should not return a nil option")
	assert.IsType(t, &contextOption{}, c, "it should actually return a *contextOption")

	cc := c.(*contextOption)
	assert.Contains(t, cc.contexts, "test", "it should set the 'test' context")
	assert.Equal(t, "data", cc.contexts["test"], "it should set the context data correctly")
}

func TestContextOption(t *testing.T) {
	c := Context("test", "data")
	assert.NotNil(t, c, "it should not return a nil option")
	
	assert.IsType(t, &contextOption{}, c, "it should actually return a *contextOption")
	cc := c.(*contextOption)

	assert.Equal(t, "contexts", c.Class(), "it should use the right option class")
	assert.Implements(t, (*MergeableOption)(nil), c, "it should implement the MergeableOption interface")

	t.Run("Merge()", func(t *testing.T) {
		t.Run("Unknown Type", func(t *testing.T) {
			out := cc.Merge(&testOption{})
			assert.Equal(t, c, out, "it should overwrite the original value")
		})
		
		t.Run("Existing Context", func(t *testing.T) {
			old := Context("test", "oldData")
			out := cc.Merge(old)

			assert.NotNil(t, out, "it should not return a nil result")
			assert.IsType(t, &contextOption{}, out, "it should return a new *contextOption")
			
			oo := out.(*contextOption)
			assert.Equal(t, map[string]interface{}{
				"test": "data",
			}, oo.contexts)
		})
		
		t.Run("Existing Different Context", func(t *testing.T) {
			old := Context("old", "oldData")
			out := cc.Merge(old)

			assert.NotNil(t, out, "it should not return a nil result")
			assert.IsType(t, &contextOption{}, out, "it should return a new *contextOption")
			
			oo := out.(*contextOption)
			assert.Equal(t, map[string]interface{}{
				"test": "data",
				"old": "oldData",
			}, oo.contexts)
		})
	})

	t.Run("MarshalJSON()", func(t *testing.T) {
		c := Context("test", "data")
		assert.Equal(t, map[string]interface{}{ "test": "data" }, testOptionsSerialize(t, c))
	})
}
