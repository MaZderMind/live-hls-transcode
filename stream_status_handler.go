package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type StreamStatusHandler struct {
	statusPageTemplateFile *template.Template
	streamStatusManager    *StreamStatusManager
	lifetimeMinutes        int32
}

func NewStreamStatusHandler(streamStatusManager *StreamStatusManager, lifetimeMinutes int32) StreamStatusHandler {
	statusPageTemplateFile, err := template.New("status-page.gohtml").ParseFiles("templates/status-page.gohtml")

	if err != nil {
		log.Fatal(err)
	}

	return StreamStatusHandler{
		statusPageTemplateFile,
		streamStatusManager,
		lifetimeMinutes,
	}
}

func (handler StreamStatusHandler) HandleStatusRequest(writer http.ResponseWriter, request *http.Request, pathMappingResult PathMappingResult) {
	if request.URL.Query()["start"] != nil {
		handler.streamStatusManager.StartStream(pathMappingResult.CalculatedPath, pathMappingResult.UrlPath)
		RelativeRedirect(writer, request, "?stream&autostart", http.StatusTemporaryRedirect)
		return
	} else if request.URL.Query()["stop"] != nil {
		handler.streamStatusManager.StopStream(pathMappingResult.CalculatedPath)
		RelativeRedirect(writer, request, "?stream", http.StatusTemporaryRedirect)
		return
	}

	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	streamInfo := handler.streamStatusManager.StreamInfo(pathMappingResult.CalculatedPath)

	streamStatus := streamInfo.DominantStatusCode()

	otherStreamInfos := handler.streamStatusManager.OtherStreamInfos(pathMappingResult.CalculatedPath)

	if err := handler.statusPageTemplateFile.Execute(writer, struct {
		UrlPath                 string
		LastAccess              time.Time
		ExpirationDate          time.Time
		NoStream                bool
		StreamInPreparation     bool
		StreamReady             bool
		StreamTranscodingFailed bool
		TranscodingFinished     bool
		OtherStreamInfos        []StreamInfo
	}{
		pathMappingResult.UrlPath,
		streamInfo.LastAccess,
		streamInfo.LastAccess.Add(time.Minute * time.Duration(handler.lifetimeMinutes)),
		streamStatus == NoStream,
		streamStatus == StreamInPreparation,
		streamStatus == StreamReady,
		streamStatus == StreamTranscodingFailed,
		streamStatus == TranscodingFinished,
		otherStreamInfos,
	}); err != nil {
		log.Printf("Template-Formatting failed: %s", err)
	}
}
