package managers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/r4stl1n/algo-benchmark-discord-bot/pkg/dto"
	"github.com/sirupsen/logrus"
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

func (serviceManager *ServiceManager) handleRegisterCommand(s *discordgo.Session, m *discordgo.MessageCreate) {

	chanCreate, chanCreateError := s.UserChannelCreate(m.Author.ID)

	if chanCreateError != nil {
		logrus.Error(chanCreateError)
		return
	}

	if serviceManager.DatabaseClient.CheckIfParticipantExist(m.Author.ID) == true {
		s.ChannelMessageSend(chanCreate.ID, "You are already registered")
		return
	}

	participantUUID, participantRegisterError := serviceManager.DatabaseClient.CreateParticipant(m.Author.ID)

	if participantRegisterError != nil {
		s.ChannelMessageSend(chanCreate.ID, "Something broke tell the owner you can't register")
	}

	s.ChannelMessageSend(chanCreate.ID, "Welcome to the algo-benchmark")
	s.ChannelMessageSend(chanCreate.ID, "Your participant ID: "+participantUUID)

}

func (serviceManager *ServiceManager) handleGiveIDCommand(s *discordgo.Session, m *discordgo.MessageCreate) {

	chanCreate, chanCreateError := s.UserChannelCreate(m.Author.ID)

	if chanCreateError != nil {
		logrus.Error(chanCreateError)
		return
	}

	if serviceManager.DatabaseClient.CheckIfParticipantExist(m.Author.ID) != true {
		s.ChannelMessageSend(chanCreate.ID, "You are not registered")
		return
	}

	participant, participantError := serviceManager.DatabaseClient.GetParticipant(m.Author.ID)

	if participantError != nil {
		logrus.Error(participantError)
		s.ChannelMessageSend(chanCreate.ID, "Something broke tell the owner you can't get your id")
	}

	s.ChannelMessageSend(chanCreate.ID, "Join Date: "+participant.CreatedAt.String())
	s.ChannelMessageSend(chanCreate.ID, "Your participant ID: "+participant.UUID)

}

func (serviceManager *ServiceManager) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	logrus.Debug("Got message from user: " + m.Author.ID + " - " + m.Content)

	if m.Content == "!register" {

		serviceManager.handleRegisterCommand(s, m)

	} else if m.Content == "!giveid" {
		serviceManager.handleGiveIDCommand(s, m)
	}

}

func (serviceManager *ServiceManager) Shutdown() {
	serviceManager.DiscordClient.Close()
}
