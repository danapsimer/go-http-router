package router

import (
	"github.com/danapsimer/go-http-matchers/extractor"
	"github.com/danapsimer/go-http-matchers/predicate"
	"net/http"
)

type when struct {
	predicate     predicate.Predicate
	trueResponse  http.Handler
	falseResponse http.Handler
}

// When can be used within a Method's config function to conditionally choose one Response or another.
func When(predicate predicate.Predicate, trueResponseBuilder func(), falseResponseBuilder func()) {

	outerHandler := CurrentHandler()
	trueResponseBuilder()
	trueHandler := CurrentHandler()

	SetCurrentHandler(outerHandler)
	falseResponseBuilder()
	falseHandler := CurrentHandler()

	SetCurrentHandler(&when{predicate, trueHandler, falseHandler})
}

func (wh *when) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if wh.predicate.Accept(request) {
		wh.trueResponse.ServeHTTP(w, request)
	} else {
		wh.falseResponse.ServeHTTP(w, request)
	}
}

type switchCase struct {
	predicate predicate.Predicate
	response  http.Handler
}

type switchCaseSet struct {
	keySupplier    extractor.Extractor
	switchCases    []*switchCase
	defaultHandler http.Handler
}

func (scs *switchCaseSet) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	key := scs.keySupplier.Extract(request)
	for _, sc := range scs.switchCases {
		if sc.predicate.Accept(key) {
			sc.response.ServeHTTP(w, request)
			return
		}
	}
	scs.defaultHandler.ServeHTTP(w, request)
}

var currentSwitch *switchCaseSet

// Switch can be used within a Method's config function to conditionally choose one of many possible responses.  The
// first Case whose predicate returns true will be selected.  Otherwise the Response defined in the Default is used.
// If there is no Default, then 404 is returned with an empty Body.
func Switch(keySupplier extractor.Extractor, cases func()) {
	handler := CurrentHandler()
	sw := &switchCaseSet{
		keySupplier: keySupplier,
		switchCases: make([]*switchCase, 0, 10),
		defaultHandler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			handler.ServeHTTP(w, request)
			w.WriteHeader(404)
		}),
	}
	outerSwitch := currentSwitch
	currentSwitch = sw
	cases()
	SetCurrentHandler(currentSwitch)
	currentSwitch = outerSwitch
}

// Case used within a Switch to define a Response that will be returned if the case's predicate is true.  The order of
// the case calls matter as the first to match will be used.
func Case(predicate predicate.Predicate, responseBuilder func()) {
	outerHandler := CurrentHandler()
	responseBuilder()
	responseHandler := CurrentHandler()
	if predicate != nil {
		currentSwitch.switchCases = append(currentSwitch.switchCases, &switchCase{predicate, responseHandler})
	} else {
		currentSwitch.defaultHandler = responseHandler
	}
	SetCurrentHandler(outerHandler)
}

// Default used to define the Response that will be returned when no other case is triggered.  The default can be placed
// anywhere but there can only be one.
func Default(responseBuilder func()) {
	Case(nil, responseBuilder)
}
