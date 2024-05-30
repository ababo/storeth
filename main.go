package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"storeth/service"
)

type config struct {
	DBConnString  string  `env:"DB_CONN_STR" default:"postgres://127.0.0.1/storeth?sslmode=disable" help:"PostgreSQL connection string"`
	EthWSEndpoint string  `env:"ETH_WS_ENDPOINT" required:"" help:"Ethereum JSON-RPC API websocket endpoint"`
	MaxNumBlocks  *uint64 `env:"MAX_NUM_BLOCKS" help:"Maximum number of stored blocks"`
	ServerAddress string  `env:"SERVER_ADDRESS" default:"127.0.0.1:9321" help:"HTTP/WS server address"`
}

func main() {
	conf := new(config)
	_ = kong.Parse(conf,
		kong.Name("storeth"),
		kong.Description("Ethereum Storage Server"))

	sqlDB, err := sql.Open("postgres", conf.DBConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()
	log.Printf("connected to database")

	db := reform.NewDB(sqlDB, postgresql.Dialect, nil)

	svc := service.NewService(db)

	go monitorEth(conf, db, svc)

	server := rpc.NewServer()
	server.RegisterName("storeth", svc)

	http.HandleFunc("/", server.ServeHTTP)
	http.Handle("/ws", server.WebsocketHandler(nil))

	log.Printf("starting http server on %s", conf.ServerAddress)
	if err := http.ListenAndServe(conf.ServerAddress, nil); err != nil {
		log.Fatal(err)
	}
}
