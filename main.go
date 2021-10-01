package main

import (
	"fmt"

	. "github.com/nsip/sif-spec-res/3.4.8"
)

func main() {
	fmt.Println(string(JSON_ATTR["Activity_0"]))
	fmt.Println(string(JSON_ATTR["Activity_1"]))
	fmt.Println(string(JSON_BOOL["Activity_0"]))
	fmt.Println(string(JSON_BOOL["Activity_1"]))
	fmt.Println(string(JSON_LIST["Activity_0"]))
	fmt.Println(string(JSON_LIST["Activity_1"]))
	fmt.Println(string(JSON_NUM["Activity_0"]))
	fmt.Println(string(JSON_NUM["Activity_1"]))
}
