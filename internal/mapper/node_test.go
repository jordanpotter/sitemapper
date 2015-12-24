package mapper

import (
	"testing"

	"golang.org/x/net/html"
)

func TestGetNodeType(t *testing.T) {
	testType := func(data string, et nodeType) {
		var n html.Node
		n.Data = data
		nodeType := getNodeType(&n)
		if nodeType != et {
			t.Errorf("Exepected %d, got %d", et, nodeType)
		}
	}

	testType("script", scriptNode)
	testType("link", linkNode)
	testType("iframe", iframeNode)
	testType("source", sourceNode)
	testType("embed", embedNode)
	testType("object", objectNode)
	testType("img", imageNode)
	testType("a", anchorNode)
	testType("unknown", unknownNode)
}

func TestGetNodeAttrValue(t *testing.T) {
	var n html.Node
	n.Attr = []html.Attribute{
		html.Attribute{"", "key1", "value1"},
		html.Attribute{"", "key2", "value2"},
		html.Attribute{"", "key3", "value3"},
	}

	for _, attr := range n.Attr {
		val, err := getNodeAttrValue(&n, attr.Key)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		} else if val != attr.Val {
			t.Errorf("Exepected %s, got %s", attr.Val, val)
		}
	}

	_, err := getNodeAttrValue(&n, "bad-key")
	if err != nodeAttrNotFound {
		t.Errorf("Expected error: %v", nodeAttrNotFound)
	}
}
