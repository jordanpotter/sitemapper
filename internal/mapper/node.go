package mapper

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

type nodeType uint32

const (
	iframeNode nodeType = iota
	sourceNode
	embedNode
	objectNode
	imageNode
	anchorNode
	scriptNode
	stylesheetNode
	icoNode
	unknownNode
)

var errNodeAttrNotFound = errors.New("node attribute not found")

func getNodeType(n *html.Node) nodeType {
	switch n.Data {
	case "iframe":
		return iframeNode
	case "source":
		return sourceNode
	case "embed":
		return embedNode
	case "object":
		return objectNode
	case "img":
		return imageNode
	case "a":
		return anchorNode
	case "script":
		return scriptNode
	case "link":
		if isStylesheetNode(n) {
			return stylesheetNode
		} else if isIcoNode(n) {
			return icoNode
		}
	}
	return unknownNode
}

func isStylesheetNode(n *html.Node) bool {
	relVal, err := getNodeAttrValue(n, "rel")
	if err != nil {
		return false
	}
	return strings.ToLower(relVal) == "stylesheet"
}

func isIcoNode(n *html.Node) bool {
	relVal, err := getNodeAttrValue(n, "rel")
	if err != nil {
		return false
	}
	return strings.ToLower(relVal) == "icon"
}

func getNodeAttrValue(n *html.Node, key string) (string, error) {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val, nil
		}
	}
	return "", errNodeAttrNotFound
}
