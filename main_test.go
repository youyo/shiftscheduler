package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/youyo/shiftscheduler/routers"
)

func TestRoute(t *testing.T) {
	router := routers.Setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 301, w.Code)
	assert.Equal(t, "<a href=\"/login\">Moved Permanently</a>.\n\n", w.Body.String())
}

func TestRouteLogin(t *testing.T) {
	router := routers.Setup()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/login", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"auth\":\"required\"}", w.Body.String())
}
