package main

import (
	"bytes"
	"html/template"
	"net/http"
	"github.com/bangarharshit/bigpipe-golang/lib"
)

// AdsPagelet is 1 more sample service simulation.
type AdsPagelet2 struct{}

// Render generates html from template. The html returned is then inserted into container by application.
// Note - Clients are responsible for handling the errors on their own and return the error dom element.

var hideLoadingBarTemplateAds2 = template.Must(template.ParseFiles("templates/removeLoadingBarAdsPagelet2.html"))

func (adsPagelet AdsPagelet2) Render(r *http.Request, cacheLookupFunc bigpipe.LookupFunc) (ret template.HTML) {
	val, err := cacheLookupFunc("localhost://ads")
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

func (adsPagelet AdsPagelet2) PreLoad() (ret template.HTML) {
	buf := bytes.NewBuffer([]byte{})
	hideLoadingBarTemplateAds2.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}

