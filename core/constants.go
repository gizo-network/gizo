package core

import (
	"os"
	"path"
)

//BlockPathProd is the path block files are saved on the disk for production
var BlockPathProd = path.Join(os.Getenv("HOME"), ".gizo", "blocks")

//IndexPathProd is the path database files are saved on the disk for production
var IndexPathProd = path.Join(os.Getenv("HOME"), ".gizo")

//BlockPathDev is the path block files are saved on the disk for development
var BlockPathDev = path.Join(os.Getenv("HOME"), ".gizo-dev", "blocks")

//IndexPathDev is the path database files are saved on the disk for development
var IndexPathDev = path.Join(os.Getenv("HOME"), ".gizo-dev")

//BlockFile is the format of block filenames
const BlockFile = "%s.blk"

//BlockBucket is the name of the bucket for the blockchain database
const BlockBucket = "blocks"

//IndexDB is the database file for the node
const IndexDB = "bc_%s.db" // node id
