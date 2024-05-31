//go:generate reform

package data

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/reform.v1"
)

// Log contains Ethereum event log data.
//
//reform:log
type Log struct {
	ID      uint            `reform:"id,pk"`
	Address common.Address  `reform:"address"`
	Block   uint64          `reform:"block"`
	Content json.RawMessage `reform:"content"`
}

// AddLog adds a new event log.
func AddLog(querier *reform.Querier, log *Log) error {
	result, err := querier.Exec(`
		UPDATE log
		   SET address = $1,
		       block = $2,
			   content = $3
		 WHERE id = (
			SELECT id
			  FROM log
			 WHERE block IS NULL
			 LIMIT 1
		 )
	`, log.Address, log.Block, log.Content)
	if err != nil {
		return err
	}

	if num, err := result.RowsAffected(); err != nil {
		return err
	} else if num > 0 {
		return nil
	}

	_, err = querier.Exec(`
		INSERT INTO log(address, block, content)
		VALUES ($1, $2, $3)
	`, log.Address, log.Block, log.Content)

	return err
}

// FindLogsFilter specifies log filtering rules.
type FindLogsFilter struct {
	Address   *common.Address
	FromBlock *uint64
	ToBlock   *uint64 // Not including ToBlock.
}

// FindLogs finds logs for specified filtering rules.
func FindLogs(querier *reform.Querier, filter *FindLogsFilter) ([]Log, error) {
	structs, err := querier.SelectAllFrom(LogTable, `
		WHERE ($1::bytea IS NULL OR address = $1) AND
		      ($2::bigint IS NULL OR block >= $2) AND
		      ($3::bigint IS NULL OR block < $3)`,
		filter.Address, filter.FromBlock, filter.ToBlock)
	if err != nil {
		return nil, err
	}

	logs := make([]Log, len(structs))
	for i, s := range structs {
		logs[i] = *s.(*Log)
	}

	return logs, nil
}
