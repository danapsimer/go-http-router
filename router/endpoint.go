package router

import (
	"github.com/danapsimer/go-http-matchers/extractor"
	"github.com/danapsimer/go-http-matchers/predicate"
	"regexp"
)

// Endpoint defines an endpoint that uses the http.ServeMux to dispatch requests.  All subsequent uses of this DSL
// element will add to the same ServeMux instance.  The ServeMux handler will have default priority.  The content of the
// configureFunc should be Method elements because this element opens a Switch on the http Method value.
func Endpoint(url string, configureFunc func()) {
	outerHandler := CurrentHandler()
	Switch(extractor.ExtractMethod(), configureFunc)
	currentRouter.Handle(url, CurrentHandler())
	SetCurrentHandler(outerHandler)
}

// DefaultPriority is the default priority for endpoint consideration.
const DefaultPriority = 100

// EndpointPattern creates an endpoint that is selected by comparing the URL path with the pattern provided.
func EndpointPattern(urlPattern string, configFunc func()) {
	pathRegex := regexp.MustCompile(urlPattern)
	EndpointForCondition(predicate.PathMatches(pathRegex), configFunc)
}

// EndpointForCondition creates an endpoint that is selected by the predicate passed.
func EndpointForCondition(predicate predicate.Predicate, configFunc func()) {
	EndpointForConditionWithPriority(DefaultPriority, predicate, configFunc)
}

// EndpointForConditionWithPriority defines an endpoint that is selected by the predicate given with the priority
// provided.
func EndpointForConditionWithPriority(priority int, predicate predicate.Predicate, configFunc func()) {
	outerHandler := CurrentHandler()
	configFunc()
	currentRouter.HandleForCondition(priority, predicate, CurrentHandler())
	SetCurrentHandler(outerHandler)
}
