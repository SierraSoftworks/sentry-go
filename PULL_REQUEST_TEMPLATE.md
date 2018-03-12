### Introduction
*Give us a short description of what your pull request does*
> This PR adds support for the new `teapot` interface in Sentry
> which allows [RFC 2324][] compliant devices to submit information
> about themselves to Sentry.

### Contains
*Tell us what your pull request includes*

- [ ] :radioactive: Breaking API changes
- [ ] :star: Features
  - [ ] :star2: New features
  - [ ] :beetle: Fixes for existing features
  - [ ] :100: Automated tests for my changes
- [ ] :books: Documentation
  - [ ] :memo: New documentation
  - [ ] :bookmark_tabs: Fixes for existing documentation
- [ ] :electric_plug: Plugins
  - [ ] :star2: New plugins
  - [ ] :beetle: Fixes for existing plugins
  - [ ] :100: Automated tests for my changes

### Description
*Give us a more in-depth description of what this pull request hopes to solve*
> In version 13.3.7 of Sentry support was added for [RFC 2324][] devices in the
> form of the new `teapot` interface. This interface allows compatible devices
> to submit information about their brew to Sentry as additional context for
> their events.
>
> This PR adds a new `Teapot()` option provider and a builder which allows
> information about a teapot to be included in a standard Sentry packet.
> It also includes extensive automated tests to ensure that both the option
> provider and builder work correctly, as well as validating that the serialized
> payload matches that expected by the Sentry API.
>
> ```go
> cl := sentry.NewClient()
> cl.Capture(
>   sentry.Teapot().WithBrew("Rooibos").WithTemperature(97.5),
> )
> ```

[RFC 2324]: http://tools.ietf.org/html/2324#section-6.5.14