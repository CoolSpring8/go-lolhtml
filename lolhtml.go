// package lolhtml provides the ability to rewrite or parse HTML on the fly,
// with a CSS-selector based API.
// It is a binding for the Rust crate lol_html.
// https://github.com/cloudflare/lol-html
package lolhtml

/*
#cgo CFLAGS:-I${SRCDIR}/build/include
#cgo LDFLAGS:-llolhtml
#cgo !windows LDFLAGS:-lm
#cgo linux,amd64 LDFLAGS:-L${SRCDIR}/build/linux-x86_64
#cgo darwin,amd64 LDFLAGS:-L${SRCDIR}/build/macos-x86_64
#cgo windows,amd64 LDFLAGS:-L${SRCDIR}/build/windows-x86_64
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"
