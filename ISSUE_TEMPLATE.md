### Problem Statement
*Tell us in as few words as possible what problem you are running into*
> I am not seeing some of the fields I provided to `sentry.Extra()` in my Sentry UI

### Expected Behaviour
*Tell us how you expected things to behave*
> I expected that all of the fields I provided in `sentry.Extra()` would be sent to
> Sentry and be shown in the user interface there.

### Environment
*Tell us about the environment you are using*

- [ ] **Go Version**: `go version go1.9.2 windows/amd64` (`go version`)
- [ ] **Sentry Version**: `8.22.0`
- [ ] **Updated `sentry-go.v1`** (`go get -u gopkg.in/SierraSoftworks/sentry-go.v1`)

### Reproduction Code
*Give us a short snippet of code that shows where you are encountering the problem*

```go
package main

import (
  "fmt"
  "os"

  "gopkg.in/SierraSoftworks/sentry-go.v1"
)

func main() {
  cl := sentry.NewClient(
    sentry.Release(fmt.Sprintf("#%s", os.Getenv("ISSUE_ID"))),
  )
  
  cl.Capture(
    sentry.Extra(map[string]interface{}{
      "field1": "this is visible",
      // ...
      "field51": "this is not visible",
    })
  )
}
```
