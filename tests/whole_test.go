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
	Age       int
}

func UserGateKeeper() (gk *govalidations.GateKeeper) {
	gk = govalidations.NewGateKeeper()

	gk.Add(govalidations.Regexp(func(object interface{}) interface{} {
		return object.(*User).Email
	}, regexp.MustCompile(`^([^@\s]+)@((?:[-a-z0-9]+\.)+[a-z]{2,})$`), "Email", "Must be a valid email"))

	gk.Add(govalidations.Presence(func(object interface{}) interface{} {
		return object.(*User).Username
	}, "Username", "Username can't be blank"))

	gk.Add(govalidations.Custom(func(object interface{}) bool {
		age := object.(*User).Age
		if age < 18 {
			return false
		}
		return true
	}, "Age", "You must be a grown man"))

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
	if !strings.Contains(body, "You must be a grown man") {
		t.Error(body)
	}
}
