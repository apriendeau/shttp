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
	mux := http.NewServeMux()
	mux.HandleFunc("/post", handler)
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Panic(err)
	}
}
