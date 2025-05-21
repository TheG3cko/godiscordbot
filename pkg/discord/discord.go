package discord

import (
	"fmt"
	"github.com/TheG3cko/godiscordbot/pkg/ollama"
	"github.com/bwmarrin/discordgo"
	"github.com/ollama/ollama/api"
	"os"
	"os/signal"
	"syscall"
)

func Client() {

	// Create a new Client session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("error creating Client session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(getmsgs)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Open a websocket connection to Client and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Client session.
	err = dg.Close()
	if err != nil {
		return
	}
}

var UserHistories []api.Message

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example, but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	merged := m.Author.Username + " " + m.Content
	if m.ChannelID == os.Getenv("CHANNEL_ID") {
		answer := ollama.AskOllama(merged, &UserHistories)
		_, err := s.ChannelMessageSend(m.ChannelID, answer)
		if err != nil {
			return
		}

	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

func getmsgs(_ *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID == os.Getenv("CHANNEL_ID") {
		fmt.Println(m.Author.Username, m.Content)
	}

}
