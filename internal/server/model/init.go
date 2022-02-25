package model

import (
	"fmt"
	"os"
	"wildlife/internal/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// database
	DB *gorm.DB
)

func InitDB() {

	dbstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(mysql.Open(dbstring), &gorm.Config{})

	if err != nil {
		log.Errf("Error connecting to database: %s", err)
		panic(err)
	}
	DB = db
	log.Logf("DB Connected")
}

func LoadUsers() (*map[string]*User, error) {
	var users []User
	var userMap = make(map[string]*User)
	if os.Getenv("MIGRATE") == "true" {
		DB.AutoMigrate(&User{})
	}
	err := DB.Find(&users).Error
	if err != nil {
		log.Errf("Error loading users: %s", err)
		return nil, err
	}
	for _, user := range users {
		userMap[user.ID] = &user
	}
	return &userMap, nil
}
