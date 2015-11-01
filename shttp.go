package shttp

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	errors "github.com/utopiansociety/errors"
)

// Err implements the error interface but also has a couple other fields including stack traces
// so you can see where exactly the error came from in the API when debugging.
type Err struct {
	Title       string `json:"title" xml:"title"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
	Stack       string `json:"stack,omitempty" xml:"stack,omitempty"`
	Err         error  `json:"error" xml:"error"`
}

// Error is for implementing the error interface
func (e Err) Error() string {
	return e.Error()
}

// Write is what you use to to respond to the req when there is no error.
// It has content negotiation for xml and json so that people can choose their
// preferred data format.
func Write(w http.ResponseWriter, r *http.Request, body interface{}, status int) (err error) {
	var c []byte
	accept := r.Header.Get("Accept")
	w.WriteHeader(status)
	switch {
	case accept == "application/xml":
		w.Header().Set("Content-Type", "application/xml")
		c, err = xml.Marshal(body)
	default:
		w.Header().Set("Content-Type", "application/json")
		c, err = json.Marshal(body)
	}
	if err != nil {
		return err
	}
	_, err = w.Write(c)
	return err
}

// Read is a helper that responds parses the req.Body into the interface of your
// choice so that you can actually do something useful with the information that
// is sent to you.
func Read(r *http.Request, body interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	t := r.Header.Get("Content-Type")
	switch {
	case t == "application/xml":
		return xml.Unmarshal(b, body)
	default:
		return json.Unmarshal(b, body)
	}
}

// Error is used for when there is any kind of error and returns a consistant
// error format that the user can parse back.
func Error(w http.ResponseWriter, r *http.Request, err error, status int) error {
	wrapped := errors.Wrap(err, 1)
	body := Err{
		Title:       http.StatusText(status),
		Description: wrapped.Error(),
		Err:         wrapped,
		Stack:       string(wrapped.StackTrace()),
	}
	return Write(w, r, body, status)
}

// Status only writes only status back to the req.
func Status(w http.ResponseWriter, r *http.Request, status int) error {
	w.WriteHeader(status)
	return nil
}
