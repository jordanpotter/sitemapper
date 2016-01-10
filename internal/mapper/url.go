package mapper

import (
	"net/url"
	"strings"
)

func getAbsoluteURL(pageURL, targetURL *url.URL) (*url.URL, error) {
	t, err := url.Parse(targetURL.String())
	if err != nil {
		return nil, err
	}
	return pageURL.ResolveReference(t), nil
}

func getHashlessURL(u *url.URL) (*url.URL, error) {
	hashIndex := strings.Index(u.String(), "#")
	if hashIndex >= 0 {
		return url.Parse(u.String()[0:hashIndex])
	}
	return url.Parse(u.String())
}

func isValidLink(linkURL *url.URL) bool {
	validScheme := (linkURL.Scheme == "" && linkURL.Host == "") ||
		linkURL.Scheme == "http" || linkURL.Scheme == "https"
	validExtension := strings.HasSuffix(linkURL.Path, ".html") ||
		strings.LastIndex(linkURL.Path, ".") <= strings.LastIndex(linkURL.Path, "/")
	return validScheme && validExtension
}

func getUniqueURLs(urls []*url.URL) []*url.URL {
	var unique []*url.URL
	seen := make(map[string]bool)
	for _, u := range urls {
		if !seen[u.String()] {
			seen[u.String()] = true
			unique = append(unique, u)
		}
	}
	return unique
}
