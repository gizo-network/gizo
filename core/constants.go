package core

import (
	"os"
	"path"
)

var BlockPath = path.Join(os.Getenv("HOME"), "gizo", "blocks")
var IndexPath = path.Join(os.Getenv("HOME"), "gizo", "index")

const BlockFile = "blk_%s.db"
const BlockBucket = "blocks"
const BlockIndex = "blk_index.db"
