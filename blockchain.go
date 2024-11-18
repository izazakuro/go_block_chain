package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

// Block
type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []*Transaction
}

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {

	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		timestamp:    time.Now().UnixNano(),
		transactions: transactions,
	}

}

func (b *Block) PrintBlock() {
	fmt.Printf("nonce:			%d\n", b.nonce)
	fmt.Printf("preHash:		%x\n", b.previousHash)
	fmt.Printf("time_Stamp:		%d\n", b.timestamp)
	for _, t := range b.transactions {
		t.Print()
	}
}

// Hash
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

func NewBlockchain(blockchainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAddress
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("-", 20), i, strings.Repeat("-", 20))
		block.PrintBlock()
		fmt.Printf("%s\n", strings.Repeat("=", 50))
	}
}

func (bc *Blockchain) AddTransaction(sender, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.senderBlockChainAddress,
				t.recipientBlockChainAddress,
				t.value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {

	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
		timestamp:    0,
	}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros

}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

// Mining
func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("MINING SUCCESS")
	return true
}

// Transaction
type Transaction struct {
	senderBlockChainAddress    string
	recipientBlockChainAddress string
	value                      float32
}

func NewTransaction(sender, recipient string, value float32) *Transaction {
	return &Transaction{
		senderBlockChainAddress:    sender,
		recipientBlockChainAddress: recipient,
		value:                      value,
	}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 50))
	fmt.Printf("sender: %s\n", t.senderBlockChainAddress)
	fmt.Printf("recipient: %s\n", t.recipientBlockChainAddress)
	fmt.Printf("value: %.1f\n", t.value)
	fmt.Printf("%s\n", strings.Repeat("-", 50))
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockChainAddress    string  `json:"sender_block_chain_address"`
		RecipientBlockChainAddress string  `json:"recipient_block_chain_address"`
		Value                      float32 `json:"value"`
	}{
		SenderBlockChainAddress:    t.senderBlockChainAddress,
		RecipientBlockChainAddress: t.recipientBlockChainAddress,
		Value:                      t.value,
	})
}

func main() {
	myAddress := "my_address"
	blockChain := NewBlockchain(myAddress)
	blockChain.Print()

	blockChain.AddTransaction("A", "B", 1.0)
	blockChain.Mining()
	blockChain.Print()

	blockChain.AddTransaction("C", "D", 2.0)
	blockChain.AddTransaction("X", "Y", 3.0)

	blockChain.Mining()

	blockChain.Print()

}
