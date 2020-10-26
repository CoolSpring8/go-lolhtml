package lolhtml_test

import (
	"bytes"
	"errors"
	"github.com/coolspring8/go-lolhtml"
	"testing"
)

func TestRewriter_NonAsciiEncoding(t *testing.T) {
	_, err := lolhtml.NewWriter(
		nil,
		nil,
		lolhtml.Config{
			Encoding: "UTF-16",
			Memory: &lolhtml.MemorySettings{
				PreallocatedParsingBufferSize: 1024,
				MaxAllowedMemoryUsage:         1<<63 - 1,
			},
			Strict: true,
		})
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "Expected ASCII-compatible encoding." {
		t.Error(err)
	}
}

func TestRewriterBuilderMemoryLimiting(t *testing.T) {
	w, err := lolhtml.NewWriter(
		nil,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					"span",
					nil,
					nil,
					nil,
				},
			},
		},
		lolhtml.Config{
			Encoding: "utf-8",
			Memory: &lolhtml.MemorySettings{
				PreallocatedParsingBufferSize: 0,
				MaxAllowedMemoryUsage:         5,
			},
			Strict: true,
		},
	)
	if err != nil {
		t.Error(err)
	}
	_, err = w.Write([]byte("<span alt='aaaaa"))
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "The memory limit has been exceeded." {
		t.Error(err)
	}
	w.Free()
}

func TestNewSelector(t *testing.T) {
	testCases := []struct {
		selector  string
		errorText string
	}{
		{"p.center", ""},
		{"p:last-child", "Unsupported pseudo-class or pseudo-element in selector."},
	}
	for _, tc := range testCases {
		t.Run(tc.selector, func(t *testing.T) {
			s, err := lolhtml.NewSelector(tc.selector)
			if err == nil {
				if tc.errorText != "" {
					t.FailNow()
				}
			} else {
				if err.Error() != tc.errorText {
					t.Error(err)
				}
				s.Free()
			}
		})
	}
}

func TestDoctypeApi(t *testing.T) {
	w, err := lolhtml.NewWriter(
		nil,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					DoctypeHandler: func(doctype *lolhtml.Doctype) lolhtml.RewriterDirective {
						name := doctype.Name()
						publicId := doctype.PublicId()
						systemId := doctype.SystemId()
						if name != "math" {
							t.Errorf("wrong doctype name %s\n", name)
						}
						if publicId != "" {
							t.Errorf("wrong doctype name %s\n", publicId)
						}
						if systemId != "http://www.w3.org/Math/DTD/mathml1/mathml.dtd" {
							t.Errorf("wrong doctype name %s\n", systemId)
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
	_, err = w.Write([]byte(`<!DOCTYPE math SYSTEM "http://www.w3.org/Math/DTD/mathml1/mathml.dtd">`))
	if err != nil {
		t.Error(err)
	}
	err = w.End()
	if err != nil {
		t.Error(err)
	}
}

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
							err = textChunk.InsertAfterAsRaw("</div>")
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

func TestDocEndApi(t *testing.T) {
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
						err = docEnd.AppendAsRaw("hello & world")
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

func TestElementApi(t *testing.T) {
	var b bytes.Buffer
	w, err:= lolhtml.NewWriter(
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

// TestNullErrorStr tests internal functions for handling a null lol_html_str_t, by calling lol_html_take_last_error()
// when there is no error.
func TestNullErrorStr(t *testing.T) {
	err := lolhtml.GetError()
	if !errors.Is(err, lolhtml.ErrCannotGetErrorMessage) {
		t.Error(err)
	}
}
