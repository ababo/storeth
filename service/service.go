package service

import (
	"sync"

	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/reform.v1"
)

// Service is a service for storing Ethereum data.
type Service struct {
	db        *reform.DB
	mutex     sync.Mutex
	newBlocks map[*rpc.Notifier]rpc.ID
}

// NewService creates a new Service instance.
func NewService(db *reform.DB) *Service {
	return &Service{
		db:        db,
		mutex:     sync.Mutex{},
		newBlocks: make(map[*rpc.Notifier]rpc.ID),
	}
}

// NotifyNewBlock notifies newBlocks subscribers about new block index.
func NotifyNewBlock(service *Service, index uint64) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	for notifier, id := range service.newBlocks {
		_ = notifier.Notify(id, index)
	}
}
