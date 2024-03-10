# go-temp-err-catcher

This package is a provides a way to detect temporary errors, automatically applying an exponential
backoff. This is useful, for example, in `Accept` loops where the underlying listener may have a
temporary failure (e.g., because you have too many open file descriptors).

The main function is `TempErrCatcher.IsTemporary`. It works by:

1. Detecting if the error is temporary (implements a `Temporary() bool` method that returns true). E.g., [`net.Error`](https://pkg.go.dev/net#Error) (ignore the warning, this library exists to handle those "few exceptions" and avoids the pitfalls with an exponential backoff).
2. Sleeping an exponentially increasing delay (the delay resets after a we fail to detect any temporary errors for a while).

Docs: https://godoc.org/github.com/jbenet/go-temp-err-catcher

## Examples

This library is meant to be used with things like `net.Lister.Accept`:

```go
import (
  tec "github.com/jbenet/go-temp-err-catcher"
)

func listen(listener net.Listener) {
  // The TempErrCatcher tracks the current exponential backoff and detects "strings"
  // of temporary errors. Use one per potential error source.
  var c tec.TempErrCatcher

  for {
    conn, err := listener.Accept()
    if err != nil && c.IsTemporary(c) {
      // IsTemporary applies the backoff internally, so there's no need to sleep here.
      continue
    }
    return conn, err
  }
}
```

You can make your errors implement `Temporary`:

```go
type errTemp struct {
  e error
}

func (e errTemp) Temporary() bool {
  return true
}

func (e errTemp) Error() string {
  return e.e.Error()
}

err := errors.New("beep boop")
var c tec.TempErrCatcher
c.IsTemporary(err)              // false
c.IsTemporary(errTemp{err}) // true
```

Or just use `ErrTemp`:

```go
err := errors.New("beep boop")
var c tec.TempErrCatcher
c.IsTemporary(err)              // false
c.IsTemporary(tec.ErrTemp{err}) // true
```


You can also define an `IsTemp` function to classify errors:

```go
var ErrSkip = errors.New("this should be skipped")
var ErrNotSkip = errors.New("this should not be skipped")

var c tec.TempErrCatcher
c.IsTemp = func(e error) bool {
  return e == ErrSkip
}

c.IsTemporary(ErrSkip) // true
c.IsTemporary(ErrNotSkip) // false
c.IsTemporary(ErrTemp) // false! no longer accepts Temporary()
```
