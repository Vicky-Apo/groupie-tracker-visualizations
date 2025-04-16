package utils

import (
	"html/template"
)

func ParseTemplates() *template.Template {
	funcMap := template.FuncMap{
		"replaceSpaces":  ReplaceSpaces,
		"cleanDate":      CleanDate,
		"formatLocation": FormatLocation,
	}
	return template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
}
