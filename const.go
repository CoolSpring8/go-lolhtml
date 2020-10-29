package lolhtml

import "C"

// RewriterDirective should returned by callback handlers, to inform the rewriter to continue or stop parsing.
type RewriterDirective int

const (
	// Let the normal parsing process continue.
	Continue RewriterDirective = iota

	// Stop the rewriter immediately. Content currently buffered is discarded, and an error is returned.
	Stop
)
