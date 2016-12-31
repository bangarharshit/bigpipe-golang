package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

// RecommendationPagelet is 1 more sample service simulation.
type RecommendationPagelet struct{}

var recommendationPageletTemplate = template.Must(template.ParseFiles("templates/recommendationpagelet.gohtml"))
var hideLoadingBarTemplateReco = template.Must(template.ParseFiles("templates/removeLoadingBarRecoPagelet.html"))

// Render generates html from template. The html returned is then inserted into container by application.
// Note - Clients are responsible for handling the errors on their own and return the error dom element.
func (recommendationPagelet RecommendationPagelet) Render(r *http.Request) (ret template.HTML) {
	time.Sleep(2 * time.Second)
	buf := bytes.NewBuffer([]byte{})
	recommendationPageletTemplate.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}

func (recommendationPagelet RecommendationPagelet) PreLoad() (ret template.HTML)  {
	buf := bytes.NewBuffer([]byte{})
	hideLoadingBarTemplateReco.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}