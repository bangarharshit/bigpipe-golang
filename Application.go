package bigpipe

import (
	"net/http"
	"sync"
	"bytes"
	"html/template"
	"fmt"
)

// Application is representation of the entire web-page in big-pipe world.
// For details check - https://www.facebook.com/notes/facebook-engineering/bigpipe-pipelining-web-pages-for-high-performance/389414033919/
// Application is composed of small components called pagelets which are rendered in parallel.
// To render the complete webpage, specify the list of pagelets with container-id in PageletsContainerMapping method.

type Application interface {
	// Render generates the basic html markup with containers for individual pagelets.
	Render(rw http.ResponseWriter, r *http.Request, servePagelet func() bool)

	// PageletsContainerMapping return the list of pagelet in the application with containerId.
	PageletsContainerMapping() map[string]Pagelet
}

// HandlerFunc is a shortcut type for handler functions.
type HandlerFunc func(rw http.ResponseWriter, r *http.Request)

func servePageletWrapper(rw http.ResponseWriter, r *http.Request, application Application) func() bool {
	return func() bool {
		return ServePagelet(rw, r, application)
	}
}
// ServeApplication is the handler for rendering the complete web-page.
// It adds application to scope by closure.
func ServeApplication(application Application) HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("content-type", "text/html")
		application.Render(rw, r, servePageletWrapper(rw, r, application))
	}
}

// ServePagelet renders individual pagelets in an application. The pagelets are rendered in separate go-routines.
// Note the following things are not implemented and
func ServePagelet(rw http.ResponseWriter, r *http.Request, application Application) (success bool){
	wg := sync.WaitGroup{}

	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		success = false
		return
	}
	// This flush is required to flush application component.
	renderBigPipeJavascript(rw)
	flusher.Flush()
	lock := &sync.Mutex{}
	for containerID, pagelet := range application.PageletsContainerMapping() {
		startRequest(rw, flusher, pagelet, &wg, r, lock, containerID)
	}
	wg.Wait()
	success = true
	return
}

func renderBigPipeJavascript(rw http.ResponseWriter) {
	buf := bytes.NewBuffer([]byte{})
	templates, err := template.ParseFiles("templates/bigpipe.html")
	if err != nil {
		return
	}
	templates.Execute(buf, nil)
	ret := template.HTML(buf.String())
	_, err2 := fmt.Fprintf(rw, "%s", ret)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	return
}
