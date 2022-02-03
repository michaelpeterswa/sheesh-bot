package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/michaelpeterswa/sheesh-bot/internal/logging"
	"github.com/michaelpeterswa/sheesh-bot/internal/settings"
	"go.uber.org/zap"
)

type OnMessage struct {
	SheeshRegex *regexp.Regexp
}

func main() {
	logger, err := logging.InitZap()
	if err != nil {
		log.Panicf("could not acquire zap logger: %s", err.Error())
	}
	logger.Info("sheesh-bot init...")

	sheeshSettings := settings.InitSettings(logger, "config/settings.yaml")

	sheeshRegex, err := createSheeshRegex()
	if err != nil {
		logger.Fatal("error creating sheesh regex", zap.Error(err))
	}

	on := OnMessage{
		SheeshRegex: sheeshRegex,
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + sheeshSettings.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(on.messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	dg.Identify.Presence.Status = sheeshSettings.Status
	dg.Identify.Presence.Game.Name = sheeshSettings.GameName

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	logger.Info("bot is now running.  press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (on *OnMessage) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	normalizedContent := strings.ToLower(m.Content)

	if on.SheeshRegex.MatchString(normalizedContent) {
		s.ChannelMessageSend(m.ChannelID, "sheesh")
	}
}

func createSheeshRegex() (*regexp.Regexp, error) {
	sheesh := `she(e+)sh`
	return regexp.Compile(sheesh)
}
