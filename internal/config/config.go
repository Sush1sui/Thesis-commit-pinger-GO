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

	botToken := os.Getenv("bot_token")
	if botToken == "" {
		return nil, fmt.Errorf("bot_token is not set in the environment variables")
	}

	appID := os.Getenv("app_id")
	if appID == "" {
		return nil, fmt.Errorf("app_id is not set in the environment variables")
	}

	publicKey := os.Getenv("public_key")
	if publicKey == "" {
		return nil, fmt.Errorf("public_key is not set in the environment variables")
	}

	githubSecret := os.Getenv("GITHUB_SECRET")
	if githubSecret == "" {
		return nil, fmt.Errorf("github_secret is not set in the environment variables")
	}

	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		return nil, fmt.Errorf("server_url is not set in the environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "6924" // Default port if not set
	}

	return &Config{
		BotToken:     botToken,
		AppID:        appID,
		PublicKey:    publicKey,
		GithubSecret: githubSecret,
		ServerURL:    serverURL,
		Port:         port,
	}, nil
}