//go:generate reform

package data

import (
	"encoding/json"

	"gopkg.in/reform.v1"
)

// Block contains Ethereum block data.
//
//reform:block
type Block struct {
	Index   uint64          `reform:"index,pk"`
	Hash    string          `reform:"hash"`
	Content json.RawMessage `reform:"content"`
}

// AddBlock adds a new block. Drops an oldest one if maxNumBlocks
// is specified and the new number of blocks exceeds this limit.
func AddBlock(db *reform.DB, block *Block, maxNumBlocks *int) error {
	return nil
}

// GetBlock retrieves block by index. Returns nil if no block is found.
func GetBlock(db *reform.DB, index uint64) (*Block, error) {
	return nil, nil
}

// GetBlock retrieves block by hash. Returns nil if no block is found.
func GetBlockByHash(db *reform.DB, hash string) (*Block, error) {
	return nil, nil
}
