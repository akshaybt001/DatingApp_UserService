package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey;unique;not null"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Password     string
	IsBlocked    bool `json:"is_blocked" gorm:"default:false"`
	ReportCount  int
	IsSubscribed bool `json:"is_subscribed" gorm:"default:false"`
	CreatedAt    time.Time
}

type Gender struct {
	Id   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name" `
}

type UserInterests struct {
	Id         int `json:"id" gorm:"primaryKey"`
	ProfileId  uuid.UUID
	Profile    Profile   `gorm:"foreignKey:ProfileId"`
	InterestId int       `json:"interest_id"`
	Interest   Interests `json:"interest" gorm:"foreignKey:InterestId;constraint:OnDelete:CASCADE"`
}

type UserGenders struct {
	Id        int `json:"id" gorm:"primaryKey"`
	ProfileId uuid.UUID
	Profile   Profile `gorm:"foreignKey:ProfileId"`
	GenderId  int     `json:"gender_id"`
	Gender    Gender  `json:"interest" gorm:"foreignKey:GenderId;constraint:OnDelete:CASCADE"`
}

type Interests struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Interest string `json:"interest"`
}

type Admin struct {
	ID       uuid.UUID
	Name     string
	Password string
	Email    string
	Phone    string
}

type Preference struct {
	Id         uuid.UUID
	ProfileId  uuid.UUID
	Profile    Profile `gorm:"foreignKey:ProfileId"`
	MinAge     int     `json:"min_age"  binding:"required" validate:"required"`
	MaxAge     int     `json:"max_age"  binding:"required" validate:"required"`
	GenderId   int
	Gender     Gender `gorm:"foreignKey:GenderId"`
	DesireCity string `json:"desirecity"  binding:"required" validate:"required"`
}
type Address struct {
	Id        uuid.UUID `gorm:"primaryKey;unique;not null"`
	Country   string
	State     string
	District  string
	City      string
	ProfileId uuid.UUID
	Profile   Profile `gorm:"foreignKey:ProfileId"`
}

type Profile struct {
	ID     uuid.UUID `gorm:"primaryKey;unique;not null"`
	UserId uuid.UUID
	User   User `gorm:"foreignKey:UserId"`
	Image  string
}
