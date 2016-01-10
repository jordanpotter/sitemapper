package mapper

import (
	"fmt"
	"net/url"
	"testing"

	"golang.org/x/net/html"
)

func TestProcessNode(t *testing.T) {
	urlStr := "https://foo.com"
	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assetURLStr2 := "image2.png"
	child3 := html.Node{
		Type: html.ElementNode,
		Data: "img",
		Attr: []html.Attribute{
			html.Attribute{Key: "src", Val: assetURLStr2},
		},
	}

	assetURLStr1 := "image1.png"
	child2 := html.Node{
		Type: html.ElementNode,
		Data: "img",
		Attr: []html.Attribute{
			html.Attribute{Key: "src", Val: assetURLStr1},
		},
		NextSibling: &child3,
	}

	linkURLStr := "path/to/site"
	child1 := html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			html.Attribute{Key: "href", Val: linkURLStr},
		},
		NextSibling: &child2,
	}

	root := html.Node{
		Type: html.ElementNode,
		Data: "div",
		Attr: []html.Attribute{
			html.Attribute{Key: "class", Val: "root"},
		},
		FirstChild: &child1,
	}

	pm := PageMap{URL: u}
	err = processNode(&pm, &root)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(pm.Links) != 1 {
		t.Errorf("Exepected number links to be 1, got %d", len(pm.Links))
	} else if len(pm.Assets) != 2 {
		t.Errorf("Exepected number assets to be 2, got %d", len(pm.Assets))
	}

	expectedLinkStr := fmt.Sprintf("%s/%s", urlStr, linkURLStr)
	if expectedLinkStr != pm.Links[0].String() {
		t.Errorf("Exepected link url to be %q, got %q", expectedLinkStr, pm.Links[0].String())
	}

	expectedAssetStr1 := fmt.Sprintf("%s/%s", urlStr, assetURLStr1)
	if expectedAssetStr1 != pm.Assets[0].String() {
		t.Errorf("Exepected asset url to be %q, got %q", expectedAssetStr1, pm.Assets[0].String())
	}

	expectedAssetStr2 := fmt.Sprintf("%s/%s", urlStr, assetURLStr2)
	if expectedAssetStr2 != pm.Assets[1].String() {
		t.Errorf("Exepected asset url to be %q, got %q", expectedAssetStr2, pm.Assets[1].String())
	}
}

func TestAddLinkURL(t *testing.T) {
	urlStr := "https://foo.com"
	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	linkURLStr := "path/to/site"
	n := html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			html.Attribute{Key: "href", Val: linkURLStr},
		},
	}

	pm := PageMap{URL: u}
	err = addLinkURL(&pm, &n)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(pm.Links) != 1 {
		t.Errorf("Exepected number links to be 1, got %d", len(pm.Links))
	} else if len(pm.Assets) != 0 {
		t.Errorf("Exepected number assets to be 0, got %d", len(pm.Assets))
	}

	expectedLinkStr := fmt.Sprintf("%s/%s", urlStr, linkURLStr)
	if expectedLinkStr != pm.Links[0].String() {
		t.Errorf("Exepected link url to be %q, got %q", expectedLinkStr, pm.Links[0].String())
	}
}

func TestAddAssetURL(t *testing.T) {
	urlStr := "https://foo.com"
	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	assetURLStr := "image.png"
	n := html.Node{
		Type: html.ElementNode,
		Data: "img",
		Attr: []html.Attribute{
			html.Attribute{Key: "src", Val: assetURLStr},
		},
	}

	pm := PageMap{URL: u}
	err = addAssetURL(&pm, &n)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(pm.Links) != 0 {
		t.Errorf("Exepected number links to be 0, got %d", len(pm.Links))
	} else if len(pm.Assets) != 1 {
		t.Errorf("Exepected number assets to be 1, got %d", len(pm.Assets))
	}

	expectedAssetStr := fmt.Sprintf("%s/%s", urlStr, assetURLStr)
	if expectedAssetStr != pm.Assets[0].String() {
		t.Errorf("Exepected asset url to be %q, got %q", expectedAssetStr, pm.Assets[0].String())
	}
}
