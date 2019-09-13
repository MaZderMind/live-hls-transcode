package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type DirectoryIndex struct {
	rootDir  string
	template *template.Template
}

func NewDirectoryIndex(rootDir string) DirectoryIndex {
	templateFile, err := template.New("directory-index.html").ParseFiles("templates/directory-index.html")

	if err != nil {
		log.Fatal(err)
	}
	return DirectoryIndex{
		rootDir,
		templateFile,
	}
}

func (i DirectoryIndex) Handle(writer http.ResponseWriter, request *http.Request) {
	calculatedPath := i.calculatePath(request.URL.Path)

	info, err := os.Stat(calculatedPath);
	if os.IsNotExist(err) {
		writer.WriteHeader(404)
		return
	} else if os.IsPermission(err) {
		writer.WriteHeader(401)
		return
	} else if !info.IsDir() {
		writer.Header().Add("Content-Type", "text/plain")
		if _, err := fmt.Fprintf(writer, "File-Content"); err != nil {
			fmt.Printf("Serving File-Content failed: %s", err)
		}
	}

	files, err := ioutil.ReadDir(calculatedPath)
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Add("Content-Type", "text/html")
	if err = i.template.Execute(writer, struct {
		Files []os.FileInfo
	}{
		files,
	}); err != nil {
		fmt.Printf("Template-Formatting failed: %s", err)
	}
}

func (i DirectoryIndex) calculatePath(urlPath string) string {
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}

	cleanPath := path.Clean(urlPath)
	return path.Join(i.rootDir, cleanPath)
}
