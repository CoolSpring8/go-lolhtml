package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"
import "unsafe"

// DocumentEnd represents the end of the document.
type DocumentEnd C.lol_html_doc_end_t

// DocumentEndHandlerFunc is a callback handler function to do something with a DocumentEnd.
type DocumentEndHandlerFunc func(*DocumentEnd) RewriterDirective

// AppendAsText appends the given content at the end of the document.
//
// The rewriter will HTML-escape the content before appending:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
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

// AppendAsHTML appends the given content at the end of the document.
// The content is appended as is.
func (d *DocumentEnd) AppendAsHTML(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_doc_end_append((*C.lol_html_doc_end_t)(d), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}
