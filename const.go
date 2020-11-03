package lolhtml

import "C"

// RewriterDirective is a "status code“ that should be returned by callback handlers, to inform the
// rewriter to continue or stop parsing.
type RewriterDirective int

const (
	// Continue lets the normal parsing process continue.
	Continue RewriterDirective = iota

	// Stop stops the rewriter immediately. Content currently buffered is discarded, and an error is returned.
	// After stopping, the Writer should not be used anymore except for Close().
	Stop
)
