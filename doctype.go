package lolhtml

/*
#include "lol_html.h"
*/
import "C"

type Doctype C.lol_html_doctype_t

type DoctypeHandlerFunc func(*Doctype) RewriterDirective

func (d *Doctype) Name() string {
	nameC := (*str)(C.lol_html_doctype_name_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return nameC.String()
}

func (d *Doctype) PublicId() string {
	nameC := (*str)(C.lol_html_doctype_public_id_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return nameC.String()
}

func (d *Doctype) SystemId() string {
	nameC := (*str)(C.lol_html_doctype_system_id_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return nameC.String()
}
