package config

import (
	"sync"
)

var (
	config *Config
	once   sync.Once
)

type DatabaseConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

type Config struct {
	database DatabaseConfig
}

func (c Config) GetDatabaseConf() DatabaseConfig {
	return c.database
}

func Load() {
	once.Do(func() {
		config = &Config{
			database: loadDatabaseConfig(),
		}
	})
}

func GetConfig() Config {
	Load()
	return *config
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		User: getEnv("DB_USER", "postgres"),
		Pass: getEnv("DB_PASS", "postgres"),
		Host: getEnv("DB_HOST", "localhost"),
		Port: getEnv("DB_PORT", "5432"),
		Name: getEnv("DB_NAME", "document_service"),
	}
}
