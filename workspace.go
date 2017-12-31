package main

import (
	"fmt"

	"github.com/gizo-network/gizo/core"
)

func main() {
	block := core.NewBlock([]byte("jobs example"), []byte("genesis block"), []byte("merkle root"))
	block.SetHash()
	// fmt.Println(helpers.MarshalBlock(*block))
	// fmt.Println(block.VerifyBlock())
	// err := block.SetHash()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(helpers.MarshalBlock(*block))

	blockchain := core.NewBlockChain()
	blockchain.AddBlock([]byte("jobs example"), []byte("merkleshash"))
	blockchain.AddBlock([]byte("jobs example"), []byte("merklehash"))

	fmt.Println(blockchain.VerifyBlockChain())

	blockchain.Blocks[1].Nonce = 40
	fmt.Println(blockchain.VerifyBlockChain())
}
