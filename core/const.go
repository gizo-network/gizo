package core

import (
	"math/big"
)

//! modify on job engine creation
var GenesisBlock *Block = NewBlock([]byte("created block engine"), []byte("0"), []byte("0"))
var MaxTreeJobs *big.Int = big.NewInt(24)
