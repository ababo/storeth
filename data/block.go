//go:generate reform

package data

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/reform.v1"
)

// Block contains Ethereum block data.
//
//reform:block
type Block struct {
	Index   uint64          `reform:"index,pk"`
	Hash    common.Hash     `reform:"hash"`
	Content json.RawMessage `reform:"content"`
}

// AddBlock adds a new block. Drops an oldest one if maxNumBlocks is
// specified and the number of blocks after addition exceeds this limit.
func AddBlock(tx *reform.TX, block *Block, maxNumBlocks *uint64) error {
	from, to, err := GetBlockRange(tx.Querier)
	if err != nil {
		return err
	}

	if maxNumBlocks == nil || to-from < *maxNumBlocks {
		_, err = tx.Exec(`
			INSERT INTO block
			VALUES ($1, $2, $3)
		`, block.Index, block.Hash, block.Content)
	} else { // Avoid VACUUM.
		_, err = tx.Exec(`
			UPDATE block
			   SET index = $1,
			       hash = $2,
				   content = $3
		     WHERE index = $4
		`, block.Index, block.Hash, block.Content, from)
	}

	return err
}

// GetBlock retrieves block by index. Returns nil if no block is found.
func GetBlock(querier *reform.Querier, index uint64) (*Block, error) {
	var block Block
	if err := querier.FindByPrimaryKeyTo(&block, index); err != nil {
		if err == reform.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &block, nil
}

// GetBlock retrieves block by hash. Returns nil if no block is found.
func GetBlockByHash(querier *reform.Querier, hash common.Hash) (*Block, error) {
	var block Block
	if err := querier.FindOneTo(&block, "hash", hash); err != nil {
		if err == reform.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &block, nil
}

// GetBlockRange returns [fromIndex, toIndex) interval of stored block indices.
func GetBlockRange(querier *reform.Querier) (fromIndex uint64, toIndex uint64, err error) {
	row := querier.QueryRow(`
		SELECT MIN(index), MAX(index)
		  FROM block
	`)

	var min, max *uint64
	if err := row.Scan(&min, &max); err != nil {
		return 0, 0, err
	}

	if min == nil {
		return 0, 0, nil
	}

	return *min, *max + 1, nil
}
