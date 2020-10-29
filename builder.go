package lolhtml

/*
#include <stdlib.h>
#include "lol_html.h"
extern void callback_sink(const char *chunk, size_t chunk_len, void *user_data);
extern lol_html_rewriter_directive_t callback_doctype(lol_html_doctype_t *doctype, void *user_data);
extern lol_html_rewriter_directive_t callback_comment(lol_html_comment_t *comment, void *user_data);
extern lol_html_rewriter_directive_t callback_text_chunk(lol_html_text_chunk_t *text_chunk, void *user_data);
extern lol_html_rewriter_directive_t callback_element(lol_html_element_t *element, void *user_data);
extern lol_html_rewriter_directive_t callback_doc_end(lol_html_doc_end_t *doc_end, void *user_data);
*/
import "C"
import (
	"unsafe"
)

// rewriterBuilder is used to build a rewriter.
type rewriterBuilder struct {
	rb       *C.lol_html_rewriter_builder_t
	pointers []unsafe.Pointer
	built    bool // this builder has built at least one writer
}

func newRewriterBuilder() *rewriterBuilder {
	return &rewriterBuilder{rb: C.lol_html_rewriter_builder_new(), pointers: nil, built: false}
}

func (rb *rewriterBuilder) Free() {
	if rb != nil {
		C.lol_html_rewriter_builder_free(rb.rb)
		if !rb.built {
			unrefPointers(rb.pointers)
		}
	}
}

func (rb *rewriterBuilder) AddDocumentContentHandlers(
	doctypeHandler DoctypeHandlerFunc,
	commentHandler CommentHandlerFunc,
	textChunkHandler TextChunkHandlerFunc,
	documentEndHandler DocumentEndHandlerFunc,
) {
	var cCallbackDoctypePointer, cCallbackCommentPointer, cCallbackTextChunkPointer, cCallbackDocumentEndPointer *[0]byte
	if doctypeHandler != nil {
		cCallbackDoctypePointer = (*[0]byte)(C.callback_doctype)
	}
	if commentHandler != nil {
		cCallbackCommentPointer = (*[0]byte)(C.callback_comment)
	}
	if textChunkHandler != nil {
		cCallbackTextChunkPointer = (*[0]byte)(C.callback_text_chunk)
	}
	if documentEndHandler != nil {
		cCallbackDocumentEndPointer = (*[0]byte)(C.callback_doc_end)
	}
	doctypeHandlerPointer := savePointer(doctypeHandler)
	commentHandlerPointer := savePointer(commentHandler)
	textChunkHandlerPointer := savePointer(textChunkHandler)
	documentEndHandlerPointer := savePointer(documentEndHandler)
	C.lol_html_rewriter_builder_add_document_content_handlers(
		rb.rb,
		cCallbackDoctypePointer,
		doctypeHandlerPointer,
		cCallbackCommentPointer,
		commentHandlerPointer,
		cCallbackTextChunkPointer,
		textChunkHandlerPointer,
		cCallbackDocumentEndPointer,
		documentEndHandlerPointer,
	)
	rb.pointers = append(
		rb.pointers,
		doctypeHandlerPointer,
		commentHandlerPointer,
		textChunkHandlerPointer,
		documentEndHandlerPointer,
	)
}

func (rb *rewriterBuilder) AddElementContentHandlers(
	selector *selector,
	elementHandler ElementHandlerFunc,
	commentHandler CommentHandlerFunc,
	textChunkHandler TextChunkHandlerFunc,
) {
	var cCallbackElementPointer, cCallbackCommentPointer, cCallbackTextChunkPointer *[0]byte
	if elementHandler != nil {
		cCallbackElementPointer = (*[0]byte)(C.callback_element)
	}
	if commentHandler != nil {
		cCallbackCommentPointer = (*[0]byte)(C.callback_comment)
	}
	if textChunkHandler != nil {
		cCallbackTextChunkPointer = (*[0]byte)(C.callback_text_chunk)
	}
	elementHandlerPointer := savePointer(elementHandler)
	commentHandlerPointer := savePointer(commentHandler)
	textChunkHandlerPointer := savePointer(textChunkHandler)
	C.lol_html_rewriter_builder_add_element_content_handlers(
		rb.rb,
		(*C.lol_html_selector_t)(selector),
		cCallbackElementPointer,
		elementHandlerPointer,
		cCallbackCommentPointer,
		commentHandlerPointer,
		cCallbackTextChunkPointer,
		textChunkHandlerPointer,
	)
	rb.pointers = append(rb.pointers, elementHandlerPointer, commentHandlerPointer, textChunkHandlerPointer)
}

func (rb *rewriterBuilder) Build(sink OutputSink, config Config) (*rewriter, error) {
	encodingC := C.CString(config.Encoding)
	defer C.free(unsafe.Pointer(encodingC))
	encodingLen := len(config.Encoding)
	memorySettingsC := C.lol_html_memory_settings_t{
		preallocated_parsing_buffer_size: C.size_t(config.Memory.PreallocatedParsingBufferSize),
		max_allowed_memory_usage:         C.size_t(config.Memory.MaxAllowedMemoryUsage),
	}
	p := savePointer(sink)
	r := C.lol_html_rewriter_build(
		rb.rb,
		encodingC,
		C.size_t(encodingLen),
		memorySettingsC,
		(*[0]byte)(C.callback_sink),
		p,
		C.bool(config.Strict),
	)
	if r != nil {
		rb.built = true
		return &rewriter{rw: r, pointers: rb.pointers}, nil
	}
	return nil, getError()
}
