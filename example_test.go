package shttp_test

import (
	"log"
	"net/http"

	"github.com/apriendeau/shttp"
)

func ExampleWrite() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]string)
		if err := shttp.Read(r, &body); err != nil {
			if err := shttp.Error(w, r, err, 422); err != nil {
				log.Panic(err)
			}
		}

		message := struct {
			Message  string            `json:"message"`
			Received map[string]string `json:"received"`
		}{"I received your information.", body}

		if err := shttp.Write(w, r, message, 200); err != nil {
			log.Panic(err)
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/post", handler)
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Panic(err)
	}
}
