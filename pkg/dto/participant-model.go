package dto

import (
	"github.com/jinzhu/gorm"
)

type ParticipantModel struct {
	gorm.Model

	UUID                  string
	AuthorID              string
	ApiKey                string
	Username              string
	Approved              bool
	ApprovedBy            string
	ShowNameInLeaderboard bool
}
