package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ol-ilyassov/booking-app/internal/config"
	"github.com/ol-ilyassov/booking-app/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})

	// Application working mode (development|production).
	testApp.InProduction = false

	// Session management.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session
	app = &testApp

	os.Exit(m.Run())
}

type myResponseWriter struct{}

func (rw *myResponseWriter) Header() http.Header {
	var h http.Header
	return h
}

func (rw *myResponseWriter) WriteHeader(statusCode int) {}

func (rw *myResponseWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
