package main

import (
	_ "embed"
	"html/template"
)

//go:embed template/index.html
var templateIndex []byte

var Templates struct {
	Index *template.Template
}

func InitialiseTemplates() {
	Templates.Index = template.Must(template.New("index").Parse(string(templateIndex)))
}
