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
	ParentHash []byte
	Hash       []byte
	Ready      bool
}

func (b *Block) genHash() {
	hash := sha256.New()
	key := strconv.Itoa(b.Index)
	key += b.Timestamp.String()
	key += b.Data
	key += string(b.ParentHash)
	hash.Write([]byte(key))
	b.Hash = hash.Sum(nil)
	b.Ready = true
}

func (b *Block) serialise() []byte {
	var jsonData []byte
	jsonData, err := json.Marshal(b)
	if err != nil {
		fmt.Println("Error serialising block")
	}
	return jsonData
}

func deserialise(ser []byte) Block {
	var block Block
	err := json.Unmarshal(ser, &block)
	if err != nil {
		fmt.Println("An unexpected error occured. Please try again")
	}
	return block
}

func (b *Block) verify(parentHash []byte) bool {
	hash := sha256.New()
	key := strconv.Itoa(b.Index)
	key += b.Timestamp.String()
	key += b.Data
	key += string(b.ParentHash)
	hash.Write([]byte(key))

	return string(hash.Sum(nil)) == string(b.Hash) && string(parentHash) == string(b.ParentHash)
}

func (b *Block) toConn(conn net.Conn) (int, error) {
	return conn.Write(b.serialise())
}

func (b *Block) Print() {
	fmt.Println("=========================")
	fmt.Printf("Index:\t\t%d\n", b.Index)
	fmt.Printf("Timestamp:\t%s\n", b.Timestamp.String())
	fmt.Printf("Data:\t\t%s\n", b.Data)
	fmt.Printf("Parent Hash:\t%064s\n", hex.EncodeToString(b.ParentHash))
	fmt.Printf("Hash:\t\t%64s\n", hex.EncodeToString(b.Hash))
}
