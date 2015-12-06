package mapper

import (
	"fmt"
	"testing"
)

func TestCreatePageMap(t *testing.T) {
	pm, err := CreatePageMap("https://digitalocean.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	for _, a := range pm.Assets {
		fmt.Println(a)
	}

	fmt.Println("==========")

	for _, l := range pm.Links {
		fmt.Println(l)
	}
}

func TestAbsoluteURL(t *testing.T) {
	testURL := func(pageURL, targetURL, expectedURL string) {
		url, err := getAbsoluteURL(pageURL, targetURL)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		} else if url != expectedURL {
			t.Errorf("Exepected (%s, %s) to be %s, got %s", pageURL, targetURL, expectedURL, url)
		}
	}

	testURL("https://foo.com", "path/to/asset.png", "https://foo.com/path/to/asset.png")
	testURL("https://foo.com", "/path/to/asset.png", "https://foo.com/path/to/asset.png")
	testURL("https://foo.com", "//bar.com/path/to/asset.png", "https://bar.com/path/to/asset.png")
	testURL("https://foo.com", "http://bar.com/path/to/asset.png", "http://bar.com/path/to/asset.png")
}

func TestIsSameHost(t *testing.T) {
	testURL := func(pageURL, targetURL string, shouldMatch bool) {
		isSameHost, err := isSameHost(pageURL, targetURL)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		} else if isSameHost != shouldMatch {
			t.Errorf("Exepected (%s, %s) to be %t, got %t", pageURL, targetURL, shouldMatch, isSameHost)
		}
	}

	testURL("https://foo.com", "https://foo.com", true)
	testURL("https://foo.com", "http://foo.com", true)
	testURL("https://foo.com", "//foo.com", true)
	testURL("https://foo.com", "https://foo.com/path/to/asset.png", true)
	testURL("https://foo.com/path/to/asset/1.png", "https://foo.com/path/to/asset/2.png", true)

	testURL("https://foo.com", "https://bar.com", false)
	testURL("https://foo.com", "http://bar.com", false)
	testURL("https://foo.com", "//bar.com", false)
	testURL("https://foo.com", "https://bar.com/path/to/asset.png", false)
	testURL("https://foo.com/path/to/asset/1.png", "https://bar.com/path/to/asset/2.png", false)
}
