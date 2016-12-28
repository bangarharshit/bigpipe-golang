package main

import (
	"github.com/bangarharshit/bigpipe-golang/lib"
	"html/template"
	"net/http"
)

// HomePageApplication sample application.
type HomePageApplication struct{}

// Data is Content for template.
type Data struct {
	PageletFunc func() bool
}

// Render generates the basic html markup with containers for individual pagelets.
func (homePageApplication HomePageApplication) Render(w http.ResponseWriter, r *http.Request, pageletFunc func() bool) {
	applicationTemplate, err := template.ParseFiles("templates/homepageapplication.html")
	if err != nil {
		panic(err)
	}
	data := Data{pageletFunc}
	err1 := applicationTemplate.Execute(w, data)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

// PageletsContainerMapping return the list of pagelet in the application with containerId.
func (homePageApplication HomePageApplication) PageletsContainerMapping() map[string]bigpipe.Pagelet {
	return map[string]bigpipe.Pagelet{
		"searchPagelet":         SearchPagelet{},
		"recommendationPagelet": RecommendationPagelet{},
		"profilePagelet":        ProfilePagelet{},
		"adsPagelet":            AdsPagelet{},
	}
}
