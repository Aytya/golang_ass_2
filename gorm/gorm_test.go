package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang2/model"
	"golang2/repo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"testing"
)

var (
	db          *gorm.DB
	userRepo    *repo.UserPostgres
	profileRepo *repo.ProfilePostgres
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	var err error
	dbURL := "user=testuser password=testpassword dbname=myproject sslmode=disable"
	db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	if err := db.AutoMigrate(&model.User{}, &model.Profile{}); err != nil {
		panic(fmt.Sprintf("Failed to migrate models: %v", err))
	}

	userRepo = repo.NewUserPostgres(db)
	profileRepo = repo.NewProfilePostgres(db)
}

func teardown() {
	sqlDB, _ := db.DB()
	sqlDB.Close()
}

func TestCreateUser(t *testing.T) {
	user := model.User{Name: "Azamat Serik", Age: 26}
	message, err := userRepo.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, "User successfully created", message)
}

func TestUpdateUser(t *testing.T) {
	updatedUser := model.User{Name: "Jane Smith", Age: 26}
	err := userRepo.UpdateUser("d8f8ca2c-247e-4fd0-99b3-8fdbcf03dba4", updatedUser)

	assert.NoError(t, err)

	retrievedUser, err := userRepo.GetUserById("d8f8ca2c-247e-4fd0-99b3-8fdbcf03dba4")
	assert.NoError(t, err)
	assert.Equal(t, "Jane Smith", retrievedUser.Name)
	assert.Equal(t, 26, retrievedUser.Age)
}

func TestDeleteUser(t *testing.T) {
	err := userRepo.DeleteUser("5db6678c-9fca-4d74-8bd6-318956053a2c")

	assert.NoError(t, err)

	retrievedUser, err := userRepo.GetUserById("5db6678c-9fca-4d74-8bd6-318956053a2c")
	assert.Error(t, err)
	assert.Empty(t, retrievedUser)
}

func TestGetAllUsers(t *testing.T) {
	retrievedUsers, err := userRepo.GetAllUsers(30, 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(retrievedUsers))
}

func TestCreateUserAndProfile(t *testing.T) {
	profile := model.Profile{Bio: "This is Aytya", ProfilePictureURL: "http://example.com/bob.jpg"}
	user := model.User{Name: "Aytya", Age: 32, ProfileId: profile.ID}
	err := profileRepo.CreateUserAndProfile(user, profile)

	assert.NoError(t, err)

	retrievedProfile, err := profileRepo.GetUsersWithProfile()
	assert.NoError(t, err)
	assert.Equal(t, 5, len(retrievedProfile))
}

func TestUpdateUserAndProfile(t *testing.T) {
	err := profileRepo.UpdateUserProfile("5caa31d5-f61a-4864-b049-bfd8573d19e2", "Updated bio", "http://example.com/bob_updated.jpg")
	assert.NoError(t, err)

	var updatedProfile model.Profile
	db.Where("user_id = ?", "5caa31d5-f61a-4864-b049-bfd8573d19e2").First(&updatedProfile)
	assert.Equal(t, "Updated bio", updatedProfile.Bio)
	assert.Equal(t, "http://example.com/bob_updated.jpg", updatedProfile.ProfilePictureURL)
}

func TestDeleteUserAndProfile(t *testing.T) {
	err := profileRepo.DeleteUserWithProfile("a60cdfcf-e269-43a3-84ce-a57bb663937d")
	assert.NoError(t, err)
}
