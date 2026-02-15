package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds application configuration.
// Loaded from file at configs/config.yaml only; no env overrides.
type Config struct {
	MongoURI   string `mapstructure:"mongo_uri"`
	DBName     string `mapstructure:"mongo_db"`
	ServerPort int    `mapstructure:"server_port"`
}

// Load reads config using AddConfigPath and SetConfigName. No env overrides.
func Load() Config {
	v := viper.New()

	v.SetDefault("mongo_uri", "mongodb://localhost:27017")
	v.SetDefault("mongo_db", "mon_go")
	v.SetDefault("server_port", 8080)

	v.AddConfigPath("./configs/")
	v.AddConfigPath("../configs/")
	v.AddConfigPath("../../configs/")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("config file not found in configs/ (run from repo root or ensure config exists in image)")
		}
		log.Fatalf("config read: %v", err)
	}
	log.Printf("loaded config from configs/config.yaml")

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("config unmarshal: %v", err)
	}
	if c.ServerPort == 0 {
		c.ServerPort = 8080
	}
	return c
}
