package cache_test

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/gizo-network/gizo/cache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/job"
)

func TestJobCache(t *testing.T) {
	godotenv.Load()
	core.RemoveDataPath()
	j1 := job.NewJob(`
		func Factorial(n){
		 if(n > 0){
		  result = n * Factorial(n-1)
		  return result
		 }
		 return 1
		}`, "Factorial")
	j1.AddExec(job.Exec{})
	j1.AddExec(job.Exec{})
	j1.AddExec(job.Exec{})
	j1.AddExec(job.Exec{})
	j2 := job.NewJob(`
			func Factorial(n){
			 if(n > 0){
			  result = n * Factorial(n-1)
			  return result
			 }
			 return 1
			}`, "Factorial")
	j2.AddExec(job.Exec{})
	j2.AddExec(job.Exec{})
	j2.AddExec(job.Exec{})
	j2.AddExec(job.Exec{})
	j3 := job.NewJob(`
				func Factorial(n){
				 if(n > 0){
				  result = n * Factorial(n-1)
				  return result
				 }
				 return 1
				}`, "Factorial")
	j3.AddExec(job.Exec{})
	j3.AddExec(job.Exec{})
	j3.AddExec(job.Exec{})
	j3.AddExec(job.Exec{})
	j3.AddExec(job.Exec{})
	j3.AddExec(job.Exec{})

	node1 := merkletree.NewNode(*j1, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node2 := merkletree.NewNode(*j2, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	node3 := merkletree.NewNode(*j3, &merkletree.MerkleNode{}, &merkletree.MerkleNode{})
	tree1 := merkletree.NewMerkleTree([]*merkletree.MerkleNode{node1, node3})
	tree2 := merkletree.NewMerkleTree([]*merkletree.MerkleNode{node2, node1})
	tree3 := merkletree.NewMerkleTree([]*merkletree.MerkleNode{node3, node2})
	bc := core.CreateBlockChain()
	blk1 := core.NewBlock(*tree1, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetNextHeight(), 10)
	bc.AddBlock(blk1)
	blk2 := core.NewBlock(*tree2, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetNextHeight(), 10)
	bc.AddBlock(blk2)
	blk3 := core.NewBlock(*tree3, bc.GetLatestBlock().GetHeader().GetHash(), bc.GetNextHeight(), 10)
	bc.AddBlock(blk3)
	c := cache.NewJobCache(bc)
	cj1, err := c.Get(j1.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, cj1)
	cj2, err := c.Get(j2.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, cj2)
	cj3, err := c.Get(j3.GetID())
	assert.NoError(t, err)
	assert.NotNil(t, cj3)
	assert.False(t, c.IsFull())
}
