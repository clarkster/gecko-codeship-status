package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"strings"
)

func main() {
	m := martini.Classic()

	m.Get("/builds/rag.json", func(req *http.Request) string {
		req.ParseForm()
		projects := req.Form["projects[]"]
		responses := make(chan bool, len(projects))
		for _, projectId := range projects {
			go fetchBuildStatus(projectId, responses)
		}
		for i := 0; i < len(projects); i++ {
			<-responses
		}
		return "{}"
	})

	http.ListenAndServe(":9999", m)
}

func fetchBuildStatus(projectId string, responses chan bool) {
	url := fmt.Sprintf("https://www.codeship.io/projects/%s/status?branch=master", projectId)
	resp, err := http.Head(url)
	if err != nil {
		// FIXME
		responses <- false
	} else {
		contentHeader := resp.Header["Content-Disposition"][0]
		success := strings.Contains(contentHeader, "status_success.png")
		responses <- success
	}
}
