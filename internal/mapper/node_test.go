package mapper

import (
	"testing"

	"golang.org/x/net/html"
)

func TestGetNodeType(t *testing.T) {
	testType := func(data string, attr []html.Attribute, et nodeType) {
		n := html.Node{Data: data, Attr: attr}
		nodeType := getNodeType(&n)
		if nodeType != et {
			t.Errorf("Exepected %d, got %d", et, nodeType)
		}
	}

	testType("iframe", []html.Attribute{}, iframeNode)
	testType("source", []html.Attribute{}, sourceNode)
	testType("embed", []html.Attribute{}, embedNode)
	testType("object", []html.Attribute{}, objectNode)
	testType("img", []html.Attribute{}, imageNode)
	testType("a", []html.Attribute{}, anchorNode)
	testType("script", []html.Attribute{}, scriptNode)
	testType("link", []html.Attribute{
		html.Attribute{Key: "rel", Val: "stylesheet"},
	}, stylesheetNode)
	testType("link", []html.Attribute{
		html.Attribute{Key: "rel", Val: "icon"},
	}, icoNode)
	testType("unknown", []html.Attribute{}, unknownNode)
}

func TestIsStylesheetNode(t *testing.T) {
	n := html.Node{Data: "link"}
	if isStylesheetNode(&n) {
		t.Errorf("Exepected node to not be a stylesheet node")
	}

	n.Attr = []html.Attribute{
		html.Attribute{Key: "rel", Val: "stylesheet"},
	}
	if !isStylesheetNode(&n) {
		t.Errorf("Exepected node to be a stylesheet node")
	}
}

func TestIsIcoNode(t *testing.T) {
	n := html.Node{Data: "link"}
	if isIcoNode(&n) {
		t.Errorf("Exepected node to not be a ico node")
	}

	n.Attr = []html.Attribute{
		html.Attribute{Key: "rel", Val: "icon"},
	}
	if !isIcoNode(&n) {
		t.Errorf("Exepected node to be a ico node")
	}
}

func TestGetNodeAttrValue(t *testing.T) {
	var n html.Node
	n.Attr = []html.Attribute{
		html.Attribute{Key: "key1", Val: "value1"},
		html.Attribute{Key: "key2", Val: "value2"},
		html.Attribute{Key: "key3", Val: "value3"},
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
	if err != errNodeAttrNotFound {
		t.Errorf("Expected error: %v", errNodeAttrNotFound)
	}
}
