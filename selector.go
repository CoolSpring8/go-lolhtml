package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"
import "unsafe"

// selector represents a parsed CSS selector.
type selector C.lol_html_selector_t

func newSelector(cssSelector string) (*selector, error) {
	selectorC := C.CString(cssSelector)
	defer C.free(unsafe.Pointer(selectorC))
	selectorLen := len(cssSelector)
	s := (*selector)(C.lol_html_selector_parse(selectorC, C.size_t(selectorLen)))
	if s != nil {
		return s, nil
	}
	return nil, getError()
}

func (s *selector) Free() {
	if s != nil {
		C.lol_html_selector_free((*C.lol_html_selector_t)(s))
	}
}
