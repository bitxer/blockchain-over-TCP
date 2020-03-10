package main

import (
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	block1 := Block{index: 1, timestamp: 1, data: "data", parentHash: []byte{0}}
	block1.genHash()
	chain := []Block{block1}
	go listen(&chain)
	wg.Wait()
}
