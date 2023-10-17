package router

import (
	"embed"
	"html/template"
	"log"
	"time"
	"io"
)

//go:embed views/*
var templatesDirFs embed.FS

type Views struct {
	templates *template.Template
}

type Meta struct {
	Title string
	Description string
}

func New() *Views {
	// funcMap := template.FuncMap{
	// 	"inc": inc,
	// }
	templates := template.Must(template.New("html_templates").ParseFS(templatesDirFs, "views/*.html"))
	for _, tmpl := range templates.Templates() {
		log.Println(tmpl.Name())
	}
	return &Views{
		templates: templates,
	}
}

type ResultPage struct {
	Meta Meta
	Paste string
	Exp time.Time
}

func (v *Views) renderResultPage(w io.Writer, page *ResultPage) error {
	err := v.templates.ExecuteTemplate(w, "result_page.html", page)
	return err
}

type CreatePage struct {
	Meta Meta
}

func (v *Views) renderCreatePage(w io.Writer, page *CreatePage) error {
	err := v.templates.ExecuteTemplate(w, "create_page.html", page)
	return err
}

var OurViews *Views

func init() {
	OurViews = New()
}
