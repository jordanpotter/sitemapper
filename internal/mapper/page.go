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

func CreatePageMap(u *url.URL) (*PageMap, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	pm := &PageMap{URL: u}
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

	link, err := getNodeAttrValue(n, "href")
	if err != nil {
		return err
	}

	linkURL, err := url.Parse(link)
	if err != nil {
		return err
	}

	if !isValidLink(linkURL) {
		return nil
	} else if !isSameHost(pm.URL, linkURL) {
		return nil
	}

	linkURL, err = getAbsoluteURL(pm.URL, linkURL)
	if err != nil {
		return err
	}

	pm.Links = append(pm.Links, linkURL)
	return nil
}

func addAssetURL(pm *PageMap, n *html.Node) error {
	var asset string
	var err error
	switch getNodeType(n) {
	case scriptNode, iframeNode, sourceNode, embedNode, imageNode:
		asset, err = getNodeAttrValue(n, "src")
	case linkNode:
		asset, err = getNodeAttrValue(n, "href")
	case objectNode:
		asset, err = getNodeAttrValue(n, "data")
	default:
		return nil
	}

	if err != nil {
		return err
	}

	assetURL, err := url.Parse(asset)
	if err != nil {
		return err
	}

	assetURL, err = getAbsoluteURL(pm.URL, assetURL)
	if err != nil {
		return err
	}

	pm.Assets = append(pm.Assets, assetURL)
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

func isValidLink(linkURL *url.URL) bool {
	return linkURL.Scheme == "" || linkURL.Scheme == "http" || linkURL.Scheme == "https"
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
