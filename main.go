package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	var arguments = ParseCliArguments()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		if _, err := fmt.Fprintf(w, "Root-Dir: %s\nArgs: %s", arguments.RootDir, r.URL.Query()); err != nil {
			log.Fatal(err)
		}
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	if err := http.ListenAndServe(":8042", nil); err != nil {
		log.Fatal(err)
	}
}
