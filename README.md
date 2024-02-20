
# Ethereum Transaction Parser**

## Overview
The Ethereum Transaction Parser is a blockchain parser designed to query transactions for subscribed addresses. It consists of three main packages:

- **internal/repository/db:** Defines the application state, including data structures and methods for manipulating the application state.

- **internal/service/scannersvc:** Implements core functionalities, interacting with the blockchain for extracting and parsing on-chain data.

- **internal/service/parsersvc:** Contains the Parser interface and the Transaction domain type, representing transaction objects in the application context.

## Running

**build**
go build -o ethparser cmd/main/main.go

**run**
./ethparser -block=[block_number] 

//block_number need to be replaced.

**Future Improvements**
**Error Handling:** Implement robust error handling for production environments.
**Logging System:** Implement a comprehensive logging system to trace operations and identify potential issues.
**Application Metrics:** Incorporate monitoring tools to gather metrics, providing insights into the application's behavior.
