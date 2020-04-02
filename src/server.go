package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"sync"
)

func add(buf []byte, chain *[]Block, conn net.Conn) {
	for {
		// lastBlock := (*chain)[len(*chain)-1]
		// block := deserialise(buf)

		// if block.verify(lastBlock.hash) {
		// 	conn.Write([]byte("Block verified and will be added to the blockchain"))
		// 	if lastBlock.index+1 < block.index {
		// 		sleep := time.Duration(lastBlock.index - block.index - 1)
		// 		time.Sleep(sleep * time.Second)
		// 	}
		// 	*chain = append(*chain, block)
		// }
		if addBlock(chain, deserialise(buf)) {
			conn.Write([]byte("Block have been added successfully"))
		}
	}
}

func search(index int, chain *[]Block, conn net.Conn) {
	var wg sync.WaitGroup
	wg.Add(len(*chain))
	stop := 0
	for _, v := range *chain {

		go func(v Block) {
			defer wg.Done()
			if stop != 1 && v.Index == index {
				stop = 1
				conn.Write([]byte{1})
				v.toConn(conn)
			}
		}(v)
	}
	wg.Wait()
	if stop == 0 {
		conn.Write([]byte{0})
	}
	conn.Close()
}

func syncchain(chain *[]Block, conn net.Conn) {
	for _, v := range *chain {
		v.toConn(conn)
	}
	conn.Write([]byte("bitxer"))
}

func listen(chain *[]Block, wg *sync.WaitGroup) {
	listener, err := net.Listen("tcp", ":"+BLOCKCHAIN_PORT)
	exit_on_error(err)

	fmt.Println("[+] Listening on port", BLOCKCHAIN_PORT)
	wg.Done()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			// defer conn.Close()
			// defer fmt.Println("")
			act := make([]byte, 1)
			_, err = conn.Read(act)
			fmt.Println("Connected to:", conn.RemoteAddr().String())

			switch act[0] {
			case 'a':
				hash := make([]byte, 512)
				_, err = conn.Read(hash)
				hash = bytes.Trim(hash, "\x00")
				go add(hash, chain, conn)
			case 'q':
				conn.Read(act)
				go search(int(act[0]), chain, conn)
			case 's':
				go syncchain(chain, conn)
			}
		}
	}
}

func exit_on_error(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
