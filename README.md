# Ethereum Storage Server

This is an example server that proactively retrieves latest N ethereum blocks together with corresponding event logs and exposes API for:
* block contents by index or hash
* index range of currently stored blocks
* event logs optionally filtered by address or/and block range
* subscription to new block indices

# Build

```sh
go install gopkg.in/reform.v1/reform
go mod download
go generate ./...
go build
```

# Launch

Here are the CLI arguments supported by the binary:

```sh
./storeth -h
Usage: storeth --eth-ws-endpoint=STRING [flags]

Ethereum Storage Server

Flags:
  -h, --help                               Show context-sensitive help.
      --db-conn-string="postgres://127.0.0.1/storeth?sslmode=disable"
                                           PostgreSQL connection string ($DB_CONN_STR)
      --eth-ws-endpoint=STRING             Ethereum JSON-RPC API websocket endpoint ($ETH_WS_ENDPOINT)
      --max-num-blocks=MAX-NUM-BLOCKS      Maximum number of stored blocks ($MAX_NUM_BLOCKS)
      --server-address="127.0.0.1:9321"    HTTP/WS server address ($SERVER_ADDRESS)
```

# Usage

The server supports two types of transport: HTTP and WS. Note, the former doesn't support subscriptions.

* Get number of currently stored blocks (not including `toBlock`):
    ```sh
    curl -X POST http://127.0.0.1:9321 -H 'Content-Type: application/json'  -d '{"jsonrpc":"2.0","method":"storeth_getBlockRange","params":[], "id": 1}'
    {"jsonrpc":"2.0","id":1,"result":{"fromIndex":19989659,"toBlock":19989660}}
    ```

* Get block by index (returned in the same format as `eth_getBlockByNumber(N, true)` does):
    ```sh
    curl -X POST http://127.0.0.1:9321 -H 'Content-Type: application/json'  -d '{"jsonrpc":"2.0","method":"storeth_getBlock","params":[{"index":19989681}], "id": 1}'
    {"jsonrpc":"2.0","id":1,"result":{"block":{"baseFeePerGas":"0x24de1d55d","blobGasUsed":"0x0","difficulty":"0x0","excessBlobGas":"0x20000", ... "amount":"0x11b58f2"}],"withdrawalsRoot":"0x8473e8e707b9b0b39241c4655f173406002e8257efba9b2b4753b53c806ab500"}}}
    ```

* Get block by hash (returned in the same format as `eth_getBlockByNumber(N, true)` does):
    ```sh
    curl -X POST http://127.0.0.1:9321 -H 'Content-Type: application/json'  -d '{"jsonrpc":"2.0","method":"storeth_getBlock","params":[{"hash":"0x66e50ad668ebc4a3bd745feeda0052fff009625009b6c5660b525df71e00a732"}], "id": 1}'
    {"jsonrpc":"2.0","id":1,"result":{"block":{"baseFeePerGas":"0x24de1d55d","blobGasUsed":"0x0","difficulty":"0x0","excessBlobGas":"0x20000", ... "amount":"0x11b58f2"}],"withdrawalsRoot":"0x8473e8e707b9b0b39241c4655f173406002e8257efba9b2b4753b53c806ab500"}}}
    ```

* Find event logs with optionally specified address or/and block range (returned in the same format as `eth_getLogs({...})` does):
    ```sh
    curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"storeth_findLogs","params":[{"fromBlock":19990045,"toBlock":19990052,"address":"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606EB48"}],"id":1}' http://127.0.0.1:9321
    {"jsonrpc":"2.0","id":1,"result":{"logs":[{"address":"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", ... "blockHash":"0x0b68f8d3e8f9bf5fa7e512b2c3cf8e77f84b9395bd38a1eec5ff60c2b42b0801","logIndex":"0x7a","removed":false}]}}
    ```
    **Note**: The `toBlock` filtering condition doesn't include events from the block above (unlike `eth_getLogs`).

* Subscribe to new block indices:
    ```sh
    websocat ws://127.0.0.1:9321/ws
    > { "id": 1, "jsonrpc": "2.0", "method": "storeth_subscribe", "params": ["newBlocks"] }
    < {"jsonrpc":"2.0","id":1,"result":"0xcf6e5c1e09de8b80f8cf9d612202bab5"}
    < {"jsonrpc":"2.0","method":"storeth_subscription","params":{"subscription":"0xcf6e5c1e09de8b80f8cf9d612202bab5","result":19990203}}
    < {"jsonrpc":"2.0","method":"storeth_subscription","params":{"subscription":"0xcf6e5c1e09de8b80f8cf9d612202bab5","result":19990204}}
    < {"jsonrpc":"2.0","method":"storeth_subscription","params":{"subscription":"0xcf6e5c1e09de8b80f8cf9d612202bab5","result":19990205}}
    ...
    ```
## Notes

* I looks like I didn't quite understand the requirements so decided to improvise.
* I will use PostgreSQL communicated by nice pseudo-ORM called "reform". Typically I avoid ORMs, but since it's not a real ORM I will stick to it (used it before quite frequently).
* Found a nice Go module called "kong" that is somehow similar to "clap" Rust crate. Previously I used the standard "flag" but after using "clap" I find the structured approach better. 
* For sake of simplicity I will ignore the following aspects that are typically addressed in production code:
    - client authentication
    - graceful server exiting
    - multiple retries for failed requests
    - logging levels (and maybe structured logging)

## Some Answers

* How would you handle security of the API?

    I would support access tokens passed in the `Authorization` header. Might be plain tokens or JWT.

* How would you improve the performance of your approach?

    I believe the chosen design gives near-optimal performance taking into account that we reuse database rows. There's always a room for improvement, just now hard to find practical bottlenecks.

* How would you adapt your design to store the same data for the entire history of Ethereum Mainnet?

    I will probably need to create tables to store additional information, including transaction receipts, uncles, contract ABIs, account states, etc. Frankly, I'm not an Ethereum expert to answer that question in detail.

* What would it take to deploy and monitor a service like this in production?
  
    The service can be deployed as a container using Kubernetes, so the standard Kubernetes toolchain will be in charge (Kubectl, Argo CD, Prometheus, Grafana, etc.). Frankly, I'm not too familiar with the DevOps field.
