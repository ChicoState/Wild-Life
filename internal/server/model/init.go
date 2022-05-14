package model

import (
	"fmt"
	"time"

	"log"
	"os"
	wlog "wildlife/internal/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// database
	DB *gorm.DB
)

// initializes the database for use later
func InitDB() error {

	// Database connection string
	// ie "user:password@tcp(localhost:3306)/dbname?charset=utf8&parseTime=True&loc=Local"
	dbstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	// Grabbed from GORM's website https://gorm.io/docs/logger.html
	// Our logger will do what this can do
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	// GORM creates a database connection pool
	db, err := gorm.Open(mysql.Open(dbstring), &gorm.Config{Logger: newLogger})
	if err != nil {
		return err
	}
	// Sets the database
	DB = db
	wlog.Logf("DB Connected")
	return nil
}

// Loads users from DB into memory as cache
func LoadUsers() (*map[string]*User, error) {
	if DB == nil {
		return nil, fmt.Errorf("DB not initialized")
	}
	// Used to store users from DB
	var users []User
	// Used to return users to cache
	var userMap = make(map[string]*User)
	// Checks if we should migrate
	if os.Getenv("MIGRATE") == "true" {
		// Migrates the database
		err := DB.AutoMigrate(&User{})
		if err != nil {
			return nil, err
		}
	}
	// Loads all users from DB
	err := DB.Find(&users).Error
	if err != nil {
		wlog.Errf("Error loading users: %s", err)
		return nil, err
	}
	// Iterates through users to load into cache
	for _, user := range users {
		userMap[user.ID] = &user
	}
	return &userMap, nil
}
