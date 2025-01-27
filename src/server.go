package main

import (
	"net"
	"os"
	"sync"
)

func server_add(buf []byte, chain *[]Block, conn net.Conn) {
	if addtoChain(chain, deserialise(buf)) {
		conn.Write([]byte{1})
		return
	}
	conn.Write([]byte{0})
	conn.Close()
}

func search(index int, chain *[]Block, conn net.Conn) {
	if index > len(*chain) {
		conn.Write([]byte{0})
		return
	} else {
		conn.Write([]byte{1})
	}
	v := (*chain)[index-1]

	v.toConn(conn)
	conn.Close()
}

func syncchain(chain *[]Block, conn net.Conn) {
	for _, v := range *chain {
		n, _ := v.toConn(conn)
		buf := make([]byte, 1)
		conn.Read(buf)
		if n != int(buf[0]) {
			printError("Mismatch length when syncing chain")
			break
		}

	}
	conn.Write([]byte("bitxer"))
	conn.Close()
}

func listen(chain *[]Block, wg *sync.WaitGroup) {
	listener, err := net.Listen("tcp", ":"+BLOCKCHAIN_PORT)
	exit_on_error(err)

	printSuccess("Listening on port", BLOCKCHAIN_PORT)
	wg.Done()
	for {
		conn, err := listener.Accept()
		if err != nil {
			printError(err.Error())
		} else {
			buf := make([]byte, 1)
			_, err = conn.Read(buf)
			if verbose {
				printInfo("Connected to:", conn.RemoteAddr().String())
			}

			switch buf[0] {
			case 'a':
				buf = make([]byte, 256)
				n, _ := conn.Read(buf)
				buf = buf[:n]
				go server_add(buf, chain, conn)
			case 'q':
				buf = make([]byte, 1)
				conn.Read(buf)
				go search(int(buf[0]), chain, conn)
			case 'l':
				lastBlock := (*chain)[len(*chain)-1]
				go search(lastBlock.Index, chain, conn)
			case 's':
				go syncchain(chain, conn)
			}
		}
	}
}

func exit_on_error(err error) {
	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}
}
