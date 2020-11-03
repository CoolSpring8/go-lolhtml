package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"
import "unsafe"

// TextChunk represents a text chunk.
type TextChunk C.lol_html_text_chunk_t

// TextChunkHandlerFunc is a callback handler function to do something with a TextChunk.
type TextChunkHandlerFunc func(*TextChunk) RewriterDirective

// Content returns the text chunk's content.
func (t *TextChunk) Content() string {
	text := (textChunkContent)(C.lol_html_text_chunk_content_get((*C.lol_html_text_chunk_t)(t)))
	return text.String()
}

// IsLastInTextNode returns whether the text chunk is the last in the text node.
func (t *TextChunk) IsLastInTextNode() bool {
	return (bool)(C.lol_html_text_chunk_is_last_in_text_node((*C.lol_html_text_chunk_t)(t)))
}

type textChunkAlter int

const (
	textChunkInsertBefore textChunkAlter = iota
	textChunkInsertAfter
	textChunkReplace
)

func (t *TextChunk) alter(content string, alter textChunkAlter, isHTML bool) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	var errCode C.int
	switch alter {
	case textChunkInsertBefore:
		errCode = C.lol_html_text_chunk_before((*C.lol_html_text_chunk_t)(t), contentC, C.size_t(contentLen), C.bool(isHTML))
	case textChunkInsertAfter:
		errCode = C.lol_html_text_chunk_after((*C.lol_html_text_chunk_t)(t), contentC, C.size_t(contentLen), C.bool(isHTML))
	case textChunkReplace:
		errCode = C.lol_html_text_chunk_replace((*C.lol_html_text_chunk_t)(t), contentC, C.size_t(contentLen), C.bool(isHTML))
	default:
		panic("not implemented")
	}
	if errCode == 0 {
		return nil
	}
	return getError()
}

// InsertBeforeAsText inserts the given content before the text chunk.
//
// The rewriter will HTML-escape the content before insertion:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (t *TextChunk) InsertBeforeAsText(content string) error {
	return t.alter(content, textChunkInsertBefore, false)
}

// InsertBeforeAsHTML inserts the given content before the text chunk.
// The content is inserted as is.
func (t *TextChunk) InsertBeforeAsHTML(content string) error {
	return t.alter(content, textChunkInsertBefore, true)
}

// InsertAfterAsText inserts the given content after the text chunk.
//
// The rewriter will HTML-escape the content before insertion:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (t *TextChunk) InsertAfterAsText(content string) error {
	return t.alter(content, textChunkInsertAfter, false)
}

// InsertAfterAsHTML inserts the given content after the text chunk.
// The content is inserted as is.
func (t *TextChunk) InsertAfterAsHTML(content string) error {
	return t.alter(content, textChunkInsertAfter, true)
}

// ReplaceAsText replace the text chunk with the supplied content.
//
// The rewriter will HTML-escape the content:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (t *TextChunk) ReplaceAsText(content string) error {
	return t.alter(content, textChunkReplace, false)
}

// ReplaceAsHTML replace the text chunk with the supplied content.
// The content is kept as is.
func (t *TextChunk) ReplaceAsHTML(content string) error {
	return t.alter(content, textChunkReplace, true)
}

// Remove removes the text chunk.
func (t *TextChunk) Remove() {
	C.lol_html_text_chunk_remove((*C.lol_html_text_chunk_t)(t))
}

// IsRemoved returns whether the text chunk is removed or not.
func (t *TextChunk) IsRemoved() bool {
	return (bool)(C.lol_html_text_chunk_is_removed((*C.lol_html_text_chunk_t)(t)))
}
