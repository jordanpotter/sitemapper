package mapper

import (
	"net/url"
	"testing"
)

// func TestCreatePageMap(t *testing.T) {
// 	url, err := url.Parse("https://digitalocean.com")
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
//
// 	_, err = CreatePageMap(url)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// }

func TestAbsoluteURL(t *testing.T) {
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
