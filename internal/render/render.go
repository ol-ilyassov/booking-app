package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/ol-ilyassov/booking-app/internal/config"
	"github.com/ol-ilyassov/booking-app/internal/models"
)

var functions template.FuncMap = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}
var app *config.AppConfig
var pathToTemplates string = "./templates"

// NewRenderer sets the config for the template package.
func NewRenderer(a *config.AppConfig) {
	app = a
}

// HumanDate returns time in YYYY-MM-DD format.
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate returns formatted date.
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// Iterate returns a slice of ints, starting at 1, goint to count.
func Iterate(count int) []int {
	var items []int
	for i := 1; i <= count; i++ {
		items = append(items, i)
	}
	return items
}

func Add(a, b int) int {
	return a + b
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)

	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}

	return td
}

// Template renders templates using html/template package.
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	// templates could be read from cache or from origin.
	var tc map[string]*template.Template
	var err error

	if app.UseCache {
		// get the template cache from the app config.
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
		// if err != nil {
		// 	fmt.Println("error, could not create templates:", err)
		// 	return err
		// }
	}

	t, ok := tc[tmpl]
	if !ok {
		// log.Println("could not get template from template cache")
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	// add default common data among web pages.
	td = AddDefaultData(td, r)

	err = t.Execute(buf, td)
	if err != nil {
		// fmt.Println("error, could not render templates:", err)
		// return err
		log.Fatal(err)

	}

	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error, writing template to browser:", err)
		return err
	}

	return nil
}

// CreateTemplateCache creates a template cache as a map.
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// fmt.Sprintf solves the problem when running program from root directory or specific folder (as handlers).
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
