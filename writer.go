package lolhtml

import (
	"bytes"
	"io"
)

// Writer takes data written to it and writes the rewritten form of that data to an
// underlying writer (see NewWriter).
type Writer struct {
	w        io.Writer
	rewriter *rewriter
	err      error
	closed   bool
}

// NewWriter returns a new Writer with Handlers and an optional Config configured.
// Writes to the returned Writer are rewritten and written to w.
//
// It is the caller's responsibility to call Close on the Writer when done.
// Writes may be buffered and not flushed until Close. There is no Flush method,
// so before using the content written by w, it is necessary to call Close
// to ensure w has finished writing.
func NewWriter(w io.Writer, handlers *Handlers, config ...Config) (*Writer, error) {
	var c Config
	var sink OutputSink
	if config != nil {
		c = config[0]
		if c.Sink != nil {
			sink = c.Sink
		} else if w == nil {
			sink = func([]byte) {}
		} else {
			sink = func(p []byte) {
				_, _ = w.Write(p)
			}
		}
	} else {
		c = newDefaultConfig()
		if w == nil {
			sink = func([]byte) {}
		} else {
			sink = func(p []byte) {
				_, _ = w.Write(p)
			}
		}
	}

	rb := newRewriterBuilder()
	var selectors []*selector
	if handlers != nil {
		for _, dh := range handlers.DocumentContentHandler {
			rb.AddDocumentContentHandlers(
				dh.DoctypeHandler,
				dh.CommentHandler,
				dh.TextChunkHandler,
				dh.DocumentEndHandler,
			)
		}
		for _, eh := range handlers.ElementContentHandler {
			s, err := newSelector(eh.Selector)
			if err != nil {
				return nil, err
			}
			selectors = append(selectors, s)
			rb.AddElementContentHandlers(
				s,
				eh.ElementHandler,
				eh.CommentHandler,
				eh.TextChunkHandler,
			)
		}
	}
	r, err := rb.Build(sink, c)
	if err != nil {
		return nil, err
	}
	rb.Free()
	for _, s := range selectors {
		s.Free()
	}

	return &Writer{w: w, rewriter: r}, nil
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}
	if len(p) == 0 {
		return 0, nil
	}
	n, err = w.rewriter.Write(p)
	if err != nil {
		w.err = err
		return
	}
	return
}

// WriteString writes a string to the Writer.
func (w *Writer) WriteString(s string) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}
	if len(s) == 0 {
		return 0, nil
	}
	n, err = w.rewriter.WriteString(s)
	if err != nil {
		w.err = err
		return
	}
	return
}

// Close closes the Writer, flushing any unwritten data to the underlying io.Writer,
// but does not close the underlying io.Writer.
// Subsequent calls to Close is a no-op.
func (w *Writer) Close() error {
	if w == nil || w.closed {
		return nil
	}
	w.closed = true
	if w.err == nil {
		w.err = w.rewriter.End()
	}
	w.rewriter.Free()
	return w.err
}

// RewriteString rewrites the given string with the provided Handlers and Config.
func RewriteString(s string, handlers *Handlers, config ...Config) (string, error) {
	var buf bytes.Buffer
	var w *Writer
	var err error
	if config != nil {
		w, err = NewWriter(&buf, handlers, config[0])
	} else {
		w, err = NewWriter(&buf, handlers)
	}
	if err != nil {
		return "", err
	}

	_, err = w.WriteString(s)
	if err != nil {
		return "", err
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
