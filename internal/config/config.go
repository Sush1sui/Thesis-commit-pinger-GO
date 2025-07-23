package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken     string
	AppID        string
	PublicKey    string
	GithubSecret string
	Port         string
	ServerURL    string
}

var (
	// GlobalConfig holds the application configuration
	GlobalConfig *Config
)

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "6924" // Default port if not set
	}
	return &Config{
		BotToken:     os.Getenv("bot_token"),
		AppID:        os.Getenv("app_id"),
		PublicKey:    os.Getenv("public_key"),
		GithubSecret: os.Getenv("github_secret"),
		Port:         port,
	}, nil
}