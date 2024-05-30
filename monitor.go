package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"storeth/data"
	"storeth/service"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/reform.v1"
)

func monitorEth(conf *config, db *reform.DB, svc *service.Service) {
	storedFrom, storedNum, err := data.GetBlockRange(db.Querier)
	if err != nil {
		log.Fatalf("failed to get stored block range: %v", err)
	}

	var latestStored *uint64 = nil
	if storedNum > 0 {
		latestStored = new(uint64)
		*latestStored = storedFrom + storedNum - 1
		log.Printf("stored blocks %d-%d", storedFrom, *latestStored)
	} else {
		log.Printf("no blocks stored yet")
	}

	client, err := ethclient.Dial(conf.EthWSEndpoint)
	if err != nil {
		log.Fatalf("failed to connect to eth client: %v", err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("failed to subscribe to new eth blocks: %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("failed to receive new eth block notification: %v", err)
		case header := <-headers:
			index := header.Number.Uint64()
			if err := processNextBlock(conf, db, svc, client, index, latestStored); err != nil {
				log.Fatalf("failed to process next block: %v", err)
			}
			latestStored = nil // We've already fetched the missing blocks.
		}
	}
}

type rpcBlock struct {
	Hash common.Hash `json:"hash"`
}

func processNextBlock(
	conf *config,
	db *reform.DB,
	svc *service.Service,
	client *ethclient.Client,
	index uint64,
	latestStored *uint64) error {
	// Fetch missing blocks.
	for latestStored != nil && *latestStored < index-1 {
		*latestStored++
		if err := addBlock(conf, db, client, *latestStored); err != nil {
			return fmt.Errorf("failed to add missing block %d: %v", *latestStored, err)
		}
		service.NotifyNewBlock(svc, *latestStored)
	}

	// Fetch the block we just got notified about.
	if err := addBlock(conf, db, client, index); err != nil {
		return fmt.Errorf("failed to add next block %d: %v", index, err)
	}

	service.NotifyNewBlock(svc, index)

	return nil
}

func addBlock(conf *config, db *reform.DB, client *ethclient.Client, index uint64) error {
	// Use lower level API to retrieve raw block and log JSON.

	var blockJSON json.RawMessage
	indexHex := fmt.Sprintf("0x%x", index)
	if err := client.Client().CallContext(context.Background(),
		&blockJSON, "eth_getBlockByNumber", indexHex, true); err != nil {
		return fmt.Errorf("failed to fetch eth block: %v", err)
	}

	var rpcBlock rpcBlock
	if err := json.Unmarshal(blockJSON, &rpcBlock); err != nil {
		return fmt.Errorf("failed to unmarshal rpc block: %v", err)
	}

	block := data.Block{
		Index:   index,
		Hash:    rpcBlock.Hash.Hex(),
		Content: blockJSON,
	}

	query := map[string]interface{}{
		"fromBlock": indexHex,
		"toBlock":   indexHex,
	}

	var logJSONs []json.RawMessage
	if err := client.Client().CallContext(
		context.Background(), &logJSONs, "eth_getLogs", query); err != nil {
		return fmt.Errorf("failed to fetch eth logs: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin tx: %v", err)
	}

	if err := data.AddBlock(tx.Querier, &block, conf.MaxNumBlocks); err != nil {
		return fmt.Errorf("failed to store block: %v", err)
	}

	for _, logJSON := range logJSONs {
		var ethLog types.Log
		if err := json.Unmarshal(logJSON, &ethLog); err != nil {
			return fmt.Errorf("failed to unmarshal log: %v", err)
		}

		log := data.Log{
			Address: ethLog.Address.Hex(),
			Block:   index,
			Content: logJSON,
		}
		if err := data.AddLog(tx.Querier, &log); err != nil {
			return fmt.Errorf("failed to store log: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %v", err)
	}

	log.Printf("added block %d", index)
	return nil
}
