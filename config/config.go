package config

import (
	"sync"
	"time"
)

var (
	config *Config
	once   sync.Once
)

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	User string
	Pass string
	Host string
	Port int
	Name string
}

type Config struct {
	server   ServerConfig
	database DatabaseConfig
}

func (c Config) GetServerConf() ServerConfig {
	return c.server
}

func (c Config) GetDatabaseConf() DatabaseConfig {
	return c.database
}

func Load() {
	once.Do(func() {
		config = &Config{
			server:   loadServerConfig(),
			database: loadDatabaseConfig(),
		}
	})
}

func GetConfig() Config {
	Load()
	return *config
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:         getEnv("SERVER_PORT", "9003"),
		Host:         getEnv("SERVER_HOST", "localhost"),
		ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 2*time.Second),
		WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 2*time.Second),
	}
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		User: getEnv("DB_USER", "postgres"),
		Pass: getEnv("DB_PASS", "postgres"),
		Host: getEnv("DB_HOST", "localhost"),
		Port: getEnvInt("DB_PORT", 5432),
		Name: getEnv("DB_NAME", "document_service"),
	}
}
