package router

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

//go:embed views/*
var templatesDirFs embed.FS

type Views struct {
	templates *template.Template
}

type Meta struct {
	Title       string
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
	Meta     Meta
	Paste    string
	Exp      time.Time
	Filename string
	Slug     string
}

func (v *Views) renderResultPage(w http.ResponseWriter, page *ResultPage) error {
	w.Header().Add("X-Content-Type-Options", "nosniff")
	err := v.templates.ExecuteTemplate(w, "result_page.html", page)
	return err
}

type CreatePage struct {
	Meta          Meta
	Expiration    time.Time
	MinExpiration time.Time
	Error         *ErrorResponse
}

func (v *Views) renderCreatePage(w http.ResponseWriter, page *CreatePage) error {
	w.Header().Add("X-Content-Type-Options", "nosniff")
	err := v.templates.ExecuteTemplate(w, "create_page.html", page)
	return err
}

type ErrorResponse struct {
	StatusCode   int
	ErrorMessage string
}

func (v *Views) renderErrorPage(w http.ResponseWriter, page *CreatePage) error {
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.WriteHeader(page.Error.StatusCode)
	err := v.templates.ExecuteTemplate(w, "create_page.html", page)
	return err
}

func (v *Views) renderErrorPageShorthand(w http.ResponseWriter, statusCode int, message string) error {
	now := time.Now()
	return v.renderErrorPage(w, &CreatePage{
		Meta: Meta{
			Title:       "Error",
			Description: fmt.Sprintf("This paste service encountered an error: %q", message),
		},
		Expiration:    now.Add(time.Hour * DEFAULT_EXPIRATION_HOURS),
		MinExpiration: now.Add(time.Hour * 24),
		Error: &ErrorResponse{
			StatusCode:   statusCode,
			ErrorMessage: message,
		},
	})
}

var OurViews *Views

func init() {
	OurViews = New()
}
