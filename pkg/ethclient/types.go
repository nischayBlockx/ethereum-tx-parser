package ethclient

type AccessListEntry struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

type RequestBody struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type JSONRPCResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Block  `json:"result"`
}

type Block struct {
	Number       string        `json:"number"`
	Hash         string        `json:"hash"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
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
