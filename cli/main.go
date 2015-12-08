package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"runtime"

	"github.com/jordanpotter/sitemapper/internal/mapper"
)

func main() {
	site := flag.String("site", "https://digitalocean.com", "entry point into site to scan")
	numWorkers := flag.Int("workers", runtime.NumCPU(), "number of workers")
	filename := flag.String("file", "sitemap.json", "file to write to")
	flag.Parse()

	siteURL, err := url.Parse(*site)
	if err != nil {
		log.Fatalln(err)
	}

	sm, err := mapper.CreateSiteMap(siteURL, *numWorkers)
	if err != nil {
		log.Fatalln(err)
	}

	b, err := json.Marshal(sm)
	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile(*filename, b, 400)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Site map written to %s", *filename)
}
