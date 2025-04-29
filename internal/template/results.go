package template

import (
	"html/template"
	"net/http"

	"github.com/dizzydwarfus/tree-builder/internal/shared"
)

func TreeResponse(w http.ResponseWriter) error {
	t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
	if err != nil {
		shared.Red(err.Error())
	}
	err = t.ExecuteTemplate(w, "T", "HTMX Test Message")
	return err
}
