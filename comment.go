package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
*/
import "C"
import "unsafe"

// Comment represents an HTML comment.
type Comment C.lol_html_comment_t

// CommentHandlerFunc is a callback handler function to do something with a Comment.
// Expected to return a RewriterDirective as instruction to continue or stop.
type CommentHandlerFunc func(*Comment) RewriterDirective

// Text returns the comment's text.
func (c *Comment) Text() string {
	textC := (str)(C.lol_html_comment_text_get((*C.lol_html_comment_t)(c)))
	defer textC.Free()
	return textC.String()
}

// SetText sets the comment's text and returns an error if there is one.
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

type commentAlter int

const (
	commentInsertBefore commentAlter = iota
	commentInsertAfter
	commentReplace
)

func (c *Comment) alter(content string, alter commentAlter, isHTML bool) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	var errCode C.int
	switch alter {
	case commentInsertBefore:
		errCode = C.lol_html_comment_before((*C.lol_html_comment_t)(c), contentC, C.size_t(contentLen), C.bool(isHTML))
	case commentInsertAfter:
		errCode = C.lol_html_comment_after((*C.lol_html_comment_t)(c), contentC, C.size_t(contentLen), C.bool(isHTML))
	case commentReplace:
		errCode = C.lol_html_comment_replace((*C.lol_html_comment_t)(c), contentC, C.size_t(contentLen), C.bool(isHTML))
	default:
		panic("not implemented")
	}
	if errCode == 0 {
		return nil
	}
	return getError()
}

// InsertBeforeAsText inserts the given content before the comment.
//
// The rewriter will HTML-escape the content before insertion:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (c *Comment) InsertBeforeAsText(content string) error {
	return c.alter(content, commentInsertAfter, false)
}

// InsertBeforeAsHTML inserts the given content before the comment.
// The content is inserted as is.
func (c *Comment) InsertBeforeAsHTML(content string) error {
	return c.alter(content, commentInsertBefore, true)
}

// InsertAfterAsText inserts the given content before the comment.
//
// The rewriter will HTML-escape the content before insertion:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (c *Comment) InsertAfterAsText(content string) error {
	return c.alter(content, commentInsertAfter, false)
}

// InsertAfterAsHTML inserts the given content before the comment.
// The content is inserted as is.
func (c *Comment) InsertAfterAsHTML(content string) error {
	return c.alter(content, commentInsertAfter, true)
}

// ReplaceAsText replace the comment with the supplied content.
//
// The rewriter will HTML-escape the content:
//
// `<` will be replaced with `&lt;`
//
// `>` will be replaced with `&gt;`
//
// `&` will be replaced with `&amp;`
func (c *Comment) ReplaceAsText(content string) error {
	return c.alter(content, commentReplace, false)
}

// ReplaceAsHTML replace the comment with the supplied content.
// The content is kept as is.
func (c *Comment) ReplaceAsHTML(content string) error {
	return c.alter(content, commentReplace, true)
}

// Remove removes the comment.
func (c *Comment) Remove() {
	C.lol_html_comment_remove((*C.lol_html_comment_t)(c))
}

// IsRemoved returns whether the comment is removed or not.
func (c *Comment) IsRemoved() bool {
	return (bool)(C.lol_html_comment_is_removed((*C.lol_html_comment_t)(c)))
}
