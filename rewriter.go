package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"
import (
	"unsafe"
)

// rewriter represents an actual HTML rewriter.
// rewriterBuilder, rewriter and selector are kept private to simplify public API.
// If you find it useful to use them publicly, please inform me.
type rewriter struct {
	rewriter *C.lol_html_rewriter_t
	pointers []unsafe.Pointer
	// TODO: unrecoverable bool
}

func (r *rewriter) Write(p []byte) (n int, err error) {
	pLen := len(p)
	// avoid 0-sized array
	if pLen == 0 {
		p = []byte("\x00")
	}
	pC := (*C.char)(unsafe.Pointer(&p[0]))
	errCode := C.lol_html_rewriter_write(r.rewriter, pC, C.size_t(pLen))
	if errCode == 0 {
		return pLen, nil
	}
	return 0, getError()
}

func (r *rewriter) WriteString(chunk string) (n int, err error) {
	chunkC := C.CString(chunk)
	defer C.free(unsafe.Pointer(chunkC))
	chunkLen := len(chunk)
	errCode := C.lol_html_rewriter_write(r.rewriter, chunkC, C.size_t(chunkLen))
	if errCode == 0 {
		return chunkLen, nil
	}
	return 0, getError()
}

func (r *rewriter) End() error {
	errCode := C.lol_html_rewriter_end(r.rewriter)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (r *rewriter) Free() {
	if r != nil {
		C.lol_html_rewriter_free(r.rewriter)
		unrefPointers(r.pointers)
	}
}
