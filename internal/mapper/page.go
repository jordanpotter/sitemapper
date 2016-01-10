package mapper

import (
	"encoding/json"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

// A PageMap contains all of the links and assets at URL.
type PageMap struct {
	URL    *url.URL
	Links  []*url.URL
	Assets []*url.URL
}

func (pm *PageMap) MarshalJSON() ([]byte, error) {
	urlsToStrings := func(urls []*url.URL) []string {
		strs := make([]string, 0, len(urls))
		for _, u := range urls {
			strs = append(strs, u.String())
		}
		return strs
	}

	return json.Marshal(struct {
		URL    string   `json:"url"`
		Links  []string `json:"links"`
		Assets []string `json:"assets"`
	}{
		URL:    pm.URL.String(),
		Links:  urlsToStrings(pm.Links),
		Assets: urlsToStrings(pm.Assets),
	})
}

// CreatePageMap creates a page map for the specified url. This is done by
// parsing the HTML for all links and assets found in the DOM tree.
func CreatePageMap(u *url.URL) (*PageMap, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

	linkURL, err = getHashlessURL(linkURL)
	if err != nil {
		return err
	}

	if !isValidLink(linkURL) {
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
	case stylesheetNode, icoNode:
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
