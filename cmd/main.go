package main

import (
	"botmanager/internal/adapter/http/route"
	"botmanager/internal/adapter/repo"
	"botmanager/internal/core/telegram"
	"botmanager/internal/setup"

	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	// load ENV data (database config, http server config and etc.)
	if err := godotenv.Load("../.env"); err != nil {
		slog.Error("environment not found")
		return
	}

	// import database
	dbConfig := setup.NewDBConfig()
	store, err := repo.NewDB(
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DbName,
	)
	if err != nil {
		slog.Error("database running failed")
		return
	}

	// set up http server
	httpConfig := setup.NewHTTPConfig()
	app := setup.NewHTTPServer()

	// init goroutines pool
	pool := telegram.NewPool()

	err = initBots(*store, pool)
	if err != nil {
		slog.Error("init bots failed")
	}

	// init routes
	route.InitRoutes(app, *store, pool)

	err = app.Listen(httpConfig.BuildIP())
	if err != nil {
		slog.Error("http server running failed")
		return
	}
}

// init all bots from DB as application runs
func initBots(store repo.Store, pool *telegram.GoroutinesPool) error {
	bots, err := store.Shop.Select()
	if err != nil {
		return err
	}

	for bot := range bots {
		g, err := telegram.New(bots[bot], store, pool)
		if err != nil {
			continue
		}
		g.Start()
	}
	return nil
}
