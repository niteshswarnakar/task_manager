package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/niteshswarnakar/task_manager/internal/model"
	"gorm.io/gorm"
)

func First[T any](db *gorm.DB, id string) (T, error) {
	var t T
	err := db.Where("id=?", id).First(&t).Error
	return t, err
}

func Pop[T any](db *gorm.DB) (T, error) {
	var t T
	err := db.Order("created_at desc").First(&t).Error
	return t, err
}

func FindPendingTasks[T any](db *gorm.DB) ([]T, error) {
	var objects []T
	err := db.Where("status=?", string(model.StatusPending)).Find(&objects).Error
	return objects, err
}

func FindByColumn[T any](db *gorm.DB, column string, value any) (T, error) {
	var t T
	err := db.Where(fmt.Sprintf("%s=?", column), value).First(&t).Error
	return t, err
}

func FindAll[T any](db *gorm.DB) ([]T, error) {
	var results []T
	err := db.Find(&results).Error
	return results, err
}

func Count[T any](tx *gorm.DB) (int64, error) {
	var object T
	var count int64

	err := tx.Model(&object).Count(&count).Error
	return count, err
}

func Create[T any](db *gorm.DB, t T) error {
	return db.Create(&t).Error
}

func Delete[T any](db *gorm.DB, id string) error {
	var t T
	return db.Model(&t).Where("id=?", id).Delete(nil).Error
}

func DeleteAll[T any](db *gorm.DB) error {
	var t T
	return db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&t).Error
}

func UpdateTask[T any](db *gorm.DB, params map[string]interface{}) error {
	var t T
	return db.Model(&t).Updates(params).Error
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
