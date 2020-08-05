package dto

import "github.com/jinzhu/gorm"

type RoiEntryModel struct {
	gorm.Model

	ParticipantUUID string
	ROIValue        float64
}
