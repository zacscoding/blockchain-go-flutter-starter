package main

import (
	"fmt"
	"github.com/gagliardetto/solana-go"
)

func main() {
	// Create a new account:
	// account := solana.NewWallet()
	account, _ := solana.WalletFromPrivateKeyBase58("2mUNxGyfxbqhQLT94nnXEjfUfCKVHAE8LYYkdrbpTrG4SGWeLhRrcq8kPBgzt2bRwE7mi4DqRbAL8hY5XyDAUxM8")
	fmt.Println("account private key:", account.PrivateKey)
	fmt.Println("account public key:", account.PublicKey())

	//// Create a new RPC client:
	//client := rpc.New(rpc.TestNet_RPC)
	//
	//// Airdrop 5 SOL to the new account:
	//out, err := client.RequestAirdrop(
	//	context.TODO(),
	//	account.PublicKey(),
	//	solana.LAMPORTS_PER_SOL*5,
	//	rpc.CommitmentFinalized,
	//)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("airdrop transaction signature:", out)
	// account private key: 2mUNxGyfxbqhQLT94nnXEjfUfCKVHAE8LYYkdrbpTrG4SGWeLhRrcq8kPBgzt2bRwE7mi4DqRbAL8hY5XyDAUxM8
	// account public key: 3rffGGn7WjZk6qVZE1qz4Nqc11CAgwr5bBTuvrYSfZGE
	// airdrop transaction signature: 3aEjM2R5Lz71iAdiZpZLtKgsnVyEuHaKJEAbJwJabG4FEEY8nBvCjEHuRmmtBqsgJj1Yct4LcWxeoS79s6DeuEzY
	// 4ETf86tK7b4W72f27kNLJLgRWi9UfJjgH4koHGUXMFtn
	// 3rffGGn7WjZk6qVZE1qz4Nqc11CAgwr5bBTuvrYSfZGE
}
