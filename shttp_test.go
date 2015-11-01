package shttp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apriendeau/shttp"
	"github.com/stretchr/testify/assert"
)

type Sample struct {
	Hello string `json:"hello" xml:"hello"`
}

func TestWriteJSON(t *testing.T) {
	assert := assert.New(t)
	var b = struct {
		Hello string `json:"hello"`
	}{"world"}
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(err)

	// test content-negotiation application/json
	r.Header.Set("Accept", "application/json")
	err = shttp.Write(w, r, b, 200)
	assert.NoError(err)
	assert.Equal(w.Code, 200)
	assert.Equal(w.Body.String(), `{"hello":"world"}`)
}

func TestWriteXML(t *testing.T) {
	assert := assert.New(t)
	b := Sample{
		Hello: "world",
	}
	// test content-negotiation application/xml
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/test", nil)
	r.Header.Set("Accept", "application/xml")
	err = shttp.Write(w, r, b, 200)
	assert.NoError(err)
	assert.Equal(w.Code, 200)
	assert.Equal(w.Body.String(), "<Sample><hello>world</hello></Sample>")
}

func TestReadJSON(t *testing.T) {
	assert := assert.New(t)
	body := strings.NewReader(`{"hello":"world"}`)
	req, err := http.NewRequest("POST", "/test", body)
	assert.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	var b Sample
	err = shttp.Read(req, &b)
	assert.Nil(err)
	assert.Equal(b.Hello, "world")
}
func TestReadXML(t *testing.T) {
	assert := assert.New(t)
	body := strings.NewReader(`<Sample><hello>world</hello></Sample>`)
	req, err := http.NewRequest("POST", "/test", body)
	assert.NoError(err)
	req.Header.Set("Content-Type", "application/xml")
	var b Sample
	err = shttp.Read(req, &b)
	assert.Nil(err)
	assert.Equal(b.Hello, "world")
}

func TestErrorJSON(t *testing.T) {
	assert := assert.New(t)
	e := errors.New("test")
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/test", nil)
	assert.NoError(err)
	err = shttp.Error(w, req, e, 200)
	assert.NoError(err)
	assert.Equal(w.Code, 200, "should be expected status code")
	assert.Contains(w.Body.String(), `"description":"test"`, "should be formatted as JSON error")
}

func TestErrorXML(t *testing.T) {
	assert := assert.New(t)
	e := errors.New("test")
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/test", nil)
	assert.NoError(err)
	req.Header.Set("Accept", "application/xml")
	err = shttp.Error(w, req, e, 503)
	assert.NoError(err)
	assert.Equal(w.Code, 503, "should be expected status code")
	assert.Contains(w.Body.String(), "<title>Service Unavailable</title>")
}

func TestStatus(t *testing.T) {
	assert := assert.New(t)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/test", nil)
	assert.NoError(err)
	err = shttp.Status(w, req, 204)
	assert.NoError(err)
	assert.Equal(w.Code, 204, "should be expected status code")
	assert.Equal(w.Body.String(), "", "should be an empty response")
}
