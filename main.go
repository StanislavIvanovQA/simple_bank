package main

import (
	"database/sql"
	"github.com/StanislavIvanovQA/simple_bank/api"
	"github.com/StanislavIvanovQA/simple_bank/config"
	db "github.com/StanislavIvanovQA/simple_bank/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal("Could not connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, cfg)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	if err = server.Start(cfg.ServerAddress); err != nil {
		log.Fatal("cannot start server", err)
	}
}
