package adapt_test

import (
	"fmt"
	"github.com/go-noodle/noodle"
	"github.com/go-noodle/adapt"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func noodleMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, noodle.WithValue(r, "testKey", "testValue"))
	}
}

func negroniMW(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Fprint(w, "HTTP>")
	next(w, r)
}

func TestNegroniContextPasses(t *testing.T) {
	is := is.New(t)

	n := noodle.New(noodleMW, adapt.Negroni(negroniMW)).Then(
		func(w http.ResponseWriter, r *http.Request) {
			val, ok := noodle.Value(r, "testKey").(string)
			is.True(ok)
			is.Equal(val, "testValue")
		},
	)
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	n(httptest.NewRecorder(), r)
}
