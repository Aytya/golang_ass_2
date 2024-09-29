package model

type User struct {
	ID        string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string  `gorm:"not null;unique" json:"name"`
	Age       int     `gorm:"not null'" json:"age"`
	ProfileId string  `json:"profileId"`
	Profile   Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
