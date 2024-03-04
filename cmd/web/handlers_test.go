package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sparrowsl/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {

	app := &application{
		// errorLog: log.New(io.Discard, "", 0),
		// infoLog:  log.New(io.Discard, "", 0),
	}

	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	result, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}
	defer result.Body.Close()

	assert.Equal(t, result.StatusCode, http.StatusOK)

	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
