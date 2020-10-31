package lolhtml_test

import (
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestDoctype_GetDoctypeFields(t *testing.T) {
	w, err := lolhtml.NewWriter(
		nil,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					DoctypeHandler: func(doctype *lolhtml.Doctype) lolhtml.RewriterDirective {
						if name := doctype.Name(); name != "math" {
							t.Errorf("wrong doctype name %s\n", name)
						}
						if publicId := doctype.PublicId(); publicId != "" {
							t.Errorf("wrong doctype name %s\n", publicId)
						}
						if systemId := doctype.SystemId(); systemId != "http://www.w3.org/Math/DTD/mathml1/mathml.dtd" {
							t.Errorf("wrong doctype name %s\n", systemId)
						}
						return lolhtml.Continue
					},
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	defer w.Free()
	_, err = w.Write([]byte(`<!DOCTYPE math SYSTEM "http://www.w3.org/Math/DTD/mathml1/mathml.dtd">`))
	if err != nil {
		t.Error(err)
	}
	err = w.End()
	if err != nil {
		t.Error(err)
	}
}

func TestDoctype_StopRewriting(t *testing.T) {
	w, err := lolhtml.NewWriter(
		nil,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					DoctypeHandler: func(d *lolhtml.Doctype) lolhtml.RewriterDirective {
						return lolhtml.Stop
					},
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}
	defer w.Free()
	_, err = w.Write([]byte("<!doctype>"))
	if err == nil {
		t.FailNow()
	}
	if err.Error() != "The rewriter has been stopped." {
		t.Error(err)
	}
}
