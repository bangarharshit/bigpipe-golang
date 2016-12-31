package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

// SearchPagelet is sample service simulation.
type SearchPagelet struct{}

var searchPageletTempalte = template.Must(template.ParseFiles("templates/searchpagelet.gohtml"))
var hideLoadingBarTemplateSearch = template.Must(template.ParseFiles("templates/removeLoadingBarSearchPagelet.html"))

// Render generates html from template. The html returned is then inserted into container by application.
// Note - Clients are responsible for handling the errors on their own and return the error dom element.
func (searchPagelet SearchPagelet) Render(r *http.Request) (ret template.HTML) {
	time.Sleep(5 * time.Second)
	buf := bytes.NewBuffer([]byte{})
	searchPageletTempalte.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}

func (searchPagelet SearchPagelet) PreLoad() (ret template.HTML) {
	buf := bytes.NewBuffer([]byte{})
	hideLoadingBarTemplateSearch.Execute(buf, nil)
	ret = template.HTML(buf.String())
	return
}