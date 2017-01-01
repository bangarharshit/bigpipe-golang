package main

import (
	"bytes"
	"html/template"
	"net/http"
	"github.com/bangarharshit/bigpipe-golang/lib"
)

// RecommendationPagelet is 1 more sample service simulation.
type RecommendationPagelet struct{}

var hideLoadingBarTemplateReco = template.Must(template.ParseFiles("templates/removeLoadingBarRecoPagelet.html"))

// Render generates html from template. The html returned is then inserted into container by application.
// Note - Clients are responsible for handling the errors on their own and return the error dom element.
func (recommendationPagelet RecommendationPagelet) Render(r *http.Request, cacheLookupFunc bigpipe.LookupFunc) (ret template.HTML) {
	val, err := cacheLookupFunc("localhost://reco")
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

func (recommendationPagelet RecommendationPagelet) PreLoad() (ret template.HTML)  {
	buf := bytes.NewBuffer([]byte{})
	hideLoadingBarTemplateReco.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}