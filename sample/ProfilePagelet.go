package main

import (
	"bytes"
	"html/template"
	"net/http"
	"github.com/bangarharshit/bigpipe-golang/lib"
)

// ProfilePagelet is 1 more sample service simulation.
type ProfilePagelet struct{}

var hideLoadingBarTemplateProfile = template.Must(template.ParseFiles("templates/removeLoadingBarProfilePagelet.html"))

// Render generates html from template. The html returned is then inserted into container by application.
// Note - Clients are responsible for handling the errors on their own and return the error dom element.
func (profilePagelet ProfilePagelet) Render(r *http.Request, cacheLookupFunc bigpipe.LookupFunc) (ret template.HTML) {
	val, err := cacheLookupFunc("localhost://profile")
	if err != nil {
		return template.HTML("") //TODO - Should return incorrect html template.
	}
	pageletCallResult := val.(PageletCallResult)
	pageletDataContainer := PageletDataContainer{pageletCallResult.name, pageletCallResult.timeToRender}
	buf := bytes.NewBuffer([]byte{})
	samplePageletTemplate.Execute(buf, pageletDataContainer)
	ret = template.HTML(buf.String())
	return
}

// PreLoad gives chance for any cleanup before the actual content is loaded.
// In this case we are removing the progress bar.
func (profilePagelet ProfilePagelet) PreLoad() (ret template.HTML)  {
	buf := bytes.NewBuffer([]byte{})
	hideLoadingBarTemplateProfile.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}