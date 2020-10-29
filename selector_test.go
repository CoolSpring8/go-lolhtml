package lolhtml_test

import (
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestNewSelector(t *testing.T) {
	testCases := []struct {
		selector  string
		errorText string
	}{
		{"p.center", ""},
		{"p:last-child", "Unsupported pseudo-class or pseudo-element in selector."},
	}
	for _, tc := range testCases {
		t.Run(tc.selector, func(t *testing.T) {
			s, err := lolhtml.NewSelector(tc.selector)
			if err == nil {
				if tc.errorText != "" {
					t.FailNow()
				}
			} else {
				if err.Error() != tc.errorText {
					t.Error(err)
				}
				s.Free()
			}
		})
	}
}
