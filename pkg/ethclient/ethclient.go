package ethclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type EthClient struct {
	endpoint string
}

func NewEthClient(URL string) *EthClient {
	return &EthClient{
		endpoint: URL,
	}
}

// BlockNumber returns the current block number. It will call
// the eth_blockNumber method of the JSON-RPC API in the given url.

func (ec *EthClient) BlockNumber() (int, error) {

	requestBody, err := json.Marshal(createRequest("eth_blockNumber", []string{}))
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(ec.endpoint, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var responseBody struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return 0, fmt.Errorf("[eth-client] Error decoding response body: %v", err)
	}

	blocknumber, err := strconv.ParseInt(responseBody.Result[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("[eth-client] Error parsing response body: %v", err)
	}

	return int(blocknumber), nil
}

// BlockByNumber retrieves information about a specific block by its number.
func (ec *EthClient) BlockByNumber(blockNumber int) (*Block, error) {
	method := "eth_getBlockByNumber"
	params := []interface{}{fmt.Sprintf("0x%x", blockNumber), true}
	requestBody, err := json.Marshal(createRequest(method, params))
	if err != nil {
		return nil, fmt.Errorf("[eth-client] Error in JSON Marshal: %v", err)
	}

	resp, err := http.Post(ec.endpoint, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("[eth-client] Error in creating Request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[eth-client] Error reading Response Body: %v", err)
	}

	var jsonResponse JSONRPCResponse
	if err := json.Unmarshal(responseBody, &jsonResponse); err != nil {
		return nil, fmt.Errorf("[eth-client] Error decoding Response Body: %v", err)
	}
	fmt.Println("Response--->", &jsonResponse)
	return &jsonResponse.Result, nil
}

// createRequest generates a JSON-RPC request.
func createRequest(method string, params interface{}) RequestBody {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return RequestBody{
		Jsonrpc: "2.0",
		ID:      int(r.Uint64()),
		Method:  method,
		Params:  params,
	}
}
