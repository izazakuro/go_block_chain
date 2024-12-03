package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/izazakuro/go_block_chain/block"
	"github.com/izazakuro/go_block_chain/utils"
	"github.com/izazakuro/go_block_chain/wallet"
)

const tempDir = "./templates"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port: port,
		gateway: gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

// Read Template API
func (ws *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(tempDir, "index.html"))
		if err != nil {
			log.Printf("Index Error : %s", err.Error())
		}
		t.Execute(w, "")
	default:
		log.Printf("Error: Invalid http method")
	}
}

// Create Wallet Info API
func (ws *WalletServer) Wallet(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		m, _ := myWallet.MarshalJSON()
		io.Writer.Write(w, m[:])
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: Invalid http method")
	}
}

// Create Transaction API
func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		var t wallet.TransactionRequest
		err := dec.Decode(&t)
		if err != nil {
			log.Printf("Error: %v", err)
			io.Writer.Write(w, utils.JsonStatus("Failed"))
			return
		}
		if !t.Validate() {
			log.Println("Error: missing field")
			return
		}

		publicKey := utils.PublickKeyFromString(*t.SenderPublicKey)
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Println("Error: Parse Float Error")
			io.Writer.Write(w, utils.JsonStatus("Failed"))
			return
		}
		value32 := float32(value)

		w.Header().Add("Conten-Type", "application/json")

		transaction := wallet.NewTransaction(privateKey, publicKey, *t.SenderBlockChainAddress, *t.RecipientBlockChainAddress, value32)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		bt := &block.TransactionRequest{
			SenderBlockChainAddress:    t.SenderBlockChainAddress,
			RecipientBlockChainAddress: t.RecipientBlockChainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      &value32,
			Signature:                  &signatureStr,
		}

		fmt.Println(*bt.RecipientBlockChainAddress, *bt.SenderBlockChainAddress, *bt.SenderPublicKey, *bt.Signature, *bt.Value)

		m, _ := json.Marshal(bt)
		buf := bytes.NewBuffer(m)

		resp, _ := http.Post(ws.Gateway()+"/transactions", "application/json", buf)
		if resp.StatusCode == 201 {
			io.Writer.Write(w, utils.JsonStatus("success"))
			return
		} else {

			io.Writer.Write(w, utils.JsonStatus("failed"))
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: Invalid http method")
	}
}

func (ws *WalletServer) WalletAmount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		endpoint := fmt.Sprintf("%s/amount", ws.Gateway())

		client := &http.Client{}
		bcsReq, _ := http.NewRequest("GET", endpoint, nil)
		q := bcsReq.URL.Query()
		q.Add("blockchain_address", blockchainAddress)
		bcsReq.URL.RawQuery = q.Encode()

		bcsResp, err := client.Do(bcsReq)
		if err != nil {
			log.Printf("Error: %v", err)
			io.Writer.Write(w, utils.JsonStatus("failed"))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if bcsResp.StatusCode == 200 {
			decoder := json.NewDecoder(bcsResp.Body)
			var bar block.AmountRespone
			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("Error: %v", err)
				io.Writer.Write(w, utils.JsonStatus("failed"))
				return
			}
			m, _ := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			})
			io.Writer.Write(w, m)

		} else {
			io.Writer.Write(w, utils.JsonStatus("failed"))
		}
	default:
		log.Println("Error: Invalid http method")
		w.WriteHeader(http.StatusBadRequest)

	}

}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/wallet/amount", ws.WalletAmount)
	http.HandleFunc("/transaction", ws.CreateTransaction)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.port)), nil))
}
