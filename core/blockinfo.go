package core

import (
	"encoding/json"

	"github.com/kpango/glg"
)

//BlockInfo - model of data written to embedded database
type BlockInfo struct {
	Header    BlockHeader
	Height    uint64
	TotalJobs uint
	FileName  string
	FileSize  int64
}

//sets blockinfo header
func (bi *BlockInfo) setHeader(bh BlockHeader) {
	bi.Header = bh
}

//GetHeader returns block header
func (bi BlockInfo) GetHeader() BlockHeader {
	return bi.Header
}

//sets height
func (bi *BlockInfo) setHeight(h uint64) {
	bi.Height = h
}

//GetHeight return height
func (bi BlockInfo) GetHeight() uint64 {
	return bi.Height
}

//sets total jobs
func (bi *BlockInfo) setTotalJobs(t uint) {
	bi.TotalJobs = t
}

//GetTotalJobs returns total jobs
func (bi BlockInfo) GetTotalJobs() uint {
	return bi.TotalJobs
}

//sets filename
func (bi *BlockInfo) setFileName(n string) {
	bi.FileName = n
}

//GetFileName returns filename
func (bi BlockInfo) GetFileName() string {
	return bi.FileName
}

//sets filename
func (bi *BlockInfo) setFileSize(s int64) {
	bi.FileSize = s
}

//GetFileSize returns the blocks filesize
func (bi BlockInfo) GetFileSize() int64 {
	return bi.FileSize
}

//Serialize returns the blockinfo in json bytes
func (bi *BlockInfo) Serialize() []byte {
	temp, err := json.Marshal(*bi)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}

//GetBlock - imports block from file into memory
func (bi BlockInfo) GetBlock() *Block {
	var temp Block
	temp.Import(bi.GetHeader().GetHash())
	return &temp
}

//DeserializeBlockInfo return blockinfo
func DeserializeBlockInfo(bi []byte) *BlockInfo {
	var temp BlockInfo
	err := json.Unmarshal(bi, &temp)
	if err != nil {
		glg.Fatal(err)
	}
	return &temp
}
