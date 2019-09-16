package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	arguments := NewCliArgumentsParser().GetCliArguments()

	pathMapper := NewPathMapper(arguments.RootDir)
	requestClassifier := NewRequestClassifier(arguments.Extensions)
	directoryIndex := NewDirectoryIndex(arguments.Extensions)
	fileHandler := NewFileHandler(arguments.RootDir)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		mappingResult := pathMapper.MapUrlPathToFilesystem(request.URL.Path)
		if mappingResult.StatError != nil {
			mappingResult.HandleError(writer)
			return
		}

		switch requestClass := requestClassifier.ClassifyRequest(request, mappingResult); requestClass {
		case DirectoryIndexRequest:
			directoryIndex.Handle(writer, request, mappingResult)
		case FileRequest:
			fileHandler.Handle(writer, request)
		case StreamStatusRequest:
			_, _ = fmt.Fprintf(writer, "Stream-Status")
		case StreamPlaylistRequest:
			_, _ = fmt.Fprintf(writer, "Stream-Playlist")
		case StreamSegmentRequest:
			_, _ = fmt.Fprintf(writer, "Stream-Segment")
		}

	})

	http.HandleFunc("/__status", func(writer http.ResponseWriter, request *http.Request) {

	})
	http.Handle("/__static/", http.StripPrefix("/__static/", http.FileServer(http.Dir("static/"))))

	fmt.Printf("Listening on %s\n", arguments.HttpBind)
	if err := http.ListenAndServe(arguments.HttpBind, nil); err != nil {
		log.Fatal(err)
	}
}
