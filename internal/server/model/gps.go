package model
type GPS struct{
	ID string `gorm:"primaryKey; not null; type varchar(36)"`
	Lognitude float64 `gorm:"type:float"`
	Latitude float64 `gorm:"type:float"`
}

