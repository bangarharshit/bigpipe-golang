package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

// AdsPagelet is 1 more sample service simulation.
type AdsPagelet struct{}

// Render generates html from template. The html returned is then inserted into container by application.
// Note - Clients are responsible for handling the errors on their own and return the error dom element.
func (adsPagelet AdsPagelet) Render(r *http.Request) (ret template.HTML) {
	time.Sleep(10 * time.Second)
	buf := bytes.NewBuffer([]byte{})
	templates, err := template.ParseFiles("templates/adspagelet.gohtml")
	if err != nil {
		return
	}
	templates.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}
