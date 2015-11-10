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

	fmt.Println(pageMap)
}
