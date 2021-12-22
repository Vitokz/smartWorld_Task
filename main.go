package main

import (
	"context"
	"github.com/Vitokz/smartWorld_Task/config"
	"github.com/Vitokz/smartWorld_Task/handler"
	"github.com/Vitokz/smartWorld_Task/internal/repository"
	"github.com/Vitokz/smartWorld_Task/internal/services/logger"
	"github.com/Vitokz/smartWorld_Task/server"
)

func main() {
	log := logger.NewLogger()

	cfg := config.Parse()
	log.Println(cfg)

	postrgresDB, err := repository.NewPgSQL(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	defer postrgresDB.Close()

	err = repository.RunPgMigrations(cfg)
	if err != nil {
		panic(err)
	}

	hdlr := handler.NewHandler(cfg, log, repository.New(postrgresDB))
	rest := server.NewServer(hdlr)

	rest.Router.Logger.Fatal(rest.Router.Start(":" + cfg.Port))
}
