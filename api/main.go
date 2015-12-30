package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jordanpotter/sitemapper/internal/mapper"
)

const port = 8000

func main() {
	log.Printf("Starting server on port %d", port)
	http.HandleFunc("/sitemap", getSiteMap)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func getSiteMap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, siteProvided := r.URL.Query()["site"]
	if !siteProvided {
		http.Error(w, "Missing query parameter \"site\"", 412)
		return
	}

	_, workersProvided := r.URL.Query()["workers"]
	if !workersProvided {
		http.Error(w, "Missing query parameter \"workers\"", 412)
		return
	}

	siteStr := r.URL.Query()["site"][0]
	u, err := url.Parse(siteStr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	numWorkersStr := r.URL.Query()["workers"][0]
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	sm, err := mapper.CreateSiteMap(u, numWorkers)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	b, err := json.Marshal(sm)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
}
