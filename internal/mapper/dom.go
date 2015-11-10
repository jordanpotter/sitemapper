package mapper

import (
	"errors"

	"golang.org/x/net/html"
)

type elementType uint32

const (
	linkElement elementType = iota
	imageElement
	unknownElement
)

var elementAttrNotFound error = errors.New("element attribute not found")

func getElementType(n *html.Node) elementType {
	switch n.Data {
	case "a":
		return linkElement
	case "img":
		return imageElement
	default:
		return unknownElement
	}
}

func getElementAttrValue(n *html.Node, key string) (string, error) {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val, nil
		}
	}
	return "", elementAttrNotFound
}
