package mapper

import (
	"log"
	"net/url"
	"sync"
)

type SiteMap struct {
	PageMaps []*PageMap `json:"pages"`
}

type workerPageResult struct {
	pm  *PageMap
	err error
}

func CreateSiteMap(u *url.URL, numWorkers int) (*SiteMap, error) {
	log.Printf("Creating site map for %q with %d workers...", u, numWorkers)
	urls := make(chan *url.URL)
	results := createWorkers(numWorkers, urls)
	pms, err := processPages(u, urls, results)
	return &SiteMap{pms}, err
}

func createWorkers(num int, urls <-chan *url.URL) <-chan *workerPageResult {
	var wg sync.WaitGroup
	results := make(chan *workerPageResult)

	wg.Add(num)
	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 0; i < num; i++ {
		go func() {
			for u := range urls {
				pm, err := CreatePageMap(u)
				results <- &workerPageResult{pm, err}
			}
			wg.Done()
		}()
	}
	return results
}

func processPages(initialURL *url.URL, urls chan<- *url.URL, results <-chan *workerPageResult) ([]*PageMap, error) {
	var pms []*PageMap
	var m sync.RWMutex
	var wg sync.WaitGroup

	go func() {
		urls <- initialURL
		wg.Add(1)
		wg.Wait()
		close(urls)
	}()

	for wr := range results {
		if wr.err != nil {
			return nil, wr.err
		}

		if hasVisitedPage(pms, &m, wr.pm.URL) {
			wg.Done()
			continue
		}

		log.Printf("Processed %s", wr.pm.URL)

		m.Lock()
		pms = append(pms, wr.pm)
		m.Unlock()

		wg.Add(len(wr.pm.Links) - 1)
		go func(links []*url.URL) {
			for _, link := range links {
				if hasVisitedPage(pms, &m, link) {
					wg.Done()
				} else {
					urls <- link
				}
			}
		}(wr.pm.Links)
	}
	return pms, nil
}

func hasVisitedPage(pms []*PageMap, m *sync.RWMutex, u *url.URL) bool {
	m.RLock()
	defer m.RUnlock()

	for _, pm := range pms {
		if pm.URL.String() == u.String() {
			return true
		}
	}
	return false
}
