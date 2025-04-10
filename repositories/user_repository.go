package repositories

import (
	"kelarin-backend/database"
	"kelarin-backend/models"
)

// CreateUser inserts a new user into the database.
func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

// GetUserByEmail returns a user by email.
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// GetUserByID returns a user by ID along with preloaded associations.
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.
		Preload("OwnedWorkspaces").
		Preload("CollabWorkspaces.User").
		First(&user, id).Error
	return &user, err
}
