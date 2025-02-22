package main

import (
	"fmt"
	"sync"

	"github.com/niteshswarnakar/task_manager/internal/constants.go"
	"github.com/niteshswarnakar/task_manager/internal/infrastructure/consumer"
	"github.com/niteshswarnakar/task_manager/internal/infrastructure/database"
	"github.com/niteshswarnakar/task_manager/internal/lib"
	"github.com/niteshswarnakar/task_manager/internal/model"
)

func main() {
	// cmd.Execute()

	db, err := database.InitDB("task.db")
	if err != nil {
		panic(err)
	}
	database.AutoMigrate(db)

	queue := lib.NewQueue[model.Task]()

	var wg sync.WaitGroup

	for i := 0; i < constants.NumberOfConsumer; i++ {
		consumerId := i
		wg.Add(1)
		consumerObj := consumer.NewConsumer(consumerId, db, queue)
		go consumerObj.PerformTask()
	}

	//producer logic will be here before wg.Wait()
	fmt.Println("Producer started")

	wg.Wait()

}
