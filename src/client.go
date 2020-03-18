package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

const (
	REMOTE_HOST = "127.0.0.1"
	REMOTE_PORT = "6444"
)

func query(hash []byte) {
	addr := REMOTE_HOST + ":" + REMOTE_PORT
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		if err.(*net.OpError).Err.(*os.SyscallError).Err == syscall.ECONNREFUSED {
			fmt.Printf("Peer at %s is not contactable\n", addr)
			return
		} else {
			exit_on_error(err)
		}
	}
	conn.Write([]byte("q"))
	conn.Write(hash)
	buf := make([]byte, 512)
	conn.Read(buf)
	block := deserialise(buf)
	fmt.Println(block)
}
