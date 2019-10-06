package main

import (
	"github.com/gobuffalo/packr/v2"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type StreamStatusHandler struct {
	statusPageTemplateFile *template.Template
	streamStatusManager    *StreamStatusManager
	lifetimeMinutes        uint32
}

func NewStreamStatusHandler(streamStatusManager *StreamStatusManager, lifetimeMinutes uint32) StreamStatusHandler {
	templates := packr.New("templates", "./templates")
	templateString, err := templates.FindString("status-page.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	statusPageTemplateFile, err := template.New("status-page.gohtml").Parse(templateString)

	if err != nil {
		log.Fatal(err)
	}

	return StreamStatusHandler{
		statusPageTemplateFile,
		streamStatusManager,
		lifetimeMinutes,
	}
}

func (handler *StreamStatusHandler) HandleStatusRequest(writer http.ResponseWriter, request *http.Request, pathMappingResult PathMappingResult) {
	if request.URL.Query()["start"] != nil {
		handler.streamStatusManager.StartStream(pathMappingResult.CalculatedPath, pathMappingResult.UrlPath)
		RelativeRedirect(writer, request, "?stream&autoplay", http.StatusTemporaryRedirect)
		return
	} else if request.URL.Query()["stop"] != nil {
		handler.streamStatusManager.StopStream(pathMappingResult.CalculatedPath)
		RelativeRedirect(writer, request, "?stream", http.StatusTemporaryRedirect)
		return
	}

	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	streamInfo := handler.streamStatusManager.StreamInfo(pathMappingResult.CalculatedPath)

	streamStatus := streamInfo.DominantStatusCode()

	otherRunningTranscoders := handler.streamStatusManager.OtherRunningTranscoders(pathMappingResult.CalculatedPath)

	dir, file := filepath.Split(pathMappingResult.UrlPath)
	if err := handler.statusPageTemplateFile.Execute(writer, struct {
		Dir  string
		File string

		LastAccess     time.Time
		ExpirationDate time.Time

		ProcessedDuration time.Duration
		TotalDuration     time.Duration
		ProcessedPercent  float64

		NoStream                bool
		StreamInPreparation     bool
		StreamReady             bool
		StreamTranscodingFailed bool
		TranscodingFinished     bool

		ShowProgress   bool
		ShowValidity   bool
		AutoplayActive bool

		OtherRunningTranscoders []StreamInfo
	}{
		dir,
		file,

		streamInfo.LastAccess,
		streamInfo.LastAccess.Add(time.Minute * time.Duration(handler.lifetimeMinutes)),

		streamInfo.ProcessedDuration(),
		streamInfo.TotalDuration(),
		streamInfo.ProcessedPercent(),

		streamStatus == NoStream,
		streamStatus == StreamInPreparation,
		streamStatus == StreamReady,
		streamStatus == StreamTranscodingFailed,
		streamStatus == TranscodingFinished,

		streamStatus == StreamInPreparation || streamStatus == StreamReady,
		streamStatus == TranscodingFinished,
		request.URL.Query()["autoplay"] != nil,

		otherRunningTranscoders,
	}); err != nil {
		log.Printf("Template-Formatting failed: %s", err)
	}
}
