package lolhtml_test

import (
	"errors"
	"github.com/coolspring8/go-lolhtml"
	"testing"
)

func TestRewriterBuilder(t *testing.T) {
	rb := lolhtml.NewRewriterBuilder()
	if rb == nil {
		t.Error("cannot get new rewriter-builder\n")
	}
	defer rb.Free()
	r, err := rb.Build(lolhtml.NewDefaultConfig())
	if err != nil {
		t.Errorf("cannot build rewriter %s\n", err)
	}
	defer r.Free()
	err = r.WriteString("<div>a<")
	if err != nil {
		t.Error(err)
	}
	err = r.WriteString("/div>")
	if err != nil {
		t.Error(err)
	}
	err = r.End()
	if err != nil {
		t.Error(err)
	}
}

func TestRewriterBuilderNonAsciiEncoding(t *testing.T) {
	rb := lolhtml.NewRewriterBuilder()
	if rb == nil {
		t.FailNow()
	}
	defer rb.Free()
	r, err := rb.Build(lolhtml.Config{
		Encoding: "UTF-16",
		Memory: &lolhtml.MemorySettings{
			PreallocatedParsingBufferSize: 0,
			MaxAllowedMemoryUsage:         16,
		},
		Sink:   func(string) {},
		Strict: true,
	})
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "Expected ASCII-compatible encoding." {
		t.Error(err)
	}
	r.Free()
}

func TestRewriterBuilderMemoryLimiting(t *testing.T) {
	rb := lolhtml.NewRewriterBuilder()
	if rb == nil {
		t.Error("cannot get new rewriter-builder\n")
	}
	defer rb.Free()
	s, err := lolhtml.NewSelector("span")
	if err != nil {
		t.Error(err)
	}
	defer s.Free()
	rb.AddElementContentHandlers(s, nil, nil, nil)
	r, err := rb.Build(lolhtml.Config{
		Encoding: "utf-8",
		Memory: &lolhtml.MemorySettings{
			PreallocatedParsingBufferSize: 0,
			MaxAllowedMemoryUsage:         5,
		},
		Sink:   func(string) {},
		Strict: true,
	})
	if err != nil {
		t.Error(err)
	}
	defer r.Free()
	err = r.WriteString("<span alt='aaaaa")
	if err.Error() != "The memory limit has been exceeded." {
		t.Error(err)
	}
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
			if err != nil && err.Error() != tc.errorText {
				t.Errorf("cannot get new selector: %s\n", err)
			}
			s.Free()
		})
	}
}

