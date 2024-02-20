package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/trust-assignment/internal/models"
)

// MemoryDb represents an in-memory database.
type MemoryDb struct {
	Db map[string][]models.Transaction // Internal storage for transactions, indexed by address
	mu *sync.RWMutex                   // Mutex for concurrent access to the database
}

// NewDB creates and returns a new instance of MemoryDb.
func NewDB() *MemoryDb {
	return &MemoryDb{
		Db: make(map[string][]models.Transaction),
		mu: &sync.RWMutex{},
	}
}

// AddSubscriber adds a new subscriber with the given address to the database.
func (m *MemoryDb) AddSubscriber(ctx context.Context, address string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	address = strings.ToLower(address)
	if _, ok := m.Db[address]; ok {
		return fmt.Errorf("[DB-error] Address already exists")
	}
	m.Db[address] = []models.Transaction{}
	return nil
}

// CheckTxns checks if transactions exist for the specified address in the database.
func (m *MemoryDb) CheckTxns(ctx context.Context, address string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	address = strings.ToLower(address)

	if _, ok := m.Db[address]; ok {
		return true, nil
	}
	return false, fmt.Errorf("[DB-error] Address does not exist")
}

// GetTxns retrieves transactions for the specified address from the database.
func (m *MemoryDb) GetTxns(ctx context.Context, address string) ([]models.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	address = strings.ToLower(address)

	if txns, ok := m.Db[address]; ok {
		result := make([]models.Transaction, len(txns))
		copy(result, txns)
		return result, nil
	}
	return nil, fmt.Errorf("[DB-error] Address not found")
}

// SaveTxns saves new transactions for multiple addresses to the database.
func (m *MemoryDb) SaveTxns(ctx context.Context, newTxs map[string][]models.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for address, txs := range newTxs {
		address = strings.ToLower(address)

		// Check if the address exists in the memory database
		if _, ok := m.Db[address]; !ok {
			return fmt.Errorf("[DB-error] Address does not exist")
		}

		// Append the new transactions to the existing transactions for the address
		m.Db[address] = append(m.Db[address], txs...)
	}

	return nil
}

// DeleteSub removes a subscriber with the specified address from the database.
func (m *MemoryDb) DeleteSub(ctx context.Context, address string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	address = strings.ToLower(address)
	delete(m.Db, address)
}

// Close deallocates the internal map to free resources.
func (m *MemoryDb) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Db = nil
}
