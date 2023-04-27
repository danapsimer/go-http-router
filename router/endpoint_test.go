package router_test

import (
	"github.com/danapsimer/go-http-matchers/predicate"
	"github.com/danapsimer/go-http-router/responder"
	. "github.com/danapsimer/go-http-router/router"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestEndpointForConditionNoMatch(t *testing.T) {
	mockery := Router(func() {
		EndpointForCondition(predicate.False(), func() {
			responder.Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 404, response.StatusCode)
}

func TestEndpointForConditionMatch(t *testing.T) {
	mockery := Router(func() {
		EndpointForCondition(predicate.True(), func() {
			responder.Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 200, response.StatusCode)
}

func TestEndpointPatternNotFound(t *testing.T) {
	mockery := Router(func() {
		EndpointPattern("/foo/.+", func() {
			responder.Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 404, response.StatusCode)
}

func TestEndpointPatternMatch(t *testing.T) {
	mockery := Router(func() {
		EndpointPattern("/fo{2}", func() {
			responder.Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 200, response.StatusCode)
}

func TestEndpointMatch(t *testing.T) {
	mockery := Router(func() {
		EndpointPattern("/foo", func() {
			responder.Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 200, response.StatusCode)
}

func TestEndpointNoMatch(t *testing.T) {
	mockery := Router(func() {
		EndpointPattern("/bar", func() {
			responder.Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 404, response.StatusCode)
}

func TestEndpointForConditionWithPriority(t *testing.T) {
	mockery := Router(func() {
		EndpointForConditionWithPriority(2, predicate.True(), func() {
			responder.Respond(201)
		})
		EndpointForConditionWithPriority(1, predicate.True(), func() {
			responder.Respond(200)
		})
	})

	// build the request
	testURL, _ := url.ParseRequestURI("http://localhost/foo")
	request := &http.Request{
		Method: "GET",
		Header: http.Header{},
		URL:    testURL,
	}

	// Execute the request
	mockWriter := httptest.NewRecorder()
	mockery.ServeHTTP(mockWriter, request)
	response := mockWriter.Result()

	// Test the result.
	assert.Equal(t, 200, response.StatusCode)
}
