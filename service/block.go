package service

import (
	"encoding/json"
	"fmt"
	"storeth/data"
)

// GetBlockArgs is a set of args for Service.GetBlock().
type GetBlockArgs struct {
	Index *uint64 `json:"index"`
	Hash  *string `json:"hash"`
}

// GetLogsResult is a result for Service.GetLogs().
type GetBlockResult struct {
	Block json.RawMessage `json:"block"`
}

// GetLogs retrieves block contents.
func (s *Service) GetBlock(args GetBlockArgs) (*GetBlockResult, error) {
	if (args.Index == nil) == (args.Hash == nil) {
		return nil, fmt.Errorf("either index or hash (but not both) to be specified")
	}

	var block *data.Block
	var err error

	if args.Index != nil {
		block, err = data.GetBlock(s.db, *args.Index)
	} else {
		block, err = data.GetBlockByHash(s.db, *args.Hash)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get block: %v", err)
	}

	if block == nil {
		return nil, fmt.Errorf("block not found")
	}

	return &GetBlockResult{Block: block.Content}, nil
}
