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
	Timestamp  int
	Data       string
	ParentHash []byte
	Hash       []byte
	Ready      bool
}

func (b *Block) genHash() {
	hash := sha256.New()
	key := strconv.Itoa(b.Index)
	key += strconv.Itoa(b.Timestamp)
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
	key += strconv.Itoa(b.Timestamp)
	key += b.Data
	key += string(b.ParentHash)
	hash.Write([]byte(key))

	return string(hash.Sum(nil)) == string(b.Hash) && string(parentHash) == string(b.ParentHash)
}

func (b *Block) toConn(conn net.Conn) (int, error) {
	return conn.Write(b.serialise())
}

func addBlock(chain *[]Block, b Block) bool {
	lastBlock := (*chain)[len(*chain)-1]

	if b.verify(lastBlock.Hash) {
		if lastBlock.Index+1 < b.Index {
			sleep := time.Duration(lastBlock.Index - b.Index - 1)
			time.Sleep(sleep * time.Second)
		}
		*chain = append(*chain, b)
	}
	return string(b.Hash) == string((*chain)[len(*chain)-1].Hash)
}

func (b *Block) Print() {
	fmt.Println("=========================")
	fmt.Printf("Index:\t\t%d\n", b.Index)
	fmt.Printf("Timestamp:\t%d\n", b.Timestamp)
	fmt.Printf("Data:\t\t%s\n", b.Data)
	fmt.Printf("Parent Hash:\t%064s\n", hex.EncodeToString(b.ParentHash))
	fmt.Printf("Hash:\t\t%64s\n", hex.EncodeToString(b.Hash))
}
