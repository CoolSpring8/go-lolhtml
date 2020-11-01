// This is a ported Go version of https://web.scraper.workers.dev/, whose source code is
// available at https://github.com/adamschwartz/web.scraper.workers.dev licensed under MIT.
//
// This translation is for demonstration purpose only, so many parts of the code are suboptimal.
//
// Sometimes you may get a "different" result, as Go's encoding/json package always sorts the
// keys of a map (When using multiple selectors), and encodes a nil slice as the null JSON value.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/coolspring8/go-lolhtml"
)

var (
	listenAddress    = ":80"
	mainPageFileName = "index.html"
)

var (
	urlHasPrefix    = regexp.MustCompile(`^[a-zA-Z]+://`)
	unifyWhitespace = regexp.MustCompile(`\s{2,}`)
)

// used to separate texts in different elements.
var textSeparator = "TEXT_SEPARATOR_TEXT_SEPARATOR"

func main() {
	log.Printf("Server started at %s", listenAddress)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	log.Println(req.URL)

	// 404
	if req.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not found"))
		return
	}

	q := req.URL.Query()

	url := q.Get("url")
	if url != "" && !urlHasPrefix.MatchString(url) {
		url = "http://" + url
	}

	selector := q.Get("selector")

	attr := q.Get("attr")

	var spaced bool
	_spaced := q.Get("spaced")
	if _spaced == "true" {
		spaced = true
	} else {
		spaced = false
	}

	var pretty bool
	_pretty := q.Get("pretty")
	if _pretty == "true" {
		pretty = true
	} else {
		pretty = false
	}

	// home page
	if url == "" && selector == "" {
		http.ServeFile(w, req, mainPageFileName)
		return
	}

	// text or attr: get text, part 1/2
	handlers := lolhtml.Handlers{}
	// matches and selectors are used by text scraper
	matches := make(map[string][]string)
	var selectors []string
	_selectors := strings.Split(selector, ",")
	for _, s := range _selectors {
		selectors = append(selectors, strings.TrimSpace(s))
	}
	// attrValue is used by attribute scraper
	var attrValue string
	if attr == "" {
		nextText := make(map[string]string)

		// sort of confusing...
		// https://stackoverflow.com/questions/27646434/reference-to-another-field-with-in-a-func-in-the-same-struct-in-golang
		// Are there better ways?
		for _, s := range selectors {
			elementContentHandler := &lolhtml.ElementContentHandler{
				Selector: s,
			}
			elementContentHandler.ElementHandler = func(ech *lolhtml.ElementContentHandler) lolhtml.ElementHandlerFunc {
				return func(e *lolhtml.Element) lolhtml.RewriterDirective {
					matches[ech.Selector] = append(matches[ech.Selector], textSeparator)
					nextText[ech.Selector] = ""
					return lolhtml.Continue
				}
			}(elementContentHandler)
			elementContentHandler.TextChunkHandler = func(ech *lolhtml.ElementContentHandler) lolhtml.TextChunkHandlerFunc {
				return func(t *lolhtml.TextChunk) lolhtml.RewriterDirective {
					nextText[ech.Selector] += t.Content()
					if t.IsLastInTextNode() {
						if spaced {
							nextText[ech.Selector] += " "
						}
						matches[ech.Selector] = append(matches[ech.Selector], nextText[ech.Selector])
						nextText[ech.Selector] = ""
					}
					return lolhtml.Continue
				}
			}(elementContentHandler)
			handlers.ElementContentHandler = append(
				handlers.ElementContentHandler,
				*elementContentHandler,
			)
		}
	} else {
		handlers = lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: selector,
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						attrValue, _ = e.AttributeValue(attr)
						return lolhtml.Stop
					},
				},
			},
		}
	}

	lolWriter, err := lolhtml.NewWriter(
		nil,
		&handlers,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error(), pretty)
		return
	}

	// fetch target page content
	resp, err := http.Get(url)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error(), pretty)
		return
	}
	if resp.StatusCode != http.StatusOK {
		sendError(w, http.StatusBadGateway, fmt.Sprintf("Status %d requesting %s", resp.StatusCode, url), pretty)
		return
	}
	defer resp.Body.Close()

	// another confusing point
	_, err = io.Copy(lolWriter, resp.Body)
	if err != nil && err.Error() != "The rewriter has been stopped." {
		sendError(w, http.StatusInternalServerError, err.Error(), pretty)
		return
	}
	if err == nil || err.Error() != "The rewriter has been stopped." {
		err = lolWriter.End()
		if err != nil {
			sendError(w, http.StatusInternalServerError, err.Error(), pretty)
			return
		}
	}
	lolWriter.Free()

	// text or attr: post-process texts, part 2/2
	if attr == "" {
		for _, s := range selectors {
			var nodeCompleteTexts []string
			nextText := ""

			for _, text := range matches[s] {
				if text == textSeparator {
					if strings.TrimSpace(nextText) != "" {
						nodeCompleteTexts = append(nodeCompleteTexts, cleanText(nextText))
						nextText = ""
					}
				} else {
					nextText += text
				}
			}

			lastText := cleanText(nextText)
			if lastText != "" {
				nodeCompleteTexts = append(nodeCompleteTexts, lastText)
			}
			matches[s] = nodeCompleteTexts
		}
	}

	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if pretty {
		enc.SetIndent("", "  ")
	}

	if attr == "" {
		err = enc.Encode(Response{Result: matches})
	} else {
		err = enc.Encode(Response{Result: attrValue})
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error(), pretty)
		return
	}
}

type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func sendError(w http.ResponseWriter, statusCode int, errorText string, pretty bool) {
	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if pretty {
		enc.SetIndent("", "  ")
	}

	err := enc.Encode(Response{Error: errorText})
	if err != nil {
		_, _ = w.Write([]byte(errorText))
	}
}

func cleanText(s string) string {
	return unifyWhitespace.ReplaceAllString(strings.TrimSpace(s), " ")
}
