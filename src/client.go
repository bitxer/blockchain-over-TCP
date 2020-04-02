package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"syscall"
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

func reqsync(chain *[]Block) {
	conn := getConn()
	if conn == nil {
		return
	}
	var lastBlock Block
	lastBlock.Ready = false
	conn.Write([]byte("s"))
	for {
		buf := make([]byte, 512)
		conn.Read(buf)
		buf = bytes.Trim(buf, "\x00")
		block := deserialise(buf)
		if (!lastBlock.Ready && block.Index == 1) || block.verify(lastBlock.Hash) {
			*chain = append(*chain, block)
		} else {
			fmt.Println("Error syncing. Check reliability")
			os.Exit(1)
		}
	}
}

// func addBlock(chain *[]Block, data string) {
// 	conn := getConn()
// 	if conn == nil {
// 		return
// 	}
// 	lastBlock := (*chain)[len(*chain)-1]
// 	block := Block{index: lastBlock.index + 1, timestamp: 1, data: data, parentHash: lastBlock.hash}

// }
