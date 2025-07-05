package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/yousifsabah0/snippets/internal/models/snippets"
	"github.com/yousifsabah0/snippets/web"
)

type templateData struct {
	CurrentYear     int
	Snippet         snippets.Snippet
	Snippets        []snippets.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func newTemplateCaceh() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(web.Files, "app/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"app/index.html",
			"app/partials/navbar.html",
			"app/partials/footer.html",
			page,
		}

		ts, err := template.New(page).Funcs(functions).ParseFS(web.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
