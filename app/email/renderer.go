// Package mail provides utilities for constructing and sending email messages.
package mail

import (
	"bytes"
	"html/template"
)

// RenderTemplate parses an HTML template file from the given path and executes it
// with the provided data, returning the rendered template as a string.
//
// Parameters:
//   - path: filesystem path to the HTML template file
//   - data: data to inject into the template during execution (typically a map or struct)
//
// Returns:
//   - the rendered template content as a string
//   - an error if parsing or executing the template fails
func RenderTemplate(path string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return "", err
	}

	return out.String(), nil
}
