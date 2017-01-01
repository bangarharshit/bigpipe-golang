package main

import (
	"github.com/bangarharshit/bigpipe-golang/lib"
	"html/template"
	"net/http"
	"time"
	"errors"
	"fmt"
)

// HomePageApplication sample application.
type HomePageApplication struct{}

// Data is Content for template.
type Data struct {
	FinishRendering bigpipe.FinishRendering
	RenderPagelet   bigpipe.RenderPagelet
}
/**
	Replace with actual http call:
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
 */
var mockHttpClientCall = func(url string) (interface{}, error) {
	fmt.Println("call to http client " + url)
	if url == "localhost://ads" {
		time.Sleep(10 * time.Second)
		return PageletCallResult{"Ads Pagelet", "Slow. Took 10s to render."}, nil
	} else if url == "localhost://profile" {
		time.Sleep(650 * time.Millisecond)
		return PageletCallResult{"Profile Pagelet", "Profile took 650 ms to render."}, nil
	} else if url == "localhost://search" {
		time.Sleep(5 * time.Second)
		return PageletCallResult{"Search Pagelet", "Search took 5s to render."}, nil
	} else if url == "localhost://reco" {
		time.Sleep(2 * time.Second)
		return PageletCallResult{"Recommendation Pagelet", "Recommendation took 2s to render."}, nil
	} else {
		time.Sleep(15 * time.Second)
		return nil, errors.New("Pagelet not available")
	}
}

type PageletCallResult struct {
	name string
	timeToRender string
}

type PageletDataContainer struct{
	Name string
	TimeToRender string
}

func (homePageApplication HomePageApplication) SetupCache(cacheContainerGenerator bigpipe.CacheContainerGenerator) {
	cacheContainerGenerator(mockHttpClientCall)
}

var samplePageletTemplate = template.Must(template.ParseFiles("templates/samplepagelet.gohtml"))
var errorTemplate = template.Must(template.ParseFiles("templates/errorpagelet.gohtml"))

// Render generates the basic html markup with containers for individual pagelets.
func (homePageApplication HomePageApplication) Render(w http.ResponseWriter, r *http.Request, finishRendering bigpipe.FinishRendering, renderPagelet bigpipe.RenderPagelet) {
	applicationTemplate, err := template.ParseFiles("templates/homepageapplication.html")
	if err != nil {
		panic(err)
	}
	data := Data{finishRendering, renderPagelet}
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
		"adsPagelet2":           AdsPagelet2{},
		"errorPagelet":		 ErrorPagelet{},
	}
}
