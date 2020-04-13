package structures

import "time"

type MetaData struct {
	Id        string     `gorm:"index:id"`
	Location  string     `gorm:"location"`
	Port      int32      `gorm:"port"`
	StartTime *time.Time `gorm:"startTime"`
	EndTime   *time.Time `gorm:"endTime"`
	Type      int32      `gorm:"type"`
}
