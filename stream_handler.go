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

func (handler *StreamHandler) HandlePlaylistRequest(writer http.ResponseWriter, request *http.Request, pathMappingResult PathMappingResult) {
	handler.streamStatusManager.UpdateLastAccess(pathMappingResult.CalculatedPath)
	streamInfo := handler.streamStatusManager.StreamInfo(pathMappingResult.CalculatedPath)

	if ! handler.ensureStreamIsReady(streamInfo, writer) {
		return
	}

	filepath := path.Join(streamInfo.TempDir, "index.m3u8")
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
	playlist.StartTime = 0.01 // must be >0.0 to make the m3u8 writer print the statement

	// Modify Segment-Names
	for index := range playlist.Segments {
		if playlist.Segments[index] == nil {
			break
		}

		playlist.Segments[index].URI = request.URL.EscapedPath() + "?stream&segment=" + url.QueryEscape(playlist.Segments[index].URI)
	}

	writer.Header().Add("Content-Type", "application/vnd.apple.mpegurl; charset=utf-8")

	_, err = fmt.Fprint(writer, playlist.String())
	if err != nil {
		log.Printf("Error writing to Socket: %s", err)
		return
	}
}

func (handler *StreamHandler) HandleSegmentRequest(writer http.ResponseWriter, request *http.Request, pathMappingResult PathMappingResult) {
	handler.streamStatusManager.UpdateLastAccess(pathMappingResult.CalculatedPath)
	streamInfo := handler.streamStatusManager.StreamInfo(pathMappingResult.CalculatedPath)

	if ! handler.ensureStreamIsReady(streamInfo, writer) {
		return
	}

	segmentFilename := request.URL.Query().Get("segment")

	segmentRequest := request
	segmentRequest.URL.Path = segmentFilename

	writer.Header().Add("Content-Type", "video/MP2T")
	fileServer := http.FileServer(http.Dir(streamInfo.TempDir))
	fileServer.ServeHTTP(writer, segmentRequest)
}

func (handler *StreamHandler) ensureStreamIsReady(streamInfo StreamInfo, writer http.ResponseWriter) bool {
	if ! streamInfo.IsReady() {
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
