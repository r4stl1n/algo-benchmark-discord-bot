package dto

import (
	"time"

	"github.com/jinzhu/gorm"
)

type RoiEntryModel struct {
	gorm.Model
	UUID            string
	ParticipantUUID string
	ROIValue        float64
	SubmissionTime  time.Time
}
