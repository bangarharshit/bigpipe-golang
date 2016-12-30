package main

import (
	"github.com/bangarharshit/bigpipe-golang/lib"
	"net/http"
	"log"
)

func main() {
	homePageApplication := HomePageApplication{}
	http.HandleFunc("/home", bigpipe.ServeApplication(homePageApplication))
	http.ListenAndServe(":3000", nil)
	log.Fatal(http.ListenAndServe(":5000", http.FileServer(http.Dir("static/"))))
}
