package lolhtml_test

import (
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestRewriter_NonAsciiEncoding(t *testing.T) {
	w, err := lolhtml.NewWriter(
		nil,
		nil,
		lolhtml.Config{
			Encoding: "UTF-16",
			Memory: &lolhtml.MemorySettings{
				PreallocatedParsingBufferSize: 0,
				MaxAllowedMemoryUsage:         16,
			},
			Strict: true,
		})
	if w != nil || err == nil {
		t.FailNow()
	}
	if err.Error() != "Expected ASCII-compatible encoding." {
		t.Error(err)
	}
}

func TestRewriter_MemoryLimiting(t *testing.T) {
	w, err := lolhtml.NewWriter(
		nil,
		&lolhtml.Handlers{
			ElementContentHandler: []lolhtml.ElementContentHandler{
				{
					"span",
					nil,
					nil,
					nil,
				},
			},
		},
		lolhtml.Config{
			Encoding: "utf-8",
			Memory: &lolhtml.MemorySettings{
				PreallocatedParsingBufferSize: 0,
				MaxAllowedMemoryUsage:         5,
			},
			Strict: true,
		},
	)
	if err != nil {
		t.Error(err)
	}
	_, err = w.Write([]byte("<span alt='aaaaa"))
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "The memory limit has been exceeded." {
		t.Error(err)
	}
	w.Free()
}
