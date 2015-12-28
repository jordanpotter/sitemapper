package mapper

import (
	"fmt"
	"net/url"
	"testing"

	"golang.org/x/net/html"
)

func TestProcessNode(t *testing.T) {
	urlStr := "https://foo.com"
	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assetURLStr2 := "image2.png"
	child3 := html.Node{
		Type: html.ElementNode,
		Data: "img",
		Attr: []html.Attribute{
			html.Attribute{Key: "src", Val: assetURLStr2},
		},
	}

	assetURLStr1 := "image1.png"
	child2 := html.Node{
		Type: html.ElementNode,
		Data: "img",
		Attr: []html.Attribute{
			html.Attribute{Key: "src", Val: assetURLStr1},
		},
		NextSibling: &child3,
	}

	linkURLStr := "path/to/site"
	child1 := html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			html.Attribute{Key: "href", Val: linkURLStr},
		},
		NextSibling: &child2,
	}

	root := html.Node{
		Type: html.ElementNode,
		Data: "div",
		Attr: []html.Attribute{
			html.Attribute{Key: "class", Val: "root"},
		},
		FirstChild: &child1,
	}

	pm := PageMap{URL: u}
	err = processNode(&pm, &root)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(pm.Links) != 1 {
		t.Errorf("Exepected number links to be 1, got %d", len(pm.Links))
	} else if len(pm.Assets) != 2 {
		t.Errorf("Exepected number assets to be 2, got %d", len(pm.Assets))
	}

	expectedLinkStr := fmt.Sprintf("%s/%s", urlStr, linkURLStr)
	if expectedLinkStr != pm.Links[0].String() {
		t.Errorf("Exepected link url to be %q, got %q", expectedLinkStr, pm.Links[0].String())
	}

	expectedAssetStr1 := fmt.Sprintf("%s/%s", urlStr, assetURLStr1)
	if expectedAssetStr1 != pm.Assets[0].String() {
		t.Errorf("Exepected asset url to be %q, got %q", expectedAssetStr1, pm.Assets[0].String())
	}

	expectedAssetStr2 := fmt.Sprintf("%s/%s", urlStr, assetURLStr2)
	if expectedAssetStr2 != pm.Assets[1].String() {
		t.Errorf("Exepected asset url to be %q, got %q", expectedAssetStr2, pm.Assets[1].String())
	}
}

func TestAddLinkURL(t *testing.T) {
	urlStr := "https://foo.com"
	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	linkURLStr := "path/to/site"
	n := html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			html.Attribute{Key: "href", Val: linkURLStr},
		},
	}

	pm := PageMap{URL: u}
	err = addLinkURL(&pm, &n)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(pm.Links) != 1 {
		t.Errorf("Exepected number links to be 1, got %d", len(pm.Links))
	} else if len(pm.Assets) != 0 {
		t.Errorf("Exepected number assets to be 0, got %d", len(pm.Assets))
	}

	expectedLinkStr := fmt.Sprintf("%s/%s", urlStr, linkURLStr)
	if expectedLinkStr != pm.Links[0].String() {
		t.Errorf("Exepected link url to be %q, got %q", expectedLinkStr, pm.Links[0].String())
	}
}

func TestAddAssetURL(t *testing.T) {
	urlStr := "https://foo.com"
	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assetURLStr := "image.png"
	n := html.Node{
		Type: html.ElementNode,
		Data: "img",
		Attr: []html.Attribute{
			html.Attribute{Key: "src", Val: assetURLStr},
		},
	}

	pm := PageMap{URL: u}
	err = addAssetURL(&pm, &n)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(pm.Links) != 0 {
		t.Errorf("Exepected number links to be 0, got %d", len(pm.Links))
	} else if len(pm.Assets) != 1 {
		t.Errorf("Exepected number assets to be 1, got %d", len(pm.Assets))
	}

	expectedAssetStr := fmt.Sprintf("%s/%s", urlStr, assetURLStr)
	if expectedAssetStr != pm.Assets[0].String() {
		t.Errorf("Exepected asset url to be %q, got %q", expectedAssetStr, pm.Assets[0].String())
	}
}

func TestGetAbsoluteURL(t *testing.T) {
	testURL := func(pageURLStr, targetURLStr, expectedURLStr string) {
		pageURL, err := url.Parse(pageURLStr)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		targetURL, err := url.Parse(targetURLStr)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		url, err := getAbsoluteURL(pageURL, targetURL)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		} else if url.String() != expectedURLStr {
			t.Errorf("Exepected (%s, %s) to be %s, got %s", pageURL, targetURL, expectedURLStr, url)
		}
	}

	testURL("https://foo.com", "path/to/asset.png", "https://foo.com/path/to/asset.png")
	testURL("https://foo.com", "/path/to/asset.png", "https://foo.com/path/to/asset.png")
	testURL("https://foo.com", "//bar.com/path/to/asset.png", "https://bar.com/path/to/asset.png")
	testURL("https://foo.com", "http://bar.com/path/to/asset.png", "http://bar.com/path/to/asset.png")
}

func TestIsValidLink(t *testing.T) {
	testURL := func(urlStr string, shouldBeValid bool) {
		u, err := url.Parse(urlStr)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		valid := isValidLink(u)
		if valid != shouldBeValid {
			t.Errorf("Exepected %q to be %t, got %t", u.String(), shouldBeValid, valid)
		}
	}

	testURL("path/to/file", true)
	testURL("path/to/file.html", true)

	testURL("http://foo.com", true)
	testURL("https://foo.com", true)

	testURL("ftp://foo.com", false)
	testURL("smtp://foo.com", false)

	testURL("https://foo.com/path/to/file", true)
	testURL("https://foo.com/path/to/file.html", true)

	testURL("https://foo.com/path/to/file.png", false)
	testURL("https://foo.com/path/to/file.html2", false)
}

func TestIsSameHost(t *testing.T) {
	testURL := func(pageURLStr, targetURLStr string, shouldMatch bool) {
		pageURL, err := url.Parse(pageURLStr)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		targetURL, err := url.Parse(targetURLStr)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
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

func TestGetUniqueURLs(t *testing.T) {
	createURL := func(str string) *url.URL {
		u, err := url.Parse(str)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		return u
	}

	testUnique := func(urls []*url.URL, expectedLen int) {
		unique := getUniqueURLs(urls)
		if len(unique) != expectedLen {
			t.Errorf("Expected length to be %d, got %d", expectedLen, len(unique))
		}

		seen := make(map[string]bool)
		for _, u := range unique {
			if seen[u.String()] {
				t.Errorf("Duplicate url %q", u.String)
			} else {
				seen[u.String()] = true
			}
		}
	}

	u1 := createURL("http://foo.com")
	u2 := createURL("http://bar.com")
	u3 := createURL("http://baz.com")

	testUnique([]*url.URL{}, 0)
	testUnique([]*url.URL{u1}, 1)
	testUnique([]*url.URL{u1, u2}, 2)
	testUnique([]*url.URL{u1, u2, u3}, 3)
	testUnique([]*url.URL{u1, u2, u3, u1, u2, u2, u1, u3}, 3)
}
