package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	arguments := NewCliArgumentsParser().GetCliArguments()
	directoryIndex := NewDirectoryIndex(arguments.RootDir, arguments.Extensions)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		directoryIndex.Handle(writer, request, func(calculatedPath string) {
			fmt.Fprintf(writer, "Serving File %s", calculatedPath)
		})
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Listening on %s\n", arguments.HttpBind)
	if err := http.ListenAndServe(arguments.HttpBind, nil); err != nil {
		log.Fatal(err)
	}
}
