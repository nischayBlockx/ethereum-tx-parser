package parser

import (
	"context"
	"log"

	"github.com/trust-assignment/internal/models"
	repo "github.com/trust-assignment/internal/repository"
	"github.com/trust-assignment/internal/service/scannersvc"
	"github.com/trust-assignment/pkg/ethclient"
)

// ParserService represents a service for parsing and managing transactions.
type ParserService struct {
	Db      repo.DBInterface           // Database interface for managing subscribers and transactions
	Scansvc *scannersvc.ScannerService // Scanner service for retrieving and updating blockchain transactions
}

func NewParser(ctx context.Context, endpoint string, startAtBlock int) *ParserService {
	data := repo.NewDB()
	ethclt := ethclient.NewEthClient(endpoint)
	scan := scannersvc.NewScanner(ctx, data, ethclt, startAtBlock)
	return &ParserService{
		Db:      data,
		Scansvc: scan,
	}
}

// Subscribe adds a new subscriber with the given address to the database.
func (p *ParserService) Subscribe(address string) bool {
	if err := p.Db.AddSubscriber(context.Background(), address); err != nil {
		log.Println("[Parser] Error subscribing address: ", err)
		return false
	}
	return true
}

// GetTransactions returns a list of inbound or outbound transactions for an address.
func (p *ParserService) GetTransactions(address string) []models.Transaction {
	txns, err := p.Db.GetTxns(context.Background(), address)
	if err != nil {
		log.Printf("[Parser] Error getting transactions for address %s: %v", address, err)
		return nil
	}
	return txns
}
