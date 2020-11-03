package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"

// AttributeIterator can be used to iterate over all attributes of an element. The only way to
// get an AttributeIterator is by calling AttributeIterator() on an Element. Note the "range" syntax is not
// applicable here, use AttributeIterator.Next() instead.
type AttributeIterator C.lol_html_attributes_iterator_t

// Attribute represents an HTML element attribute. Obtained by calling Next() on an AttributeIterator.
type Attribute C.lol_html_attribute_t

// Free frees the memory held by the AttributeIterator.
func (ai *AttributeIterator) Free() {
	C.lol_html_attributes_iterator_free((*C.lol_html_attributes_iterator_t)(ai))
}

// Next advances the iterator and returns next attribute.
// Returns nil if the iterator has been exhausted.
func (ai *AttributeIterator) Next() *Attribute {
	return (*Attribute)(C.lol_html_attributes_iterator_next((*C.lol_html_attributes_iterator_t)(ai)))
}

// Name returns the name of the attribute.
func (a *Attribute) Name() string {
	nameC := (str)(C.lol_html_attribute_name_get((*C.lol_html_attribute_t)(a)))
	defer nameC.Free()
	return nameC.String()
}

// Value returns the value of the attribute.
func (a *Attribute) Value() string {
	valueC := (str)(C.lol_html_attribute_value_get((*C.lol_html_attribute_t)(a)))
	defer valueC.Free()
	return valueC.String()
}
