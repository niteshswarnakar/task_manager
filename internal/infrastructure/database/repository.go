package database

import (
	"errors"

	"gorm.io/gorm"
)

func First[T any](db *gorm.DB, id string) (T, error) {
	var t T
	err := db.First(&t, id).Error
	return t, err
}

func FindAll[T any](db *gorm.DB) ([]T, error) {
	var results []T
	err := db.Find(&results).Error
	return results, err
}

func UpdateTask[T any](db *gorm.DB, params map[string]interface{}) error {
	var t T
	return db.Model(&t).Updates(params).Error
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
