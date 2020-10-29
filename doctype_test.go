package lolhtml_test

import (
	"testing"

	"github.com/coolspring8/go-lolhtml"
)

func TestDoctypeApi(t *testing.T) {
	w, err := lolhtml.NewWriter(
		nil,
		&lolhtml.Handlers{
			DocumentContentHandler: []lolhtml.DocumentContentHandler{
				{
					DoctypeHandler: func(doctype *lolhtml.Doctype) lolhtml.RewriterDirective {
						name := doctype.Name()
						publicId := doctype.PublicId()
						systemId := doctype.SystemId()
						if name != "math" {
							t.Errorf("wrong doctype name %s\n", name)
						}
						if publicId != "" {
							t.Errorf("wrong doctype name %s\n", publicId)
						}
						if systemId != "http://www.w3.org/Math/DTD/mathml1/mathml.dtd" {
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
