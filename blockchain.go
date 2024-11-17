package main

import (
	"fmt"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	previousHash string
	timestamp    int64
	transactions []string
}

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

func NewBlock(nonce int, previousHash string) *Block {

	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		timestamp:    time.Now().UnixNano(),
	}

}

func (b *Block) PrintBlock() {
	fmt.Printf("nonce:			%d\n", b.nonce)
	fmt.Printf("preHash:		%s\n", b.previousHash)
	fmt.Printf("time_Stamp:		%d\n", b.timestamp)
	fmt.Printf("transaction:		%s\n", b.transactions)
}

func NewBlockchain() *Blockchain {
	bc := new(Blockchain)
	bc.CreateBlock(0, "init hash")
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash string) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("-", 20), i, strings.Repeat("-", 20))
		block.PrintBlock()
		fmt.Printf("%s\n", strings.Repeat("=", 50))
	}
}

func main() {
	blockChain := NewBlockchain()
	blockChain.Print()
	blockChain.CreateBlock(5, "hash 1")
	blockChain.Print()
	blockChain.CreateBlock(2, "hash 2")
	blockChain.Print()

}
