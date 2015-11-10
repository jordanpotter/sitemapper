package mapper

import (
	"fmt"
	"net/http"

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

	pm := new(PageMap)
	err = processDOMElement(root, pm)
	return pm, err
}

func processDOMElement(n *html.Node, pm *PageMap) error {
	if n.Type == html.ElementNode {
		switch getElementType(n) {
		case linkElement:
			link, err := getElementAttrValue(n, "href")
			fmt.Println(link, err)
		case imageElement:
			src, err := getElementAttrValue(n, "src")
			fmt.Println(src, err)
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		processDOMElement(child, pm)
	}
	return nil
}
