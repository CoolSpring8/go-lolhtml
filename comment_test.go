package lolhtml_test

import (
	"bytes"
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestComment_GetSetText(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					CommentHandler: func(comment *lolhtml.Comment) lolhtml.RewriterDirective {
						if text := comment.Text(); text != "Hey 42" {
							t.Errorf("wrong text %s\n", text)
						}
						if err := comment.SetText("Yo"); err != nil {
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
	if _, err = w.Write([]byte("<!--Hey 42-->")); err != nil {
		t.Error(err)
	}
	if err = w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<!--Yo-->"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestComment_Replace(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					CommentHandler: func(c *lolhtml.Comment) lolhtml.RewriterDirective {
						if err := c.ReplaceAsHtml("<repl>"); err != nil {
							t.Error(err)
						}
						if !c.IsRemoved() {
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
	if _, err := w.Write([]byte("<div><!--hello--></div>")); err != nil {
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

func TestComment_InsertAfter(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					CommentHandler: func(c *lolhtml.Comment) lolhtml.RewriterDirective {
						if err := c.InsertAfterAsHtml("<after>"); err != nil {
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
	if _, err := w.Write([]byte("<div><!--hello--></div>")); err != nil {
		t.Error(err)
	}
	if err := w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<div><!--hello--><after></div>"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestComment_Remove(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					CommentHandler: func(c *lolhtml.Comment) lolhtml.RewriterDirective {
						if c.IsRemoved() {
							t.FailNow()
						}
						c.Remove()
						if !c.IsRemoved() {
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
	if _, err := w.Write([]byte("<<!--0_0-->>")); err != nil {
		t.Error(err)
	}
	if err := w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<>"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestComment_InsertBeforeAndAfter(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					CommentHandler: func(c *lolhtml.Comment) lolhtml.RewriterDirective {
						if err := c.InsertBeforeAsHtml("<div>"); err != nil {
							t.Error(err)
						}
						if err := c.InsertAfterAsText("</div>"); err != nil {
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
	if _, err := w.Write([]byte("<!--Hey 42-->")); err != nil {
		t.Error(err)
	}
	if err := w.End(); err != nil {
		t.Error(err)
	}
	wantedText := "<div><!--Hey 42-->&lt;/div&gt;"
	if finalText := buf.String(); finalText != wantedText {
		t.Errorf("want %s got %s \n", wantedText, finalText)
	}
}

func TestComment_StopRewriting(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					CommentHandler: func(c *lolhtml.Comment) lolhtml.RewriterDirective {
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
	_, err = w.Write([]byte("<div><!-- foo --></div>"))
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "The rewriter has been stopped." {
		t.Error(err)
	}
}

func TestComment_StopRewritingWithSelector(t *testing.T) {
	var buf bytes.Buffer
	w, err := lolhtml.NewWriter(
		&buf,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "*",
					CommentHandler: func(c *lolhtml.Comment) lolhtml.RewriterDirective {
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
	_, err = w.Write([]byte("<div><!-- foo --></div>"))
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "The rewriter has been stopped." {
		t.Error(err)
	}
}
