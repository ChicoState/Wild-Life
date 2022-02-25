package model

import (
	"errors"

	"github.com/google/uuid"
)

// GORM model
type User struct {
	// ID string is a UUID in struct but varchar(36) in DB
	ID string `gorm:"primaryKey;not null;type:varchar(36)"`
	// User name from Google as varchar(64)
	Name string `gorm:"type:varchar(64)"`
	// User email from Google as varchar(64)
	Email string `gorm:"type:varchar(64);unique"`
}

// Creates a new user
func NewUser(name string, email string) *User {
	var user *User
	//find user in DB by email
	DB.Where("email = ?", email).First(user)
	//if user doesn't exist, create new user
	if user == nil || user.ID == "" {
		// Creates a new UUID as a string
		id := uuid.NewString()
		return &User{id, name, email}
	}
	return nil
}

// Returns error if user couldn't be created
func (u *User) Create() error {
	// Creates a new user in the DB
	err := DB.Create(u).Error
	return err
}

// Returns error if user couldn't be saved
func (u *User) Save() error {
	// Saves the user to the DB
	err := DB.Save(u).Error
	return err
}

// Returns error if user couldn't be deleted
func (u *User) Delete() error {
	// Deletes the user from the DB
	err := DB.Delete(u).Error
	return err
}

// Returns user by email
func FindUserByEmail(email string) (User, error) {
	var user *User
	//find user in DB by email
	DB.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&user)
	if user == nil || user.ID == "" {
		return User{}, errors.New("User not found")
	}
	return *user, nil
}
