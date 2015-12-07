package mapper

import (
	"net/url"
	"testing"
)

func TestCreatePageMap(t *testing.T) {
	url, err := url.Parse("https://digitalocean.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	_, err = CreatePageMap(url)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestAbsoluteURL(t *testing.T) {
	testURL := func(pageURLStr, targetURLStr, expectedURLStr string) {
		pageURL, err := url.Parse(pageURLStr)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		targetURL, err := url.Parse(targetURLStr)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		url, err := getAbsoluteURL(pageURL, targetURL)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		} else if url.String() != expectedURLStr {
			t.Errorf("Exepected (%s, %s) to be %s, got %s", pageURL, targetURL, expectedURLStr, url)
		}
	}

	testURL("https://foo.com", "path/to/asset.png", "https://foo.com/path/to/asset.png")
	testURL("https://foo.com", "/path/to/asset.png", "https://foo.com/path/to/asset.png")
	testURL("https://foo.com", "//bar.com/path/to/asset.png", "https://bar.com/path/to/asset.png")
	testURL("https://foo.com", "http://bar.com/path/to/asset.png", "http://bar.com/path/to/asset.png")
}

func TestIsSameHost(t *testing.T) {
	testURL := func(pageURLStr, targetURLStr string, shouldMatch bool) {
		pageURL, err := url.Parse(pageURLStr)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		targetURL, err := url.Parse(targetURLStr)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if isSameHost(pageURL, targetURL) != shouldMatch {
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
