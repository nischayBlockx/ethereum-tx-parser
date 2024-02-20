package parser

import (
	"github.com/trust-assignment/internal/models"
)

type ParserServiceInterface interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []models.Transaction
}
