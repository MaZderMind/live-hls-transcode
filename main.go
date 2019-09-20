package main

import (
	"github.com/gobuffalo/packr/v2"
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
	streamStatusHandler := NewStreamStatusHandler(&statusManager, arguments.LifetimeMinutes)
	streamHandler := NewStreamHandler(&statusManager, arguments.RootDir)

	cleanup := NewCleanup(&statusManager, arguments.LifetimeMinutes)
	cleanup.Start()

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

	bootstrap := packr.New("bootstrap", "./frontend/node_modules/bootstrap/dist/css/")
	http.Handle("/___frontend/bootstrap/", http.StripPrefix("/___frontend/bootstrap/", http.FileServer(bootstrap)))

	jquery := packr.New("jquery", "frontend/node_modules/jquery/dist")
	http.Handle("/___frontend/jquery/", http.StripPrefix("/___frontend/jquery/", http.FileServer(jquery)))

	frontend := packr.New("frontend", "frontend/code")
	http.Handle("/___frontend/", http.StripPrefix("/___frontend/", http.FileServer(frontend)))

	log.Printf("Listening on %s\n", arguments.HttpBind)
	if err := http.ListenAndServe(arguments.HttpBind, nil); err != nil {
		log.Fatal(err)
	}
}
