package main

import (
	"github.com/gobuffalo/packr/v2"
	"html/template"
	"log"
)

func readTemplate(filename string) *template.Template {
	templates := packr.New("templates", "./templates")
	parsedTemplate := template.New("")
	parsedTemplate = addTemplateFile(templates, filename, parsedTemplate)
	parsedTemplate = addTemplateFile(templates, "base.gohtml", parsedTemplate)

	return parsedTemplate
}

func addTemplateFile(templates *packr.Box, filename string, parsedTemplate *template.Template) *template.Template {
	templateString, err := templates.FindString(filename)
	if err != nil {
		log.Fatal(err)
	}

	parsedTemplate, err = parsedTemplate.New(filename).Parse(templateString)
	if err != nil {
		log.Fatal(err)
	}

	return parsedTemplate
}
