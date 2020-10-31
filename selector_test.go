package lolhtml_test

import (
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestNewSelector_UnsupportedSelector(t *testing.T) {
	s, err := lolhtml.NewSelector("p:last-child")
	if s != nil || err == nil {
		t.FailNow()
	}
	if err.Error() != "Unsupported pseudo-class or pseudo-element in selector." {
		t.Error(err)
	}
}
