package lolhtml_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

const dataDir = "testdata"

const ChunkSize = 1024

func BenchmarkNewWriter(b *testing.B) {
	benchmarks := []struct {
		category string
		name     string
		handlers *lolhtml.Handlers
	}{
		{
			"Parsing",
			"TagScanner",
			nil,
		},
		{
			"Parsing",
			"Lexer",
			&lolhtml.Handlers{
				DocumentContentHandler: []lolhtml.DocumentContentHandler{
					{
						DoctypeHandler: func(d *lolhtml.Doctype) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
				},
			},
		},
		{
			"Parsing",
			"TextRewritableUnitParsingAndDecoding",
			&lolhtml.Handlers{
				DocumentContentHandler: []lolhtml.DocumentContentHandler{
					{
						TextChunkHandler: func(c *lolhtml.TextChunk) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
				},
			},
		},
		{
			"Rewriting",
			"ModificationOfTagsOfAnElementWithLotsOfContent",
			&lolhtml.Handlers{
				ElementContentHandler: []lolhtml.ElementContentHandler{
					{
						Selector: "body",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							err := e.SetTagName("body1")
							if err != nil {
								b.Fatal(err)
							}
							err = e.InsertAfterEndTagAsText("test")
							if err != nil {
								b.Fatal(err)
							}
							return lolhtml.Continue
						},
					},
				},
			},
		},
		{
			"Rewriting",
			"RemoveContentOfAnElement",
			&lolhtml.Handlers{
				ElementContentHandler: []lolhtml.ElementContentHandler{
					{
						Selector: "ul",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							err := e.SetInnerContentAsText("")
							if err != nil {
								b.Fatal(err)
							}
							return lolhtml.Continue
						},
					},
				},
			},
		},
		{
			"SelectorMatching",
			"MatchAllSelector",
			&lolhtml.Handlers{
				ElementContentHandler: []lolhtml.ElementContentHandler{
					{
						Selector: "*",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
				},
			},
		},
		{
			"SelectorMatching",
			"TagNameSelector",
			&lolhtml.Handlers{
				ElementContentHandler: []lolhtml.ElementContentHandler{
					{
						Selector: "div",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
				},
			},
		},
		{
			"SelectorMatching",
			"ClassSelector",
			&lolhtml.Handlers{
				ElementContentHandler: []lolhtml.ElementContentHandler{
					{
						Selector: ".note",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
				},
			},
		},
		{
			"SelectorMatching",
			"AttributeSelector",
			&lolhtml.Handlers{
				ElementContentHandler: []lolhtml.ElementContentHandler{
					{
						Selector: "[href]",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
				},
			},
		},
		{
			"SelectorMatching",
			"MultipleSelectors",
			&lolhtml.Handlers{
				ElementContentHandler: []lolhtml.ElementContentHandler{
					{
						Selector: "ul",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
					{
						Selector: "ul > li",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
					{
						Selector: "table > tbody td dfn",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
					{
						Selector: "body table > tbody tr",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
					{
						Selector: "body [href]",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
					{
						Selector: "div img",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
					{
						Selector: "div.note span",
						ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
							return lolhtml.Continue
						},
					},
				},
			},
		},
	}

	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		b.Fatal("benchmark data files not found", err)
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(filepath.Join(dataDir, file.Name()))
		if err != nil {
			b.Fatal("cannot read benchmark data files", err)
		}

		for _, bm := range benchmarks {
			b.Run(fmt.Sprintf("%s-%s-%s", bm.category, bm.name, file.Name()), func(b *testing.B) {
				b.SetBytes(int64(len(data)))
				b.ReportAllocs()
				runtime.GC()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					w, err := lolhtml.NewWriter(nil, bm.handlers)
					if err != nil {
						b.Fatal(err)
					}

					r := bytes.NewReader(data)
					copyBuf := make([]byte, ChunkSize)
					_, err = io.CopyBuffer(w, r, copyBuf)
					if err != nil {
						b.Fatal(err)
					}

					err = w.End()
					if err != nil {
						b.Fatal(err)
					}
					w.Free()
				}
			})
		}
	}
}
