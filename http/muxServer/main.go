package main

import (
	"fmt"
	"net/http"
)

type textHandler struct {
	responseText string
}

func (th *textHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, th.responseText)
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, r.URL.Path)
	}))

	thwelcome := &textHandler{"TextHandle!"}

	mux.Handle("/text", thwelcome)

	http.ListenAndServe(":8000", mux)

}
