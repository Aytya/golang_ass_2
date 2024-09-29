package model

type Profile struct {
	ID                string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID            string `gorm:"foreignKey;not null;unique"`
	Bio               string `gorm:"type:text"`
	ProfilePictureURL string
}
