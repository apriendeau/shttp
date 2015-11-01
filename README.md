# shttp (simple-http)

So after a while while programming with Go, I discovered that
golang's `net/http` is very nicely designed for writing web apps
but I wanted to add helpers on top of it to cover my API use cases.
The [godoc](https://godoc.org/github.com/apriendeau/shttp) is more
useful than here but here is a quick example.

```go
package main

import (
	"log"
	"net/http"

	"github.com/apriendeau/shttp"
)

func handler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/post", handler)
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Panic(err)
	}
}
```
