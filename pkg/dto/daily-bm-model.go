package dto

import (
	"time"

	"github.com/jinzhu/gorm"
)

type DailyBmEntryModel struct {
	gorm.Model
	UUID     string
	ROIValue float64
	Date     time.Time
}
