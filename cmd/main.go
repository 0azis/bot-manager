package main

import (
	// "botmanager/internal/repos"
	"botmanager/internal/models"
	"botmanager/internal/routes"
	"botmanager/internal/setup"
	"botmanager/internal/tools"

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
	store, err := setup.NewDB(
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
	pool := models.NewGoroutinesPool()

	// init routes and bots
	routes.InitRoutes(app, store, pool)
	err = tools.InitBots(store, pool)
	if err != nil {
		slog.Error("init bots failed")
		return 
	}

	err = app.Listen(httpConfig.BuildIP())
	if err != nil {
		slog.Error("http server running failed")
		return
	}
}
