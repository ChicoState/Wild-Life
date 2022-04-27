package model

type Plant struct {
	ID string `gorm:"primaryKey;not null; type:varchar(36)"`
	Gps GPS `gorm: "foreignKey::GpsID"`
	GpsID string `gorm:"not null; type:varchar(36)"`
	Name string `gorm:"not null; type:varchar(255)"`
	Birth string `gorm:"not null; type:datatime"`
}