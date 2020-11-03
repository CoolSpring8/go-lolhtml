# go-lolhtml

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/coolspring8/go-lolhtml/Go) [![codecov](https://codecov.io/gh/CoolSpring8/go-lolhtml/branch/main/graph/badge.svg)](https://codecov.io/gh/CoolSpring8/go-lolhtml) [![Go Report Card](https://goreportcard.com/badge/github.com/coolspring8/go-lolhtml)](https://goreportcard.com/report/github.com/coolspring8/go-lolhtml) [![PkgGoDev](https://pkg.go.dev/badge/github.com/coolspring8/go-lolhtml)](https://pkg.go.dev/github.com/coolspring8/go-lolhtml)

Go bindings for the Rust crate [cloudflare/lol-html](https://github.com/cloudflare/lol-html/), the *Low Output Latency streaming HTML rewriter/parser with CSS-selector based API*, talking via cgo.

**Status:** 

**All abilities provided by lol_html's c-api are available**, except for customized user data in handlers. The original tests included in c-api package have also been translated to examine this binding's functionality.

The code is at its early stage and **breaking changes might be introduced**. If you have any ideas on how the public API can be better structured, feel free to open a PR or an issue.

   * [go-lolhtml](#go-lolhtml)
      * [Installation](#installation)
      * [Features](#features)
      * [Getting Started](#getting-started)
      * [Examples](#examples)
      * [Documentation](#documentation)
      * [Other Bindings](#other-bindings)
      * [Versioning](#versioning)
      * [Help Wanted!](#help-wanted)
      * [License](#license)
      * [Disclaimer](#disclaimer)

## Installation

For Linux/macOS/Windows x86_64 platform users, installation is as simple as a single `go get` command:

```shell
$ go get github.com/coolspring8/go-lolhtml
```

Installing Rust is not a necessary step. That's because lol-html could be prebuilt into static libraries, stored and shipped in `/build` folder, so that cgo can handle other compilation matters naturally and smoothly, without intervention.

For other platforms, you will have to compile it yourself.

## Features

- Fast: A Go (cgo) wrapper built around the highly-optimized Rust HTML parsing crate lol_html.
- Easy to use: Utilizing Go's idiomatic I/O methods, [lolhtml.Writer](https://pkg.go.dev/github.com/coolspring8/go-lolhtml#Writer) implements [io.Writer](https://golang.org/pkg/io/#Writer) interface.

## Getting Started

Now let's initialize a project and create `main.go`:

```go
package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/coolspring8/go-lolhtml"
)

func main() {
	chunk := []byte("Hello, <span>World</span>!")
	r := bytes.NewReader(chunk)
	w, err := lolhtml.NewWriter(
		// output to stdout
		os.Stdout,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					Selector: "span",
					ElementHandler: func(e *lolhtml.Element) lolhtml.RewriterDirective {
						err := e.SetInnerContentAsText("LOL-HTML")
						if err != nil {
							log.Fatal(err)
						}
						return lolhtml.Continue
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// copy from the bytes reader to lolhtml writer
	_, err = io.Copy(w, r)
	if err != nil {
		log.Fatal(err)
	}

	// explicitly close the writer and flush the remaining content
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
	// Output: Hello, <span>LOL-HTML</span>!
}
```

The above program creates a new Writer configured to rewrite all texts in `span` tags to "LOL-HTML". It takes the chunk `Hello, <span>World</span>!` as input, and prints the result to standard output.

And the result is `Hello, <span>LOL-HTML</span>!` .

## Examples

example_test.go contains two examples.

For more detailed examples, please visit the `/examples` subdirectory.

- defer-scripts

  Usage: curl -NL https://git.io/JeOSZ | go run main.go

- mixed-content-rewriter

  Usage: curl -NL https://git.io/JeOSZ | go run main.go

- web-scraper

  A ported Go version of https://web.scraper.workers.dev/.

## Documentation

Available at [pkg.go.dev](https://pkg.go.dev/github.com/coolspring8/go-lolhtml).

## Other Bindings

- Rust (native), C, JavaScript - [cloudflare/lol-html](https://github.com/cloudflare/lol-html/)
- Lua - [jdesgats/lua-lolhtml](https://github.com/jdesgats/lua-lolhtml/)

## Versioning

This package does not really follow [Semantic Versioning](https://semver.org/). The current strategy is to follow lol_html's major and minor version, and the patch version number is reserved for this binding's updates, for Go Modul to upgrade correctly.

## Help Wanted!

There are a few interesting things at [Projects](https://github.com/coolspring8/go-lolhtml/projects/1) panel that I have considered but is not yet implemented. Other contributions and suggestions are also welcome!

## License

BSD 3-Clause "New" or "Revised" License

## Disclaimer

This is an unofficial binding.

Cloudflare is a registered trademark of Cloudflare, Inc. Cloudflare names used in this project are for identification purposes only. The project is not associated in any way with Cloudflare Inc.