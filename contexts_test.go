package sentry

import (
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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

func TestContexts(t *testing.T) {
	Convey("Contexts", t, func() {
		Convey("RuntimeContext()", func() {
			Convey("Should return a context option", func() {
				So(RuntimeContext("go", runtime.Version()), ShouldHaveSameTypeAs, Context("runtime", nil))
			})

			Convey("Should set the context type to 'runtime'", func() {
				c := RuntimeContext("go", runtime.Version())
				cc, ok := c.(*contextOption)
				So(ok, ShouldBeTrue)
				So(cc.contexts, ShouldContainKey, "runtime")
			})

			Convey("Should set the context correctly", func() {
				c := RuntimeContext("go", runtime.Version())
				cc, ok := c.(*contextOption)
				So(ok, ShouldBeTrue)
				So(cc.contexts["runtime"], ShouldResemble, map[string]string{
					"name":    "go",
					"version": runtime.Version(),
				})
			})
		})

		Convey("OSContext()", func() {
			osInfo := OSContextInfo{
				Version:       "CentOS 7.3",
				Build:         "centos7.3.1611",
				KernelVersion: "3.10.0-514",
				Rooted:        false,
			}

			Convey("Should return a context option", func() {
				So(OSContext(&osInfo), ShouldHaveSameTypeAs, Context("os", nil))
			})

			Convey("Should set the context type to 'os'", func() {
				c := OSContext(&osInfo)
				cc, ok := c.(*contextOption)
				So(ok, ShouldBeTrue)
				So(cc.contexts, ShouldContainKey, "os")
			})

			Convey("Should set the context correctly", func() {
				c := OSContext(&osInfo)
				cc, ok := c.(*contextOption)
				So(ok, ShouldBeTrue)
				So(cc.contexts["os"], ShouldResemble, &osInfo)
			})
		})

		Convey("DeviceContext()", func() {
			deviceInfo := DeviceContextInfo{
				Architecture: "arm",
				BatteryLevel: 100,
				Family:       "Samsung Galaxy",
				Model:        "Samsung Galaxy S8",
				ModelID:      "SM-G95550",
				Name:         "Samsung Galaxy S8",
				Orientation:  "portrait",
			}

			Convey("Should return a context option", func() {
				So(DeviceContext(&deviceInfo), ShouldHaveSameTypeAs, Context("device", nil))
			})

			Convey("Should set the context type to 'device'", func() {
				c := DeviceContext(&deviceInfo)
				cc, ok := c.(*contextOption)
				So(ok, ShouldBeTrue)
				So(cc.contexts, ShouldContainKey, "device")
			})

			Convey("Should set the context correctly", func() {
				c := DeviceContext(&deviceInfo)
				cc, ok := c.(*contextOption)
				So(ok, ShouldBeTrue)
				So(cc.contexts["device"], ShouldResemble, &deviceInfo)
			})
		})

		Convey("Context()", func() {
			c := Context("test", "data")
			So(c, ShouldNotBeNil)
			So(c, ShouldHaveSameTypeAs, &contextOption{})

			cc := c.(*contextOption)
			So(cc.contexts, ShouldContainKey, "test")
			So(cc.contexts["test"], ShouldEqual, "data")
		})

		Convey("ContextOption", func() {
			c := Context("test", "data")
			So(c, ShouldNotBeNil)

			Convey("Should have the correct Class()", func() {
				So(c.Class(), ShouldEqual, "contexts")
			})

			Convey("Should implement MergableOption interface", func() {
				So(c, ShouldImplement, (*MergeableOption)(nil))
			})

			Convey("Merge()", func() {
				cc, ok := c.(*contextOption)
				So(ok, ShouldBeTrue)

				Convey("Should overwrite if it cannot identify the old type", func() {
					out := cc.Merge(&testOption{})
					So(out, ShouldEqual, c)
				})

				Convey("Should overwriting old fields", func() {
					old := Context("test", "oldData")
					out := cc.Merge(old)
					So(out, ShouldNotBeNil)
					So(out, ShouldHaveSameTypeAs, &contextOption{})

					So(out.(*contextOption).contexts, ShouldResemble, map[string]interface{}{
						"test": "data",
					})
				})

				Convey("Should add new fields", func() {
					old := Context("old", "data")
					out := cc.Merge(old)
					So(out, ShouldNotBeNil)
					So(out, ShouldHaveSameTypeAs, &contextOption{})

					So(out.(*contextOption).contexts, ShouldResemble, map[string]interface{}{
						"old":  "data",
						"test": "data",
					})
				})
			})

			Convey("MarshalJSON", func() {
				c := Context("test", "data")
				So(testOptionsSerialize(c), ShouldResemble, map[string]interface{}{
					"test": "data",
				})
			})
		})
	})
}
