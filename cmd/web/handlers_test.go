package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.tmgasek.net/internal/assert"
)

func TestPing(t *testing.T) {
	// Init new httptest.ResponseRecorder.
	rr := httptest.NewRecorder()

	// Init new dummy http.Request.
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(rr, r)

	// Get the response generated by the ping handler.
	rs := rr.Result()

	// Check if the status code was 200.
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// Check if the response body was "OK".
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}