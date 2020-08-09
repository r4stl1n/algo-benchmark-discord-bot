package managers

import (
	"time"

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
	databaseClient.AutoMigrate(&dto.DailyBmEntryModel{})

	return &DatabaseManager{
		gormClient: databaseClient,
	}, nil
}

func (databaseManager *DatabaseManager) CreateParticipant(authorID string, username string) (dto.ParticipantModel, error) {

	if databaseManager.CheckIfParticipantExist(authorID) != false {
		return dto.ParticipantModel{}, nil
	}

	newUUID := uuid.NewV4().String()
	apiKeyUUID := uuid.NewV4().String()
	participantModel := dto.ParticipantModel{
		UUID:                  newUUID,
		AuthorID:              authorID,
		ApiKey:                apiKeyUUID,
		Username:              username,
		ShowNameInLeaderboard: false,
	}

	createError := databaseManager.gormClient.Create(&participantModel).Error

	if createError != nil {
		return dto.ParticipantModel{}, createError
	}

	return participantModel, nil
}

func (databaseManager *DatabaseManager) CheckIfParticipantExist(authorID string) bool {

	participantModel := new(dto.ParticipantModel)

	findError := databaseManager.gormClient.Find(&participantModel, "author_id = ?", authorID).Error

	if findError != nil {
		return false
	}

	return true
}

func (databaseManager *DatabaseManager) CheckIfParticipantExistByUUID(participantUUID string) bool {

	participantModel := new(dto.ParticipantModel)

	findError := databaseManager.gormClient.Find(&participantModel, "uuid = ?", participantUUID).Error

	if findError != nil {
		return false
	}

	return true
}

func (databaseManager *DatabaseManager) ApproveParticipantByUUID(participantUUID string, approvedBy string) bool {

	participantModel := new(dto.ParticipantModel)

	findError := databaseManager.gormClient.Find(&participantModel, "uuid = ?", participantUUID).Error

	if findError != nil {
		return false
	}

	participantModel.Approved = true
	participantModel.ApprovedBy = approvedBy

	databaseManager.gormClient.Save(&participantModel)

	return true
}

func (databaseManager *DatabaseManager) GetParticipant(authorID string) (*dto.ParticipantModel, error) {

	participantModel := new(dto.ParticipantModel)

	findError := databaseManager.gormClient.Find(&participantModel, "author_id = ?", authorID).Error

	if findError != nil {
		return participantModel, findError
	}

	return participantModel, nil
}

func (databaseManager *DatabaseManager) GetParticipantByUUID(uuid string) (*dto.ParticipantModel, error) {

	participantModel := new(dto.ParticipantModel)

	findError := databaseManager.gormClient.Find(&participantModel, "uuid = ?", uuid).Error

	if findError != nil {
		return participantModel, findError
	}

	return participantModel, nil
}

func (databaseManager *DatabaseManager) ShowNameInLeaderboardParticipantByUUID(participantUUID string, showName bool) bool {

	participantModel := new(dto.ParticipantModel)

	findError := databaseManager.gormClient.Find(&participantModel, "uuid = ?", participantUUID).Error

	if findError != nil {
		return false
	}

	participantModel.ShowNameInLeaderboard = showName

	databaseManager.gormClient.Save(&participantModel)

	return true
}

func (databaseManager *DatabaseManager) CreateRoiEntry(participantUUID string, roiValue float64) (string, error) {

	newUUID := uuid.NewV4().String()

	roiEntryModel := dto.RoiEntryModel{
		UUID:            newUUID,
		ParticipantUUID: participantUUID,
		ROIValue:        roiValue,
		SubmissionTime:  time.Now().UTC(),
	}

	createError := databaseManager.gormClient.Create(&roiEntryModel).Error

	if createError != nil {
		return "", createError
	}

	return newUUID, nil
}

func (databaseManager *DatabaseManager) GetRoiEntriesForToday() ([]dto.RoiEntryModel, error) {
	roiEntries := []dto.RoiEntryModel{}

	currentTime := time.Now().UTC()

	newDateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)

	findError := databaseManager.gormClient.Find(&roiEntries, "submission_time  > ?", newDateTime).Error

	if findError != nil {
		return roiEntries, findError
	}

	if len(roiEntries) == 0 {
		return roiEntries, nil
	}

	return roiEntries, nil

}

func (databaseManager *DatabaseManager) GetLatestEntryForParticipant(participantUUID string) (bool, dto.RoiEntryModel, error) {
	roiEntries := []dto.RoiEntryModel{}

	findError := databaseManager.gormClient.Find(&roiEntries, "participant_uuid = ?", participantUUID).Error

	if findError != nil {
		return false, dto.RoiEntryModel{}, findError
	}

	if len(roiEntries) == 0 {
		return false, dto.RoiEntryModel{}, nil
	}

	return true, roiEntries[len(roiEntries)-1], nil
}

func (databaseManager *DatabaseManager) GetDailyBmForToday() (dto.DailyBmEntryModel, error) {

	dailyBmModel := dto.DailyBmEntryModel{}

	currentTime := time.Now().UTC()

	newDateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)

	findError := databaseManager.gormClient.Find(&dailyBmModel, " date = ?", newDateTime).Error

	if findError != nil {
		return dailyBmModel, findError
	}

	return dailyBmModel, nil
}

func (databaseManager *DatabaseManager) GetAllDailyBmEntries() ([]dto.DailyBmEntryModel, error) {

	roiEntries := []dto.DailyBmEntryModel{}

	findError := databaseManager.gormClient.Find(&roiEntries).Error

	if findError != nil {
		return roiEntries, findError
	}

	return roiEntries, nil
}

func (databaseManager *DatabaseManager) CreateDailyBmEntry(startRoi float64) error {

	newUUID := uuid.NewV4().String()

	currentTime := time.Now().UTC()

	newDateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)

	dailyBmEntryModel := dto.DailyBmEntryModel{
		UUID:     newUUID,
		ROIValue: startRoi,
		Date:     newDateTime,
	}

	createError := databaseManager.gormClient.Create(&dailyBmEntryModel).Error

	if createError != nil {
		return createError
	}

	return nil
}

func (databaseManager *DatabaseManager) UpdateDailyBmEntry(uuid string, newValue float64) error {

	dailyBmModel := dto.DailyBmEntryModel{}

	findError := databaseManager.gormClient.Find(&dailyBmModel, " uuid = ?", uuid).Error

	if findError != nil {
		return findError
	}

	dailyBmModel.ROIValue = newValue

	databaseManager.gormClient.Save(&dailyBmModel)

	return nil

}
