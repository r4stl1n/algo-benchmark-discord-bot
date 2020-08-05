package managers

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/r4stl1n/algo-benchmark-discord-bot/pkg/dto"
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
