package mapper

import (
	"testing"

	"golang.org/x/net/html"
)

func TestGetElementType(t *testing.T) {
	testType := func(data string, et elementType) {
		var n html.Node
		n.Data = data
		elemType := getElementType(&n)
		if elemType != et {
			t.Errorf("Exepected %d, got %d", et, elemType)
		}
	}

	testType("a", linkElement)
	testType("img", imageElement)
	testType("unknown", unknownElement)
}

func TestGetElementAttrValue(t *testing.T) {
	var n html.Node
	attr1 := html.Attribute{"", "key1", "value1"}
	attr2 := html.Attribute{"", "key2", "value2"}
	attr3 := html.Attribute{"", "key3", "value3"}
	n.Attr = []html.Attribute{attr1, attr2, attr3}

	for _, attr := range n.Attr {
		val, err := getElementAttrValue(&n, attr.Key)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		} else if val != attr.Val {
			t.Errorf("Exepected %s, got %s", attr.Val, val)
		}
	}

	_, err := getElementAttrValue(&n, "bad-key")
	if err != elementAttrNotFound {
		t.Errorf("Expected error: %v", elementAttrNotFound)
	}
}
