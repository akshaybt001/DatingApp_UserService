package concurrency

import (
	"fmt"
	"log"

	"github.com/akshaybt001/DatingApp_UserService/internal/service"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

type CronJob struct {
	service *service.UserService
	DB      *gorm.DB
}

func NewCronJob(service *service.UserService, db *gorm.DB) *CronJob {
	return &CronJob{
		service: service,
		DB:      db,
	}
}
func (c *CronJob) Start() {
	cron := cron.New()
	err := cron.AddFunc("00 00 * * *", func() {
		c.UpdateLikeCount()
	})
	if err != nil {
		log.Print("error scheduling cron job ", err)
	}
	cron.Start()
}

func (c *CronJob) UpdateLikeCount() error {
	err := c.DB.Exec("UPDATE users SET like_count = 3").Error
	if err != nil {
		return fmt.Errorf("failed to update like_count: %w", err)
	}

	fmt.Println("Like count for all users set to 3 successfully")
	return nil
}
