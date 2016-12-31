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

var adsPageletTemplate = template.Must(template.ParseFiles("templates/adspagelet.gohtml"))
var hideLoadingBarTemplateAds = template.Must(template.ParseFiles("templates/removeLoadingBarAdsPagelet.html"))

func (adsPagelet AdsPagelet) Render(r *http.Request) (ret template.HTML) {
	time.Sleep(10 * time.Second)
	buf := bytes.NewBuffer([]byte{})
	adsPageletTemplate.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}

func (adsPagelet AdsPagelet) PreLoad() (ret template.HTML) {
	buf := bytes.NewBuffer([]byte{})
	hideLoadingBarTemplateAds.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}

