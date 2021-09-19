package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type PlayerHandler struct {
	template *template.Template
}

func NewPlayerHandler() PlayerHandler {
	return PlayerHandler{
		readTemplate("player.gohtml"),
	}
}

func (playerHandler *PlayerHandler) Handle(writer http.ResponseWriter, request *http.Request, mappingResult PathMappingResult) {
	writer.Header().Add("Content-Type", "text/html; charset=utf-8")

	dir, file := filepath.Split(mappingResult.UrlPath)
	if err := playerHandler.template.Execute(writer, struct {
		Dir  string
		File string
	}{
		dir,
		file,
	}); err != nil {
		log.Printf("Template-Formatting failed: %s", err)
	}
}
