package handlers

import (
	"net/http"
	"text/template"

	"github.com/pkg/errors"
)

type HTMLHandler interface {
	Index() http.HandlerFunc
}

type htmlHandler struct {
}

type UserInfo struct {
	Name  string
	Email string
}

func NewHTMLHandler() HTMLHandler {
	return &htmlHandler{}
}

func (h *htmlHandler) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, "./templates/home.page.html", nil)
	}
}

func render(w http.ResponseWriter, path string, data map[string]string) {
	tmpl, err := template.ParseFiles(path, "./templates/main.layout.html")
	if err != nil {
		http.Error(w, errors.Wrap(err, "failed to render html page").Error(), http.StatusInternalServerError)
		return
	}
	if data != nil {
		tmpl.ExecuteTemplate(w, "main", data)
		return
	}
	tmpl.Execute(w, nil)
}
