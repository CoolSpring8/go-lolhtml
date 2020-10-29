package lolhtml_test

import (
	"errors"
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

// TestNullErrorStr tests internal functions for handling a null lol_html_str_t, by calling lol_html_take_last_error()
// when there is no error.
func TestNullErrorStr(t *testing.T) {
	err := lolhtml.GetError()
	if !errors.Is(err, lolhtml.ErrCannotGetErrorMessage) {
		t.Error(err)
	}
}
