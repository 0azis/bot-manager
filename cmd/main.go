package main

import (
	"botmanager/internal/adapter/http/route"
	"botmanager/internal/adapter/redis"
	"botmanager/internal/adapter/repo"
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
	store, err := repo.NewDB(
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
	route.InitRoutes(app, *store, pool)

	err = app.Listen(config.Http.BuildIP())
	if err != nil {
		slog.Error("http server running failed")
	}
}

// init all bots from DB as application runs
func initBots(homeBotToken string, redisDB redis.RedisInterface, store repo.Store, pool *goroutine.GoroutinesPool) error {
	bots, err := store.Shop.Select()
	if err != nil {
		return err
	}

	homeBot, err := goroutine.NewHomeBot(homeBotToken, pool, redisDB)
	if err != nil {
		return err
	}
	homeBot.InitHomeHandlers()
	homeBot.Start()

	for bot := range bots {
		shopBot, err := goroutine.NewShopBot(bots[bot].Token, pool, store)
		if err != nil {
			continue
		}
		shopBot.InitShopHandlers()
		shopBot.Start()
	}
	return nil
}
