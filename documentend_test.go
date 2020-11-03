package lolhtml_test

import (
	"bytes"
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestDocumentEnd_AppendToEmptyDoc(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					DocumentEndHandler: func(docEnd *lolhtml.DocumentEnd) lolhtml.RewriterDirective {
						if err := docEnd.AppendAsHTML("<!--appended text-->"); err != nil {
							t.Error(err)
						}
						if err := docEnd.AppendAsText("hello & world"); err != nil {
							t.Error(err)
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

	if _, err = w.Write([]byte("")); err != nil {
		t.Error(err)
	}
	if err = w.Close(); err != nil {
		t.Error(err)
	}
	wantedText := "<!--appended text-->hello &amp; world"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestDocumentEnd_AppendAtEnd(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					DocumentEndHandler: func(docEnd *lolhtml.DocumentEnd) lolhtml.RewriterDirective {
						if err := docEnd.AppendAsHTML("<!--appended text-->"); err != nil {
							t.Error(err)
						}
						if err := docEnd.AppendAsText("hello & world"); err != nil {
							t.Error(err)
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

	if _, err = w.Write([]byte("<html><div>Hello</div></html>")); err != nil {
		t.Error(err)
	}
	if err = w.Close(); err != nil {
		t.Error(err)
	}
	wantedText := "<html><div>Hello</div></html><!--appended text-->hello &amp; world"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}
