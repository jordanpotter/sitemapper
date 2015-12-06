package mapper

import (
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type PageMap struct {
	Links  []string
	Assets []string
}

func CreatePageMap(url string) (*PageMap, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	links := make(map[string]bool)
	assets := make(map[string]bool)
	processNode(root, url, links, assets)

	pm := new(PageMap)
	for url := range links {
		pm.Links = append(pm.Links, url)
	}
	for url := range assets {
		pm.Assets = append(pm.Assets, url)
	}
	return pm, nil
}

func processNode(n *html.Node, pageURL string, links, assets map[string]bool) error {
	if n.Type == html.ElementNode {
		err := addLinkURL(n, pageURL, links)
		if err != nil {
			return err
		}

		err = addAssetURL(n, pageURL, assets)
		if err != nil {
			return err
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		processNode(child, pageURL, links, assets)
	}
	return nil
}

func addLinkURL(n *html.Node, pageURL string, links map[string]bool) error {
	if getNodeType(n) != anchorNode {
		return nil
	}

	url, err := getNodeAttrValue(n, "href")
	if err != nil {
		return err
	}

	isSameHost, err := isSameHost(pageURL, url)
	if err != nil {
		return err
	} else if !isSameHost {
		return nil
	}

	absURL, err := getAbsoluteURL(pageURL, url)
	if err != nil {
		return err
	}

	links[absURL] = true
	return nil
}

func addAssetURL(n *html.Node, pageURL string, assets map[string]bool) error {
	var url string
	var err error
	switch getNodeType(n) {
	case scriptNode, iframeNode, sourceNode, embedNode, imageNode:
		url, err = getNodeAttrValue(n, "src")
	case linkNode:
		url, err = getNodeAttrValue(n, "href")
	case objectNode:
		url, err = getNodeAttrValue(n, "data")
	default:
		return nil
	}

	if err != nil {
		return err
	}

	absURL, err := getAbsoluteURL(pageURL, url)
	if err != nil {
		return err
	}

	assets[absURL] = true
	return nil
}

func getAbsoluteURL(pageURL, targetURL string) (string, error) {
	parsedPageURL, err := url.Parse(pageURL)
	if err != nil {
		return "", err
	}

	parsedTargetURL, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}

	if parsedTargetURL.Scheme == "" {
		parsedTargetURL.Scheme = parsedPageURL.Scheme
	}
	if parsedTargetURL.Host == "" {
		parsedTargetURL.Host = parsedPageURL.Host
	}
	return parsedTargetURL.String(), nil
}

func isSameHost(pageURL, targetURL string) (bool, error) {
	if strings.HasPrefix(targetURL, "/") {
		return true, nil
	}

	parsedPageURL, err := url.Parse(pageURL)
	if err != nil {
		return false, err
	}

	parsedTargetURL, err := url.Parse(targetURL)
	if err != nil {
		return false, err
	}

	return parsedPageURL.Host == parsedTargetURL.Host, nil
}
