package lolhtml_test

import (
	"bytes"
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestCommentApi(t *testing.T) {
	var b bytes.Buffer
	w, err := lolhtml.NewWriter(
		&b,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					CommentHandler: func(comment *lolhtml.Comment) lolhtml.RewriterDirective {
						text := comment.Text()
						if text != "Hey 42" {
							t.Errorf("wrong text %s\n", text)
						}
						err := comment.SetText("Yo")
						if err != nil {
							t.Errorf("set text error %s\n", err)
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
	_, err = w.Write([]byte("<!--Hey 42-->"))
	if err != nil {
		t.Error(err)
	}
	err = w.End()
	if err != nil {
		t.Error(err)
	}
	finalText := b.String()
	if finalText != "<!--Yo-->" {
		t.Errorf("wrong output %s\n", finalText)
	}
}
