package exporter

import (
	"bytes"
	"strings"
	"text/template"
)

func RenderTemplate(Name string, templateData string, data map[string]interface{}) string {
	tmpl, err := template.New(Name).Parse(templateData)
	if err != nil {
		panic(err)
	}
	return renderTemplateToString(tmpl, data)
}

func RenderEnvClause(env map[string]string) string {
	clauses := make([]string, 0, len(env))
	for name, value := range env {
		clauses = append(clauses, name+"="+value)
	}
	return strings.Join(clauses, " ")
}

func renderTemplateToString(template *template.Template, data interface{}) string {
	buffer := new(bytes.Buffer)

	err := template.Execute(buffer, data)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}
