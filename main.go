package main

import (
	"log"
	"net/http"
)

func main() {
	arguments := NewCliArgumentsParser().GetCliArguments()

	pathMapper := NewPathMapper(arguments.RootDir)
	requestClassifier := NewRequestClassifier(arguments.Extensions)
	directoryIndex := NewDirectoryIndex(arguments.Extensions)
	fileHandler := NewFileHandler(arguments.RootDir)

	statusManager := NewStreamStatusManager(arguments.TempDir)
	streamStatusHandler := NewStreamStatusHandler(&statusManager)
	streamHandler := NewStreamHandler(&statusManager, arguments.RootDir)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("%s %s", request.Method, request.URL)

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
			streamStatusHandler.HandleStatusRequest(writer, request, mappingResult)
		case StreamPlaylistRequest:
			streamHandler.HandlePlaylistRequest(writer, request, mappingResult)
		case StreamSegmentRequest:
			streamHandler.HandleSegmentRequest(writer, request, mappingResult)
		}

	})

	http.Handle("/__static/", http.StripPrefix("/__static/", http.FileServer(http.Dir("static/"))))

	log.Printf("Listening on %s\n", arguments.HttpBind)
	if err := http.ListenAndServe(arguments.HttpBind, nil); err != nil {
		log.Fatal(err)
	}
}
