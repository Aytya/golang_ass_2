package repo

import (
	"fmt"
	"golang2/model"
	"gorm.io/gorm"
)

type ProfilePostgres struct {
	db *gorm.DB
}

func NewProfilePostgres(db *gorm.DB) *ProfilePostgres {
	return &ProfilePostgres{db: db}
}

func (p *ProfilePostgres) CreateUserAndProfile(user model.User, profile model.Profile) error {
	if err := p.db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	profile.UserID = user.ID
	if err := p.db.Create(&profile).Error; err != nil {
		return fmt.Errorf("failed to create profile: %v", err)
	}

	user.ProfileId = profile.ID
	if err := p.db.Model(&user).Updates(model.User{ProfileId: user.ProfileId}).Error; err != nil {
		return fmt.Errorf("failed to update user with profile ID: %v", err)
	}

	return nil
}

func (p *ProfilePostgres) GetUsersWithProfile() ([]model.Profile, error) {
	var profile []model.Profile

	err := p.db.Find(&profile).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users with profiles: %v", err)
	}

	return profile, nil
}

func (p *ProfilePostgres) UpdateUserProfile(userID string, bio string, profilePictureURL string) error {
	err := p.db.Model(&model.Profile{}).Where("user_id = ?", userID).Updates(model.Profile{
		Bio:               bio,
		ProfilePictureURL: profilePictureURL,
	}).Error
	if err != nil {
		return fmt.Errorf("failed to update user profile: %v", err)
	}

	return nil
}

func (p *ProfilePostgres) DeleteUserWithProfile(userID string) error {
	if err := p.db.Where("user_id = ?", userID).Delete(&model.Profile{}).Error; err != nil {
		return fmt.Errorf("failed to delete user profile: %v", err)
	}

	if err := p.db.Where("id = ?", userID).Delete(&model.User{}).Error; err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}
