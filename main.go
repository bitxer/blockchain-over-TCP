package main

import (
	"fmt"
)

func main() {
	block1 := Block{index: 1, timestamp: 1, data: "data"}
	h := block1.genHash()
	fmt.Printf("%x", h)
}
