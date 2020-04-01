package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

var REMOTE_HOST string
var BLOCKCHAIN_PORT string

func main() {
	fmt.Print("Enter peer address: ")
	reader := bufio.NewReader(os.Stdin)
	REMOTE_HOST, _ = reader.ReadString('\n')
	REMOTE_HOST = strings.Replace(REMOTE_HOST, "\n", "", -1)

	fmt.Print("Enter blockchain port: ")
	BLOCKCHAIN_PORT, _ = reader.ReadString('\n')
	BLOCKCHAIN_PORT = strings.Replace(BLOCKCHAIN_PORT, "\n", "", -1)

	var wg sync.WaitGroup
	wg.Add(1)
	block1 := Block{index: 1, timestamp: 1, data: "data", parentHash: []byte{0}}
	block1.genHash()
	chain := []Block{block1}
	go listen(&chain)
	go query(block1.hash)
	wg.Wait()
}
