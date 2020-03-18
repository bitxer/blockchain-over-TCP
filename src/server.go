package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	LISTEN_HOST = "127.0.0.1"
	LISTEN_PORT = "6443"
)

func add(buf []byte, chain *[]Block, conn net.Conn) {
	for {
		lastBlock := (*chain)[len(*chain)-1]
		block := deserialise(buf)

		if block.verify(lastBlock.hash) {
			conn.Write([]byte("Block verified and will be added to the blockchain"))
			if lastBlock.index+1 < block.index {
				sleep := time.Duration(lastBlock.index - block.index - 1)
				time.Sleep(sleep * time.Second)
			}
			*chain = append(*chain, block)
		}
	}
}

func search(buf []byte, chain *[]Block, conn net.Conn) {
	var wg sync.WaitGroup
	wg.Add(len(*chain))
	target := string(buf)
	stop := 0
	for _, v := range *chain {

		go func(v Block) {
			defer wg.Done()
			switch stop {
			case 0:
				if string(v.hash) == target {
					stop = 1
					ser := v.serialise().Bytes()
					conn.Write(ser)
				}
			}
		}(v)
	}
	wg.Wait()
	if stop == 0 {
		_ = binary.Write(conn, binary.LittleEndian, []byte("Not found"))
	}
}

func listen(chain *[]Block) {
	listener, err := net.Listen("tcp", ":"+LISTEN_PORT)
	exit_on_error(err)

	fmt.Println("Listening on port", LISTEN_PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			defer conn.Close()
			defer fmt.Println("")
			act := make([]byte, 1)
			_, err = conn.Read(act)
			fmt.Printf("Connected to: %s\n", conn.RemoteAddr().String())
			hash := make([]byte, 512)
			_, err = conn.Read(hash)
			hash = bytes.Trim(hash, "\x00")
			switch act[0] {
			case 'a':
				go add(hash, chain, conn)
				break
			case 'q':
				go search(hash, chain, conn)
				break
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
