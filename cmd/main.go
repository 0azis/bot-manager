package main

import (
	"botmanager/internal/adapter/primary/http"
	"botmanager/internal/adapter/secondary/database"
	"botmanager/internal/adapter/secondary/redis"
	"botmanager/internal/core/goroutine"
	"botmanager/internal/setup"

	"log/slog"

	"github.com/joho/godotenv"
)

func main() {
	// load ENV data (database config, http server config and etc.)
	if err := godotenv.Load("../.env"); err != nil {
		slog.Error("environment not found")
	}

	// get application config
	config := setup.New()

	// import database
	store, err := database.NewDB(
		config.Store.User,
		config.Store.Password,
		config.Store.Host,
		config.Store.Port,
		config.Store.DbName,
	)
	if err != nil {
		slog.Error("database running failed")
	}

	redisDb, err := redis.NewRedis(config.Redis.BuildIP())
	if err != nil {
		slog.Error("redis running failed")
	}

	// init goroutines pool
	pool := goroutine.NewPool()

	err = initBots(config.HomeBot.Token, redisDb, *store, pool)
	if err != nil {
		slog.Error("init bots failed")
	}

	// set up http server
	app := setup.NewHTTPServer()

	// init routes
	http.InitRoutes(app, *store, pool, config.HomeBot.Token)

	err = app.Listen(config.Http.BuildIP())
	if err != nil {
		slog.Error("http server running failed")
	}
}

// init all bots from DB as application runs
func initBots(homeBotToken string, redisDB redis.RedisInterface, store database.Store, pool *goroutine.GoroutinesPool) error {
	bots, err := store.Shop.Select()
	if err != nil {
		return err
	}

	homeBot, err := goroutine.New(homeBotToken, pool, store, redisDB)
	if err != nil {
		return err
	}
	homeBot.InitHomeHandlers()
	homeBot.Start()

	for bot := range bots {
		shopBot, err := goroutine.New(bots[bot].Token, pool, store, redisDB)
		if err != nil {
			continue
		}
		shopBot.InitShopHandlers()
		shopBot.Start()
	}
	return nil
}
