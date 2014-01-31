package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/leocassarani/gecko-codeship-status/geckoboard"
	"log"
	"net/http"
	"os"
	"strings"
)

type BuildStatus int

const (
	RedBuild BuildStatus = iota
	GreenBuild
	UnknownStatus
)

func main() {
	m := martini.Classic()

	m.Get("/builds/rag.json", func(req *http.Request, res http.ResponseWriter) []byte {
		req.ParseForm()
		projects := req.Form["projects[]"]
		builds := make(chan BuildStatus, len(projects))

		rag := &geckoboard.RAG{
			RedText:   "Failing builds",
			AmberText: "Unknown status",
			GreenText: "Green builds",
		}

		for _, projectId := range projects {
			go fetchBuildStatus(projectId, builds)
		}

		for i := 0; i < len(projects); i++ {
			status := <-builds
			switch status {
			case RedBuild:
				rag.RedValue += 1
			case UnknownStatus:
				rag.AmberValue += 1
			case GreenBuild:
				rag.GreenValue += 1
			}
		}

		json, err := json.Marshal(rag)
		if err != nil {
			log.Fatal(err)
		}

		res.Header()["Content-Type"] = []string{"application/json"}
		return json
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	http.ListenAndServe(":"+port, m)
}

func fetchBuildStatus(projectId string, builds chan BuildStatus) {
	url := fmt.Sprintf("https://www.codeship.io/projects/%s/status?branch=master", projectId)
	resp, err := http.Head(url)
	if err != nil {
		builds <- UnknownStatus
		return
	}

	contentDisposition := resp.Header["Content-Disposition"]
	if len(contentDisposition) > 0 {
		builds <- buildStatus(contentDisposition[0])
	} else {
		builds <- UnknownStatus
	}
}

func buildStatus(contentDisposition string) BuildStatus {
	if strings.Contains(contentDisposition, "status_success.png") {
		return GreenBuild
	}
	if strings.Contains(contentDisposition, "status_error.png") {
		return RedBuild
	}
	return UnknownStatus
}
