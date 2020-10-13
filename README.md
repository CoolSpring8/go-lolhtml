# lolhtml

[![Go Report Card](https://goreportcard.com/badge/github.com/coolspring8/lolhtml)](https://goreportcard.com/report/github.com/coolspring8/lolhtml)

Go bindings for [cloudflare/lol-html](https://github.com/cloudflare/lol-html/), the *Low Output Latency streaming HTML rewriter/parser with CSS-selector based API.*

**Status:** All abilities provided by C-API implemented, except for customized user data in handlers. The code is at its early stage and the API is therefore subject to change. If you have any ideas on how API can be better structured, feel free to open a PR or an issue.

## Installation

Rust is required to build the lol-html library.

For Linux:

```bash
git clone --recursive https://github.com/coolspring8/lolhtml.git
cargo build --release --manifest-path ./lol-html/c-api/ --target-dir ./
go intall
```

For Windows users, as Rust relies on MSVC toolchain by default, one more step is needed between `cargo build` and `go install`: create a `.a` file from compiled artifacts. This snippet works for me:

```powershell
gendef ./release/lolhtml.dll
dlltool --as-flags=--64 -m i386:x86-64 -k --output-lib ./lolhtml.a --input-def lolhtml.def
cp ./release/lolhtml.dll ./
```

Now let's initialize a project and create `main.go`:

```go
package main

import (
    "fmt"
    "github.com/coolspring8/lolhtml"
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
	r.WriteString("<p>Hello <span>")
	r.WriteString("World</span>!</p>")
	r.End()
}
```

The above program takes chunked input `<p>Hello <span>World</span>!</p>`, rewrites texts in `span` tags to "LOL-HTML" and prints the result to standard output. The result is ``<p>Hello <span>LOL-HTML</span>!</p>`` .

## Documentation

Available at [pkg.go.dev](https://pkg.go.dev/github.com/coolspring8/lolhtml). (WIP)

## Known Issue

- For now, to use `Rewriter.End()` without causing panic, you will probably need to assign a stub `DocEndHandler` function when calling `AddDocumentContentHandlers()`.

## Other Bindings

- Rust (native), C, JavaScript - [cloudflare/lol-html](https://github.com/cloudflare/lol-html/)
- Lua - [jdesgats/lua-lolhtml](https://github.com/jdesgats/lua-lolhtml/)

## License

BSD 3-Clause "New" or "Revised" License