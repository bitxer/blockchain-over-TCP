package main

import (
	"bytes"
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
		printInfo("Block of index", string(index), "not found")
		// fmt.Printf("not found")
		return
	} else {
		buf = make([]byte, 256)
		n, _ := conn.Read(buf)
		buf = buf[:n]
		block := deserialise(buf)
		block.Print()
	}
	conn.Close()
}

func reqsync(chain *[]Block) {
	conn := getConn()
	if conn == nil {
		return
	}
	conn.Write([]byte("s"))

	if len(*chain) == 0 {
		conn.Write([]byte{byte(0)})
	} else {
		lastBlock := (*chain)[(len(*chain))-1]
		conn.Write([]byte{byte(lastBlock.Index)})
	}

	stop := false
	stopcode := []byte("bitxer")
	stopcode_len := len(stopcode)
	for !stop {
		buf := make([]byte, 256)
		n, _ := conn.Read(buf)
		buf = buf[:n]
		stop = n == stopcode_len && bytes.Compare(buf, stopcode) == 0
		if !stop {
			block := deserialise(buf)
			addtoChain(chain, block)
			conn.Write([]byte{byte(n)})
		}
	}
	conn.Close()
}

func add(chain *[]Block, data string) {
	reqsync(chain)

	conn := getConn()
	if conn == nil {
		return
	}

	index := 0
	lhash := ""

	// if len(*chain) > 0 {
	lastBlock := (*chain)[len(*chain)-1]
	index = lastBlock.Index
	lhash = lastBlock.Hash
	// } else {
	// 	printInfo("Chain not initialised")
	// 	printInfo("Adding block as start of chain block")
	// 	lhash = []byte{0}
	// }
	block := Block{Index: index + 1, Timestamp: time.Now(), Data: data, ParentHash: lhash}
	block.genHash()
	// addtoChain(chain, block)
	conn.Write([]byte("a"))
	buf := make([]byte, 1)
	// conn.Read(buf)

	// if buf[0] == 1 {
	block.toConn(conn)
	conn.Read(buf)
	if buf[0] == 1 {
		printSuccess("Successfully add block to chain")
	} else {
		printError("Block added unsuccessfully due to validation errors")
	}
	conn.Close()
	// }
}

func querylast() {
	conn := getConn()
	if conn == nil {
		return
	}
	conn.Write([]byte("l"))

	buf := make([]byte, 1)
	conn.Read(buf)
	if buf[0] == 0 {
		printInfo("Chain not initialised")
		return
	} else {
		buf = make([]byte, 256)
		n, _ := conn.Read(buf)
		buf = buf[:n]
		block := deserialise(buf)
		block.Print()
	}
	conn.Close()
}
