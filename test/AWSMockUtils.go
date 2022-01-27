package test

import (
	"net/http"
	"net/http/httptest"
)

func MockServer() *httptest.Server {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	return mock
}
