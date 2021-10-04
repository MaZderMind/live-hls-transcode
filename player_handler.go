package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type PlayerHandler struct {
	template            *template.Template
	streamStatusManager *StreamStatusManager
}

func NewPlayerHandler(streamStatusManager *StreamStatusManager) PlayerHandler {
	return PlayerHandler{
		readTemplate("player.gohtml"),
		streamStatusManager,
	}
}

func (handler *PlayerHandler) Handle(writer http.ResponseWriter, request *http.Request, mappingResult PathMappingResult) {
	streamInfo := handler.streamStatusManager.StreamInfo(mappingResult.CalculatedPath)

	isStream := request.URL.Query()["stream"] != nil
	if isStream {
		streamStatus := streamInfo.DominantStatusCode()
		switch streamStatus {
		case NoStream:
		case StreamTranscodingFailed:
		case StreamInPreparation:
			streamInfoUrl := mappingResult.UrlPath + "?stream"
			http.Redirect(writer, request, streamInfoUrl, http.StatusSeeOther)
			return
		}
	}

	writer.Header().Add("Content-Type", "text/html; charset=utf-8")

	playbackUrl := mappingResult.UrlPath
	if isStream {
		playbackUrl += "?stream&playlist"
	}

	dir, file := filepath.Split(mappingResult.UrlPath)
	if err := handler.template.Execute(writer, struct {
		Dir         string
		File        string
		Url         string
		PlaybackUrl string
	}{
		dir,
		file,
		mappingResult.UrlPath,
		playbackUrl,
	}); err != nil {
		log.Printf("Template-Formatting failed: %s", err)
	}
}
