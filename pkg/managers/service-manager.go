package managers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/r4stl1n/algo-benchmark-discord-bot/pkg/dto"
)

type ServiceManager struct {
	Config         *dto.ConfigStruct
	DiscordClient  *discordgo.Session
	DatabaseClient *DatabaseManager
}

func CreateServiceManager(config *dto.ConfigStruct, databaseClient *DatabaseManager) *ServiceManager {

	return &ServiceManager{
		Config:         config,
		DatabaseClient: databaseClient,
	}

}

func (serviceManager *ServiceManager) Initalize() error {

	discordClient, discordClientError := discordgo.New("Bot " + serviceManager.Config.BotToken)

	if discordClientError != nil {
		return discordClientError
	}

	serviceManager.DiscordClient = discordClient

	serviceManager.DiscordClient.AddHandler(serviceManager.messageHandler)

	serviceManager.DiscordClient.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	discordWebsocketError := serviceManager.DiscordClient.Open()

	if discordWebsocketError != nil {
		return discordWebsocketError
	}

	return nil
}

func (serviceManager *ServiceManager) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

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
func (serviceManager *ServiceManager) Shutdown() {
	serviceManager.DiscordClient.Close()
}
