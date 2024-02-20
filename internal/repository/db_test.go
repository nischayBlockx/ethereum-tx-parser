package repository

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/trust-assignment/internal/models"
)

func TestMemoryDb(t *testing.T) {
	db := NewDB()
	defer db.Close()
	// Test AddSubscriber
	address := "0x1234567890abcdef"
	err := db.AddSubscriber(context.Background(), address)
	if err != nil {
		t.Errorf("AddSubscriber failed: %v", err)
	}

	// Try adding the same address again
	err = db.AddSubscriber(context.Background(), address)
	if err == nil {
		t.Error("AddSubscriber should return an error for existing address, but it did not")
	}

	// Test CheckTxns for existing address
	exists, _ := db.CheckTxns(context.Background(), address)
	if !exists {
		t.Error("CheckTxns failed for existing address.")
	}

	// Test GetTxns for existing address
	transactions, err := db.GetTxns(context.Background(), address)
	if err != nil || len(transactions) != 0 {
		t.Errorf("GetTxns failed for existing address. Expected: [], Got: %v", transactions)
	}

	// Test SaveTxns
	newTxs := map[string][]models.Transaction{
		address: {
			{Hash: "tx1", From: "0xabcdef", To: "0x123456", Value: big.NewInt(100)},
			{Hash: "tx2", From: "0x123456", To: "0xabcdef", Value: big.NewInt(50)},
		},
	}
	err = db.SaveTxns(context.Background(), newTxs)
	if err != nil {
		t.Errorf("SaveTxns failed: %v", err)
	}

	// Test GetTxns for existing address after saving transactions
	transactions, err = db.GetTxns(context.Background(), address)
	expectedTransactions := []models.Transaction{
		{Hash: "tx1", From: "0xabcdef", To: "0x123456", Value: big.NewInt(100)},
		{Hash: "tx2", From: "0x123456", To: "0xabcdef", Value: big.NewInt(50)},
	}
	if err != nil || !reflect.DeepEqual(transactions, expectedTransactions) {
		t.Errorf("GetTxns failed after saving transactions. Expected: %v, Got: %v", expectedTransactions, transactions)
	}

	// Test DeleteSub
	db.DeleteSub(context.Background(), address)
	exists, _ = db.CheckTxns(context.Background(), address)
	if exists {
		t.Errorf("DeleteSub failed. The address should have been deleted, but it still exists")
	}
}
