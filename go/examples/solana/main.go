package main

import (
	"context"
	"encoding/json"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
)

func main() {
	//endpoint := "http://localhost:8899"
	endpoint := rpc.MainNetBeta_RPC
	cli := rpc.New(endpoint)

	checkBlocks(cli)
}

func checkBlocks(cli *rpc.Client) {
	height, err := cli.GetBlockHeight(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Check block height: %d", height)

	for i := 0; i < 1; i++ {
		bn := height - uint64(i)
		log.Printf("Check block: %d", bn)

		block, err := cli.GetBlockWithOpts(context.Background(), bn, &rpc.GetBlockOpts{
			TransactionDetails: "none",
		})
		if err != nil {
			log.Printf("> Error: %v", err)
			continue
		}
		log.Printf(">\n%s", toJson(block, true))
	}
}

func toJson(v interface{}, pretty bool) string {
	var (
		b   []byte
		err error
	)
	if pretty {
		b, err = json.MarshalIndent(v, "", "  ")
	} else {
		b, err = json.Marshal(v)
	}
	if err != nil {
		err.Error()
	}
	return string(b)
}
