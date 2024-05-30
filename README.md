# Ethereum Storage Server

## Notes

* I looks like I didn't quite understand the requirements so decided to improvise.
* The service proactively retrieves and stores latest N ethereum blocks as well as the corresponding event logs. The exposed API will support JSON-RPC (to reuse `go-ethereum/rpc`) and allow to:
    - Get block contents by index or hash.
    - Get range of currently stored blocks.
    - Get event logs for a specified range of blocks.
    - Subscribe to "new-block" events when connected via websocket.
* I will use PostgreSQL communicated by nice pseudo-ORM called "reform". Typically I avoid ORMs, but since it's not a real ORM I will stick to it (used it before quite frequently).
* Found a nice Go module called "kong" that is somehow similar to "clap" Rust crate. Previously I used the standard "flag" but after using "clap" I find the structured approach better. 
* For sake of simplicity I will ignore the following aspects that are typically addressed in production code:
    - client authentication
    - graceful server exiting
    - failed request retries
