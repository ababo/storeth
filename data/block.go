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

// AddBlock adds a new block. Drops an oldest one if maxNumBlocks is
// specified and the number of blocks after addition exceeds this limit.
func AddBlock(querier *reform.Querier, block *Block, maxNumBlocks *uint64) error {
	return nil
}

// GetBlock retrieves block by index. Returns nil if no block is found.
func GetBlock(querier *reform.Querier, index uint64) (*Block, error) {
	return nil, nil
}

// GetBlock retrieves block by hash. Returns nil if no block is found.
func GetBlockByHash(querier *reform.Querier, hash string) (*Block, error) {
	return nil, nil
}

// GetBlockRange returns a range of stored blocks.
func GetBlockRange(querier *reform.Querier) (fromIndex uint64, numBlocks uint64, err error) {
	return 0, 0, nil
}
