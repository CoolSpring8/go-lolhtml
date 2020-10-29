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
					var buf bytes.Buffer
					w, err := lolhtml.NewWriter(&buf, bm.handlers)
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

// BenchmarkNewRewriterBuilder tests the cost of not caching the RewriterBuilder.
// In current implementation, every call to NewWriter involves creating a RewriterBuilder.
// Ideally, if multiple documents are to be processed at the same time, a RewriterBuilder
// is created in advance and kept in memory, and Rewriter-s could be built from it.
// This benchmark aims to find out the extent of its impact.
func BenchmarkNewRewriterBuilder(b *testing.B) {
	b.Run("Builder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = lolhtml.NewRewriterBuilder()
		}
	})
	b.Run("BuilderWithFree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rb := lolhtml.NewRewriterBuilder()
			rb.Free()
		}
	})
	b.Run("BuilderWithEmptyDocumentHandlerAndFree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rb := lolhtml.NewRewriterBuilder()
			rb.AddDocumentContentHandlers(nil, nil, nil, nil)
			rb.Free()
		}
	})
	b.Run("BuilderWithEmptyElementHandlerAndFree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rb := lolhtml.NewRewriterBuilder()
			s, _ := lolhtml.NewSelector("*")
			rb.AddElementContentHandlers(s, nil, nil, nil)
			rb.Free()
			s.Free()
		}
	})
	b.Run("BuilderWithElementHandlerAndFree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rb := lolhtml.NewRewriterBuilder()
			s, _ := lolhtml.NewSelector("*")
			rb.AddElementContentHandlers(
				s,
				func(e *lolhtml.Element) lolhtml.RewriterDirective { return lolhtml.Continue },
				nil,
				nil,
			)
			rb.Free()
			s.Free()
		}
	})
	b.Run("BuilderWithElementHandlerAndBuildAndFree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rb := lolhtml.NewRewriterBuilder()
			s, _ := lolhtml.NewSelector("*")
			rb.AddElementContentHandlers(
				s,
				func(e *lolhtml.Element) lolhtml.RewriterDirective { return lolhtml.Continue },
				nil,
				nil,
			)
			_, _ = rb.Build(func([]byte) {}, lolhtml.NewDefaultConfig())
			rb.Free()
			s.Free()
		}
	})
	b.Run("BuildMultipleRewriterFromOneBuilder", func(b *testing.B) {
		rb := lolhtml.NewRewriterBuilder()
		s, _ := lolhtml.NewSelector("*")
		rb.AddElementContentHandlers(
			s,
			func(e *lolhtml.Element) lolhtml.RewriterDirective { return lolhtml.Continue },
			nil,
			nil,
		)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = rb.Build(func([]byte) {}, lolhtml.NewDefaultConfig())
		}
		b.StopTimer()
		rb.Free()
		s.Free()
	})
	b.Run("Writer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = lolhtml.NewWriter(nil,
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
			)
		}
	})
}
