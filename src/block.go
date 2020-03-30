package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"net"
	"strconv"
)

type Block struct {
	index      int
	timestamp  int
	data       string
	parentHash []byte
	hash       []byte
	ready      bool
}

func (b *Block) genHash() {
	hash := sha256.New()
	key := strconv.Itoa(b.index)
	key += strconv.Itoa(b.timestamp)
	key += b.data
	key += string(b.parentHash)
	hash.Write([]byte(key))
	b.hash = hash.Sum(nil)
	b.ready = true
}

func (b *Block) serialise() *bytes.Buffer {
	ex := make(map[string]string)
	ex["index"] = strconv.Itoa(b.index)
	ex["timestamp"] = strconv.Itoa(b.timestamp)
	ex["data"] = b.data
	ex["parentHash"] = string(b.parentHash)
	ex["hash"] = string(b.hash)
	buf := new(bytes.Buffer)

	en := gob.NewEncoder(buf)
	err := en.Encode(ex)
	if err != nil {
		panic(err)
	}
	return buf
}

func deserialise(ser []byte) Block {
	buf := bytes.NewBuffer(ser)

	im := make(map[string]string)
	de := gob.NewDecoder(buf)

	err := de.Decode(&im)
	if err != nil {
		panic(err)
	}

	index, err := strconv.Atoi(im["index"])
	if err != nil {
		panic(err)
	}

	timestamp, err := strconv.Atoi(im["timestamp"])
	if err != nil {
		panic(err)
	}

	data := im["data"]
	parentHash := []byte(im["parentHash"])
	hash := []byte(im["hash"])
	return Block{index: index, timestamp: timestamp, data: data, parentHash: parentHash, hash: hash}
}

func (b *Block) verify(parentHash []byte) bool {
	hash := sha256.New()
	key := strconv.Itoa(b.index)
	key += strconv.Itoa(b.timestamp)
	key += b.data
	key += string(b.parentHash)
	hash.Write([]byte(key))

	return string(hash.Sum(nil)) == string(b.hash) && string(parentHash) == string(b.parentHash)
}

func toConn(b Block, conn net.Conn) (int, error) {
	return conn.Write(b.serialise().Bytes())
}
