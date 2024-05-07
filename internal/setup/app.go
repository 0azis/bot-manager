package setup

type Config struct {
	Store   *db
	Http    *http
	Redis   *redisDb
	HomeBot *homeBot
}

func New() *Config {
	config := new(Config)

	store := dbConfig()
	config.Store = store

	http := httpConfig()
	config.Http = http

	redisDb := redisConfig()
	config.Redis = redisDb

	homeBot := homeBotConfig()
	config.HomeBot = homeBot

	return config
}
