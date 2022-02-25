package model

import "github.com/google/uuid"

type User struct {
	ID    string `gorm:"primaryKey;not null;autoIncrement;type:varchar(36)"`
	Name  string `gorm:"type:varchar(255)"`
	Email string `gorm:"type:varchar(255);unique"`
}

func AddUser(name string, email string) User {
	//check db for user with email
	id := uuid.NewString()
	return NewUser(id, name, email)
}

func NewUser(id string, name string, email string) User {
	return User{id, name, email}
}

func (u *User) Save() error {
	err := DB.Create(u).Error
	return err
}

func (u *User) Delete() error {
	err := DB.Delete(u).Error
	return err
}
