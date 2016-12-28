package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

// ProfilePagelet is 1 more sample service simulation.
type ProfilePagelet struct{}

// Render generates html from template. The html returned is then inserted into container by application.
// Note - Clients are responsible for handling the errors on their own and return the error dom element.
func (profilePagelet ProfilePagelet) Render(r *http.Request) (ret template.HTML) {
	time.Sleep(650 * time.Millisecond)
	buf := bytes.NewBuffer([]byte{})
	templates, err := template.ParseFiles("templates/profilepagelet.gohtml")
	if err != nil {
		return
	}
	templates.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}
