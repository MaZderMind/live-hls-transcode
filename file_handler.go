package main

import (
	"net/http"
)

type FileHandler struct {
	fileServerHandler http.Handler
}

func NewFileHandler(rootDir string) FileHandler {
	return FileHandler{
		http.FileServer(http.Dir(rootDir)),
	}
}

func (fileHandler *FileHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	fileHandler.fileServerHandler.ServeHTTP(writer, request)
}
