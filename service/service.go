package service

import "gopkg.in/reform.v1"

// Service is a service for storing Ethereum data.
type Service struct {
	db *reform.DB
}

// NewService creates a new Service instance.
func NewService(db *reform.DB) *Service {
	return &Service{db}
}
