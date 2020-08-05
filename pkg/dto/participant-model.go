package dto

import (
	"github.com/jinzhu/gorm"
)

type ParticipantModel struct {
	gorm.Model

	UUID     string
	AuthorID string
}
