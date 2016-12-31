package bigpipe

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Sample testApplication with 0 pagelet
type TestApplication struct{}

// Sample testApplication with 1 pagelet
type TestApplicationWithPagelet struct{}

// Sample pagelet
type TestPagelet struct{}

var testApplicationContent = "test string"

// Render just put test string inside responsewriter for testing.
func (testApplication *TestApplication) Render(rw http.ResponseWriter, r *http.Request, servePagelet func() bool, renderPagelet func(pageletId string) template.HTML) {
	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, testApplicationContent)
}

// PageletsContainerMapping return empty pagelet for testing.
func (testApplication *TestApplication) PageletsContainerMapping() map[string]Pagelet {
	return map[string]Pagelet{}
}

// Render just put test string inside response writer for testing.
func (testApplication *TestApplicationWithPagelet) Render(rw http.ResponseWriter, r *http.Request, servePagelet func() bool, renderPagelet func(pageletId string) template.HTML) {
	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, "test string")
	// It is executed in a blocking way by templates. Simulating the same by invoking it manually.
	servePagelet()
}

// PageletsContainerMapping returns pagelet for testing.
func (testApplication *TestApplicationWithPagelet) PageletsContainerMapping() map[string]Pagelet {
	return map[string]Pagelet{
		"testPagelet": TestPagelet{},
	}
}

func (testPagelet TestPagelet) Render(r *http.Request) (ret template.HTML) {
	return template.HTML("pagelet test content")
}

// For testing handler check - https://elithrar.github.io/article/testing-http-handlers-go/
func TestServeApplicationRendersApplicationForClientSideRendering(t *testing.T) {
	testApplication := &TestApplication{}
	handlerFunc := ServeApplication(testApplication, true)
	http.HandleFunc("/test", handlerFunc)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	if rr.Body.String() != testApplicationContent {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), testApplicationContent)
	}
}

func TestServePageletShouldRenderPageletForServerSideRendering(t *testing.T) {
	testApplicationWithPagelet := &TestApplicationWithPagelet{}
	handlerFunc := ServeApplication(testApplicationWithPagelet, true)
	http.HandleFunc("/testwithpagelet", handlerFunc)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/testwithpagelet", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	renderedText := "test string" + // Application template
		// Bigpipe glue
		"<script type=\"text/javascript\">" +
		"function renderInDom(value, containerId) {" +
		"document.getElementById(containerId).innerHTML = value;" +
		"}" +
		"</script>" +
		"<script type=\"text/javascript\">" +
		"renderInDom(\"pagelet test content\", \"testPagelet\")" + // Pagelet container and value.
		"</script>"
	// Check the response body is what we expect.
	if rr.Body.String() != renderedText {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), renderedText)
	}

}
