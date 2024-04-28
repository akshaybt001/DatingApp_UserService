package db

import (
	"github.com/akshaybt001/DatingApp_UserService/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(connectTo string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectTo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&entities.User{})
	db.AutoMigrate(&entities.Admin{})
	db.AutoMigrate(&entities.Gender{})
	db.AutoMigrate(&entities.Interests{})
	db.AutoMigrate(&entities.UserInterests{})
	db.AutoMigrate(&entities.Address{})
	return db, nil

}
