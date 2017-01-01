package bigpipe

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

type RenderPagelet func(pageletId string) template.HTML
type FinishRendering func() bool

// Application is representation of the entire web-page in big-pipe world.
// For details check - https://www.facebook.com/notes/facebook-engineering/bigpipe-pipelining-web-pages-for-high-performance/389414033919/
// Application is composed of small components called pagelets which are rendered in parallel.
// To render the complete webpage, specify the list of pagelets with container-id in PageletsContainerMapping method.
type Application interface {
	// Render generates the basic html markup with containers for individual pagelets.
	Render(rw http.ResponseWriter, r *http.Request, finishRendering FinishRendering, renderPagelet RenderPagelet)

	// PageletsContainerMapping return the list of pagelet in the application with containerId.
	PageletsContainerMapping() map[string]Pagelet

	SetupCache(cacheContainerGenerator CacheContainerGenerator)
}

func servePageletWrapper(
	rw http.ResponseWriter,
	channelTemplateMapping map[string]PageletChannelContainer,
	clientSideRendering bool,
	flusher http.Flusher) func() bool {
	return func() bool {
		return ServePagelet(rw, channelTemplateMapping, clientSideRendering, flusher)
	}
}

// ServeApplication is the handler for rendering the complete web-page.
// It adds application to scope by closure.
func ServeApplication(application Application, clientSideRendering bool) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var cacheContainer CacheContainer
		application.SetupCache(newCache(&cacheContainer))
		rw.Header().Set("content-type", "text/html")
		if (cacheContainer.f == nil) {
			panic("function not implemented")
		}
		channelTemplateMapping := startPageletRendering(application, r, cacheContainer.GetValueForKey)
		flusher, ok := rw.(http.Flusher)
		if !ok {
			http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}
		application.Render(rw, r, servePageletWrapper(rw, channelTemplateMapping, clientSideRendering, flusher), renderPagelet(clientSideRendering, channelTemplateMapping, flusher, rw))
	}
}

// ServePagelet renders individual pagelets in an application. The pagelets are rendered in separate go-routines.
// Note the following things are not implemented and
func ServePagelet(rw http.ResponseWriter, channelTemplateMapping map[string]PageletChannelContainer, isClientSideRendering bool, flusher http.Flusher) (success bool) {
	if !isClientSideRendering {
		success = true
		return
	}
	// This flush is required to flush application component.
	renderBigPipeJavascript(rw)
	flusher.Flush()
	clientSideRender(rw, flusher, channelTemplateMapping)
	success = true
	return
}

func startPageletRendering(application Application, r *http.Request, cacheLookupFunc LookupFunc) map[string]PageletChannelContainer {
	channelTemplateMap := make(map[string]PageletChannelContainer)
	for containerID, pagelet := range application.PageletsContainerMapping() {
		channelTemplateMap[containerID] = PageletChannelContainer{pagelet, startRequest(r, pagelet, cacheLookupFunc)}
	}
	return channelTemplateMap
}

func renderPagelet(isClientSideRendering bool, channelTemplateMapping map[string]PageletChannelContainer, flusher http.Flusher, rw http.ResponseWriter) func(pageletId string) template.HTML {
	return func(pageletId string) template.HTML {
		if isClientSideRendering {
			return generateContainerDiv(pageletId)
		}
		flusher.Flush()
		pageletContentChannel := channelTemplateMapping[pageletId]
		pageletContent := <-pageletContentChannel.pageletChannelTemplate

		preLoadConent := channelTemplateMapping[pageletId].pagelet.PreLoad()
		_, err1 := fmt.Fprintf(rw, "%s", preLoadConent)
		if err1 != nil {
			fmt.Println(err1)
		}

		return pageletContent
	}
}

func renderBigPipeJavascript(rw http.ResponseWriter) {
	buf := bytes.NewBuffer([]byte{})
	bigPipeTemplate.Execute(buf, nil)
	ret := template.HTML(buf.String())
	_, err2 := fmt.Fprintf(rw, "%s", ret)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	return
}

var bigPipe = "<script type=\"text/javascript\">" +
	"function renderInDom(value, containerId) {" +
	"document.getElementById(containerId).innerHTML = value;" +
	"}" +
	"</script>"

var bigPipeTemplate = template.Must(template.New("bigpipe").Parse(bigPipe))
