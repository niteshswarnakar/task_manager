package main

import (
	"sync"

	"github.com/niteshswarnakar/task_manager/internal/constants.go"
	"github.com/niteshswarnakar/task_manager/internal/infrastructure/consumer"
	"github.com/niteshswarnakar/task_manager/internal/infrastructure/database"
	"github.com/niteshswarnakar/task_manager/internal/lib"
	"github.com/niteshswarnakar/task_manager/internal/logger"
	"github.com/niteshswarnakar/task_manager/internal/model"
)

func main() {
	// cmd.Execute()

	appLogger := logger.NewAppLogger()

	appLogger.Info("Task Manager started")

	db, err := database.InitDB("task.db")
	if err != nil {
		appLogger.Panic(err)
	}

	err = database.AutoMigrate(db)
	if err != nil {
		appLogger.Panic(err)
	}

	queue := lib.NewQueue[model.Task]()

	var wg sync.WaitGroup

	for i := 0; i < constants.NumberOfConsumer; i++ {
		consumerId := i
		wg.Add(1)
		consumerObj := consumer.NewConsumer(consumerId, db, queue, appLogger)
		go consumerObj.PerformTask()
	}

	//producer logic will be here before wg.Wait()
	appLogger.Info("Producer started")

	wg.Wait()

}
