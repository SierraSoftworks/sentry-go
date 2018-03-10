package sentry

func ExampleBreadcrumbs() {
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

func ExampleBreadcrumbsContext() {
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

	//Output:

}
