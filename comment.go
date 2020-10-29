package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"
import "unsafe"

type Comment C.lol_html_comment_t

type CommentHandlerFunc func(*Comment) RewriterDirective

func (c *Comment) Text() string {
	textC := (str)(C.lol_html_comment_text_get((*C.lol_html_comment_t)(c)))
	defer textC.Free()
	return textC.String()
}

func (c *Comment) SetText(text string) error {
	textC := C.CString(text)
	defer C.free(unsafe.Pointer(textC))
	textLen := len(text)
	errCode := C.lol_html_comment_text_set((*C.lol_html_comment_t)(c), textC, C.size_t(textLen))
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (c *Comment) InsertBeforeAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_comment_before((*C.lol_html_comment_t)(c), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (c *Comment) InsertBeforeAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_comment_before((*C.lol_html_comment_t)(c), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (c *Comment) InsertAfterAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_comment_after((*C.lol_html_comment_t)(c), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (c *Comment) InsertAfterAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_comment_after((*C.lol_html_comment_t)(c), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (c *Comment) ReplaceAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_comment_replace((*C.lol_html_comment_t)(c), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (c *Comment) ReplaceAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_comment_replace((*C.lol_html_comment_t)(c), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (c *Comment) Remove() {
	C.lol_html_comment_remove((*C.lol_html_comment_t)(c))
}

func (c *Comment) IsRemoved() bool {
	return (bool)(C.lol_html_comment_is_removed((*C.lol_html_comment_t)(c)))
}
