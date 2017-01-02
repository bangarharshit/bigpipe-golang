package main

import (
	"bytes"
	"html/template"
	"net/http"
	"github.com/bangarharshit/bigpipe-golang/lib"
)

// ErrorPagelet is 1 more sample service simulation.
// It demonstrates how to handle errors.
type ErrorPagelet struct{}

var hideLoadingBarTemplateError = template.Must(template.ParseFiles("templates/removeLoadingBarErrorPagelet.html"))

// Render generates html from template. The html returned is then inserted into container by application.
// Note - Clients are responsible for handling the errors on their own and return the error dom element.
func (errorPagelet ErrorPagelet) Render(r *http.Request, cacheLookupFunc bigpipe.LookupFunc) (ret template.HTML) {
	val, err := cacheLookupFunc("localhost://error")
	buf := bytes.NewBuffer([]byte{})
	if err != nil {
		errorTemplate.Execute(buf, nil)
		ret = template.HTML(buf.String())
		return
	}
	pageletCallResult := val.(PageletCallResult)
	pageletDataContainer := PageletDataContainer{pageletCallResult.name, pageletCallResult.timeToRender}
	samplePageletTemplate.Execute(buf, pageletDataContainer)
	ret = template.HTML(buf.String())
	return
}

// PreLoad gives chance for any cleanup before the actual content is loaded.
// In this case we are removing the progress bar.
func (errorPagelet ErrorPagelet) PreLoad() (ret template.HTML) {
	buf := bytes.NewBuffer([]byte{})
	hideLoadingBarTemplateError.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}

