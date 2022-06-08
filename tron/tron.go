package tron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
	"io/ioutil"
	"net/http"
	"os"
)

// GenerateWallet creates new wallet end save credentials if it is possible
func GenerateWallet(keystorePassword string, derivationPath string) (Wallet, error) {
	var w Wallet
	bip39.SetWordList(wordlists.English)
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return w, err
	}
	mnemonic, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed(mnemonic, keystorePassword)
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return w, err
	}
	path := hdwallet.MustParseDerivationPath(derivationPath)
	account, err := wallet.Derive(path, false)
	if err != nil {
		return w, err
	}

	w = Wallet{account, wallet, mnemonic, keystorePassword}
	w.saveCredentials()
	return w, nil
}

func (w *Wallet) signTx(txHashBytes []byte) (signature []byte, err error) {
	privateKeyECDSA, err := w.Wallet.PrivateKey(w.Account)
	if err != nil {
		return nil, err
	}
	signature, err = crypto.Sign(txHashBytes, privateKeyECDSA)
	if err != nil {
		return nil, err
	}
	return
}

func (w *Wallet) saveCredentials() {
	privKeyHex, err := w.Wallet.PrivateKeyHex(w.Account)
	if err != nil {
		fmt.Println("Failed to save credentials. Error: ", err)
		return
	}

	pubKeyHex, err := w.Wallet.PublicKeyHex(w.Account)
	if err != nil {
		fmt.Println("Failed to save credentials. Error: ", err)
		return
	}
	file, err := os.Create("credentials.txt")
	if err != nil {
		fmt.Println("Failed to save credentials. Error: ", err)
		return
	}
	defer file.Close()
	credentials := fmt.Sprintf("pubKeyHex: %s\nprivKeyHex: %s\nmnemonic: %s\nkeystorePassword: %s\n",
		pubKeyHex, privKeyHex, w.mnemonic, w.keystorePassword)
	_, err = file.WriteString(credentials)
	if err != nil {
		fmt.Println("Failed to save credentials. Error: ", err)
		return
	}
	fmt.Println("Credentials successfully saved in credentials.txt")
}

// Transfer Execute transaction
func Transfer(passPhrase, toAddress string, amount int32) (txID string, err error) {
	data := TransferRequest{passPhrase, toAddress, amount}
	tReq, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	resp, err := http.Post("https://api.shasta.trongrid.io/wallet/easytransfer", "application/json", bytes.NewBuffer(tReq))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	var tResp TransferResponse
	err = json.Unmarshal(body, &tResp)
	if !tResp.Result.Result {
		return "", fmt.Errorf("something went wrong")
	}
	return tResp.Transaction.TxID, nil
}

type Wallet struct {
	Account          accounts.Account
	Wallet           *hdwallet.Wallet
	mnemonic         string
	keystorePassword string
}

type TransferResponse struct {
	Result      Result      `json:"result"`
	Transaction Transaction `json:"transaction"`
}
type Result struct {
	Result bool `json:"result"`
}
type Transaction struct {
	TxID string `json:"txID"`
}
type TransferRequest struct {
	PassPhrase string `json:"passPhrase"`
	ToAddress  string `json:"toAddress"`
	Amount     int32  `json:"amount"`
}
