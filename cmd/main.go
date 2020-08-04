package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/r4stl1n/algo-benchmark-discord-bot/util"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	configStruct := util.ConfigStruct{}
	fmt.Println("Loading bot token")

	configFile, configLoadError := ioutil.ReadFile("config.json")

	if configLoadError != nil {
		fmt.Println("Could not load the config file")
		fmt.Println("Creating new config file, modify and run again")
		marshaledStruct, marshaledStructError := json.MarshalIndent(configStruct, "", "")
		if marshaledStructError != nil {
			panic("No idea how the heck this happened")
		}

		writeFileError := ioutil.WriteFile("config.json", marshaledStruct, 0644)

		if writeFileError != nil {
			panic("Could not write file out")
		}
		return
	}

	configUnmarshalError := json.Unmarshal(configFile, &configStruct)

	if configUnmarshalError != nil {
		panic("Could not unmarshall config file")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + configStruct.BotToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
