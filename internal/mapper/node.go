package mapper

import (
	"errors"

	"golang.org/x/net/html"
)

type nodeType uint32

const (
	scriptNode nodeType = iota
	linkNode
	iframeNode
	sourceNode
	embedNode
	objectNode
	imageNode
	anchorNode
	unknownNode
)

var nodeAttrNotFound error = errors.New("node attribute not found")

func getNodeType(n *html.Node) nodeType {
	switch n.Data {
	case "script":
		return scriptNode
	case "link":
		return linkNode
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
	default:
		return unknownNode
	}
}

func getNodeAttrValue(n *html.Node, key string) (string, error) {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val, nil
		}
	}
	return "", nodeAttrNotFound
}