func TestDoctypeApi(t *testing.T) {
	rb := lolhtml.NewRewriterBuilder()
	defer rb.Free()
	rb.AddDocumentContentHandlers(
		func(doctype *lolhtml.Doctype) lolhtml.RewriterDirective {
			name := doctype.GetName()
			publicId := doctype.GetPublicId()
			systemId := doctype.GetSystemId()
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
		nil,
		nil,
		nil,
	)
	r, err := rb.Build(lolhtml.NewDefaultConfig())
	if err != nil {
		t.Error(err)
	}
	defer r.Free()
	err = r.WriteString(`<!DOCTYPE math SYSTEM "http://www.w3.org/Math/DTD/mathml1/mathml.dtd">`)
	if err != nil {
		t.Error(err)
	}
	err = r.End()
	if err != nil {
		t.Error(err)
	}
}

func TestCommentApi(t *testing.T) {
	rb := lolhtml.NewRewriterBuilder()
	defer rb.Free()
	rb.AddDocumentContentHandlers(
		nil,
		func(comment *lolhtml.Comment) lolhtml.RewriterDirective {
			text := comment.GetText()
			if text != "Hey 42" {
				t.Errorf("wrong text %s\n", text)
			}
			err := comment.SetText("Yo")
			if err != nil {
				t.Errorf("set text error %s\n", err)
			}
			return lolhtml.Continue
		},
		nil,
		nil,
	)
	var finalText string
	r, err := rb.Build(lolhtml.Config{
		Encoding: "utf-8",
		Memory: &lolhtml.MemorySettings{
			PreallocatedParsingBufferSize: 1024,
			MaxAllowedMemoryUsage:         1<<63 - 1,
		},
		Sink: func(s string) {
			finalText += s
		},
		Strict: false,
	})
	if err != nil {
		t.Error(err)
	}
	defer r.Free()
	err = r.WriteString("<!--Hey 42-->")
	if err != nil {
		t.Error(err)
	}
	err = r.End()
	if err != nil {
		t.Error(err)
	}
	if finalText != "<!--Yo-->" {
		t.Errorf("wrong output %s\n", finalText)
	}
}

func TestTextChunkApi(t *testing.T) {
	rb := lolhtml.NewRewriterBuilder()
	defer rb.Free()
	rb.AddDocumentContentHandlers(
		nil,
		nil,
		func(textChunk *lolhtml.TextChunk) lolhtml.RewriterDirective {
			content := textChunk.GetContent()
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
		nil,
	)
	var finalText string
	r, err := rb.Build(lolhtml.Config{
		Encoding: "utf-8",
		Memory: &lolhtml.MemorySettings{
			PreallocatedParsingBufferSize: 1024,
			MaxAllowedMemoryUsage:         1<<63 - 1,
		},
		Sink: func(s string) {
			finalText += s
		},
		Strict: false,
	})
	if err != nil {
		t.Error(err)
	}
	defer r.Free()
	err = r.WriteString("Hey 42")
	if err != nil {
		t.Error(err)
	}
	err = r.End()
	if err != nil {
		t.Error(err)
	}
	if finalText != "<div>Hey 42&lt;/div&gt;" {
		t.Errorf("wrong output %s\n", finalText)
	}
}

func TestDocEndApi(t *testing.T) {
	rb := lolhtml.NewRewriterBuilder()
	defer rb.Free()
	rb.AddDocumentContentHandlers(
		nil,
		nil,
		nil,
		func(docEnd *lolhtml.DocEnd) lolhtml.RewriterDirective {
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
	)
	var finalText string
	r, err := rb.Build(lolhtml.Config{
		Encoding: "utf-8",
		Memory: &lolhtml.MemorySettings{
			PreallocatedParsingBufferSize: 1024,
			MaxAllowedMemoryUsage:         1<<63 - 1,
		},
		Sink: func(s string) {
			finalText += s
		},
		Strict: false,
	})
	if err != nil {
		t.Error(err)
	}
	defer r.Free()
	err = r.WriteString("")
	if err != nil {
		t.Error(err)
	}
	err = r.End()
	if err != nil {
		t.Error(err)
	}
	if finalText != "<!--appended text-->hello &amp; world" {
		t.Errorf("wrong output %s\n", finalText)
	}
}

func TestElementApi(t *testing.T) {
	rb := lolhtml.NewRewriterBuilder()
	defer rb.Free()
	s, err := lolhtml.NewSelector("*")
	if err != nil {
		t.Error(err)
	}
	defer s.Free()
	rb.AddElementContentHandlers(
		s,
		func(element *lolhtml.Element) lolhtml.RewriterDirective {
			name := element.GetTagName()
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
		nil,
		nil,
	)
	var finalText string
	r, err := rb.Build(lolhtml.Config{
		Encoding: "utf-8",
		Memory: &lolhtml.MemorySettings{
			PreallocatedParsingBufferSize: 1024,
			MaxAllowedMemoryUsage:         1<<63 - 1,
		},
		Sink: func(s string) {
			finalText += s
		},
		Strict: false,
	})
	if err != nil {
		t.Error(err)
	}
	defer r.Free()
	err = r.WriteString("Hi <div>")
	if err != nil {
		t.Error(err)
	}
	err = r.End()
	if err != nil {
		t.Error(err)
	}
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
