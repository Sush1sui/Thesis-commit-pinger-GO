package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sush1sui/thesis-bot-pinger-go/internal/config"
	"github.com/bwmarrin/discordgo"
)

var Session *discordgo.Session

func StartBot() {

	s, err := discordgo.New("Bot " + config.GlobalConfig.BotToken)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	} else {
		Session = s
		fmt.Println("Discord session created successfully")
	}

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildPresences | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages

	s.AddHandler(func(sess *discordgo.Session, ready *discordgo.Ready) {
		sess.UpdateStatusComplex(discordgo.UpdateStatusData{
			Status: "idle",
			Activities: []*discordgo.Activity {
				{
					Name: "to repo commits",
					Type: discordgo.ActivityTypeListening,
				},
			},
		})
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("error opening Discord session: %v", err)
	}

	defer s.Close()

	fmt.Println("Bot is now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Shutting down bot gracefully...")
}