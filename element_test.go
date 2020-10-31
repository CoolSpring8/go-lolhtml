package lolhtml_test

import (
	"bytes"
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestElement_ModifyTagName(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "*",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						wantName := "div"
						if name := e.TagName(); name != wantName {
							t.Errorf("got %s want %s\n", name, wantName)
						}
						err := e.SetTagName("")
						if err == nil {
							t.FailNow()
						}
						if err.Error() != "Tag name can't be empty." {
							t.Error(err)
						}
						if err = e.SetTagName("span"); err != nil {
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
	if _, err = w.Write([]byte("Hi <div>")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "Hi <span>"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestElement_ModifyAttributes(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "*",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						has, err := e.HasAttribute("foo")
						if err != nil {
							t.Error(err)
						}
						if !has {
							t.FailNow()
						}
						has, err = e.HasAttribute("Bar")
						if err != nil {
							t.Error(err)
						}
						if has {
							t.FailNow()
						}

						a, err := e.AttributeValue("foo")
						if err != nil {
							t.Error(err)
						}
						wantValue := "42"
						if a != wantValue {
							t.Errorf("got %s; want %s", a, wantValue)
						}
						a, err = e.AttributeValue("Bar")
						if err != nil {
							t.Error(err)
						}
						if a != "" {
							t.Errorf("got %s; want empty", a)
						}

						if err := e.SetAttribute("Bar", "hey"); err != nil {
							t.Error(err)
						}

						if err := e.RemoveAttribute("foo"); err != nil {
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
	if _, err = w.Write([]byte("<span foo=42>")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<span bar=\"hey\">"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestElement_InsertContentAroundElement(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "*",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						if err := e.InsertBeforeStartTagAsText("&before"); err != nil {
							t.Error(err)
						}
						if err := e.InsertAfterStartTagAsHtml("<!--prepend-->"); err != nil {
							t.Error(err)
						}
						if err := e.InsertBeforeEndTagAsHtml("<!--append-->"); err != nil {
							t.Error(err)
						}
						if err := e.InsertAfterEndTagAsText("&after"); err != nil {
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
	if _, err = w.Write([]byte("<div>Hi</div>")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "&amp;before<div><!--prepend-->Hi<!--append--></div>&amp;after"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestElement_SetInnerContent(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "div",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						if err := e.SetInnerContentAsText("hey & ya"); err != nil {
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
	if _, err = w.Write([]byte("<div><span>42</span></div>")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<div>hey &amp; ya</div>"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestElement_Replace(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "div",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						if err := e.ReplaceAsHtml("hey & ya"); err != nil {
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
	if _, err = w.Write([]byte("<div><span>42</span></div><h1>Hello<div>good bye</div></h1><h2>Hello2</h2>")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "hey & ya<h1>Hellohey & ya</h1><h2>Hello2</h2>"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestElement_Remove(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "h1",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						if e.IsRemoved() {
							t.FailNow()
						}
						e.Remove()
						if !e.IsRemoved() {
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
	if _, err = w.Write([]byte("<div><span>42</span></div><h1>Hello</h1><h2>Hello2</h2>")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<div><span>42</span></div><h2>Hello2</h2>"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestElement_RemoveElementAndKeepContent(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "h2",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						if e.IsRemoved() {
							t.FailNow()
						}
						e.RemoveAndKeepContent()
						if !e.IsRemoved() {
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
	if _, err = w.Write([]byte("<div><span>42<h2>Hello1</h2></span></div><h1>Hello</h1><h2>Hello2</h2>")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<div><span>42Hello1</span></div><h1>Hello</h1>Hello2"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestElement_GetEmptyElementAttribute(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "span",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						has, err := e.HasAttribute("foo")
						if err != nil {
							t.Error(err)
						}
						if !has {
							t.FailNow()
						}
						value, err := e.AttributeValue("foo")
						if err != nil {
							t.Error(err)
						}
						if value != "" {
							t.Errorf("got %s; want empty", value)
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
	if _, err = w.Write([]byte("<span foo>")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<span foo>"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestElement_IterateAttributes(t *testing.T) {
	w, err := lolhtml.NewWriter(
		nil,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "*",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						ai := e.AttributeIterator()

						a := ai.Next()
						if name := a.Name(); name != "foo" {
							t.Errorf("got %s; want foo", name)
						}
						if value := a.Value(); value != "42" {
							t.Errorf("got %s; want foo", value)
						}

						a = ai.Next()
						if name := a.Name(); name != "bar" {
							t.Errorf("got %s; want bar", name)
						}
						if value := a.Value(); value != "1337" {
							t.Errorf("got %s; want 1337", value)
						}

						a = ai.Next()
						if a != nil {
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
	if _, err = w.Write([]byte("<div foo=42 bar='1337'>")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
}

func TestElement_AssertNsIsHtml(t *testing.T) {
	w, err := lolhtml.NewWriter(
		nil,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "script",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						wantedText := "http://www.w3.org/1999/xhtml"
						if ns := e.NamespaceUri(); ns != wantedText {
							t.Errorf("got %s; want %s", ns, wantedText)
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
	if _, err = w.Write([]byte("<script></script>")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
}

func TestElement_AssertNsIsSvg(t *testing.T) {
	w, err := lolhtml.NewWriter(
		nil,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "script",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						wantedText := "http://www.w3.org/2000/svg"
						if ns := e.NamespaceUri(); ns != wantedText {
							t.Errorf("got %s; want %s", ns, wantedText)
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
	if _, err = w.Write([]byte("<svg><script></script></svg>")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
}

func TestElement_StopRewriting(t *testing.T) {
	w, err := lolhtml.NewWriter(
		nil,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "span",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
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
	_, err = w.Write([]byte("<span foo>"))
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "The rewriter has been stopped." {
		t.Error(err)
	}
}
