package main

import (
	"net/http"
	"net/url"
)

func RelativeRedirect(writer http.ResponseWriter, request *http.Request, relativeUrl string, httpCode int) {
	parsedUrl, err := url.Parse(relativeUrl)
	if err != nil {
		return
	}
	resolvedUrl := request.URL.ResolveReference(parsedUrl).String()
	http.Redirect(writer, request, resolvedUrl, httpCode)
}
