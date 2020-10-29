package lolhtml_test

import (
	"bytes"
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestElementApi(t *testing.T) {
	var b bytes.Buffer
	w, err := lolhtml.NewWriter(
		&b,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "*",
					ElementHandler: func(element *lolhtml.Element) lolhtml.RewriterDirective {
						name := element.TagName()
						if name != "div" {
							t.Errorf("get wrong tag name %s\n", name)
						}
						err := element.SetTagName("")
						if err != nil && err.Error() != "Tag name can't be empty." {
							t.Error(err)
						}
						err = element.SetTagName("span")
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
	_, err = w.Write([]byte("Hi <div>"))
	if err != nil {
		t.Error(err)
	}
	err = w.End()
	if err != nil {
		t.Error(err)
	}
	finalText := b.String()
	if finalText != "Hi <span>" {
		t.Errorf("wrong output %s", finalText)
	}
}
