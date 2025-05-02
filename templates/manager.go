package templates

import (
	"embed"
	"text/template"
	"fmt"
)

//go:embed *
var templates embed.FS

func GetTemplate(name string) (*template.Template, error) {
	tmpl, err := templates.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("error reading template file: %w", err)
	}

	return template.New(name).Parse(string(tmpl))
}