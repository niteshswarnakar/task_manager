package database

import (
	"github.com/niteshswarnakar/task_manager/internal/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Task{})
}
