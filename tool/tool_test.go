package tool

import (
	"fmt"
	"testing"
)

func TestGetAllVer(t *testing.T) {
	fmt.Println(GetAllVer("v", ""))
}

func TestGetAttrPaths(t *testing.T) {
	GetAttrPaths("3.4.9")
}
