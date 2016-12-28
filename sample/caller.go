package main

import (
	"github.com/bangarharshit/bigpipe-golang/lib"
	"net/http"
)

func main() {
	homePageApplication := HomePageApplication{}
	http.HandleFunc("/", bigpipe.ServeApplication(homePageApplication))
	http.ListenAndServe(":3000", nil)
}
