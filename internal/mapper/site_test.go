package mapper

import (
	"net/url"
	"sync"
	"testing"
)

func TestProcessPagesInitialURL(t *testing.T) {
	u, err := url.Parse("https://foo.com")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	urls := make(chan *url.URL)
	results := make(chan *workerPageResult)
	close(results)
	_, err = processPages(u, urls, results)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	firstURL := <-urls
	if firstURL.String() != u.String() {
		t.Errorf("Expected initial url to be %q, got %q", u.String(), firstURL.String())
	}
}

func TestProcessPagesAddURLs(t *testing.T) {
	u, err := url.Parse("https://foo.com")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	circularLink, err := url.Parse("https://foo.com")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	unvisitedLink, err := url.Parse("https://foo.com/link/two")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	firstResult := workerPageResult{
		pm: &PageMap{
			URL:   u,
			Links: []*url.URL{circularLink, unvisitedLink},
		},
	}

	urls := make(chan *url.URL)
	results := make(chan *workerPageResult, 1)
	results <- &firstResult
	close(results)

	_, err = processPages(u, urls, results)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	firstURL := <-urls
	if firstURL.String() == u.String() {
		firstURL = <-urls
	}

	if firstURL.String() != unvisitedLink.String() {
		t.Errorf("Expected new url to be %q, got %q", unvisitedLink.String(), firstURL.String())
	}
}

func TestHasVisitedPage(t *testing.T) {
	testURL := func(pageURLs []string, target string, shouldBeVisited bool) {
		var pms []*PageMap
		for _, pageURL := range pageURLs {
			u, err := url.Parse(pageURL)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			pms = append(pms, &PageMap{URL: u})
		}

		targetURL, err := url.Parse(target)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		visited := hasVisitedPage(pms, new(sync.RWMutex), targetURL)
		if visited != shouldBeVisited {
			t.Errorf("Expected visited to be %t, got %t", shouldBeVisited, visited)
		}
	}

	testURL([]string{"https://foo.com", "https://bar.com"}, "https://foo.com", true)
	testURL([]string{"https://foo.com", "https://bar.com"}, "https://bar.com", true)
	testURL([]string{"https://foo.com", "https://bar.com"}, "https://baz.com", false)
}
