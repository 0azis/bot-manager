package main

import (
	"botmanager/internal/routes"
	"botmanager/internal/setup"
	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	// load ENV data (database config, http server config and etc.)
	if err := godotenv.Load("../.env"); err != nil {
		slog.Error("environment not found")
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
	}

	// set up http server
	httpConfig := setup.NewHTTPConfig()
	app := setup.NewHTTPServer()

	// init routes and bots
	routes.InitRoutes(app, store)
	err = setup.InitBots(store)
	if err != nil {
		slog.Error("init bots failed")
	}

	err = app.Listen(httpConfig.BuildIP())
	if err != nil {
		slog.Error("http server running failed")
	}
}
