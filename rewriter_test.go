package lolhtml_test

import (
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestRewriter_NonAsciiEncoding(t *testing.T) {
	_, err := lolhtml.NewWriter(
		nil,
		nil,
		lolhtml.Config{
			Encoding: "UTF-16",
			Memory: &lolhtml.MemorySettings{
				PreallocatedParsingBufferSize: 1024,
				MaxAllowedMemoryUsage:         1<<63 - 1,
			},
			Strict: true,
		})
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "Expected ASCII-compatible encoding." {
		t.Error(err)
	}
}
