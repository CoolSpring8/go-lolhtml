package lolhtml

/*
#include "lol_html.h"
*/
import "C"
import (
	"unsafe"
)

// Config defines settings for the rewriter.
type Config struct {
	// defaults to "utf-8".
	Encoding string
	// defaults to PreallocatedParsingBufferSize: 1024, MaxAllowedMemoryUsage: 1<<63 - 1.
	Memory *MemorySettings
	// defaults to func([]byte) {}. In other words, totally discard output.
	Sink OutputSink
	// defaults to true. If true, bail out for security reasons when ambiguous.
	Strict bool
}

func newDefaultConfig() Config {
	return Config{
		Encoding: "utf-8",
		Memory: &MemorySettings{
			PreallocatedParsingBufferSize: 1024,
			MaxAllowedMemoryUsage:         1<<63 - 1,
		},
		Sink:   func([]byte) {},
		Strict: true,
	}
}

// MemorySettings sets the memory limitations for the rewriter.
type MemorySettings struct {
	PreallocatedParsingBufferSize int // defaults to 1024
	MaxAllowedMemoryUsage         int // defaults to 1<<63 -1
}

// OutputSink is a callback function where output is written to. A byte slice is passed each time,
// representing a chunk of output.
//
// Exported for special usages which require each output chunk to be identified and processed
// individually. For most common uses, NewWriter would be more convenient.
type OutputSink func([]byte)

// DocumentContentHandler is a group of handlers that would be applied to the whole HTML document.
type DocumentContentHandler struct {
	DoctypeHandler     DoctypeHandlerFunc
	CommentHandler     CommentHandlerFunc
	TextChunkHandler   TextChunkHandlerFunc
	DocumentEndHandler DocumentEndHandlerFunc
}

// ElementContentHandler is a group of handlers that would be applied to the content matched by
// the given selector.
type ElementContentHandler struct {
	Selector         string
	ElementHandler   ElementHandlerFunc
	CommentHandler   CommentHandlerFunc
	TextChunkHandler TextChunkHandlerFunc
}

// Handlers contain DocumentContentHandlers and ElementContentHandlers. Can contain arbitrary numbers
// of them, including zero (nil slice).
type Handlers struct {
	DocumentContentHandler []DocumentContentHandler
	ElementContentHandler  []ElementContentHandler
}

//export callbackSink
func callbackSink(chunk *C.char, chunkLen C.size_t, userData unsafe.Pointer) {
	c := C.GoBytes(unsafe.Pointer(chunk), C.int(chunkLen))
	cb := restorePointer(userData).(OutputSink)
	cb(c)
}

//export callbackDoctype
func callbackDoctype(doctype *Doctype, userData unsafe.Pointer) RewriterDirective {
	cb := restorePointer(userData).(DoctypeHandlerFunc)
	return cb(doctype)
}

//export callbackComment
func callbackComment(comment *Comment, userData unsafe.Pointer) RewriterDirective {
	cb := restorePointer(userData).(CommentHandlerFunc)
	return cb(comment)
}

//export callbackTextChunk
func callbackTextChunk(textChunk *TextChunk, userData unsafe.Pointer) RewriterDirective {
	cb := restorePointer(userData).(TextChunkHandlerFunc)
	return cb(textChunk)
}

//export callbackElement
func callbackElement(element *Element, userData unsafe.Pointer) RewriterDirective {
	cb := restorePointer(userData).(ElementHandlerFunc)
	return cb(element)
}

//export callbackDocumentEnd
func callbackDocumentEnd(documentEnd *DocumentEnd, userData unsafe.Pointer) RewriterDirective {
	cb := restorePointer(userData).(DocumentEndHandlerFunc)
	return cb(documentEnd)
}
