package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/minhnghia2k3/simple_bank/api"
	db "github.com/minhnghia2k3/simple_bank/db/sqlc"
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(serverAddress); err != nil {
		log.Fatal("error when starting server:", err)
	}
}
