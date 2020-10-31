package lolhtml_test

import (
	"bytes"
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestTextChunk_InsertBeforeAndAfter(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					TextChunkHandler: func(tc *lolhtml.TextChunk) lolhtml.RewriterDirective {
						content := tc.Content()
						if len(content) > 0 {
							if content != "Hey 42" {
								t.Errorf("got %s, want Hey 42", content)
							}
							if tc.IsLastInTextNode() {
								t.Error("text chunk last in text node flag incorrect, expected false, got true")
							}
							if tc.IsRemoved() {
								t.Error("text chunk removed flag incorrect, expected false, got true")
							}
							if err := tc.InsertBeforeAsHtml("<div>"); err != nil {
								t.Error(err)
							}
							if err := tc.InsertAfterAsText("</div>"); err != nil {
								t.Error(err)
							}
						} else {
							if !tc.IsLastInTextNode() {
								t.Error("text chunk last in text node flag incorrect, expected true, got false")
							}
						}
						return lolhtml.Continue
					},
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	defer w.Free()
	if _, err := w.Write([]byte("Hey 42")); err != nil {
		t.Error(err)
	}
	if err := w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<div>Hey 42&lt;/div&gt;"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestTextChunk_Replace(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					TextChunkHandler: func(tc *lolhtml.TextChunk) lolhtml.RewriterDirective {
						if len(tc.Content()) > 0 {
							if err := tc.ReplaceAsHtml("<repl>"); err != nil {
								t.Error(err)
							}
							if !tc.IsRemoved() {
								t.FailNow()
							}
						}
						return lolhtml.Continue
					},
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	defer w.Free()
	if _, err := w.Write([]byte("<div>Hello</div>")); err != nil {
		t.Error(err)
	}
	if err := w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<div><repl></div>"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestTextChunk_InsertAfter(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					TextChunkHandler: func(tc *lolhtml.TextChunk) lolhtml.RewriterDirective {
						if len(tc.Content()) > 0 {
							if err := tc.InsertAfterAsHtml("<after>"); err != nil {
								t.Error(err)
							}
						}
						return lolhtml.Continue
					},
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	defer w.Free()
	if _, err := w.Write([]byte("<div>hello</div>")); err != nil {
		t.Error(err)
	}
	if err := w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<div>hello<after></div>"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestTextChunk_Remove(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					TextChunkHandler: func(tc *lolhtml.TextChunk) lolhtml.RewriterDirective {
						if tc.IsRemoved() {
							t.FailNow()
						}
						tc.Remove()
						if !tc.IsRemoved() {
							t.FailNow()
						}
						return lolhtml.Continue
					},
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	defer w.Free()
	if _, err := w.Write([]byte("<span>0_0</span>")); err != nil {
		t.Error(err)
	}
	if err := w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<span></span>"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestTextChunk_StopRewriting(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					TextChunkHandler: func(tc *lolhtml.TextChunk) lolhtml.RewriterDirective {
						return lolhtml.Stop
					},
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	defer w.Free()
	_, err = w.Write([]byte("42"))
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "The rewriter has been stopped." {
		t.Error(err)
	}
}

func TestTextChunk_StopRewritingWithSelector(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "*",
					TextChunkHandler: func(tc *lolhtml.TextChunk) lolhtml.RewriterDirective {
						return lolhtml.Stop
					},
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	defer w.Free()
	_, err = w.Write([]byte("<div>42</div>"))
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "The rewriter has been stopped." {
		t.Error(err)
	}
}
