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

	//Парсинг конфига
	cfg := config.Parse()
	log.Println(cfg)

	postrgresDb, err := repository.NewPgSQL(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	defer postrgresDb.Close()

	hdlr := handler.NewHandler(cfg, log, repository.New(postrgresDb))
	rest := server.NewServer(hdlr)

	rest.Router.Logger.Fatal(rest.Router.Start(":" + cfg.Port))
}
