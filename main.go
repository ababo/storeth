package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/lib/pq"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"storeth/service"
)

type config struct {
	DBConnString  string `env:"DB_CONN_STR" default:"postgres://127.0.0.1/storeth" help:"PostgreSQL connection string"`
	ServerAddress string `env:"SERVER_ADDRESS" default:"127.0.0.1:9321" help:"HTTP/WS server address"`
}

func main() {
	var conf config
	_ = kong.Parse(&conf)

	sqlDB, err := sql.Open("postgres", conf.DBConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	db := reform.NewDB(sqlDB, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))

	svc := service.NewService(db)

	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Second)
			service.NotifyNewBlock(svc, uint64(i))
		}
	}()

	server := rpc.NewServer()
	server.RegisterName("storeth", svc)

	http.HandleFunc("/", server.ServeHTTP)
	http.Handle("/ws", server.WebsocketHandler(nil))

	log.Println("Starting HTTP server...")
	if err := http.ListenAndServe(conf.ServerAddress, nil); err != nil {
		log.Fatal(err)
	}
}
