package lolhtml

/*
#include <stdio.h>
#include "lol_html.h"

extern void callbackSink(const char *chunk, size_t chunk_len, void *);

extern lol_html_rewriter_directive_t callbackDoctype(lol_html_doctype_t *doctype, void *user_data);

extern lol_html_rewriter_directive_t callbackComment(lol_html_comment_t *comment, void *user_data);

extern lol_html_rewriter_directive_t callbackTextChunk(lol_html_text_chunk_t *text_chunk, void *user_data);

extern lol_html_rewriter_directive_t callbackElement(lol_html_element_t *element, void *user_data);

extern lol_html_rewriter_directive_t callbackDocEnd(lol_html_doc_end_t *doc_end, void *user_data);

void callback_sink(const char *chunk, size_t chunk_len, void *user_data) {
    return callbackSink(chunk, chunk_len, user_data);
}

lol_html_rewriter_directive_t callback_doctype(lol_html_doctype_t *doctype, void *user_data) {
    return callbackDoctype(doctype, user_data);
}

lol_html_rewriter_directive_t callback_comment(lol_html_comment_t *comment, void *user_data) {
    return callbackComment(comment, user_data);
}

lol_html_rewriter_directive_t callback_text_chunk(lol_html_text_chunk_t *text_chunk, void *user_data) {
    return callbackTextChunk(text_chunk, user_data);
}

lol_html_rewriter_directive_t callback_element(lol_html_element_t *element, void *user_data){
    return callbackElement(element, user_data);
}

lol_html_rewriter_directive_t callback_doc_end(lol_html_doc_end_t *doc_end, void *user_data) {
    return callbackDocEnd(doc_end, user_data);
}
*/
import "C"
