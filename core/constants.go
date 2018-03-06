package core

import (
	"os"
	"path"
)

var BlockPathProd = path.Join(os.Getenv("HOME"), ".gizo", "blocks")
var IndexPathProd = path.Join(os.Getenv("HOME"), ".gizo")

var BlockPathDev = path.Join(os.Getenv("HOME"), ".gizo-dev", "blocks")
var IndexPathDev = path.Join(os.Getenv("HOME"), ".gizo-dev")

const BlockFile = "%s.blk" //format of block filenames
const BlockBucket = "blocks"
const IndexDB = "bc_%s.db" // node id
