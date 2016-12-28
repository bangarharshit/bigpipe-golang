package bigpipe

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

// Pagelet is the single unit of rendering in big-pipe world.
// Pagelets in an application are rendered in parallel.
// TODO: Pagelets doesn't support their own css and js. Bigpipe supports it - https://www.facebook.com/notes/facebook-engineering/bigpipe-pipelining-web-pages-for-high-performance/389414033919/
// TODO: Add a cache between different pagelets to dedupe network calls.
// TODO: Better error handling. Add context for better request handling.
type Pagelet interface {
	// Render generates html from template. The html returned is then inserted into container by application.
	// Note - Clients are responsible for handling the errors on their own and return the error dom element.
	Render(r *http.Request) (ret template.HTML)
}

func startRequest(rw http.ResponseWriter, flusher http.Flusher, pagelet Pagelet, wg *sync.WaitGroup, r *http.Request, lock *sync.Mutex, containerID string) {
	wg.Add(1)
	go func(rw http.ResponseWriter, flusher http.Flusher, pagelet Pagelet, wg *sync.WaitGroup, containerID string) {
		defer wg.Done()
		ret := pagelet.Render(r)
		buf := bytes.NewBuffer([]byte{})
		data := struct {
			ContainerID string
			Data        template.HTML
		}{containerID, ret}
		applicationGlueTemplate.Execute(buf, data)
		ret1 := template.HTML(buf.String())
		_, err2 := fmt.Fprintf(rw, "%s", ret1)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		lock.Lock()
		defer lock.Unlock()
		flusher.Flush()
	}(rw, flusher, pagelet, wg, containerID)
}

var applicationGlueScript = "<script type=\"text/javascript\">" +
	"renderInDom({{.Data}}, {{.ContainerID}})" +
	"</script>"

var applicationGlueTemplate = template.Must(template.New("applicationGlue").Parse(applicationGlueScript))
