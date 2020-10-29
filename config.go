package lolhtml

/*
#include "lol_html.h"
*/
import "C"
import (
	"github.com/mattn/go-pointer"
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

type MemorySettings struct {
	PreallocatedParsingBufferSize int
	MaxAllowedMemoryUsage         int
}

// OutputSink takes each chunked output as a byte slice and processes it.
type OutputSink func([]byte)

type DocumentContentHandler struct {
	DoctypeHandler     DoctypeHandlerFunc
	CommentHandler     CommentHandlerFunc
	TextChunkHandler   TextChunkHandlerFunc
	DocumentEndHandler DocumentEndHandlerFunc
}

type ElementContentHandler struct {
	Selector         string
	ElementHandler   ElementHandlerFunc
	CommentHandler   CommentHandlerFunc
	TextChunkHandler TextChunkHandlerFunc
}

type Handlers struct {
	DocumentContentHandler []DocumentContentHandler
	ElementContentHandler  []ElementContentHandler
}

//export callbackSink
func callbackSink(chunk *C.char, chunkLen C.size_t, userData unsafe.Pointer) {
	c := C.GoBytes(unsafe.Pointer(chunk), C.int(chunkLen))
	cb := pointer.Restore(userData).(OutputSink)
	cb(c)
}

//export callbackDoctype
func callbackDoctype(doctype *Doctype, userData unsafe.Pointer) RewriterDirective {
	cb := pointer.Restore(userData).(DoctypeHandlerFunc)
	return cb(doctype)
}

//export callbackComment
func callbackComment(comment *Comment, userData unsafe.Pointer) RewriterDirective {
	cb := pointer.Restore(userData).(CommentHandlerFunc)
	return cb(comment)
}

//export callbackTextChunk
func callbackTextChunk(textChunk *TextChunk, userData unsafe.Pointer) RewriterDirective {
	cb := pointer.Restore(userData).(TextChunkHandlerFunc)
	return cb(textChunk)
}

//export callbackElement
func callbackElement(element *Element, userData unsafe.Pointer) RewriterDirective {
	cb := pointer.Restore(userData).(ElementHandlerFunc)
	return cb(element)
}

//export callbackDocumentEnd
func callbackDocumentEnd(documentEnd *DocumentEnd, userData unsafe.Pointer) RewriterDirective {
	cb := pointer.Restore(userData).(DocumentEndHandlerFunc)
	return cb(documentEnd)
}
