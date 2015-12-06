package mapper

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type PageMap struct {
	URL    *url.URL
	Links  []*url.URL
	Assets []*url.URL
}

func CreatePageMap(url *url.URL) (*PageMap, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	pm := &PageMap{URL: url}
	processNode(pm, root)
	pm.Links = getUniqueURLs(pm.Links)
	pm.Assets = getUniqueURLs(pm.Assets)
	return pm, nil
}

func processNode(pm *PageMap, n *html.Node) error {
	if n.Type == html.ElementNode {
		err := addLinkURL(pm, n)
		if err != nil {
			return err
		}

		err = addAssetURL(pm, n)
		if err != nil {
			return err
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		processNode(pm, child)
	}
	return nil
}

func addLinkURL(pm *PageMap, n *html.Node) error {
	if getNodeType(n) != anchorNode {
		return nil
	}

	urlStr, err := getNodeAttrValue(n, "href")
	if err != nil {
		return err
	}

	url, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	if !isSameHost(pm.URL, url) {
		return nil
	}

	absURL, err := getAbsoluteURL(pm.URL, url)
	if err != nil {
		return err
	}

	pm.Links = append(pm.Links, absURL)
	return nil
}

func addAssetURL(pm *PageMap, n *html.Node) error {
	var urlStr string
	var err error
	switch getNodeType(n) {
	case scriptNode, iframeNode, sourceNode, embedNode, imageNode:
		urlStr, err = getNodeAttrValue(n, "src")
	case linkNode:
		urlStr, err = getNodeAttrValue(n, "href")
	case objectNode:
		urlStr, err = getNodeAttrValue(n, "data")
	default:
		return nil
	}

	if err != nil {
		return err
	}

	url, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	absURL, err := getAbsoluteURL(pm.URL, url)
	if err != nil {
		return err
	}

	pm.Assets = append(pm.Assets, absURL)
	return nil
}

func getAbsoluteURL(pageURL, targetURL *url.URL) (*url.URL, error) {
	absURL, err := url.Parse(targetURL.String())
	if err != nil {
		return nil, err
	}

	if absURL.Scheme == "" {
		absURL.Scheme = pageURL.Scheme
	}
	if absURL.Host == "" {
		absURL.Host = pageURL.Host
	}
	return absURL, nil
}

func isSameHost(pageURL, targetURL *url.URL) bool {
	return targetURL.Host == "" || pageURL.Host == targetURL.Host
}

func getUniqueURLs(urls []*url.URL) []*url.URL {
	var unique []*url.URL
	seen := make(map[string]bool)
	for _, url := range urls {
		if !seen[url.String()] {
			seen[url.String()] = true
			unique = append(unique, url)
		}
	}
	return unique
}
