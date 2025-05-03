package emaildomain

import (
	"bytes"
	"html/template"
	"log"
)

// RenderHTMLTemplate renders an HTML template with dynamic data
func RenderHTMLTemplate(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		return "", err
	}

	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return "", err
	}

	return rendered.String(), nil
}
