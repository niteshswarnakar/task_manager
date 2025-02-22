package main

import (
	"github.com/niteshswarnakar/task_manager/internal/infrastructure/database"
)

func main() {
	// cmd.Execute()

	db, err := database.InitDB("task.db")
	if err != nil {
		panic(err)
	}
	database.AutoMigrate(db)
}
