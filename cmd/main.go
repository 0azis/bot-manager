package main

import (
	"botmanager/internal/models/goroutine"
	"botmanager/internal/repos"
	"botmanager/internal/routes"
	"botmanager/internal/setup"
	"fmt"

	// "fmt"
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
	store, err := repos.NewDB(
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
	pool := goroutine.NewPool() 

	err = initBots(store, pool)
	if err != nil {
		slog.Warn("init bots failed")
	}

	// init routes and bots
	routes.InitRoutes(app, store, pool)

	err = app.Listen(httpConfig.BuildIP())
	if err != nil {
		slog.Error("http server running failed")
		return
	}
}

// init all bots from DB as application runs
func initBots(store repos.Store, pool goroutine.GoroutinesPool) error {
	bots, err := store.Shop().Select()
	if err != nil {
		return err
	}

	for bot := range bots {
		g, err := goroutine.New(bots[bot], store, pool)	
		if err != nil {
			continue
		}	
		g.Start()
		fmt.Println(pool)
	}
	return nil
}
