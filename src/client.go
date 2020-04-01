package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"syscall"
)

func getConn() net.Conn {
	addr := REMOTE_HOST + ":" + BLOCKCHAIN_PORT
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		if err.(*net.OpError).Err.(*os.SyscallError).Err == syscall.ECONNREFUSED {
			fmt.Printf("Peer at %s is not contactable\n", addr)
			return nil
		} else {
			exit_on_error(err)
		}
	}
	return conn

}
func query(hash []byte) {
	conn := getConn()
	if conn == nil {
		return
	}
	conn.Write([]byte("q"))
	conn.Write(hash)
	buf := make([]byte, 512)
	conn.Read(buf)
	block := deserialise(buf)
	fmt.Println(block)
}

func reqsync(chain *[]Block) {
	conn := getConn()
	if conn == nil {
		return
	}
	var lastBlock Block
	lastBlock.ready = false
	conn.Write([]byte("s"))
	for {
		buf := make([]byte, 512)
		conn.Read(buf)
		buf = bytes.Trim(buf, "\x00")
		block := deserialise(buf)
		if (!lastBlock.ready && block.index == 1) || block.verify(lastBlock.hash) {
			*chain = append(*chain, block)
		} else {
			fmt.Println("Error syncing. Check reliability")
			os.Exit(1)
		}
	}
}
