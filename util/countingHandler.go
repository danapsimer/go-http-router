package util

import "net/http"

type CountingHandler int

func (ch *CountingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	*ch++
}

func (ch *CountingHandler) Reset() {
	*ch = 0
}
