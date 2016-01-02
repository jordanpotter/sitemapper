package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jordanpotter/sitemapper/internal/mapper"
)

func main() {
	port := flag.Int("port", 8000, "port to serve the API")
	staticPath := flag.String("static", "gui", "path to static files to serve")
	flag.Parse()

	log.Printf("Starting server on port %d", *port)
	http.HandleFunc("/sitemap", getSiteMap)
	http.Handle("/", http.FileServer(http.Dir(*staticPath)))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func getSiteMap(w http.ResponseWriter, r *http.Request) {
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
