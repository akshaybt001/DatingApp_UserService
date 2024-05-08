package initializer

import (
	"github.com/akshaybt001/DatingApp_UserService/concurrency"
	"github.com/akshaybt001/DatingApp_UserService/internal/adapters"
	"github.com/akshaybt001/DatingApp_UserService/internal/service"
	"github.com/akshaybt001/DatingApp_UserService/internal/usecases"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *service.UserService {
	repo := adapters.NewUserAdapter(db)
	usecase := usecases.NewUserUseCase(repo)
	service := service.NewUserService(repo, usecase)
	concurrency := concurrency.NewCronJob(service, db)
	concurrency.Start()

	return service
}
