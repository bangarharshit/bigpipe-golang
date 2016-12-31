package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

// ProfilePagelet is 1 more sample service simulation.
type ProfilePagelet struct{}

var profilePageletTemplate = template.Must(template.ParseFiles("templates/profilepagelet.gohtml"))
var hideLoadingBarTemplateProfile = template.Must(template.ParseFiles("templates/removeLoadingBarProfilePagelet.html"))

// Render generates html from template. The html returned is then inserted into container by application.
// Note - Clients are responsible for handling the errors on their own and return the error dom element.
func (profilePagelet ProfilePagelet) Render(r *http.Request) (ret template.HTML) {
	time.Sleep(650 * time.Millisecond)
	buf := bytes.NewBuffer([]byte{})
	profilePageletTemplate.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}

func (profilePagelet ProfilePagelet) PreLoad() (ret template.HTML)  {
	buf := bytes.NewBuffer([]byte{})
	hideLoadingBarTemplateProfile.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}