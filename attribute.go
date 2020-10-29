package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"

// AttributeIterator cannot be iterated by "range" syntax. You should use AttributeIterator.Next() instead.
type AttributeIterator C.lol_html_attributes_iterator_t
type Attribute C.lol_html_attribute_t

func (ai *AttributeIterator) Free() {
	C.lol_html_attributes_iterator_free((*C.lol_html_attributes_iterator_t)(ai))
}

func (ai *AttributeIterator) Next() *Attribute {
	return (*Attribute)(C.lol_html_attributes_iterator_next((*C.lol_html_attributes_iterator_t)(ai)))
}

func (a *Attribute) Name() string {
	nameC := (str)(C.lol_html_attribute_name_get((*C.lol_html_attribute_t)(a)))
	defer nameC.Free()
	return nameC.String()
}

func (a *Attribute) Value() string {
	valueC := (str)(C.lol_html_attribute_value_get((*C.lol_html_attribute_t)(a)))
	defer valueC.Free()
	return valueC.String()
}
