package lolhtml

// some string types passed by c api, and their helper functions

/*
#include "lol_html.h"
*/
import "C"

type str C.lol_html_str_t

// textChunkContent does not need to be de-allocated manually.
type textChunkContent C.lol_html_text_chunk_content_t

func (s *str) Free() {
	if s != nil {
		C.lol_html_str_free(*(*C.lol_html_str_t)(s))
	}
}

// String is a helper function that translates the underlying-library-defined lol_html_str_t data to Go string.
// It is the caller's responsibility to arrange for lol_html_str_t to be freed,
// by calling str.Free() or lol_html_str_free().
// Potential issue: lol_html_str_t->len from size_t (uint) to int (int32) on 32-bit machines?
func (s *str) String() string {
	if s == nil {
		return ""
	}
	return C.GoStringN(s.data, C.int(s.len))
}

func (s *textChunkContent) String() string {
	//var nullTextChunkContent textChunkContent
	//if s == nullTextChunkContent {
	//	return ""
	//}
	if s == nil {
		return ""
	}
	return C.GoStringN(s.data, C.int(s.len))
}
