package utils

import (
	"time"

	"kelarin-backend/database"
	"kelarin-backend/models"
)

// IncrementStreak increments the user's streak by 1, but only once per day.
func IncrementStreak(userID uint) error {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	now := time.Now().UTC()

	if user.LastStreakAt != nil {
		last := user.LastStreakAt.UTC()
		if now.Year() == last.Year() && now.YearDay() == last.YearDay() {
			return nil
		}
	}

	user.Streak++
	user.LastStreakAt = &now

	return database.DB.Save(&user).Error
}

// HasStreakToday checks whether the user has already performed a streak-worthy activity today.
func HasStreakToday(user *models.User) bool {
	if user.LastStreakAt == nil {
		return false
	}

	now := time.Now().UTC()
	last := user.LastStreakAt.UTC()

	return now.Year() == last.Year() && now.YearDay() == last.YearDay()
}
