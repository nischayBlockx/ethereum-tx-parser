package models

import "math/big"

type RequestBody struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type Block struct {
	Number       string        `json:"number"`
	Hash         string        `json:"hash"`
	Transactions []Transaction `json:"transactions"`
}

type AccessListEntry struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

// these details we get from the blockchain in BlockByNumber function
type RawTransaction struct {
	ChainID          string            `json:"chainId"`
	BlockNumber      string            `json:"blockNumber"`
	BlockHash        string            `json:"-"`
	Hash             string            `json:"hash"`
	Nonce            string            `json:"nonce"`
	From             string            `json:"from"`
	To               string            `json:"to"`
	Value            string            `json:"value"`
	Gas              string            `json:"gas"`
	GasPrice         string            `json:"gasPrice"`
	Input            string            `json:"input"`
	Type             string            `json:"-"`
	R                string            `json:"-"`
	S                string            `json:"-"`
	V                string            `json:"-"`
	TransactionIndex string            `json:"-"`
	AccessList       []AccessListEntry `json:"-"`
}

type Transaction struct {
	ChainID     *big.Int `json:"chainId"`
	BlockNumber *big.Int `json:"blockNumber"`
	Hash        string   `json:"hash"`
	Nonce       *big.Int `json:"nonce"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	Value       *big.Int `json:"value"`
	Gas         *big.Int `json:"gas"`
	GasPrice    *big.Int `json:"gasPrice"`
	Input       string   `json:"input"`
}
