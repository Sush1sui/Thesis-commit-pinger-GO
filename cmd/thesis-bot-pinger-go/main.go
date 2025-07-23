package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sush1sui/thesis-bot-pinger-go/internal/bot"
	"github.com/Sush1sui/thesis-bot-pinger-go/internal/common"
	"github.com/Sush1sui/thesis-bot-pinger-go/internal/config"
	"github.com/Sush1sui/thesis-bot-pinger-go/internal/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println("Error initializing configuration:", err)
	} else {
		config.GlobalConfig = cfg
		fmt.Println("Configuration loaded successfully")
	}

	addr := fmt.Sprintf(":%s", config.GlobalConfig.Port)
	router := server.NewRouter()
	fmt.Printf("Server is listening on Port: %s\n", config.GlobalConfig.Port)

	go func() {
		if err := http.ListenAndServe(addr, router); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	go bot.StartBot()

	go common.PingServerLoop(config.GlobalConfig.ServerURL)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Shutting down server gracefully...")
}