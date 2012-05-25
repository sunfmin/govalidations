package tests

import (
	"fmt"
	"github.com/sunfmin/govalidations"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

type User struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	Bio       string
}

func UserGateKeeper() (gk *govalidations.GateKeeper) {
	gk = govalidations.NewGateKeeper()

	gk.Add(govalidations.FormatValidator(func(object interface{}) interface{} {
		return object.(*User).Email
	}, regexp.MustCompile(`^([^@\s]+)@((?:[-a-z0-9]+\.)+[a-z]{2,})$`), "Email", "Must be a valid email"))

	gk.Add(govalidations.PresenceValidator(func(object interface{}) interface{} {
		return object.(*User).Username
	}, "Username", "Username can't be blank"))

	return
}

func theMux() (sm *http.ServeMux) {
	sm = http.NewServeMux()

	tpl := template.Must(template.ParseGlob("validate.html"))

	gk := UserGateKeeper()

	sm.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		u := &User{
			Username: "",
			Email:    "fake",
		}

		vd := gk.Validate(u)
		if vd.HasError() {
			tpl.Execute(w, vd)
			return
		}

		fmt.Fprintln(w, "Yeah!")
	})

	return
}

func TestRenderErrors(t *testing.T) {
	ts := httptest.NewServer(theMux())
	defer ts.Close()

	r, _ := http.Get(ts.URL + "/validate")

	b, _ := ioutil.ReadAll(r.Body)
	body := string(b)
	if !strings.Contains(body, "Must be a valid email") {
		t.Error(body)
	}
}
