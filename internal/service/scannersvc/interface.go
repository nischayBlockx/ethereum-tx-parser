package scannersvc

import "time"

type ScannerServiceInterface interface {
	// Run starts the block scanning process.In case of no pending
	// blocks to be scanned it will return 0.
	Run() (int, error)
	StartScan(interval time.Duration)
	// GetCurrentBlock returns the last scanned block.
	GetCurrentBlock() int
}
