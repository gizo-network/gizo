package core

import (
	"os"
	"path"
)

//TODO: allow dynamic data path
var BlockPath = path.Join(os.Getenv("HOME"), ".gizo", "blocks")
var IndexPath = path.Join(os.Getenv("HOME"), ".gizo")

const BlockFile = "%s.blk" // block hash
const BlockBucket = "blocks"
const IndexDB = "bc_%s.db" // node id
