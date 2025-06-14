package domain

import (
	"bytes"
	"html/template"
)

// RenderHTMLTemplate renders an HTML template with dynamic data
func RenderHTMLTemplate(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, data)
	if err != nil {
		return "", err
	}

	return rendered.String(), nil
}
