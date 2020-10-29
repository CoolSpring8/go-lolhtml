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

type Element C.lol_html_element_t

type ElementHandlerFunc func(*Element) RewriterDirective

func (e *Element) TagName() string {
	tagNameC := (str)(C.lol_html_element_tag_name_get((*C.lol_html_element_t)(e)))
	defer tagNameC.Free()
	return tagNameC.String()
}

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

func (e *Element) NamespaceUri() string {
	namespaceUriC := C.lol_html_element_namespace_uri_get((*C.lol_html_element_t)(e))
	return C.GoString(namespaceUriC)
}

func (e *Element) AttributeIterator() *AttributeIterator {
	return (*AttributeIterator)(C.lol_html_attributes_iterator_get((*C.lol_html_element_t)(e)))
}

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

func (e *Element) InsertBeforeStartTagAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_before((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) InsertBeforeStartTagAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_before((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) InsertAfterStartTagAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_prepend((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) InsertAfterStartTagAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_prepend((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) InsertBeforeEndTagAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_append((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) InsertBeforeEndTagAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_append((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) InsertAfterEndTagAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_after((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) InsertAfterEndTagAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_after((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) SetInnerContentAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_set_inner_content((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) SetInnerContentAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_set_inner_content((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) ReplaceAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_replace((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) ReplaceAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_element_replace((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) Remove() {
	C.lol_html_element_remove((*C.lol_html_element_t)(e))
}

func (e *Element) RemoveAndKeepContent() {
	C.lol_html_element_remove_and_keep_content((*C.lol_html_element_t)(e))
}

func (e *Element) IsRemoved() bool {
	return (bool)(C.lol_html_element_is_removed((*C.lol_html_element_t)(e)))
}
