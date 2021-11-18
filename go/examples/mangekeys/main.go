package main

import (
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

func main() {
	_, priv, err := ed25519.GenerateKey(crypto_rand.Reader)
	if err != nil {
		panic(err)
	}
	account := solana.Wallet{PrivateKey: solana.PrivateKey(priv)}
	fmt.Println("account private key:", account.PrivateKey.String())
	fmt.Println("account public key:", account.PublicKey())

	pk := account.PrivateKey.String()
	privateKey, err := base58.Decode(pk)

	parsed := solana.Wallet{PrivateKey: privateKey}
	fmt.Println("account private key:", account.PrivateKey.String())
	fmt.Println("account public key:", parsed.PublicKey())

	//account := solana.NewWallet()
	//fmt.Println("account private key:", account.PrivateKey.String())
	//fmt.Println("account public key:", account.PublicKey())
	//
	//// Output
	//// account private key: 5So4WjCytJi4SM9J91n4hHURiiKneKck7oUdr3RnKVKELivdoiSwaECkhHn8kDayGyKZHLXBwFXWJcEATpxbXRWE
	//// account public key: Cp71tMR9icw3m5xLoudXcGuYb99VbJfRK6q5EmjHgXDg
	//
	//pk := account.PrivateKey.String()
	//privateKey, err := base58.Decode(pk)
	//if err != nil {
	//	panic(err)
	//}
	//
	//parsed := solana.Wallet{PrivateKey: privateKey}
	//fmt.Println("account private key:", account.PrivateKey.String())
	//fmt.Println("account public key:", parsed.PublicKey())
}
