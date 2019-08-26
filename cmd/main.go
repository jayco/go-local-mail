package main

import (
	"context"
	"os"

	"github.com/jayco/go-local-email/internal/server"
	"github.com/jayco/go-local-email/internal/store"
)

func main() {
	var port, dbConn string
	var ok bool

	if port, ok = os.LookupEnv("PORT"); !ok {
		port = "8080"
	}

	if dbConn, ok = os.LookupEnv("DB_CONN"); !ok {
		dbConn = "mongodb://localhost:27017"
	}

	ctx := context.Background()
	db := store.NewMailDB(ctx, &dbConn)

	server.Serve(&port, db)
}
