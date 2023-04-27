package router_test

import (
	. "github.com/danapsimer/go-http-matchers/predicate"
	. "github.com/danapsimer/go-http-router/responder"
	. "github.com/danapsimer/go-http-router/router"
	"github.com/danapsimer/go-http-router/util"

	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	handler := Router(func() {
		Header("SNAFU", "BAZ")
		pathPattern := regexp.MustCompile("/foo/bar/snafu.*")
		EndpointForCondition(
			And(PathMatches(pathPattern), MethodIs("GET")),
			func() {
				Header("FOO", "SNAFU")
				RespondWithFile(200, "../testdata/ok.json")
			})
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Header("FOO", "BAR")
				RespondWithFile(500, "../testdata/error.json")
			})
		})
		Endpoint("/foo/bar/", func() {
			Method("GET", func() {
				Header("FOO", "BAZ")
				RespondWithFile(200, "../testdata/ok.json")
			})
		})
	})

	mockWriter := httptest.NewRecorder()
	mockRequest := httptest.NewRequest("GET", "/foo/bar", nil)

	handler.ServeHTTP(mockWriter, mockRequest)

	assert.Equal(t, 500, mockWriter.Code)
	assert.Equal(t, "{\"error\": \"This is an error\"}", mockWriter.Body.String())
	assert.Equal(t, "BAZ", mockWriter.Header().Get("SNAFU"))

	mockWriter = httptest.NewRecorder()
	mockRequest = httptest.NewRequest("GET", "/foo/bar/snafu", nil)

	handler.ServeHTTP(mockWriter, mockRequest)

	assert.Equal(t, 200, mockWriter.Code)
	assert.Equal(t, "{\"ok\": \"everything is ok!\"}", mockWriter.Body.String())
	assert.Equal(t, "SNAFU", mockWriter.Header().Get("FOO"))

	mockWriter = httptest.NewRecorder()
	mockRequest = httptest.NewRequest("GET", "/foo/bar/fubar", nil)

	handler.ServeHTTP(mockWriter, mockRequest)

	assert.Equal(t, 200, mockWriter.Code)
	assert.Equal(t, "{\"ok\": \"everything is ok!\"}", mockWriter.Body.String())
	assert.Equal(t, "BAZ", mockWriter.Header().Get("FOO"))
}

func BenchmarkServeHTTP(b *testing.B) {
	handler := Router(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Header("FOO", "BAR")
				RespondWithFile(500, "../testdata/error.json")
			})
		})
		EndpointPattern("/foo/bar/snafu", func() {

		})
		Endpoint("/foo/bar/", func() {
			Method("GET", func() {
				Header("FOO", "BAR")
				RespondWithFile(200, "ok.json")
			})
		})
	})

	for i := 0; i < b.N; i++ {
		mockWriter := httptest.NewRecorder()
		mockRequest := httptest.NewRequest("GET", "/foo/bar/snafu", nil)
		handler.ServeHTTP(mockWriter, mockRequest)
	}
}

func TestDecorateHandler(t *testing.T) {
	preCounter := util.CountingHandler(0)
	counter := util.CountingHandler(0)
	postCounter := util.CountingHandler(0)
	r := Router(func() {
		Endpoint("/foo/bar/snafu", func() {
			Method("GET", func() {
				DecorateHandlerBefore(&counter)
				DecorateHandler(&preCounter, &postCounter)
			})
		})
	})
	mockWriter := httptest.NewRecorder()
	mockRequest := httptest.NewRequest("GET", "/foo/bar/snafu", nil)
	r.ServeHTTP(mockWriter, mockRequest)
	assert.Equal(t, 1, int(preCounter))
	assert.Equal(t, 1, int(counter))
	assert.Equal(t, 1, int(postCounter))
}

func TestDecorateHandlerBefore(t *testing.T) {
	preCounter := util.CountingHandler(0)
	counter := util.CountingHandler(0)
	r := Router(func() {
		Endpoint("/foo/bar/snafu", func() {
			Method("GET", func() {
				DecorateHandlerBefore(&counter)
				DecorateHandlerBefore(&preCounter)
			})
		})
	})
	mockWriter := httptest.NewRecorder()
	mockRequest := httptest.NewRequest("GET", "/foo/bar/snafu", nil)
	r.ServeHTTP(mockWriter, mockRequest)
	assert.Equal(t, 1, int(preCounter))
	assert.Equal(t, 1, int(counter))
}

func TestDecorateHandlerAfter(t *testing.T) {
	postCounter := util.CountingHandler(0)
	counter := util.CountingHandler(0)
	r := Router(func() {
		Endpoint("/foo/bar/snafu", func() {
			Method("GET", func() {
				DecorateHandlerBefore(&counter)
				DecorateHandlerAfter(&postCounter)
			})
		})
	})
	mockWriter := httptest.NewRecorder()
	mockRequest := httptest.NewRequest("GET", "/foo/bar/snafu", nil)
	r.ServeHTTP(mockWriter, mockRequest)
	assert.Equal(t, 1, int(postCounter))
	assert.Equal(t, 1, int(counter))
}
