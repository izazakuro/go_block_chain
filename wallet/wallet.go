package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {
	// Create Publickey and Privatekey
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &privateKey.PublicKey

	// SHA-256 on publickey
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)

	// RIPEMD-160 on SHA-256
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	// Add version byte in front of RIPEMD-160
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	// SHA-256 on vd4
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	// SHA-256 on digets5
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	// take 4bytes for checksum
	chsum := digest6[:4]

	// Add chsum
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])

	//Convert into base58
	address := base58.Encode(dc8)
	w.blockchainAddress = address

	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}
