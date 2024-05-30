package service

import (
	"encoding/json"
	"fmt"
	"storeth/data"
)

// GetLogsArgs is a set of args for Service.GetLogs().
type GetLogsArgs struct {
	Address   string  `json:"address"`
	FromBlock *uint64 `json:"fromBlock"`
	NumBlocks *uint64 `json:"numBlocks"`
}

// GetLogsResult is a result for Service.GetLogs().
type GetLogsResult struct {
	Logs []json.RawMessage `json:"logs"`
}

// GetLogs retrieves event logs for a given address.
func (s *Service) GetLogs(args GetLogsArgs) (*GetLogsResult, error) {
	logs, err := data.FindLogs(s.db.Querier, args.Address, args.FromBlock, args.NumBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to find logs: %v", err)
	}

	contents := make([]json.RawMessage, len(logs))
	for _, val := range logs {
		contents = append(contents, val.Content)
	}

	return &GetLogsResult{Logs: contents}, nil
}
