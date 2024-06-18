package service

import (
	"encoding/json"
	"fmt"
	"storeth/data"

	"github.com/ethereum/go-ethereum/common"
)

// GetLogsArgs is a set of args for Service.GetLogs().
type GetLogsArgs struct {
	Address   *common.Address `json:"address"`
	FromBlock *uint64         `json:"fromBlock"`
	ToBlock   *uint64         `json:"toBlock"` // Not including ToBlock.
}

// GetLogsResult is a result for Service.GetLogs().
type GetLogsResult struct {
	Logs []json.RawMessage `json:"logs"`
}

// GetLogs retrieves event logs for a given address.
func (s *Service) FindLogs(args GetLogsArgs) (*GetLogsResult, error) {
	logs, err := data.FindLogs(
		s.db.Querier,
		&data.FindLogsFilter{
			Address:   args.Address,
			FromBlock: args.FromBlock,
			ToBlock:   args.ToBlock,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to find logs: %v", err)
	}

	contents := make([]json.RawMessage, len(logs))
	for i, val := range logs {
		contents[i] = val.Content
	}

	return &GetLogsResult{Logs: contents}, nil
}
