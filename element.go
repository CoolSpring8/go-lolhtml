package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Element represents an HTML element.
type Element C.lol_html_element_t

// ElementHandlerFunc is a callback handler function to do something with an Element.
type ElementHandlerFunc func(*Element) RewriterDirective

// TagName gets the element's tag name.
func (e *Element) TagName() string {
	tagNameC := (str)(C.lol_html_element_tag_name_get((*C.lol_html_element_t)(e)))
	defer tagNameC.Free()
	return tagNameC.String()
}

// SetTagName sets the element's tag name.
func (e *Element) SetTagName(name string) error {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	nameLen := len(name)
	errCode := C.lol_html_element_tag_name_set((*C.lol_html_element_t)(e), nameC, C.size_t(nameLen))
	if errCode == 0 {
		return nil
	}
	return getError()
}

// NamespaceURI gets the element's namespace URI.
func (e *Element) NamespaceURI() string {
	// don't need to be freed
	namespaceURIC := C.lol_html_element_namespace_uri_get((*C.lol_html_element_t)(e))
	return C.GoString(namespaceURIC)
}

// AttributeIterator returns a pointer to an AttributeIterator. Can be used to iterate
// over all attributes of the element.
func (e *Element) AttributeIterator() *AttributeIterator {
	return (*AttributeIterator)(C.lol_html_attributes_iterator_get((*C.lol_html_element_t)(e)))
}

// AttributeValue returns the value of the attribute on this element.
func (e *Element) AttributeValue(name string) (string, error) {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	nameLen := len(name)
	valueC := (*str)(C.lol_html_element_get_attribute((*C.lol_html_element_t)(e), nameC, C.size_t(nameLen)))
	defer valueC.Free()
	// always check error, so not using getError()
	errC := (*str)(C.lol_html_take_last_error())
	defer errC.Free()
	errMsg := errC.String()
	if errMsg != "" {
		return "", errors.New(errMsg)
	}
	return valueC.String(), nil
}

// HasAttribute returns whether the element has the attribute of this name or not.
func (e *Element) HasAttribute(name string) (bool, error) {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	nameLen := len(name)
	codeC := C.lol_html_element_has_attribute((*C.lol_html_element_t)(e), nameC, C.size_t(nameLen))
	if codeC == 1 {
		return true, nil
	} else if codeC == 0 {
		return false, nil
	}
	return false, getError()
}

// SetAttribute updates or creates the attribute with name and value on the element.
func (e *Element) SetAttribute(name string, value string) error {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	nameLen := len(name)
	valueC := C.CString(value)
	defer C.free(unsafe.Pointer(valueC))
	valueLen := len(value)
	errCode := C.lol_html_element_set_attribute(
		(*C.lol_html_element_t)(e),
		nameC,
		C.size_t(nameLen),
		valueC,
		C.size_t(valueLen),
	)
	if errCode == 0 {
		return nil
	}
	return getError()
}

// RemoveAttribute removes the attribute with the name from the element.
func (e *Element) RemoveAttribute(name string) error {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	nameLen := len(name)
	errCode := C.lol_html_element_remove_attribute((*C.lol_html_element_t)(e), nameC, C.size_t(nameLen))
	if errCode == 0 {
		return nil
	}
	return getError()
}

type elementAlter int

const (
	elementInsertBeforeStartTag elementAlter = iota
	elementInsertAfterStartTag
	elementInsertBeforeEndTag
	elementInsertAfterEndTag
	elementSetInnerContent
	elementReplace
)

func (e *Element) alter(content string, alter elementAlter, isHTML bool) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	var errCode C.int
	switch alter {
	case elementInsertBeforeStartTag:
		errCode = C.lol_html_element_before((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), C.bool(isHTML))
	case elementInsertAfterStartTag:
		errCode = C.lol_html_element_prepend((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), C.bool(isHTML))
	case elementInsertBeforeEndTag:
		errCode = C.lol_html_element_append((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), C.bool(isHTML))
	case elementInsertAfterEndTag:
		errCode = C.lol_html_element_after((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), C.bool(isHTML))
	case elementSetInnerContent:
		errCode = C.lol_html_element_set_inner_content((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), C.bool(isHTML))
	case elementReplace:
		errCode = C.lol_html_element_replace((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), C.bool(isHTML))
	default:
		panic("not implemented")
	}
	if errCode == 0 {
		return nil
	}
	return getError()
}

