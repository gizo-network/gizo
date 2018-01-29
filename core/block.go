package core

//FIXME: unexport all to avoid data modification
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gizo-network/gizo/core/merkletree"
	"github.com/gizo-network/gizo/helpers"

	"github.com/kpango/glg"
)

var (
	ErrUnableToExport   = errors.New("Unable to export block")
	ErrHashModification = errors.New("Attempt to modify hash value of block")
)

type Block struct {
	Header       BlockHeader              `json:"header"`
	Jobs         []*merkletree.MerkleNode `json:"jobs"`
	Height       uint64                   `json:"height"`
	ReceivedAt   int64                    `json:"received_at"`   //time it was received
	ReceivedFrom interface{}              `json:"received_from"` //node it received from
	By           string                   `json:"by"`            // id of node that generated block
}

func (b Block) GetHeader() BlockHeader {
	return b.Header
}

func (b *Block) SetHeader(h BlockHeader) {
	b.Header = h
}

func (b Block) GetJobs() []*merkletree.MerkleNode {
	return b.Jobs
}

func (b *Block) SetJobs(j []*merkletree.MerkleNode) {
	b.Jobs = j
}

func (b Block) GetHeight() uint64 {
	return b.Height
}

func (b *Block) SetHeight(h uint64) {
	b.Height = h
}

//FIXME: implement block status
func NewBlock(tree merkletree.MerkleTree, pHash []byte, height uint64, difficulty uint8) *Block {
	block := &Block{
		Header: BlockHeader{
			Timestamp:     time.Now().Unix(),
			PrevBlockHash: pHash,
			MerkleRoot:    tree.GetRoot(),
			Difficulty:    big.NewInt(int64(difficulty)),
		},
		Jobs:   tree.GetLeafNodes(),
		Height: height,
	}
	pow := NewPOW(block)
	pow.Run() //mines block
	err := block.export()
	if err != nil {
		glg.Fatal(err)
	}
	return block
}

//writes block on disk
func (b Block) export() error {
	InitializeDataPath()
	if b.IsEmpty() {
		return ErrUnableToExport
	}
	bBytes := b.Serialize()
	err := ioutil.WriteFile(path.Join(BlockPath, fmt.Sprintf(BlockFile, hex.EncodeToString(b.Header.GetHash()[:]))), []byte(helpers.Encode64(bBytes)), os.FileMode(0555))
	if err != nil {
		glg.Fatal(err)
	}
	return nil
}

// reads block into memory
func (b *Block) Import(hash []byte) {
	if b.IsEmpty() == false {
		glg.Warn("Overwriting umempty block")
	}
	read, err := ioutil.ReadFile(path.Join(BlockPath, fmt.Sprintf(BlockFile, hex.EncodeToString(hash))))
	if err != nil {
		glg.Fatal(err) //FIXME: handle block doesn't exist
	}
	bBytes := helpers.Decode64(string(read))
	temp, err := DeserializeBlock(bBytes)
	if err != nil {
		glg.Fatal(err)
	}
	b.SetHeader(temp.GetHeader())
	b.SetHeight(temp.GetHeight())
	b.SetJobs(temp.GetJobs())
}

func (b Block) FileStats() os.FileInfo {
	info, err := os.Stat(path.Join(BlockPath, fmt.Sprintf(BlockFile, hex.EncodeToString(b.Header.GetHash()))))
	if os.IsNotExist(err) {
		glg.Fatal("Block file doesn't exist")
	}
	return info
}

func (b *Block) IsEmpty() bool {
	return b.Header.GetHash() == nil
}

//Serialize returns bytes of block
func (b *Block) Serialize() []byte {
	temp, err := json.Marshal(*b)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}

//DeserializeBlock returns block of bytes
func DeserializeBlock(b []byte) (*Block, error) {
	var temp Block
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return nil, err
	}
	return &temp, nil
}

func (b *Block) VerifyBlock() bool {
	tree := merkletree.MerkleTree{Root: b.GetHeader().GetMerkleRoot(), LeafNodes: b.GetJobs()}
	mBytes, err := tree.Serialize()
	if err != nil {
		glg.Fatal(err)
	}
	data := bytes.Join(
		[][]byte{
			b.GetHeader().GetPrevBlockHash(),
			[]byte(strconv.FormatInt(b.GetHeader().GetTimestamp(), 10)),
			mBytes,
			[]byte(strconv.FormatInt(int64(b.GetHeader().GetNonce()), 10)),
			[]byte(strconv.FormatInt(int64(b.GetHeight()), 10)),
			[]byte(strconv.FormatInt(b.GetHeader().GetDifficulty().Int64(), 10)),
		},
		[]byte{},
	)
	hash := sha256.Sum256(data)
	return bytes.Equal(hash[:], b.GetHeader().GetHash())
}

//DeleteFile deletes block file on disk
func (b Block) DeleteFile() {
	err := os.Remove(path.Join(BlockPath, b.FileStats().Name()))
	if err != nil {
		glg.Fatal(err)
	}
}
