package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
      _"github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:dics@localhost:5532/simple_bank?sslmode=disable"
	serverAddress= "localhost:8080"
)

func main() {
	// This is a simple Go program that prints "Hello, World!" to the console.
	println("Server Iniciando")


	conn, err:= sql.Open(dbDriver, dbSource)
	
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store:= db.NewStore(conn)
	server:= api.NewServer(&store)

	err= server.Start(serverAddress)

	if err != nil {
		log.Fatal("Server not running:", err)
	}


}