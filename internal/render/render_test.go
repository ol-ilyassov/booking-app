package render

import (
	"net/http"
	"testing"

	"github.com/ol-ilyassov/booking-app/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "test case 1")

	result := AddDefaultData(&td, r)
	if result.Flash != "test case 1" {
		t.Error("flash value of 'test case 1' not found in session")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var w myResponseWriter

	err = RenderTemplate(&w, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error, writing template to browser:", err)
	}

	app.UseCache = true

	err = RenderTemplate(&w, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("error, rendered template that does not exists")
	}
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
