package main

import (
	"github.com/bangarharshit/bigpipe-golang/lib"
	"log"
	"net/http"
)

func main() {
	homePageApplication := HomePageApplication{}
	go func() {
		log.Fatal(http.ListenAndServe(":5000", http.FileServer(http.Dir("static/"))))
	}()
	http.HandleFunc("/home", bigpipe.ServeApplication(homePageApplication, false))
	http.ListenAndServe(":3000", nil)
}
