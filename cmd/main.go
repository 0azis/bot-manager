package main

import (
	// "botmanager/internal/repos"
	"botmanager/internal/routes"
	"botmanager/internal/setup"
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

	// repo := repos.NewSubscriberRepo(store)
	// subs, _ := repo.Select()
	// fmt.Println(subs)

	// set up http server
	httpConfig := setup.NewHTTPConfig()
	app := setup.NewHTTPServer()

	// init routes and bots
	routes.InitRoutes(app, store)
	err = setup.InitBots(store)
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
