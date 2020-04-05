package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Block struct {
	Index      int
	Timestamp  time.Time
	Data       string
	ParentHash string
	Hash       string
	Ready      bool
}

func (b *Block) genHash() {
	hash := sha256.New()
	key := strconv.Itoa(b.Index)
	key += b.Timestamp.String()
	key += b.Data
	key += string(b.ParentHash)
	hash.Write([]byte(key))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
	b.Ready = true
}

func (b *Block) verify(parentHash string) bool {
	hash := sha256.New()
	key := strconv.Itoa(b.Index)
	key += b.Timestamp.String()
	key += b.Data
	key += string(b.ParentHash)
	hash.Write([]byte(key))

	return hex.EncodeToString(hash.Sum(nil)) == b.Hash && parentHash == b.ParentHash
}

func (b *Block) serialise() []byte {
	var jsonData []byte
	jsonData, err := json.Marshal(b)
	if err != nil {
		printError("Error serialising block")
	}
	return jsonData
}

func deserialise(ser []byte) Block {
	var block Block
	err := json.Unmarshal(ser, &block)
	if err != nil {
		printError("An unexpected error occured. Please try again")
	}
	return block
}

func (b *Block) toConn(conn net.Conn) (int, error) {
	return conn.Write(b.serialise())
}

func (b *Block) Print() {
	fmt.Println("=========================")
	fmt.Printf("Index:\t\t%d\n", b.Index)
	fmt.Printf("Timestamp:\t%s\n", b.Timestamp.String())
	fmt.Printf("Data:\t\t%s\n", b.Data)
	fmt.Printf("Parent Hash:\t%064s\n", b.ParentHash)
	fmt.Printf("Hash:\t\t%64s\n", b.Hash)
}
