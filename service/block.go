package service

import (
	"context"
	"encoding/json"
	"fmt"
	"storeth/data"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// GetBlockArgs is a set of args for Service.GetBlock().
type GetBlockArgs struct {
	Index *uint64      `json:"index"`
	Hash  *common.Hash `json:"hash"`
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
		block, err = data.GetBlock(s.db.Querier, *args.Index)
	} else {
		block, err = data.GetBlockByHash(s.db.Querier, *args.Hash)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get block: %v", err)
	}

	if block == nil {
		return nil, fmt.Errorf("block not found")
	}

	return &GetBlockResult{Block: block.Content}, nil
}

// GetBlockRangeResult is a result for Service.GetBlockRange().
type GetBlockRangeResult struct {
	FromIndex uint64 `json:"fromIndex"`
	ToBlock   uint64 `json:"toBlock"`
}

// GetBlockRange returns a range of stored blocks.
func (s *Service) GetBlockRange() (*GetBlockRangeResult, error) {
	from, to, err := data.GetBlockRange(s.db.Querier)
	if err != nil {
		return nil, fmt.Errorf("failed to get block range: %v", err)
	}
	return &GetBlockRangeResult{FromIndex: from, ToBlock: to}, nil
}

// NewBlocks creates a subscription to new block notifications.
func (s *Service) NewBlocks(ctx context.Context) (*rpc.Subscription, error) {
	notifier, ok := rpc.NotifierFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to create subscription")
	}

	sub := notifier.CreateSubscription()

	s.mutex.Lock()
	s.newBlocks[notifier] = sub.ID
	s.mutex.Unlock()

	go func() {
		<-sub.Err()
		s.mutex.Lock()
		delete(s.newBlocks, notifier)
		s.mutex.Unlock()
	}()

	return sub, nil
}
