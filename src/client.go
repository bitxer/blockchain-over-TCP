package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

func getConn() net.Conn {
	conn, err := net.Dial("tcp", REMOTE_HOST)
	if err != nil {
		if err.(*net.OpError).Err.(*os.SyscallError).Err == syscall.ECONNREFUSED {
			fmt.Printf("Peer at %s is not contactable\n", REMOTE_HOST)
			return nil
		} else {
			exit_on_error(err)
		}
	}
	return conn

}
func query(index int) {
	conn := getConn()
	if conn == nil {
		return
	}
	conn.Write([]byte("q"))

	conn.Write([]byte{byte(index)})
	buf := make([]byte, 1)
	conn.Read(buf)
	if buf[0] == 0 {
		fmt.Printf("not found")
		return
	} else {
		buf = make([]byte, 512)
		n, _ := conn.Read(buf)
		buf = buf[:n]
		block := deserialise(buf)
		// fmt.Printf("Queried block has '%s' in data\n", block.Data)
		block.Print()
	}
	// fmt.Println(block)
}

func reqsync(gChain *[]Block) {
	conn := getConn()
	if conn == nil {
		return
	}
	chain := *gChain
	lastBlock := chain[(len(chain))-1]
	conn.Write([]byte("s"))
	conn.Write([]byte{byte(lastBlock.Index)})

}

func add(chain *[]Block, data string) {
	conn := getConn()
	if conn == nil {
		return
	}
	lastBlock := (*chain)[len(*chain)-1]
	block := Block{Index: lastBlock.Index + 1, Timestamp: time.Now(), Data: data, ParentHash: lastBlock.Hash}
	fmt.Println(block)

}