// InsertBeforeStartTagAsText inserts the given content before the element's start tag.
//
// The rewriter will HTML-escape the content before insertion:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (e *Element) InsertBeforeStartTagAsText(content string) error {
	return e.alter(content, elementInsertBeforeStartTag, false)
}

// InsertBeforeStartTagAsHTML inserts the given content before the element's start tag.
// The content is inserted as is.
func (e *Element) InsertBeforeStartTagAsHTML(content string) error {
	return e.alter(content, elementInsertBeforeStartTag, true)
}

// InsertAfterStartTagAsText inserts (prepend) the given content after the element's start tag.
//
// The rewriter will HTML-escape the content before insertion:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (e *Element) InsertAfterStartTagAsText(content string) error {
	return e.alter(content, elementInsertAfterStartTag, false)
}

// InsertAfterStartTagAsHTML inserts (prepend) the given content after the element's start tag.
// The content is inserted as is.
func (e *Element) InsertAfterStartTagAsHTML(content string) error {
	return e.alter(content, elementInsertAfterStartTag, true)
}

// InsertBeforeEndTagAsText inserts (append) the given content after the element's end tag.
//
// The rewriter will HTML-escape the content before insertion:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (e *Element) InsertBeforeEndTagAsText(content string) error {
	return e.alter(content, elementInsertBeforeEndTag, false)
}

// InsertBeforeEndTagAsHTML inserts (append) the given content before the element's end tag.
// The content is inserted as is.
func (e *Element) InsertBeforeEndTagAsHTML(content string) error {
	return e.alter(content, elementInsertBeforeEndTag, true)
}

// InsertAfterEndTagAsText inserts the given content after the element's end tag.
//
// The rewriter will HTML-escape the content before insertion:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (e *Element) InsertAfterEndTagAsText(content string) error {
	return e.alter(content, elementInsertAfterEndTag, false)
}

// InsertAfterEndTagAsHTML inserts the given content after the element's end tag.
// The content is inserted as is.
func (e *Element) InsertAfterEndTagAsHTML(content string) error {
	return e.alter(content, elementInsertAfterEndTag, true)
}

// SetInnerContentAsText overwrites the element's inner content.
//
// The rewriter will HTML-escape the content:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (e *Element) SetInnerContentAsText(content string) error {
	return e.alter(content, elementSetInnerContent, false)
}

// SetInnerContentAsHTML overwrites the element's inner content.
// The content is kept as is.
func (e *Element) SetInnerContentAsHTML(content string) error {
	return e.alter(content, elementSetInnerContent, true)
}

// ReplaceAsText replace the whole element with the supplied content.
//
// The rewriter will HTML-escape the content:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (e *Element) ReplaceAsText(content string) error {
	return e.alter(content, elementReplace, false)
}

// ReplaceAsHTML replace the whole element with the supplied content.
// The content is kept as is.
func (e *Element) ReplaceAsHTML(content string) error {
	return e.alter(content, elementReplace, true)
}

// Remove completely removes the element.
func (e *Element) Remove() {
	C.lol_html_element_remove((*C.lol_html_element_t)(e))
}

// RemoveAndKeepContent removes the element but keeps the inner content.
func (e *Element) RemoveAndKeepContent() {
	C.lol_html_element_remove_and_keep_content((*C.lol_html_element_t)(e))
}

// IsRemoved returns whether the element is removed or not.
func (e *Element) IsRemoved() bool {
	return (bool)(C.lol_html_element_is_removed((*C.lol_html_element_t)(e)))
}
