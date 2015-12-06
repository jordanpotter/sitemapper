package mapper

import (
	"fmt"
	"testing"
)

func TestCreatePageMap(t *testing.T) {
	pageMap, err := CreatePageMap("https://digitalocean.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	for _, a := range pageMap.Assets {
		fmt.Println(a)
	}

	fmt.Println("==========")

	for _, l := range pageMap.Links {
		fmt.Println(l)
	}
}
