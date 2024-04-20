package setup

import "os"

// Db Config
type dbConfig struct {
	User     string // db_user
	Password string // db_user password
	DbName   string // name of the main database
	Host     string // host of the database (localhost by default)
	Port     string // port of the database (5433 by default)
}

// HTTP Config
type httpConfig struct {
	host string // host of the server (localhost by default)
	port string // port of the server (8000 by default)
}

// constructor for full socket IP address
func (hc httpConfig) BuildIP() string {
	return hc.host + ":" + hc.port
}

// compile dbConfig* struct
func NewDBConfig() *dbConfig {
	return &dbConfig{
		User:     getEnv("DATABASE_USER", ""),
		Password: getEnv("DATABASE_PASSWORD", ""),
		DbName:   getEnv("DATABASE_DB", ""),
		Host:     getEnv("DATABASE_HOST", "localhost"),
		Port:     getEnv("DATABASE_PORT", "5433"),
	}
}

// compiler httpConfig* struct
func NewHTTPConfig() *httpConfig {
	return &httpConfig{
		host: getEnv("HTTP_HOST", "0.0.0.0"),
		port: getEnv("HTTP_PORT", "8000"),
	}
}

// getter for environment variables
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
