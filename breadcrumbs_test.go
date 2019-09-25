package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleDefaultBreadcrumbs() {
	// We can change the maximum number of breadcrumbs to be stored
	DefaultBreadcrumbs().WithSize(5)

	DefaultBreadcrumbs().NewDefault(nil).WithMessage("This is an example")
	DefaultBreadcrumbs().NewDefault(map[string]interface{}{
		"example": true,
	}).WithMessage("It should give you an idea of how you can use breadcrumbs in your app")

	DefaultBreadcrumbs().
		NewNavigation("introduction", "navigation").
		WithCategory("navigation").
		WithMessage("You can use them to represent navigations from one page to another")

	DefaultBreadcrumbs().
		NewNavigation("navigation", "http").
		WithCategory("navigation").
		WithMessage("Or to represent changes in the state of your application's workflows")

	DefaultBreadcrumbs().
		NewHTTPRequest("GET", "https://example.com/api/v1/awesome", 200, "OK").
		WithLevel(Debug).
		WithMessage("I think we can agree that they're pretty awesome")

	NewClient().Capture(Message("Finally, we send the event with all our breadcrumbs included"))
}

func ExampleBreadcrumbs() {
	rootClient := NewClient()
	DefaultBreadcrumbs().NewDefault(nil).WithMessage("Breadcrumb in the default context")

	breadcrumbs := NewBreadcrumbsList(10)
	contextClient := rootClient.With(Breadcrumbs(breadcrumbs))
	breadcrumbs.NewDefault(nil).WithMessage("Breadcrumb in the private context")

	// Will include only the first breadcrumb
	rootClient.Capture(
		Message("Event in default context"),
		Logger("default"),
	)

	// Will include only the second breadcrumb
	contextClient.Capture(
		Message("Event in private context"),
		Logger("private"),
	)
}

func TestBreadcrumbs(t *testing.T) {
	t.Run("Options Providers", func(t *testing.T) {
		assert.NotNil(t, testGetOptionsProvider(t, &breadcrumbsList{}), "Breadcrumbs should be registered as a default options provider")
	})

	t.Run("DefaultBreadcrumbs()", func(t *testing.T) {
		assert.NotNil(t, DefaultBreadcrumbs())
		assert.Implements(t, (*BreadcrumbsList)(nil), DefaultBreadcrumbs())
	})

	t.Run("Breadcrumbs()", func(t *testing.T) {
		assert.Nil(t, Breadcrumbs(nil), "it should return a nil option if it receives a nil breadcrumbs list")

		l := NewBreadcrumbsList(3)
		assert.NotNil(t, l, "it should create a breadcrumbs list")

		b := Breadcrumbs(l)
		assert.NotNil(t, b, "it should return an option when the list is not nil")

		assert.Equal(t, "breadcrumbs", b.Class(), "it should use the correct option class")
	})
}

