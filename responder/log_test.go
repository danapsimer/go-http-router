package responder_test

import (
	"bytes"
	. "github.com/bluesoftdev/go-http-router/responder"
	"github.com/bluesoftdev/go-http-router/router"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

func TestLogRequest(t *testing.T) {
	r := router.Router(func() {
		router.Endpoint("/foo", func() {
			router.Method("GET", func() {
				LogRequest()
				RespondWithString(200, "The quick brown fox jumped over the lazy dogs.")
			})
		})
	})
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}
	mockWriter := httptest.NewRecorder()
	var str bytes.Buffer
	log.SetOutput(&str)
	r.ServeHTTP(mockWriter, request)
	logMessage := str.String()
	assert.Regexp(t, regexp.MustCompile("(?s:.*Request:\\s*GET /foo HTTP/0.0\\s*Host: localhost.*)"), logMessage)
}
