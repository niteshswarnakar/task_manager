package consumer

import (
	"fmt"

	"github.com/niteshswarnakar/task_manager/internal/infrastructure/database"
	"github.com/niteshswarnakar/task_manager/internal/lib"
	"github.com/niteshswarnakar/task_manager/internal/model"
	"gorm.io/gorm"
)

type Consumer struct {
	consumerId int
	db         *gorm.DB
	queue      lib.Queue[model.Task]
}

func (c Consumer) PerformTask() {
	for {
		fmt.Println("Task started by consumer", c.consumerId)
		task := c.queue.Get()

		_, err := database.First[model.Task](c.db, task.ID)
		if database.IsNotFound(err) {
			panic(err)
		}

		err = database.UpdateTask[model.Task](c.db.Where("id = ?", task.ID), map[string]interface{}{"status": model.StatusCompleted})
		if err != nil {
			panic(err)
		}
	}
}

func NewConsumer(consumerId int, db *gorm.DB, queue lib.Queue[model.Task]) Consumer {
	fmt.Println("Consumer created with id", consumerId)
	return Consumer{
		consumerId: consumerId,
		db:         db,
		queue:      queue,
	}
}