func TestNewBreadcrumbsList(t *testing.T) {
	l := NewBreadcrumbsList(3)
	assert.NotNil(t, l, "it should return a non-nil list")

	assert.Implements(t, (*Option)(nil), l, "it should implement the option interface")

	ll, ok := l.(*breadcrumbsList)
	assert.True(t, ok, "it should actually be a *breadcrumbsList")

	assert.Equal(t, 3, ll.MaxLength, "it should have the right max length")
	assert.Equal(t, 0, ll.Length, "it should start with no breadcrumbs")
	assert.Nil(t, ll.Head, "it should have no head to start with")
	assert.Nil(t, ll.Tail, "it should have no tail to start with")

	t.Run("NewDefault(nil)", func(t *testing.T) {
		b := l.NewDefault(nil)
		assert.NotNil(t, b, "it should return a non-nil breadcrumb")
		assert.Implements(t, (*Breadcrumb)(nil), b, "the breadcrumb should implement the Breadcrumb interface")

		assert.NotNil(t, ll.Tail, "the list's tail should no longer be nil")
		assert.Equal(t, ll.Tail.Value, b, "the list's tail should now be the new breadcrumb")

		bb, ok := b.(*breadcrumb)
		assert.True(t, ok, "it should actually be a *breadcrumb object")
		assert.Equal(t, "", bb.Type, "it should use the default breadcrumb type")
		assert.Equal(t, map[string]interface{}{}, bb.Data, "it should use the passed breadcrumb data")
	})

	t.Run("NewDefault(data)", func(t *testing.T) {
		data := map[string]interface{}{
			"test": true,
		}

		b := l.NewDefault(data)
		assert.NotNil(t, b, "it should return a non-nil breadcrumb")
		assert.Implements(t, (*Breadcrumb)(nil), b, "the breadcrumb should implement the Breadcrumb interface")

		assert.NotNil(t, ll.Tail, "the list's tail should no longer be nil")
		assert.Equal(t, ll.Tail.Value, b, "the list's tail should now be the new breadcrumb")

		bb, ok := b.(*breadcrumb)
		assert.True(t, ok, "it should actually be a *breadcrumb object")
		assert.Equal(t, "", bb.Type, "it should use the default breadcrumb type")
		assert.Equal(t, data, bb.Data, "it should use the passed breadcrumb data")
	})

	t.Run("NewNavigation()", func(t *testing.T) {
		b := l.NewNavigation("/from", "/to")
		assert.NotNil(t, b, "it should return a non-nil breadcrumb")
		assert.Implements(t, (*Breadcrumb)(nil), b, "the breadcrumb should implement the Breadcrumb interface")

		assert.NotNil(t, ll.Tail, "the list's tail should no longer be nil")
		assert.Equal(t, ll.Tail.Value, b, "the list's tail should now be the new breadcrumb")

		bb, ok := b.(*breadcrumb)
		assert.True(t, ok, "it should actually be a *breadcrumb object")
		assert.Equal(t, "navigation", bb.Type, "it should use the default breadcrumb type")
		assert.Equal(t, map[string]interface{}{
			"from": "/from",
			"to":   "/to",
		}, bb.Data, "it should use the correct breadcrumb data")
	})

	t.Run("NewHTTPRequest()", func(t *testing.T) {
		b := l.NewHTTPRequest("GET", "/test", 200, "OK")
		assert.NotNil(t, b, "it should return a non-nil breadcrumb")
		assert.Implements(t, (*Breadcrumb)(nil), b, "the breadcrumb should implement the Breadcrumb interface")

		assert.NotNil(t, ll.Tail, "the list's tail should no longer be nil")
		assert.Equal(t, ll.Tail.Value, b, "the list's tail should now be the new breadcrumb")

		bb, ok := b.(*breadcrumb)
		assert.True(t, ok, "it should actually be a *breadcrumb object")
		assert.Equal(t, "http", bb.Type, "it should use the default breadcrumb type")
		assert.Equal(t, map[string]interface{}{
			"method":      "GET",
			"url":         "/test",
			"status_code": 200,
			"reason":      "OK",
		}, bb.Data, "it should use the correct breadcrumb data")
	})

	t.Run("WithSize()", func(t *testing.T) {
		cl := l.WithSize(5)
		assert.Equal(t, l, cl, "it should return the list so that the call is chainable")
		assert.Equal(t, 5, ll.MaxLength, "it should update the lists's max size")

		var b Breadcrumb
		for i := 0; i < ll.MaxLength*2; i++ {
			b = l.NewDefault(map[string]interface{}{
				"index": i,
			})
		}

		assert.Equal(t, ll.MaxLength, ll.Length, "the list should cap out at its max size")

		l.WithSize(1)
		assert.Equal(t, 1, ll.Length, "the list should be resized to the new max size")

		assert.Equal(t, b, ll.Head.Value, "the head of the list should be the last breadcrumb which was added")
		assert.Equal(t, b, ll.Tail.Value, "the tail of the list should be the last breadcrumb which was added")
	})

	t.Run("append()", func(t *testing.T) {
		l.WithSize(0).WithSize(3)

		var b Breadcrumb
		for i := 0; i < 10; i++ {
			b = l.NewDefault(map[string]interface{}{
				"index": i,
			})
		}

		assert.Equal(t, 3, ll.Length, "it should evict values to ensure that the length remains capped")
		assert.Equal(t, b, ll.Tail.Value, "it should add new breadcrumbs at the end of the list")
	})

	t.Run("list()", func(t *testing.T) {
		l.WithSize(0).WithSize(3)
		assert.Equal(t, 0, ll.Length, "should start with an empty breadcrumbs list")

		ol := ll.list()
		assert.NotNil(t, ol, "should not return nil if the list is empty")
		assert.Len(t, ol, 0, "it should return an empty list")

		for i := 0; i < 10; i++ {
			l.NewDefault(map[string]interface{}{"index": i})
		}
		assert.Equal(t, 3, ll.Length, "should now have three breadcrumbs in the list")

		ol = ll.list()
		assert.NotNil(t, ol, "should not return nil if the list is non-empty")
		assert.Len(t, ol, 3, "should return the maximum number of items if the list is full")

		for i, item := range ol {
			assert.IsType(t, &breadcrumb{}, item, "every list item should be a *breadcrumb")
			assert.Equal(t, i+7, item.(*breadcrumb).Data["index"], "the items should be in the right order")
		}
	})

	t.Run("MarshalJSON()", func(t *testing.T) {
		l.WithSize(0).WithSize(5).NewDefault(map[string]interface{}{"test": true})

		data := testOptionsSerialize(t, Breadcrumbs(l))
		assert.NotNil(t, data, "should not return a nil result")
		assert.IsType(t, []interface{}{}, data, "should return a JSON array")
	})
}
