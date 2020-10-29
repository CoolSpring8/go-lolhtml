package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"
import "unsafe"

type DocumentEnd C.lol_html_doc_end_t
type DocumentEndHandlerFunc func(*DocumentEnd) RewriterDirective

func (d *DocumentEnd) AppendAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_doc_end_append((*C.lol_html_doc_end_t)(d), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (d *DocumentEnd) AppendAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_doc_end_append((*C.lol_html_doc_end_t)(d), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}
