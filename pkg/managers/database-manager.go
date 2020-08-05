package managers

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/r4stl1n/algo-benchmark-discord-bot/pkg/dto"
	"github.com/satori/go.uuid"
)

type DatabaseManager struct {
	gormClient *gorm.DB
}

func CreateDatabaseManager(databaseName string) (*DatabaseManager, error) {

	databaseClient, databaseClientError := gorm.Open("sqlite3", databaseName)

	if databaseClientError != nil {
		return &DatabaseManager{}, databaseClientError
	}

	databaseClient.AutoMigrate(&dto.ParticipantModel{})
	databaseClient.AutoMigrate(&dto.RoiEntryModel{})

	return &DatabaseManager{
		gormClient: databaseClient,
	}, nil
}

func (databaseManager *DatabaseManager) CheckIfParticipantExist(authorID string) bool {

	participantModel := new(dto.ParticipantModel)

	findError := databaseManager.gormClient.Find(&participantModel, "author_id = ?", authorID).Error

	if findError != nil {
		return false
	}

	return true
}

func (databaseManager *DatabaseManager) CreateParticipant(authorID string) (string, error) {

	if databaseManager.CheckIfParticipantExist(authorID) != false {
		return "", nil
	}

	newUUID := uuid.NewV4().String()

	participantModel := dto.ParticipantModel{
		UUID:     newUUID,
		AuthorID: authorID,
	}

	createError := databaseManager.gormClient.Create(&participantModel).Error

	if createError != nil {
		return "", createError
	}

	return newUUID, nil
}

func (databaseManager *DatabaseManager) GetParticipant(authorID string) (*dto.ParticipantModel, error) {

	participantModel := new(dto.ParticipantModel)

	findError := databaseManager.gormClient.Find(&participantModel, "author_id = ?", authorID).Error

	if findError != nil {
		return participantModel, findError
	}

	return participantModel, nil
}

func (databaseManager *DatabaseManager) CreateRoiEntry(participantUUID string, roiValue float64) (string, error) {

	newUUID := uuid.NewV4().String()

	roiEntryModel := dto.RoiEntryModel{
		UUID:            newUUID,
		ParticipantUUID: participantUUID,
		ROIValue:        roiValue,
	}

	createError := databaseManager.gormClient.Create(&roiEntryModel).Error

	if createError != nil {
		return "", createError
	}

	return newUUID, nil
}
