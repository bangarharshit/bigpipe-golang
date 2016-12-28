package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

// RecommendationPagelet is 1 more sample service simulation.
type RecommendationPagelet struct{}

// Render generates html from template. The html returned is then inserted into container by application.
// Note - Clients are responsible for handling the errors on their own and return the error dom element.
func (recommendationPagelet RecommendationPagelet) Render(r *http.Request) (ret template.HTML) {
	time.Sleep(2 * time.Second)
	buf := bytes.NewBuffer([]byte{})
	templates, err := template.ParseFiles("templates/recommendationpagelet.gohtml")
	if err != nil {
		return
	}
	templates.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}
