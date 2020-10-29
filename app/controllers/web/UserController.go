package web

import (
	"errors"
	"fiber-boilerplate/app/models"
	"fiber-boilerplate/database"
)

// Return a single user as JSON
func FindUserByUsername(db *database.Database, username string) (*models.User, error) {
	User := new(models.User)
	if response := db.Where("name = ?", username).First(&User); response.Error != nil {
		return nil, response.Error
	}
	if User.ID == 0 {
		return User, errors.New("user not found")
	}
	// Match role to user
	if User.RoleID != 0 {
		Role := new(models.Role)
		if res := db.Find(&Role, User.RoleID); res.Error != nil {
			return User, errors.New("error when retrieving the role of the user")
		}
		if Role.ID != 0 {
			User.Role = *Role
		}
	}
	return User, nil
}

// Return a single user as JSON
func FindUserByID(db *database.Database, id int64) (*models.User, error) {
	User := new(models.User)
	if response := db.Where("id = ?", id).First(&User); response.Error != nil {
		return nil, response.Error
	}
	if User.ID == 0 {
		return User, errors.New("user not found")
	}
	// Match role to user
	if User.RoleID != 0 {
		Role := new(models.Role)
		if res := db.Find(&Role, User.RoleID); res.Error != nil {
			return User, errors.New("error when retrieving the role of the user")
		}
		if Role.ID != 0 {
			User.Role = *Role
		}
	}
	return User, nil
}
