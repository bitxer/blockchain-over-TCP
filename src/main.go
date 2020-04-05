package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var REMOTE_HOST string
var BLOCKCHAIN_PORT string
var verbose = false

func main() {
	printPrompt("Enter peer address:")
	reader := bufio.NewReader(os.Stdin)
	REMOTE_HOST, _ = reader.ReadString('\n')
	REMOTE_HOST = strings.Replace(REMOTE_HOST, "\n", "", -1)

	printPrompt("Enter blockchain port:")
	BLOCKCHAIN_PORT, _ = reader.ReadString('\n')
	BLOCKCHAIN_PORT = strings.Replace(BLOCKCHAIN_PORT, "\n", "", -1)

	printPrompt("Run listener on verbose? [Y/n]:")
	vopt, _ := reader.ReadString('\n')
	vopt = strings.Replace(vopt, "\n", "", -1)
	verbose = vopt[0] == 'Y'

	printPrompt("Would you like to initialise the blockchain? [Y/n]:")
	choice, _ := reader.ReadString('\n')
	chain := []Block{}

	switch choice[0] {
	case 'Y':
		block1 := Block{Index: 1, Timestamp: time.Now(), Data: "Genesis", ParentHash: ""}
		block1.genHash()
		chain = []Block{block1}
		printSuccess("Chain initialised with genesis block")
	case 'n':
		printInfo("Skipped chain initialisation")
	default:
		printError("Unknown input")
		printInfo("Skipped chain inistialisation")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go listen(&chain, &wg)
	wg.Wait()

	for {
		fmt.Println("=========================")
		fmt.Println("What woulld you like to do?")
		fmt.Println("1. Query Block")
		fmt.Println("2. Query Last Block")
		fmt.Println("3. Add Block")
		fmt.Println("4. Sync")
		fmt.Println("5. Print all blocks")
		fmt.Print("> ")
		option, _ := reader.ReadString('\n')
		option = strings.Replace(option, "\n", "", -1)
		if option == "" {
			continue
		}

		switch {
		case option[0] == '1':
			// Query for block at index
			printPrompt("Query for block of index:")
			buf, _ := reader.ReadString('\n')
			buf = strings.Replace(buf, "\n", "", -1)
			index, _ := strconv.Atoi(buf)
			query(index)
		case option[0] == '2':
			// Query for last block
			querylast()
		case option[0] == '3':
			// Create block with data
			printPrompt("Data to be contained in block:")
			buf, _ := reader.ReadString('\n')
			buf = strings.Replace(buf, "\n", "", -1)
			add(&chain, buf)
		case option[0] == '4':
			// Sync block chain with remote peer
			reqsync(&chain)
		case option[0] == '5':
			// Print all block chain to screen
			for _, v := range chain {
				v.Print()
			}
		default:
			printError("Invalid option specified. Please choose a valid option")
		}
	}
}

func addtoChain(chain *[]Block, b Block) bool {
	if len(*chain) > 0 {
		lastBlock := (*chain)[len(*chain)-1]
		if b.verify(lastBlock.Hash) && lastBlock.Index+1 == b.Index {
			return false
		}
		*chain = append(*chain, b)
	} else {
		*chain = append(*chain, b)
	}
	return b.Hash == (*chain)[len(*chain)-1].Hash
}
