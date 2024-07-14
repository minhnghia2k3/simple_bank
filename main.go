package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/minhnghia2k3/simple_bank/api"
	db "github.com/minhnghia2k3/simple_bank/db/sqlc"
	"github.com/minhnghia2k3/simple_bank/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(config.ServerAddress); err != nil {
		log.Fatal("error when starting server:", err)
	}
}
