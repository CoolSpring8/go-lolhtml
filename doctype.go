package lolhtml

/*
#include "lol_html.h"
*/
import "C"

// Doctype represents the document's doctype.
type Doctype C.lol_html_doctype_t

// DoctypeHandlerFunc is a callback handler function to do something with a Comment.
type DoctypeHandlerFunc func(*Doctype) RewriterDirective

// Name returns doctype name.
func (d *Doctype) Name() string {
	nameC := (*str)(C.lol_html_doctype_name_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return nameC.String()
}

// PublicID returns doctype public ID.
func (d *Doctype) PublicID() string {
	nameC := (*str)(C.lol_html_doctype_public_id_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return nameC.String()
}

// SystemID returns doctype system ID.
func (d *Doctype) SystemID() string {
	nameC := (*str)(C.lol_html_doctype_system_id_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return nameC.String()
}
