package lolhtml

import (
	"bytes"
	"io"
)

type Writer struct {
	w io.Writer
	r *rewriter
}

// NewWriter returns a new Writer with Handlers and Config configured, writing to w.
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

	return &Writer{w, r}, nil
}

func (w Writer) Write(p []byte) (n int, err error) {
	return w.r.Write(p)
}

func (w Writer) WriteString(s string) (n int, err error) {
	return w.r.WriteString(s)
}

func (w *Writer) Free() {
	if w != nil {
		w.r.Free()
	}
}

func (w *Writer) End() error {
	return w.r.End()
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
	defer w.Free()

	_, err = w.WriteString(s)
	if err != nil {
		return "", err
	}

	err = w.End()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
