package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type HttpServer struct {
	Addr string `yaml:"address"`
}

type Config struct {
	Env         string `yaml:"env" required:"true"`
	StoragePath string `yaml:"storage_path"`
	HttpServer  `yaml:"http_server"`
}

func resolveConfigPath(explicitPath string) string {
	if explicitPath != "" {
		return explicitPath
	}

	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	return "config/local.yml"
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env file: %v", err)
	}

	flags := flag.String("config", "", "path to the config file")
	flag.Parse()

	configPath := resolveConfigPath(*flags)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error())
	}

	return &cfg
}
