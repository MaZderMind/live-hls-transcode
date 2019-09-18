package main

import (
	"bufio"
	"fmt"
	"github.com/grafov/m3u8"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

type StreamHandler struct {
	streamStatusManager *StreamStatusManager
}

func NewStreamHandler(streamStatusManager *StreamStatusManager, rootDir string) StreamHandler {
	return StreamHandler{
		streamStatusManager,
	}
}

func (handler StreamHandler) HandlePlaylistRequest(writer http.ResponseWriter, request *http.Request, pathMappingResult PathMappingResult) {
	streamStatus := handler.streamStatusManager.GetStreamStatus(pathMappingResult.CalculatedPath)

	if ! handler.assureStreamIsReady(streamStatus, writer) {
		return
	}

	tempdir := handler.streamStatusManager.GetStreamTempdir(pathMappingResult.CalculatedPath)
	filepath := path.Join(tempdir, "index.m3u8")
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("Error loading Playlist-File %s: %s", filepath, err)
		return
	}

	genericPlaylist, listType, err := m3u8.DecodeFrom(bufio.NewReader(file), true)
	if err != nil {
		log.Printf("Error decoding Playlist-File %s: %s", filepath, err)
		return
	}

	if listType != m3u8.MEDIA {
		log.Printf("Playlist-File %s is not a MEDIA-Playlist", filepath)
		return
	}

	// Set Start-Time
	playlist := genericPlaylist.(*m3u8.MediaPlaylist)
	playlist.StartTime = 0.01

	// Modify Segment-Names
	for index := range playlist.Segments {
		if playlist.Segments[index] == nil {
			break
		}

		playlist.Segments[index].URI = request.URL.Path + "?stream&segment=" + url.QueryEscape(playlist.Segments[index].URI)
	}

	writer.Header().Add("Content-Type", "application/vnd.apple.mpegurl; charset=utf-8")
	_, err = fmt.Fprint(writer, playlist.String())
	if err != nil {
		log.Printf("Error writing to Socket: %s", err)
		return
	}
}

func (handler StreamHandler) assureStreamIsReady(streamStatus StreamStatus, writer http.ResponseWriter) bool {
	if streamStatus != StreamReady && streamStatus != TranscodingFinished {
		writer.Header().Add("Content-Type", "text/plain")

		_, err := fmt.Fprint(writer, "Stream not Ready")
		if err != nil {
			log.Printf("Error writing to Socket: %s", err)
			return false
		}

		return false
	}

	return true
}

func (handler StreamHandler) HandleSegmentRequest(writer http.ResponseWriter, request *http.Request, pathMappingResult PathMappingResult) {
	streamStatus := handler.streamStatusManager.GetStreamStatus(pathMappingResult.CalculatedPath)

	if ! handler.assureStreamIsReady(streamStatus, writer) {
		return
	}

	tempdir := handler.streamStatusManager.GetStreamTempdir(pathMappingResult.CalculatedPath)
	segmentFilename := request.URL.Query().Get("segment")

	segmentRequest := request
	segmentRequest.URL.Path = segmentFilename

	writer.Header().Add("Content-Type", "video/MP2T")
	writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", segmentFilename))
	fileServer := http.FileServer(http.Dir(tempdir))
	fileServer.ServeHTTP(writer, segmentRequest)
}
