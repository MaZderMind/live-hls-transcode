package main

import (
	"log"
	"net/http"

	"github.com/gobuffalo/packr/v2"
)

func main() {
	arguments := NewCliArgumentsParser().GetCliArguments()

	pathMapper := NewPathMapper(arguments.RootDir)
	requestClassifier := NewRequestClassifier(arguments.TranscodeExtensions, arguments.PlayerExtensions)
	directoryIndex := NewDirectoryIndex(arguments.TranscodeExtensions, arguments.PlayerExtensions)
	fileHandler := NewFileHandler(arguments.RootDir)
	playerHandler := NewPlayerHandler()

	statusManager := NewStreamStatusManager(arguments.TempDir, arguments.MinimalTranscodeDurationSeconds)
	streamStatusHandler := NewStreamStatusHandler(&statusManager, arguments.LifetimeMinutes)
	streamHandler := NewStreamHandler(&statusManager, arguments.RootDir)

	cleanup := NewCleanup(&statusManager, arguments.LifetimeMinutes)
	cleanup.Start()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if len(request.Header.Get("range")) > 0 {
			log.Printf("%s %s [%s]", request.Method, request.URL, request.Header.Get("range"))
		} else {
			log.Printf("%s %s", request.Method, request.URL)
		}

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
		case PlayerRequest:
			playerHandler.Handle(writer, request, mappingResult)
		case StreamStatusRequest:
			streamStatusHandler.HandleStatusRequest(writer, request, mappingResult)
		case StreamPlaylistRequest:
			streamHandler.HandlePlaylistRequest(writer, request, mappingResult)
		case StreamSegmentRequest:
			streamHandler.HandleSegmentRequest(writer, request, mappingResult)
		}
	})

	configureStaticCodePacks()

	log.Printf("Listening on %s\n", arguments.HttpBind())
	if err := http.ListenAndServe(arguments.HttpBind(), nil); err != nil {
		log.Fatal(err)
	}
}

func configureStaticCodePacks() {
	bootstrap := packr.New("bootstrap", "frontend/node_modules/bootstrap/dist/css")
	http.Handle("/___frontend/bootstrap/", http.StripPrefix("/___frontend/bootstrap/", http.FileServer(bootstrap)))

	jquery := packr.New("jquery", "frontend/node_modules/jquery/dist")
	http.Handle("/___frontend/jquery/", http.StripPrefix("/___frontend/jquery/", http.FileServer(jquery)))

	videojs := packr.New("videojs", "frontend/node_modules/video.js/dist")
	http.Handle("/___frontend/video.js/", http.StripPrefix("/___frontend/video.js/", http.FileServer(videojs)))

	fontAwesomeCss := packr.New("fontAwesomeCss", "frontend/node_modules/@fortawesome/fontawesome-free/css")
	http.Handle("/___frontend/font-awesome/css/", http.StripPrefix("/___frontend/font-awesome/css", http.FileServer(fontAwesomeCss)))

	fontAwesomeWebfonts := packr.New("fontAwesomeWebfonts", "frontend/node_modules/@fortawesome/fontawesome-free/webfonts")
	http.Handle("/___frontend/font-awesome/webfonts/", http.StripPrefix("/___frontend/font-awesome/webfonts", http.FileServer(fontAwesomeWebfonts)))

	frontend := packr.New("frontend", "frontend/code")
	http.Handle("/___frontend/", http.StripPrefix("/___frontend/", http.FileServer(frontend)))
}
