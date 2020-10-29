package lolhtml_test

import (
	"bytes"
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestTextChunkApi(t *testing.T) {
	var b bytes.Buffer
	w, err := lolhtml.NewWriter(
		&b,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					TextChunkHandler: func(textChunk *lolhtml.TextChunk) lolhtml.RewriterDirective {
						content := textChunk.Content()
						if len(content) > 0 {
							if textChunk.IsLastInTextNode() {
								t.Error("text chunk last in text node flag incorrect, expected false, got true")
							}
							if textChunk.IsRemoved() {
								t.Error("text chunk removed flag incorrect, expected false, got true")
							}
							err := textChunk.InsertBeforeAsHtml("<div>")
							if err != nil {
								t.Error(err)
							}
							err = textChunk.InsertAfterAsText("</div>")
							if err != nil {
								t.Error(err)
							}
						} else {
							if !textChunk.IsLastInTextNode() {
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
	_, err = w.Write([]byte("Hey 42"))
	if err != nil {
		t.Error(err)
	}
	err = w.End()
	if err != nil {
		t.Error(err)
	}
	finalText := b.String()
	if finalText != "<div>Hey 42&lt;/div&gt;" {
		t.Errorf("wrong output %s\n", finalText)
	}
}
