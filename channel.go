package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	LOCAL_HOST = "127.0.0.1"
	LOCAL_PORT = "6442"

	LISTEN_HOST = "127.0.0.1"
	LISTEN_PORT = "6443"

	REMOTE_HOST = LISTEN_HOST
	REMOTE_PORT = LISTEN_PORT
)

func add(buf []byte, chain *[]Block, conn net.Conn) {
	for {
		lastBlock := (*chain)[len(*chain)-1]
		block := deserialise(buf)

		if block.verify(lastBlock.hash) {
			_ = binary.Write(conn, binary.LittleEndian, "Block verified and will be added to the blockchain")
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
					ser := v.serialise()
					_ = binary.Write(conn, binary.LittleEndian, ser)
				}
			}
		}(v)
	}
	wg.Wait()
	if stop == 0 {
		_ = binary.Write(conn, binary.LittleEndian, "Not found")
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
			buf := make([]byte, 1)
			_, err = conn.Read(buf)
			fmt.Printf("Connected to: %s\n", conn.RemoteAddr().String())
			err = binary.Read(conn, binary.LittleEndian, &buf)
			act := buf[0]
			buf = make([]byte, 512)
			err = binary.Read(conn, binary.LittleEndian, &buf)
			switch act {
			case 'a':
				go add(buf, chain, conn)
				break
			case 'q':
				go search(buf, chain, conn)
				break
			}
		}
	}
}

func query() {
	conn, err := net.Dial("tcp", REMOTE_HOST+":"+REMOTE_PORT)
	exit_on_error(err)
	err = binary.Write(conn, binary.LittleEndian, "a")
	exit_on_error(err)
}

func exit_on_error(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
