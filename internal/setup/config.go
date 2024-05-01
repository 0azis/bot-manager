package setup

import (
	"os"
)

// Db Config
type db struct {
	User     string // db_user
	Password string // db_user password
	DbName   string // name of the main database
	Host     string // host of the database (localhost by default)
	Port     string // port of the database (5433 by default)
}

func dbConfig() *db {
	return &db{
		User:     getEnv("DATABASE_USER", ""),
		Password: getEnv("DATABASE_PASSWORD", ""),
		DbName:   getEnv("DATABASE_DB", ""),
		Host:     getEnv("DATABASE_HOST", "localhost"),
		Port:     getEnv("DATABASE_PORT", "5433"),
	}
}

// HTTP Config
type http struct {
	Host string // host of the server (localhost by default)
	Port string // port of the server (8000 by default)
}

// constructor for full socket IP address
func (hc http) BuildIP() string {
	return hc.Host + ":" + hc.Port
}

// compile httpConfig* struct
func httpConfig() *http {
	return &http{
		Host: getEnv("HTTP_HOST", "0.0.0.0"),
		Port: getEnv("HTTP_PORT", "8000"),
	}
}

type redisDb struct {
	Host string
	Port string
}

// compile redis* struct
func redisConfig() *redisDb {
	return &redisDb{
		Host: getEnv("REDIS_HOST", "127.0.0.1"),
		Port: getEnv("REDIS_PORT", "6739"),
	}
}

// constructor for full socket IP address
func (rd redisDb) BuildIP() string {
	return rd.Host + ":" + rd.Port
}

type homeBotToken struct {
	Token string
}

func homeBotConfig() *homeBotToken {
	return &homeBotToken{
		Token: getEnv("TOKEN", ""),
	}
}

// getter for environment variables
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
