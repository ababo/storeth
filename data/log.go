//go:generate reform

package data

import (
	"encoding/json"

	"gopkg.in/reform.v1"
)

// Log contains Ethereum event log data.
//
//reform:log
type Log struct {
	Address string          `reform:"address"`
	Block   uint64          `reform:"block"`
	Content json.RawMessage `reform:"content"`
}

// AddLog adds a new event log.
func AddLog(db *reform.DB, log *Log) error {
	return nil
}

// FindLogs finds logs for a specified address and optionally for a block index range.
func FindLogs(db *reform.DB, address string, fromBlock *uint64, toBlock *uint64) ([]Log, error) {
	return nil, nil
}
