package router_test

import (
	. "github.com/bluesoftdev/go-http-matchers/extractor"
	. "github.com/bluesoftdev/go-http-matchers/predicate"
	. "github.com/bluesoftdev/go-http-router/router"
	"github.com/bluesoftdev/go-http-router/util"

	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWhen(t *testing.T) {
	trueCounter := util.CountingHandler(0)
	falseCounter := util.CountingHandler(0)
	handler := Router(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				When(PredicateFunc(func(r interface{}) bool {
					return r.(*http.Request).Header.Get("Accept") == "application/json"
				}), func() {
					DecorateHandler(&trueCounter, NoopHandler)
				}, func() {
					DecorateHandler(&falseCounter, NoopHandler)
				})
			})
		})
	})

	assert.NotNil(t, handler)
	testReq := httptest.NewRequest("GET", "http://localhost/foo/bar", nil)

	mockWriter := httptest.NewRecorder()
	handler.ServeHTTP(mockWriter, testReq)
	assert.Equal(t, 0, int(trueCounter))
	assert.Equal(t, 1, int(falseCounter))

	falseCounter.Reset()
	trueCounter.Reset()

	mockWriter = httptest.NewRecorder()
	testReq.Header.Add("Accept", "application/json")
	handler.ServeHTTP(mockWriter, testReq)
	assert.Equal(t, 1, int(trueCounter))
	assert.Equal(t, 0, int(falseCounter))
}

func TestSwitch(t *testing.T) {
	case1Counter := util.CountingHandler(0)
	case2Counter := util.CountingHandler(0)
	defaultCounter := util.CountingHandler(0)
	handler := Router(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Switch(ExtractorFunc(func(r interface{}) interface{} {
					return r.(*http.Request).Header.Get("Accept")
				}), func() {
					Case(PredicateFunc(func(acceptHeader interface{}) bool {
						return acceptHeader.(string) == "application/json"
					}), func() {
						DecorateHandler(&case1Counter, NoopHandler)
					})
					Case(PredicateFunc(func(acceptHeader interface{}) bool {
						return acceptHeader.(string) == "application/xml"
					}), func() {
						DecorateHandler(&case2Counter, NoopHandler)
					})
					Default(func() {
						DecorateHandler(&defaultCounter, NoopHandler)
					})
				})
			})
		})
	})

	assert.NotNil(t, handler)

	testReq, err := http.NewRequest("GET", "http://localhost/foo/bar", nil)

	assert.NoError(t, err)

	mockWriter := httptest.NewRecorder()
	testReq.Header.Set("Accept", "application/json")
	handler.ServeHTTP(mockWriter, testReq)
	assert.Equal(t, 1, int(case1Counter))
	assert.Equal(t, 0, int(case2Counter))
	assert.Equal(t, 0, int(defaultCounter))

	case1Counter.Reset()

	mockWriter = httptest.NewRecorder()
	testReq.Header.Set("Accept", "application/xml")
	handler.ServeHTTP(mockWriter, testReq)
	assert.Equal(t, 0, int(case1Counter))
	assert.Equal(t, 1, int(case2Counter))
	assert.Equal(t, 0, int(defaultCounter))

	case2Counter.Reset()

	mockWriter = httptest.NewRecorder()
	testReq.Header.Set("Accept", "application/pdf")
	handler.ServeHTTP(mockWriter, testReq)
	assert.Equal(t, 0, int(case1Counter))
	assert.Equal(t, 0, int(case2Counter))
	assert.Equal(t, 1, int(defaultCounter))
}

func TestSwitchWithoutDefault(t *testing.T) {
	case1Counter := util.CountingHandler(0)
	case2Counter := util.CountingHandler(0)
	handler := Router(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Switch(ExtractorFunc(func(r interface{}) interface{} {
					return r.(*http.Request).Header.Get("Accept")
				}), func() {
					Case(PredicateFunc(func(acceptHeader interface{}) bool {
						return acceptHeader.(string) == "application/json"
					}), func() {
						DecorateHandler(&case1Counter, NoopHandler)
					})
					Case(PredicateFunc(func(acceptHeader interface{}) bool {
						return acceptHeader.(string) == "application/xml"
					}), func() {
						DecorateHandler(&case2Counter, NoopHandler)
					})
				})
			})
		})
	})

	assert.NotNil(t, handler)

	testReq, err := http.NewRequest("GET", "http://localhost/foo/bar", nil)

	assert.NoError(t, err)

	mockWriter := httptest.NewRecorder()
	testReq.Header.Set("Accept", "application/json")
	handler.ServeHTTP(mockWriter, testReq)
	assert.Equal(t, 1, int(case1Counter))
	assert.Equal(t, 0, int(case2Counter))

	case1Counter.Reset()
	case2Counter.Reset()

	mockWriter = httptest.NewRecorder()
	testReq.Header.Set("Accept", "application/xml")
	handler.ServeHTTP(mockWriter, testReq)
	assert.Equal(t, 0, int(case1Counter))
	assert.Equal(t, 1, int(case2Counter))

	case1Counter.Reset()
	case2Counter.Reset()

	mockWriter = httptest.NewRecorder()
	testReq.Header.Set("Accept", "application/pdf")
	handler.ServeHTTP(mockWriter, testReq)
	result := mockWriter.Result()
	assert.Equal(t, 404, result.StatusCode)
	assert.Equal(t, 0, int(case1Counter))
	assert.Equal(t, 0, int(case2Counter))

}
