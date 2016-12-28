package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

// SearchPagelet is sample service simulation.
type SearchPagelet struct{}

// Render generates html from template. The html returned is then inserted into container by application.
// Note - Clients are responsible for handling the errors on their own and return the error dom element.
func (searchPagelet SearchPagelet) Render(r *http.Request) (ret template.HTML) {
	time.Sleep(5 * time.Second)
	buf := bytes.NewBuffer([]byte{})
	templates, err := template.ParseFiles("templates/searchpagelet.gohtml")
	if err != nil {
		return
	}
	templates.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}
