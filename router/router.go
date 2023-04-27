package router

import (
	"github.com/danapsimer/go-http-matchers/predicate"
	"net/http"
	"sort"
)

type routingHandler struct {
	priority  int
	predicate predicate.Predicate
	handler   http.Handler
}

type byPriority []*routingHandler

func (a byPriority) Len() int           { return len(a) }
func (a byPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPriority) Less(i, j int) bool { return a[i].priority < a[j].priority }

type router struct {
	mux            *http.ServeMux
	handlers       byPriority
	currentHandler http.Handler
}

func (m *router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	for _, h := range m.handlers {
		if h.predicate.Accept(request) {
			h.handler.ServeHTTP(w, request)
			return
		}
	}
	w.WriteHeader(404)
}

func (m *router) Handle(path string, handler http.Handler) {
	if m.mux == nil {
		m.mux = http.NewServeMux()
		m.HandleForCondition(DefaultPriority, predicate.PredicateFunc(func(r interface{}) bool {
			_, p := m.mux.Handler(r.(*http.Request))
			return p != ""
		}), m.mux)
	}
	m.mux.Handle(path, handler)
}

func (m *router) HandleForCondition(priority int, predicate predicate.Predicate, handler http.Handler) {
	m.handlers = append(m.handlers, &routingHandler{priority, predicate, handler})
}

var (
	currentRouter *router
)

// Router contains the top level dispatcher.  This method establishes the root handler and the configFunc is called to
// create handlers for the various routes.  Once the config method returns some clean up actions will occur and the
// router handler will be returned.
func Router(configFunc func()) http.Handler {
	currentRouter = &router{handlers: make(byPriority, 0, 10)}
	currentRouter.currentHandler = NoopHandler
	defer func() { currentRouter = nil }()
	configFunc()
	sort.Stable(currentRouter.handlers)
	return currentRouter
}

// CurrentHandler returns the current handler that should be decorated with any additional behaviors.
func CurrentHandler() http.Handler {
	return currentRouter.currentHandler
}

func SetCurrentHandler(handler http.Handler) {
	currentRouter.currentHandler = handler
}

// NoopHandler is a handler that does nothing.
var NoopHandler http.HandlerFunc = func(w http.ResponseWriter, request *http.Request) {
}

// DecorateHandler is used by DSL methods to inject pre & post actions to the current handler.  For instance, the
// Header(string,string) function adds a preHandler that adds a Header to the ResponseWriter.  To use this function
// to create new DSL Methods, follow this pattern:
//
//	   func Header(name, value string) {
//	     return DecorateHandler(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
//					w.Header.Add(name,value)
//	   	}), NoopHandler)
//	   }
func DecorateHandler(preHandler, postHandler http.Handler) {
	delegate := CurrentHandler()
	SetCurrentHandler(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		preHandler.ServeHTTP(w, request)
		delegate.ServeHTTP(w, request)
		postHandler.ServeHTTP(w, request)
	}))
}

// DecorateHandlerBefore is like DecorateHandler but only applies a Decoration before the handler is called.
func DecorateHandlerBefore(preHandler http.Handler) {
	DecorateHandler(preHandler, NoopHandler)
}

// DecorateHandlerAfter is like DecorateHandler but only applies a Decoration after the handler is called.
func DecorateHandlerAfter(postHandler http.Handler) {
	DecorateHandler(NoopHandler, postHandler)
}
