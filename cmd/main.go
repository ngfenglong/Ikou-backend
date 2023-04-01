package main

import (
	"log"
	"os"

	"github.com/ngfenglong/ikou-backend/api"
	"github.com/ngfenglong/ikou-backend/api/config"
	"github.com/ngfenglong/ikou-backend/api/store"

	_ "github.com/go-sql-driver/mysql"
)

const version = "1.0.0"

func main() {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	err := config.LoadConfig(".")
	if err != nil {
		errorLog.Fatal(err)
	}

	dataStore, err := store.NewStore(*config.C)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &api.Application{
		Config:   *config.C,
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: errorLog,
		Version:  version,
		Store:    dataStore,
	}

	server := api.NewServer(app)
	if err := server.Serve(); err != nil {
		log.Fatalf("Error running the server: %v\n", err)
	}
}
