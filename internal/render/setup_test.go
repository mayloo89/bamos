package render

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/mayloo89/bamos/internal/config"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false
	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())

}

type testWriter struct{}

func (tw *testWriter) Header() http.Header {
	return http.Header{}
}

func (tw *testWriter) WriteHeader(i int) {}

func (tw *testWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
