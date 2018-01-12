package core

import (
	"os"
	"path"
)

//TODO: allow dynamic data path
var BlockPath = path.Join(os.Getenv("HOME"), "gizo", "blocks")
var IndexPath = path.Join(os.Getenv("HOME"), "gizo", "index")

const BlockFile = "blk_%s.db" //holds blocks
const BlockBucket = "blocks"
const BlockIndex = "blk_index.db"
