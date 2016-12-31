package bigpipe

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
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
	PreLoad() (ret template.HTML)
}

type PageletChannelContainer struct {
	pagelet Pagelet
	pageletChannelTemplate <- chan template.HTML
}

func clientSideRender(
	rw http.ResponseWriter,
	flusher http.Flusher,
	templateChannelMapping map[string]PageletChannelContainer) {
	cases := make([]reflect.SelectCase, len(templateChannelMapping))
	idContainerMapping := make(map[int]string)
	index := 0
	for i, ch := range templateChannelMapping {
		cases[index] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch.pageletChannelTemplate)}
		idContainerMapping[index] = i
		index = index + 1
	}

	remaining := len(cases)
	for remaining > 0 {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			cases[chosen].Chan = reflect.ValueOf(nil)
			continue
		}
		remaining --
		buf := bytes.NewBuffer([]byte{})
		ret := template.HTML(value.String())
		data := struct {
			ContainerID string
			Data        template.HTML
		}{idContainerMapping[chosen], ret}

		preLoadConent := templateChannelMapping[idContainerMapping[chosen]].pagelet.PreLoad()
		_, err1 := fmt.Fprintf(rw, "%s", preLoadConent)
		if err1 != nil {
			fmt.Println(err1)
			return
		}

		applicationGlueTemplate.Execute(buf, data)
		ret1 := template.HTML(buf.String())
		_, err2 := fmt.Fprintf(rw, "%s", ret1)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		flusher.Flush()
	}
}

func startRequest(r *http.Request, pagelet Pagelet) <-chan template.HTML {
	pageletChannel := make(chan template.HTML)
	go func() {
		pageletChannel <- pagelet.Render(r)
	}()
	return pageletChannel

}

func generateContainerDiv(containerID string) template.HTML {
	buf := bytes.NewBuffer([]byte{})
	data := struct {
		ContainerID string
	}{containerID}
	err := clientSideRendingContainerTemplate.Execute(buf, data)
	if err != nil {
		// todo: error handling
	}
	return template.HTML(buf.String())
}

var applicationGlueScript = "<script type=\"text/javascript\">" +
	"renderInDom({{.Data}}, {{.ContainerID}})" +
	"</script>"

var applicationGlueTemplate = template.Must(template.New("applicationGlue").Parse(applicationGlueScript))

var clientSideRenderingContainer = "<div id={{.ContainerID}}></div>"

var clientSideRendingContainerTemplate = template.Must(template.New("clientSideRenderingContainerTemplate").Parse(clientSideRenderingContainer))
