package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"
import "unsafe"

type TextChunk C.lol_html_text_chunk_t

type TextChunkHandlerFunc func(*TextChunk) RewriterDirective

func (t *TextChunk) Content() string {
	text := (textChunkContent)(C.lol_html_text_chunk_content_get((*C.lol_html_text_chunk_t)(t)))
	return textChunkContentToGoString(text)
}

func (t *TextChunk) IsLastInTextNode() bool {
	return (bool)(C.lol_html_text_chunk_is_last_in_text_node((*C.lol_html_text_chunk_t)(t)))
}

func (t *TextChunk) InsertBeforeAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_text_chunk_before((*C.lol_html_text_chunk_t)(t), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (t *TextChunk) InsertBeforeAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_text_chunk_before((*C.lol_html_text_chunk_t)(t), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (t *TextChunk) InsertAfterAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_text_chunk_after((*C.lol_html_text_chunk_t)(t), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (t *TextChunk) InsertAfterAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_text_chunk_after((*C.lol_html_text_chunk_t)(t), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (t *TextChunk) ReplaceAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_text_chunk_replace((*C.lol_html_text_chunk_t)(t), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (t *TextChunk) ReplaceAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_text_chunk_replace((*C.lol_html_text_chunk_t)(t), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (t *TextChunk) Remove() {
	C.lol_html_text_chunk_remove((*C.lol_html_text_chunk_t)(t))
}

func (t *TextChunk) IsRemoved() bool {
	return (bool)(C.lol_html_text_chunk_is_removed((*C.lol_html_text_chunk_t)(t)))
}
