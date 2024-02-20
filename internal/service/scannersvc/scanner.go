package scannersvc

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/trust-assignment/internal/models"
	repo "github.com/trust-assignment/internal/repository"
	"github.com/trust-assignment/pkg/ethclient"
)

type ScannerService struct {
	ctx              context.Context
	Db               repo.DBInterface
	Client           *ethclient.EthClient
	lastScannedBlock int
	once             sync.Once
}

func NewScanner(ctx context.Context, db repo.DBInterface, client *ethclient.EthClient, startAt int) *ScannerService {
	fmt.Println("[Scanner] Scanner set to start at block: ", startAt)
	return &ScannerService{
		ctx:              ctx,
		Db:               db,
		Client:           client,
		lastScannedBlock: startAt,
	}
}

// StartScan spawn a goroutine that will run the block scanning process
// at the given interval. It will stop the process when a signal is received
// on the shutdown channel. It will spawn only one goroutine for each Blockscan
// instance.
func (s *ScannerService) StartScan(interval time.Duration) {
	fmt.Println("[Scanner Service] StartScan method called")
	s.once.Do(func() {
		go func() {
			ticker := time.NewTicker(interval)
			for {
				select {
				case <-s.ctx.Done():
					ticker.Stop()
					fmt.Println("[Scanner] stopping blockscan")
					return
				case <-ticker.C:
					ticker.Stop()
					for scannedBlock, err := s.Run(s.ctx); scannedBlock != 0 || err != nil; scannedBlock, err = s.Run(s.ctx) {
						if err != nil {
							fmt.Println(fmt.Errorf("[Scanner] error scanning block: %s", err))
							continue
						}
					}
					ticker.Reset(interval)
					fmt.Printf("[Scanner] last scanned block %d\n", s.GetCurrentBlock())
				}
			}
		}()
	})
}

// Run starts the block scanning process. It will return the number
// of the last scanned block and an error if any. In case of no pending
// blocks to be scanned it will return 0.
func (s *ScannerService) Run(ctx context.Context) (int, error) {
	headBlock, err := s.Client.BlockNumber() //Step1.  get the current latest block
	fmt.Println("Headblock", headBlock)
	if err != nil {
		fmt.Println("[Scanner] Error querying head block : ", err)
		return 0, err
	}

	nextBlock := nextBlock(s.lastScannedBlock, headBlock) // // Step2. Get the next block
	fmt.Println("Nextblock", nextBlock)
	if nextBlock == 0 {
		return 0, nil
	}

	txs, err := s.ScanBlock(ctx, nextBlock) // // Step3. Get the transaction of next block
	if err != nil {
		fmt.Println("[Scanner] Error scanning block: ", err)
		return 0, err
	}

	s.Db.SaveTxns(ctx, txs)
	s.lastScannedBlock = nextBlock

	return s.lastScannedBlock, nil
}

func nextBlock(lastScannedBlock, headBlock int) int {
	if lastScannedBlock == headBlock {
		return 0
	}
	if lastScannedBlock == 0 {
		return headBlock
	}
	next := lastScannedBlock + 1
	return next
}

func (s *ScannerService) ScanBlock(ctx context.Context, blockNumber int) (map[string][]models.Transaction, error) {
	block, err := s.Client.BlockByNumber(blockNumber) // Step1. Get All the transactions of block number
	if err != nil {
		fmt.Println("[Scanner] Error querying block: ", err)
		return nil, err
	}
	fmt.Println("[Scanner] Block Details", block.Number)
	fmt.Println("[Scanner] Block HAsh", block.Hash)
	newTxs := s.Pull(ctx, parseTxs(block.Transactions))
	if len(newTxs) == 0 {
		return nil, err
	}
	return newTxs, nil
}

// parseTxs converts a list of ethclient.Transaction into a list of
// service.Transaction.
func parseTxs(txs []ethclient.Transaction) []models.Transaction {
	transactions := make([]models.Transaction, len(txs))
	for i, v := range txs {
		transactions[i] = ParseTx(v)
	}
	return transactions
}

func (s *ScannerService) Pull(ctx context.Context, txs []models.Transaction) map[string][]models.Transaction {
	result := make(map[string][]models.Transaction)
	for _, tx := range txs {
		if ok, _ := s.Db.CheckTxns(ctx, tx.From); ok {
			result[tx.From] = append(result[tx.From], tx)
		}
		if ok, _ := s.Db.CheckTxns(ctx, tx.To); ok {
			result[tx.To] = append(result[tx.To], tx)
		}
	}
	return result
}

func ParseTx(tx ethclient.Transaction) models.Transaction {
	return models.Transaction{
		ChainID:     decodeHexString(tx.ChainID),
		BlockNumber: decodeHexString(tx.BlockNumber),
		Hash:        tx.Hash,
		Nonce:       decodeHexString(tx.Nonce),
		From:        tx.From,
		To:          tx.To,
		Value:       decodeHexString(tx.Value),
		Gas:         decodeHexString(tx.Gas),
		GasPrice:    decodeHexString(tx.GasPrice),
		Input:       tx.Input,
	}
}

// GetCurrentBlock returns the last scanned block.
func (s *ScannerService) GetCurrentBlock() int {
	return s.lastScannedBlock
}

func decodeHexString(hexStr string) *big.Int {
	hexStr = strings.TrimPrefix(hexStr, "0x")

	n := new(big.Int)
	n.SetString(hexStr, 16)
	return n
}
