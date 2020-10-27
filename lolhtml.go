// package lolhtml provides the ability to rewrite or parse HTML on the fly,
// with CSS-selector based API.
// It is a binding for Rust crate lol_html.
// https://github.com/cloudflare/lol-html
package lolhtml

/*
#cgo CFLAGS:-I${SRCDIR}/build/include
#cgo LDFLAGS:-llolhtml
#cgo !windows LDFLAGS:-lm
#cgo linux,amd64 LDFLAGS:-L${SRCDIR}/build/linux-x86_64
#cgo darwin,amd64 LDFLAGS:-L${SRCDIR}/build/macos-x86_64
#cgo windows,amd64 LDFLAGS:-L${SRCDIR}/build/windows-x86_64
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
	"bytes"
	"errors"
	"github.com/mattn/go-pointer"
	"io"
	"unsafe"
)

var ErrCannotGetErrorMessage = errors.New("cannot get error message from underlying lol-html lib")

// RewriterDirective should returned by callback handlers, to inform the rewriter to continue or stop parsing.
type RewriterDirective int

const (
	// Let the normal parsing process continue.
	Continue RewriterDirective = iota

	// Stop the rewriter immediately. Content currently buffered is discarded, and an error is returned.
	Stop
)

// rewriterBuilder is used to build a rewriter.
type rewriterBuilder C.lol_html_rewriter_builder_t

// rewriter represents an actual HTML rewriter.
// rewriterBuilder, rewriter and selector are kept private to simplify public API.
// If you find it useful to use them publicly, please inform me.
type rewriter C.lol_html_rewriter_t

// selector represents a parsed CSS selector.
type selector C.lol_html_selector_t

type Doctype C.lol_html_doctype_t
type DocumentEnd C.lol_html_doc_end_t
type Comment C.lol_html_comment_t
type TextChunk C.lol_html_text_chunk_t
type Element C.lol_html_element_t
// AttributeIterator cannot be iterated by "range" syntax. You should use AttributeIterator.Next() instead.
type AttributeIterator C.lol_html_attributes_iterator_t
type Attribute C.lol_html_attribute_t

type str C.lol_html_str_t
// textChunkContent does not need to be de-allocated manually.
type textChunkContent C.lol_html_text_chunk_content_t

// OutputSink takes each chunked output as a byte slice.
type OutputSink func([]byte)

type DoctypeHandlerFunc func(*Doctype) RewriterDirective
type CommentHandlerFunc func(*Comment) RewriterDirective
type TextChunkHandlerFunc func(*TextChunk) RewriterDirective
type ElementHandlerFunc func(*Element) RewriterDirective
type DocumentEndHandlerFunc func(*DocumentEnd) RewriterDirective

// Config defines settings for the rewriter.
type Config struct {
	// defaults to "utf-8".
	Encoding string
	// defaults to PreallocatedParsingBufferSize: 1024, MaxAllowedMemoryUsage: 1<<63 - 1.
	Memory   *MemorySettings
	// defaults to func([]byte) {}. In other words, totally discard output.
	Sink     OutputSink
	// defaults to true. If true, bail out for security reasons when ambiguous.
	Strict   bool
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

// RewriteString rewrites the given string with the provided Handlers and Config.
func RewriteString(s string, handlers *Handlers, config ...Config) (string, error) {
	var buf bytes.Buffer
	var w *Writer
	var err error
	if config != nil {
		w, err = NewWriter(&buf, handlers, config[0])
	} else {
		w, err = NewWriter(&buf, handlers)
	}
	if err != nil {
		return "", err
	}
	defer w.Free()

	_, err = w.WriteString(s)
	if err != nil {
		return "", err
	}

	err = w.End()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

type Writer struct {
	w io.Writer
	r *rewriter
}

// NewWriter returns a new Writer with Handlers and Config configured, writing to w.
func NewWriter(w io.Writer, handlers *Handlers, config ...Config) (*Writer, error) {
	var c Config
	var sink OutputSink
	if config != nil {
		c = config[0]
		if c.Sink != nil {
			sink = c.Sink
		} else if w == nil {
			sink = func([]byte) {}
		} else {
			sink = func(p []byte) {
				_, _ = w.Write(p)
			}
		}
	} else {
		c = newDefaultConfig()
		if w == nil {
			sink = func([]byte) {}
		} else {
			sink = func(p []byte) {
				_, _ = w.Write(p)
			}
		}
	}

	rb := newRewriterBuilder()
	var selectors []*selector
	if handlers != nil {
		for _, dh := range handlers.DocumentContentHandler {
			rb.AddDocumentContentHandlers(
				dh.DoctypeHandler,
				dh.CommentHandler,
				dh.TextChunkHandler,
				dh.DocumentEndHandler,
			)
		}
		for _, eh := range handlers.ElementContentHandler {
			s, err := newSelector(eh.Selector)
			if err != nil {
				return nil, err
			}
			selectors = append(selectors, s)
			rb.AddElementContentHandlers(
				s,
				eh.ElementHandler,
				eh.CommentHandler,
				eh.TextChunkHandler,
			)
		}
	}
	r, err := rb.Build(sink, c)
	if err != nil {
		return nil, err
	}
	rb.Free()
	for _, s := range selectors {
		s.Free()
	}

	return &Writer{w, r}, nil
}

func (w Writer) Write(p []byte) (n int, err error) {
	return w.r.Write(p)
}

func (w Writer) WriteString(s string) (n int, err error) {
	return w.r.WriteString(s)
}

func (w *Writer) Free() {
	if w != nil {
		w.r.Free()
	}
}

func (w *Writer) End() error {
	return w.r.End()
}

func newRewriterBuilder() *rewriterBuilder {
	return (*rewriterBuilder)(C.lol_html_rewriter_builder_new())
}

func (rb *rewriterBuilder) Free() {
	if rb != nil {
		C.lol_html_rewriter_builder_free((*C.lol_html_rewriter_builder_t)(rb))
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
	doctypeHandlerPointer := pointer.Save(doctypeHandler)
	commentHandlerPointer := pointer.Save(commentHandler)
	textChunkHandlerPointer := pointer.Save(textChunkHandler)
	documentEndHandlerPointer := pointer.Save(documentEndHandler)
	C.lol_html_rewriter_builder_add_document_content_handlers(
		(*C.lol_html_rewriter_builder_t)(rb),
		cCallbackDoctypePointer,
		doctypeHandlerPointer,
		cCallbackCommentPointer,
		commentHandlerPointer,
		cCallbackTextChunkPointer,
		textChunkHandlerPointer,
		cCallbackDocumentEndPointer,
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
	elementHandlerPointer := pointer.Save(elementHandler)
	commentHandlerPointer := pointer.Save(commentHandler)
	textChunkHandlerPointer := pointer.Save(textChunkHandler)
	C.lol_html_rewriter_builder_add_element_content_handlers(
		(*C.lol_html_rewriter_builder_t)(rb),
		(*C.lol_html_selector_t)(selector),
		cCallbackElementPointer,
		elementHandlerPointer,
		cCallbackCommentPointer,
		commentHandlerPointer,
		cCallbackTextChunkPointer,
		textChunkHandlerPointer,
	)
}

func (rb *rewriterBuilder) Build(sink OutputSink, config Config) (*rewriter, error) {
	encodingC := C.CString(config.Encoding)
	defer C.free(unsafe.Pointer(encodingC))
	encodingLen := len(config.Encoding)
	memorySettingsC := C.lol_html_memory_settings_t{
		preallocated_parsing_buffer_size: C.size_t(config.Memory.PreallocatedParsingBufferSize),
		max_allowed_memory_usage:         C.size_t(config.Memory.MaxAllowedMemoryUsage),
	}
	p := pointer.Save(sink)
	r := (*rewriter)(C.lol_html_rewriter_build(
		(*C.lol_html_rewriter_builder_t)(rb),
		encodingC,
		C.size_t(encodingLen),
		memorySettingsC,
		(*[0]byte)(C.callback_sink),
		p,
		C.bool(config.Strict),
	))
	if r != nil {
		return r, nil
	}
	return nil, getError()
}

func (r *rewriter) Write(p []byte) (n int, err error) {
	pLen := len(p)
	// avoid 0-sized array
	if pLen == 0 {
		p = []byte("\x00")
	}
	pC := (*C.char)(unsafe.Pointer(&p[0]))
	errCode := C.lol_html_rewriter_write((*C.lol_html_rewriter_t)(r), pC, C.size_t(pLen))
	if errCode == 0 {
		return pLen, nil
	}
	return 0, getError()
}

func (r *rewriter) WriteString(chunk string) (n int, err error) {
	chunkC := C.CString(chunk)
	defer C.free(unsafe.Pointer(chunkC))
	chunkLen := len(chunk)
	errCode := C.lol_html_rewriter_write((*C.lol_html_rewriter_t)(r), chunkC, C.size_t(chunkLen))
	if errCode == 0 {
		return chunkLen, nil
	}
	return 0, getError()
}

func (r *rewriter) End() error {
	errCode := C.lol_html_rewriter_end((*C.lol_html_rewriter_t)(r))
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (r *rewriter) Free() {
	if r != nil {
		C.lol_html_rewriter_free((*C.lol_html_rewriter_t)(r))
	}
}

func (d *Doctype) Name() string {
	nameC := (*str)(C.lol_html_doctype_name_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return strToGoString(nameC)
}

func (d *Doctype) PublicId() string {
	nameC := (*str)(C.lol_html_doctype_public_id_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return strToGoString(nameC)
}

func (d *Doctype) SystemId() string {
	nameC := (*str)(C.lol_html_doctype_system_id_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return strToGoString(nameC)
}

//func (d* Doctype) SetUserData(){}

//func (d* Doctype) UserData(){}

func (c *Comment) Text() string {
	textC := (str)(C.lol_html_comment_text_get((*C.lol_html_comment_t)(c)))
	defer textC.Free()
	return strToGoString2(textC)
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

func (e *Element) TagName() string {
	tagNameC := (str)(C.lol_html_element_tag_name_get((*C.lol_html_element_t)(e)))
	defer tagNameC.Free()
	return strToGoString2(tagNameC)
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
	errC := (*str)(C.lol_html_take_last_error())
	defer errC.Free()
	errMsg := strToGoString(errC)
	if errMsg != "" {
		return "", errors.New(errMsg)
	}
	return strToGoString(valueC), nil
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

func (ai *AttributeIterator) Free() {
	C.lol_html_attributes_iterator_free((*C.lol_html_attributes_iterator_t)(ai))
}

func (ai *AttributeIterator) Next() *Attribute {
	return (*Attribute)(C.lol_html_attributes_iterator_next((*C.lol_html_attributes_iterator_t)(ai)))
}

func (a *Attribute) Name() string {
	nameC := (str)(C.lol_html_attribute_name_get((*C.lol_html_attribute_t)(a)))
	defer nameC.Free()
	return strToGoString2(nameC)
}

func (a *Attribute) Value() string {
	valueC := (str)(C.lol_html_attribute_value_get((*C.lol_html_attribute_t)(a)))
	defer valueC.Free()
	return strToGoString2(valueC)
}

func (d *DocumentEnd) AppendAsText(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_doc_end_append((*C.lol_html_doc_end_t)(d), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (d *DocumentEnd) AppendAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_doc_end_append((*C.lol_html_doc_end_t)(d), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func newSelector(cssSelector string) (*selector, error) {
	selectorC := C.CString(cssSelector)
	defer C.free(unsafe.Pointer(selectorC))
	selectorLen := len(cssSelector)
	s := (*selector)(C.lol_html_selector_parse(selectorC, C.size_t(selectorLen)))
	if s != nil {
		return s, nil
	}
	return nil, getError()
}

func (s *selector) Free() {
	if s != nil {
		C.lol_html_selector_free((*C.lol_html_selector_t)(s))
	}
}

func (s *str) Free() {
	if s != nil {
		C.lol_html_str_free(*(*C.lol_html_str_t)(s))
	}
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

// strToGoString is a helper function that translates the underlying-library-defined lol_html_str_t data to Go string.
// It is the caller's responsibility to arrange for lol_html_str_t to be freed,
// by calling str.Free() or lol_html_str_free().
// Potential issue: lol_html_str_t->len from size_t (uint) to int (int32) on 32-bit machines?
// TODO: rename to String() to implement Stringer?
func strToGoString(s *str) string {
	if s == nil {
		return ""
	}
	return C.GoStringN(s.data, C.int(s.len))
}

// strToGoString2 is similar to strToGoString, except for the function argument.
func strToGoString2(s str) string {
	var nullStr str
	if s == nullStr {
		return ""
	}
	return C.GoStringN(s.data, C.int(s.len))
}

func textChunkContentToGoString(s textChunkContent) string {
	var nullTextChunkContent textChunkContent
	if s == nullTextChunkContent {
		return ""
	}
	return C.GoStringN(s.data, C.int(s.len))
}

// getError is a helper function that gets error message for the last function call.
// You should make sure there is an error when calling this, or the function interprets
// the NULL error message obtained as ErrCannotGetErrorMessage.
func getError() error {
	errC := (*str)(C.lol_html_take_last_error())
	defer errC.Free()
	if errMsg := strToGoString(errC); errMsg != "" {
		return errors.New(errMsg)
	}
	return ErrCannotGetErrorMessage
}
