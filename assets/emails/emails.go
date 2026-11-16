package emails

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed *.html
var FS embed.FS

func Render(name string, data any) (string, error) {
	templ, err := template.ParseFS(FS, name)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = templ.Execute(&buffer, data)
	if err != nil {
		return "", err
	}
	renderedHTML := buffer.String()

	return renderedHTML, nil

}
