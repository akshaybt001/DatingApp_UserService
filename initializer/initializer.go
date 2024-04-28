package initializer

import (
	"github.com/akshaybt001/DatingApp_UserService/internal/adapters"
	"github.com/akshaybt001/DatingApp_UserService/internal/service"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB)*service.UserService{
	repo:=adapters.NewUserAdapter(db)
	service:=service.NewUserService(repo)
	return service
}