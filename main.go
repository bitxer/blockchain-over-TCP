package main

import (
	"fmt"
)

func main() {
	block1 := Block{index: 1, timestamp: 1, data: "data", parentHash: []byte{0}}
	block1.genHash()
	chain := []Block{block1}
	chain = chain
	be := block1.serialise()

	fmt.Printf("%s\n", be)

	fmt.Printf("%x\n", block1)
	block2 := deserialise(be)
	fmt.Printf("%x\n", block2)
}
