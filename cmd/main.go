package main

import (
	"botmanager/internal/setup"
	"github.com/joho/godotenv"
	"log/slog"
)

func main() {
	// load ENV data (database config, http server config and etc.)
	if err := godotenv.Load("../.env"); err != nil {
		slog.Error("environment not found")
	}

	// import database 
	dbConfig := setup.NewDBConfig()
	_, err := setup.NewDB(
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

	err = app.Listen(httpConfig.BuildIP())
	if err != nil {
		slog.Error("http server running failed")
	}
}
