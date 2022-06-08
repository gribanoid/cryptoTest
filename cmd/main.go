package main

import (
	"fmt"
	"test/tron"
)

func main() {
	keystorePassword := "password"
	derivationPath := "m/44'/195'/0'/0/0"
	wallet, err := tron.GenerateWallet(keystorePassword, derivationPath)
	if err != nil {
		fmt.Printf("An error occurred while creating the wallet.\nError: %v", err)
		return
	}

	txID, err := tron.Transfer("7465737470617373776f7264", "41D1E7A6BC354106CB410E65FF8B181C600FF14292", 10)
	if err != nil {
		fmt.Printf("An error occurred while transfering.\nError: %v", err)
		return
	}
	fmt.Println(wallet, txID)
}
