package model

type GPS struct {
	// ID string is a UUID in struct but varchar(36) in DB
	ID        string  `gorm:"primaryKey;not null;type:varchar(36)"`
	Longitude float64 `gorm:"type:float"`
	Latitude  float64 `gorm:"type:float"`
}
