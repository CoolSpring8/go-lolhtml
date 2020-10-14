package lolhtml

/*
#cgo CFLAGS: -I/usr/local/include/lolhtml
#cgo LDFLAGS: /usr/local/lib/lolhtml/liblolhtml.so
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
	"errors"
	"github.com/mattn/go-pointer"
	"unsafe"
)

var ErrCannotGetErrorMessage = errors.New("cannot get error message from underlying lol-html lib")

// RewriterDirective as declared in include/lol_html.h:84
type RewriterDirective int

// RewriterDirective enumeration from include/lol_html.h:84
const (
	Continue RewriterDirective = iota
	Stop
)

// RewriterBuilder as declared in include/lol_html.h:22
type RewriterBuilder C.lol_html_rewriter_builder_t

// Rewriter as declared in include/lol_html.h:23
type Rewriter C.lol_html_rewriter_t

// Doctype as declared in include/lol_html.h:24
type Doctype C.lol_html_doctype_t

// DocEnd as declared in include/lol_html.h:25
type DocEnd C.lol_html_doc_end_t

// Comment as declared in include/lol_html.h:26
type Comment C.lol_html_comment_t

// TextChunk as declared in include/lol_html.h:27
type TextChunk C.lol_html_text_chunk_t

// Element as declared in include/lol_html.h:28
type Element C.lol_html_element_t

// AttributeIterator as declared in include/lol_html.h:29
type AttributeIterator C.lol_html_attributes_iterator_t

// Attribute as declared in include/lol_html.h:30
type Attribute C.lol_html_attribute_t

// Selector as declared in include/lol_html.h:31
type Selector C.lol_html_selector_t

// str as declared in include/lol_html.h:45
type str C.lol_html_str_t

// TextChunkContent as declared in include/lol_html.h:60
type TextChunkContent C.lol_html_text_chunk_content_t

type OutputSink func(string)

// DoctypeHandler type as declared in include/lol_html.h:86
type DoctypeHandler func(*Doctype) RewriterDirective

// CommentHandler type as declared in include/lol_html.h:91
type CommentHandler func(*Comment) RewriterDirective

// TextChunkHandler type as declared in include/lol_html.h:96
type TextChunkHandler func(*TextChunk) RewriterDirective

// ElementHandler type as declared in include/lol_html.h:101
type ElementHandler func(*Element) RewriterDirective

// DocEndHandler type as declared in include/lol_html.h:106
type DocEndHandler func(*DocEnd) RewriterDirective

type Config struct {
	Encoding string
	Memory   *MemorySettings
	Sink     OutputSink
	//UserData interface{}
	Strict bool
}

func NewDefaultConfig() Config {
	return Config{
		Encoding: "utf-8",
		Memory: &MemorySettings{
			PreallocatedParsingBufferSize: 1024,
			MaxAllowedMemoryUsage:         1<<63 - 1,
		},
		Sink:   func(string) {},
		Strict: true,
	}
}

type MemorySettings struct {
	PreallocatedParsingBufferSize int
	MaxAllowedMemoryUsage         int
}

func NewRewriterBuilder() *RewriterBuilder {
	return (*RewriterBuilder)(C.lol_html_rewriter_builder_new())
}

func (rb *RewriterBuilder) Free() {
	if rb != nil {
		C.lol_html_rewriter_builder_free((*C.lol_html_rewriter_builder_t)(rb))
	}
}

// TODO: BUG? For now, to use *Rewriter.End() without causing panic, you will probably need to assign
// a stub handler function to it.
func (rb *RewriterBuilder) AddDocumentContentHandlers(
	doctypeHandler DoctypeHandler,
	commentHandler CommentHandler,
	textChunkHandler TextChunkHandler,
	docEndHandler DocEndHandler,
) {
	doctypeHandlerPointer := pointer.Save(doctypeHandler)
	commentHandlerPointer := pointer.Save(commentHandler)
	textChunkHandlerPointer := pointer.Save(textChunkHandler)
	docEndHandlerPointer := pointer.Save(docEndHandler)
	C.lol_html_rewriter_builder_add_document_content_handlers(
		(*C.lol_html_rewriter_builder_t)(rb),
		(*[0]byte)(C.callback_doctype),
		doctypeHandlerPointer,
		(*[0]byte)(C.callback_comment),
		commentHandlerPointer,
		(*[0]byte)(C.callback_text_chunk),
		textChunkHandlerPointer,
		(*[0]byte)(C.callback_doc_end),
		docEndHandlerPointer,
	)
}

func (rb *RewriterBuilder) AddElementContentHandlers(
	selector *Selector,
	elementHandler ElementHandler,
	commentHandler CommentHandler,
	textChunkHandler TextChunkHandler,
) {
	commentHandlerPointer := pointer.Save(commentHandler)
	elementHandlerPointer := pointer.Save(elementHandler)
	textChunkHandlerPointer := pointer.Save(textChunkHandler)
	C.lol_html_rewriter_builder_add_element_content_handlers(
		(*C.lol_html_rewriter_builder_t)(rb),
		(*C.lol_html_selector_t)(selector),
		(*[0]byte)(C.callback_element),
		elementHandlerPointer,
		(*[0]byte)(C.callback_comment),
		commentHandlerPointer,
		(*[0]byte)(C.callback_text_chunk),
		textChunkHandlerPointer,
	)
}

func (rb *RewriterBuilder) Build(config Config) (*Rewriter, error) {
	encodingC := C.CString(config.Encoding)
	defer C.free(unsafe.Pointer(encodingC))
	encodingLen := len(config.Encoding)
	memorySettingsC := C.lol_html_memory_settings_t{
		preallocated_parsing_buffer_size: C.size_t(config.Memory.PreallocatedParsingBufferSize),
		max_allowed_memory_usage:         C.size_t(config.Memory.MaxAllowedMemoryUsage),
	}
	p := pointer.Save(config.Sink)
	r := (*Rewriter)(C.lol_html_rewriter_build(
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

//func (r *Rewriter) Write(b [] byte) error {}

func (r *Rewriter) WriteString(chunk string) error {
	chunkC := C.CString(chunk)
	defer C.free(unsafe.Pointer(chunkC))
	chunkLen := len(chunk)
	errCode := C.lol_html_rewriter_write((*C.lol_html_rewriter_t)(r), chunkC, C.size_t(chunkLen))
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (r *Rewriter) End() error {
	errCode := C.lol_html_rewriter_end((*C.lol_html_rewriter_t)(r))
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (r *Rewriter) Free() {
	if r != nil {
		C.lol_html_rewriter_free((*C.lol_html_rewriter_t)(r))
	}
}

func (d *Doctype) GetName() string {
	nameC := (*str)(C.lol_html_doctype_name_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return strToGoString(nameC)
}

func (d *Doctype) GetPublicId() string {
	nameC := (*str)(C.lol_html_doctype_public_id_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return strToGoString(nameC)
}

func (d *Doctype) GetSystemId() string {
	nameC := (*str)(C.lol_html_doctype_system_id_get((*C.lol_html_doctype_t)(d)))
	defer nameC.Free()
	return strToGoString(nameC)
}

//func (d* Doctype) SetUserData(){}

//func (d* Doctype) GetUserData(){}

func (c *Comment) GetText() string {
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

func (c *Comment) InsertBeforeAsRaw(content string) error {
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

func (c *Comment) InsertAfterAsRaw(content string) error {
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

func (c *Comment) ReplaceAsRaw(content string) error {
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

func (t *TextChunk) GetContent() string {
	text := (TextChunkContent)(C.lol_html_text_chunk_content_get((*C.lol_html_text_chunk_t)(t)))
	return textChunkContentToGoString(text)
}

func (t *TextChunk) IsLastInTextNode() bool {
	return (bool)(C.lol_html_text_chunk_is_last_in_text_node((*C.lol_html_text_chunk_t)(t)))
}

func (t *TextChunk) InsertBeforeAsRaw(content string) error {
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

func (t *TextChunk) InsertAfterAsRaw(content string) error {
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

func (t *TextChunk) ReplaceAsRaw(content string) error {
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

func (e *Element) GetTagName() string {
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

func (e *Element) GetNamespaceUri() string {
	namespaceUriC := C.lol_html_element_namespace_uri_get((*C.lol_html_element_t)(e))
	return C.GoString(namespaceUriC)
}

func (e *Element) GetAttributeIterator() *AttributeIterator {
	return (*AttributeIterator)(C.lol_html_attributes_iterator_get((*C.lol_html_element_t)(e)))
}

func (e *Element) GetAttributeValue(name string) (string, error) {
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

func (e *Element) InsertBeforeStartTagAsRaw(content string) error {
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
	errCode := C.lol_html_element_prepend((*C.lol_html_element_t)(e), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (e *Element) InsertAfterStartTagAsRaw(content string) error {
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

func (e *Element) InsertBeforeEndTagAsRaw(content string) error {
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

func (e *Element) InsertAfterEndTagAsRaw(content string) error {
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

func (e *Element) SetInnerContentAsRaw(content string) error {
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

func (e *Element) ReplaceAsRaw(content string) error {
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

func (a *Attribute) GetName() string {
	nameC := (str)(C.lol_html_attribute_name_get((*C.lol_html_attribute_t)(a)))
	defer nameC.Free()
	return strToGoString2(nameC)
}

func (a *Attribute) GetValue() string {
	valueC := (str)(C.lol_html_attribute_value_get((*C.lol_html_attribute_t)(a)))
	defer valueC.Free()
	return strToGoString2(valueC)
}

func (d *DocEnd) AppendAsRaw(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_doc_end_append((*C.lol_html_doc_end_t)(d), contentC, C.size_t(contentLen), false)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func (d *DocEnd) AppendAsHtml(content string) error {
	contentC := C.CString(content)
	defer C.free(unsafe.Pointer(contentC))
	contentLen := len(content)
	errCode := C.lol_html_doc_end_append((*C.lol_html_doc_end_t)(d), contentC, C.size_t(contentLen), true)
	if errCode == 0 {
		return nil
	}
	return getError()
}

func NewSelector(selector string) (*Selector, error) {
	selectorC := C.CString(selector)
	defer C.free(unsafe.Pointer(selectorC))
	selectorLen := len(selector)
	s := (*Selector)(C.lol_html_selector_parse(selectorC, C.size_t(selectorLen)))
	if s != nil {
		return s, nil
	}
	return nil, getError()
}

func (s *Selector) Free() {
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
	c := C.GoStringN(chunk, C.int(chunkLen))
	cb := pointer.Restore(userData).(OutputSink)
	cb(c)
}

//export callbackDoctype
func callbackDoctype(doctype *Doctype, userData unsafe.Pointer) RewriterDirective {
	cb := pointer.Restore(userData).(DoctypeHandler)
	return cb(doctype)
}

//export callbackComment
func callbackComment(comment *Comment, userData unsafe.Pointer) RewriterDirective {
	cb := pointer.Restore(userData).(CommentHandler)
	return cb(comment)
}

//export callbackTextChunk
func callbackTextChunk(textChunk *TextChunk, userData unsafe.Pointer) RewriterDirective {
	cb := pointer.Restore(userData).(TextChunkHandler)
	return cb(textChunk)
}

//export callbackElement
func callbackElement(element *Element, userData unsafe.Pointer) RewriterDirective {
	cb := pointer.Restore(userData).(ElementHandler)
	return cb(element)
}

//export callbackDocEnd
func callbackDocEnd(docEnd *DocEnd, userData unsafe.Pointer) RewriterDirective {
	cb := pointer.Restore(userData).(DocEndHandler)
	return cb(docEnd)
}

// strToGoString is a helper function that translates the underlying-library-defined lol_html_str_t data to Go string.
// It is the caller's responsibility to arrange for lol_html_str_t to be freed,
// by calling str.Free() or lol_html_str_free().
// Potential issue: lol_html_str_t->len from size_t (uint) to int (int32) on 32-bit machines?
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

func textChunkContentToGoString(s TextChunkContent) string {
	var nullTextChunkContent TextChunkContent
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
