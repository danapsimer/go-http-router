[![Go Lang Version](https://img.shields.io/badge/go-1.20-00ADD8.svg?style=plastic)](http://golang.com)
[![Go Doc](https://img.shields.io/badge/godoc-reference-00ADD8.svg?style=plastic)](https://godoc.org/github.com/bluesoftdev/mockery)
[![Go Report Card](https://goreportcard.com/badge/github.com/bluesoftdev/mockery?style=plastic)](https://goreportcard.com/report/github.com/bluesoftdev/mockery)
[![codecov](https://img.shields.io/codecov/c/github/bluesoftdev/mockery.svg?style=plastic)](https://codecov.io/gh/bluesoftdev/mockery)
[![CircleCI](https://img.shields.io/circleci/project/github/bluesoftdev/mockery.svg?style=plastic)](https://circleci.com/gh/bluesoftdev/mockery/tree/master)

# Go Http Router
Originally developed as part of the [Mocker](https://github.com/bluesoftdev/mockery) project, this library
provides a versitile, secure, and fast router for building web services.  In addition to being used by
Mockery, it is the heart of the API Gateway product [Argonath](https://github.com/bluesoftdev/argonath)

# Getting Started

Here is an example Router.

``` golang
package mockery_test

import (
	"log"
	"net/http"
	. "github.com/bluesoftdev/go-http-router/router
	. "github.com/bluesoftdev/go-http-matchers/extractor"
	. "github.com/bluesoftdev/go-http-matchers/predicate"
)

func main() {
	mockery := Router(func() {
		Endpoint("/foo/bar", func() {
			Method("GET", func() {
				Header("Content-Type", "application/json")
				Header("FOO", "BAR")
				RespondWithFile(500, "./error.json")
			})
		})
		Endpoint("/foo/bar/", func() {
			Method("GET", func() {
				Header("Content-Type", "application/json")
				Header("FOO", "BAR")
				RespondWithFile(200, "./ok.json")
			})
		})
		Endpoint("/snafu/", func() {
			Method("GET", func() {
				Header("Content-Type", "application/xml")
				Header("Cache-Control", "no-cache")
				Header("Access-Control-Allow-Origin", "*")
				Switch(ExtractQueryParameter("foo"), func() {
					Case(StringEquals("bar"), func() {
						RespondWithFile(http.StatusOK, "response.xml")
					})
					Default(func() {
						RespondWithFile(http.StatusBadRequest, "error.xml")
					})
				})
			})
		})
	})

	log.Fatal(http.ListenAndServe(":8080", mockery))
}
```

# Contributing

see [Contributing](CONTRIBUTING.md)
