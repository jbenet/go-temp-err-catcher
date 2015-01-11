package temperrcatcher

import (
	"fmt"
	"testing"
	"time"
)

var (
	ErrTemp  = ErrTemporary{fmt.Errorf("ErrTemp")}
	ErrSkip  = fmt.Errorf("ErrSkip")
	ErrOther = fmt.Errorf("ErrOther")
)

func testTec(t *testing.T, c TempErrCatcher, errs map[error]bool) {
	for e, expected := range errs {
		if c.IsTemporary(e) != expected {
			t.Error("expected %s to be %v", e, expected)
		}
	}
}

func TestNil(t *testing.T) {
	var c TempErrCatcher
	testTec(t, c, map[error]bool{
		ErrTemp:  true,
		ErrSkip:  false,
		ErrOther: false,
	})
}

func TestWait(t *testing.T) {
	var c TempErrCatcher
	worked := make(chan time.Duration, 3)
	c.Wait = func(t time.Duration) {
		worked <- t
	}
	testTec(t, c, map[error]bool{
		ErrTemp:  true,
		ErrSkip:  false,
		ErrOther: false,
	})

	// should've called it once
	select {
	case <-worked:
	default:
		t.Error("did not call our Wait func")
	}

	// should've called it ONLY once
	select {
	case <-worked:
		t.Error("called our Wait func more than once")
	default:
	}
}

func TestTemporary(t *testing.T) {
	var c TempErrCatcher
	testTec(t, c, map[error]bool{
		ErrTemp:  true,
		ErrSkip:  false,
		ErrOther: false,
	})
}

func TestFunc(t *testing.T) {
	var c TempErrCatcher
	c.IsTemp = func(e error) bool {
		return e == ErrSkip
	}
	testTec(t, c, map[error]bool{
		ErrTemp:  false,
		ErrSkip:  true,
		ErrOther: false,
	})
}
