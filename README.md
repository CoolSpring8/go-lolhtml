# go-lolhtml

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/coolspring8/go-lolhtml/Go)![Codecov](https://img.shields.io/codecov/c/github/coolspring8/go-lolhtml)[![Go Report Card](https://goreportcard.com/badge/github.com/coolspring8/go-lolhtml)](https://goreportcard.com/report/github.com/coolspring8/go-lolhtml)[![PkgGoDev](https://pkg.go.dev/badge/CoolSpring8/go-lolhtml)](https://pkg.go.dev/CoolSpring8/go-lolhtml)

Go bindings for the Rust library [cloudflare/lol-html](https://github.com/cloudflare/lol-html/), the *Low Output Latency streaming HTML rewriter/parser with CSS-selector based API*, talking via cgo.

**Status:** All abilities provided by C-API implemented, except for customized user data in handlers. The code is at its early stage and the API is therefore subject to change. If you have any ideas on how API can be better structured, feel free to open a PR or an issue.

## Installation

For Linux/macOS/Windows x86_64 platforms, installation is as simple as a single `go get`:

```shell
$ go get github.com/coolspring8/go-lolhtml
```

There is no need for you to install Rust. That's because lib-lolhtml could be prebuilt into static libraries, stored and shipped in `/build` folder, so that cgo can handle other matters naturally and smoothly.

(For other platforms, you'll have to compile it yourself.)

## Getting Started

Now let's initialize a project and create `main.go`:

```go
package main

import (
    "fmt"
    "github.com/coolspring8/go-lolhtml"
)

func main() {
	rb := lolhtml.NewRewriterBuilder()
	defer rb.Free()
	s, _ := lolhtml.NewSelector("span")
	defer s.Free()
	rb.AddElementContentHandlers(
		s,
		func(e *lolhtml.Element) lolhtml.RewriterDirective {
			e.SetInnerContentAsRaw("LOL-HTML")
			return lolhtml.Continue
		},
		nil,
		func(*lolhtml.TextChunk) lolhtml.RewriterDirective {
			return lolhtml.Continue
		},
	)
	r, _ := rb.Build(
		lolhtml.Config{
			Encoding: "utf-8",
			Memory: &lolhtml.MemorySettings{
				PreallocatedParsingBufferSize: 1024,
				MaxAllowedMemoryUsage:         1<<63 - 1,
			},
			Sink:   func(s string) { fmt.Print(s) },
			Strict: true,
		},
	)
	defer r.Free()
	r.WriteString("Hello, <span>")
	r.WriteString("World</span>!")
	r.End()
}
```

The above program takes chunked input `Hello, <span>World</span>!`, rewrites texts in `span` tags to "LOL-HTML" and prints the result to standard output.

And the result is `Hello, <span>LOL-HTML</span>!` .

## Documentation

Available at [pkg.go.dev](https://pkg.go.dev/github.com/coolspring8/go-lolhtml). (WIP)

## Known Issues

- For now, to use `Rewriter.End()` without causing panic, you will probably need to assign a stub `DocEndHandler` function when calling `AddDocumentContentHandlers()`.

## Other Bindings

- Rust (native), C, JavaScript - [cloudflare/lol-html](https://github.com/cloudflare/lol-html/)
- Lua - [jdesgats/lua-lolhtml](https://github.com/jdesgats/lua-lolhtml/)

## License

BSD 3-Clause "New" or "Revised" License