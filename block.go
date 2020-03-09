package main

import (
	"crypto/sha256"
	"strconv"
)

type Block struct {
	index      int
	timestamp  int
	data       string
	parentHash []byte
	hash       []byte
}

func (b *Block) genHash() []byte {
	hash := sha256.New()
	key := strconv.Itoa(b.index)
	key += strconv.Itoa(b.timestamp)
	key += b.data
	key += string(b.parentHash)
	hash.Write([]byte(key))
	return hash.Sum(nil)
}
