package repository

import (
	"context"

	"github.com/trust-assignment/internal/models"
)

type DBInterface interface {
	AddSubscriber(ctx context.Context, address string) error
	SaveTxns(ctx context.Context, txns map[string][]models.Transaction) error
	CheckTxns(ctx context.Context, address string) (bool, error)
	GetTxns(ctx context.Context, address string) ([]models.Transaction, error)
	DeleteSub(ctx context.Context, address string)
}
