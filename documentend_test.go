package lolhtml_test

import (
	"bytes"
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestDocumentEndApi(t *testing.T) {
	var b bytes.Buffer
	w, err := lolhtml.NewWriter(
		&b,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					DocumentEndHandler: func(docEnd *lolhtml.DocumentEnd) lolhtml.RewriterDirective {
						err := docEnd.AppendAsHtml("<!--appended text-->")
						if err != nil {
							t.Error(err)
						}
						err = docEnd.AppendAsText("hello & world")
						if err != nil {
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
	defer w.Free()
	_, err = w.Write([]byte(""))
	if err != nil {
		t.Error(err)
	}
	err = w.End()
	if err != nil {
		t.Error(err)
	}
	finalText := b.String()
	if finalText != "<!--appended text-->hello &amp; world" {
		t.Errorf("wrong output %s\n", finalText)
	}
}
