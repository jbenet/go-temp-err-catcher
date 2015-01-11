// Package temperrcatcher provides a TempErrCatcher object,
// which implements simple error-retrying functionality.
package temperrcatcher

import (
	"time"
)

// Temporary is an interface errors can implement to
// ensure they are correctly classified by the default
// TempErrCatcher classifier
type Temporary interface {
	Temporary() bool
}

// ErrIsTemporary returns whether an error is Temporary(),
// iff it implements the Temporary interface.
func ErrIsTemporary(e error) bool {
	te, ok := e.(Temporary)
	return ok && te.Temporary()
}

// TempErrCatcher catches temporary errors for you. It then sleeps
// for a bit before returning (you should then try again). This may
// seem odd, but it's exactly what net/http does:
// http://golang.org/src/net/http/server.go?s=51504:51550#L1728
//
// You can set a few options in TempErrCatcher. They all have defaults
// so a zero TempErrCatcher is ready to be used:
//
//  var c tec.TempErrCatcher
//  c.IsTemporary(tempErr)
//
type TempErrCatcher struct {
	IsTemp func(error) bool    // the classifier to use. default: ErrIsTemporary
	Wait   func(time.Duration) // the wait func to call. default: time.Sleep
	Max    time.Duration       // the maximum time to wait. default: time.Second
	delay  time.Duration
}

func (tec *TempErrCatcher) init() {
	if tec.Max == 0 {
		tec.Max = time.Second
	}
	if tec.IsTemp == nil {
		tec.IsTemp = ErrIsTemporary
	}
	if tec.Wait == nil {
		tec.Wait = time.Sleep
	}
}

// IsTemporary wa
func (tec *TempErrCatcher) IsTemporary(e error) bool {
	tec.init()
	if tec.IsTemp(e) {
		if tec.delay == 0 {
			tec.delay = 5 * time.Millisecond
		} else {
			tec.delay *= 2
		}
		if tec.delay > tec.Max {
			tec.delay = tec.Max
		}
		tec.Wait(tec.delay)
		return true
	}
	return false
}

// ErrTemporary wraps any error and implements Temporary function.
//
//   err := errors.New("beep boop")
//   var c tec.TempErrCatcher
//   c.IsTemporary(err)              // false
//   c.IsTemporary(tec.ErrTemp{err}) // true
//
type ErrTemporary struct {
	Err error
}

func (e ErrTemporary) Temporary() bool {
	return true
}

func (e ErrTemporary) Error() string {
	return e.Err.Error()
}

func (e ErrTemporary) String() string {
	return e.Error()
}
