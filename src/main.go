package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
	block1 := Block{Index: 1, Timestamp: 1, Data: "data", ParentHash: []byte{0}}
	block1.genHash()
	chain := []Block{block1}
	go listen(&chain, &wg)
	wg.Wait()

	for {
		fmt.Println("What woulld you like to do?")
		fmt.Println("1. Query Block")
		fmt.Println("2. Add Block")
		fmt.Println("3. Sync")
		fmt.Println("4. Print all blocks")
		fmt.Print("> ")
		option, _ := reader.ReadString('\n')
		option = strings.Replace(option, "\n", "", -1)

		switch {
		case option[0] == '1':
			fmt.Print("Query for block of index: ")
			buf, _ := reader.ReadString('\n')
			buf = strings.Replace(option, "\n", "", -1)
			index, _ := strconv.Atoi(buf)
			query(index)
		case option[0] == '2':
			// add()
			fmt.Println("Add Block")
		case option[0] == '3':
			fmt.Println("Sync chain")
		case option[0] == '4':
			fmt.Print("Printing all blocks in chain")
		default:
			fmt.Println("[-] Invalid option specified. Please choose a valid option")
		}
	}
}
