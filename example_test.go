package shttp_test

import (
	"log"
	"net/http"

	"github.com/apriendeau/shttp"
)

func ExampleWrite() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		message := struct {
			Message string `json:"message"`
		}{"I received your request."}

		if err := shttp.Write(w, r, message, 200); err != nil {
			log.Panic(err)
		}
	}

	http.HandleFunc("/post", handler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Panic(err)
	}
}

func ExampleRead() {
	// Simple mimic server
	handler := func(w http.ResponseWriter, r *http.Request) {
		// parsing for example sake
		body := make(map[string]interface{})
		if err := shttp.Read(r, &body); err != nil {
			log.Panic(err)
		}
		// send body back to mimic
		if err := shttp.Write(w, r, body, 200); err != nil {
			log.Panic(err)
		}
	}

	http.HandleFunc("/post", handler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Panic(err)
	}
}
