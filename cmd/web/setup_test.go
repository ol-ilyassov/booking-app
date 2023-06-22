package main

import (
	"net/http"
	"os"
	"testing"
)

// TestMain is a logic place to run all the tests within the package.
// It is used to setup the environment before tests and make operations after tests.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

// cmd/web $: go test -v
// cmd/web $: go test -cover
// cmd/web $: go test -coverprofile=coverage.out && go tool cover -html=coverage.out
