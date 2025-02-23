package consumer

import (
	"fmt"
	"time"

	"github.com/niteshswarnakar/task_manager/internal/constants.go"
	"github.com/niteshswarnakar/task_manager/internal/infrastructure/database"
	"github.com/niteshswarnakar/task_manager/internal/lib"
	"github.com/niteshswarnakar/task_manager/internal/logger"
	"github.com/niteshswarnakar/task_manager/internal/model"
	"gorm.io/gorm"
)

type Consumer struct {
	consumerId  int
	db          *gorm.DB
	queue       lib.Queue[model.Task]
	logger      logger.AppLogger
	stopChannel chan struct{}
}

func (c Consumer) PerformTask() {
	c.logger.Info(fmt.Sprintf("worker %v started, waiting for task!", c.consumerId))
	taskChannel := c.queue.GetChannel()

	for {
		select {
		case task, ok := <-taskChannel:
			if !ok {
				c.logger.Info(fmt.Sprintf("worker %v stopped !ok", c.consumerId))
				return
			}
			// To mimic actual task is happening
			time.Sleep(time.Second)

			c.logger.Info(fmt.Sprintf("worker %v performing '%s' task", c.consumerId, task.Title))
			time.Sleep(2 * time.Second)

			err := database.UpdateTask[model.Task](c.db.Where("id = ?", task.ID), map[string]interface{}{"status": model.StatusCompleted})
			if err != nil {
				c.logger.Error("Consumer: PerformTask: Update Task", err)
			}

			c.logger.Info(fmt.Sprintf("worker %v doing ...", c.consumerId))
			time.Sleep(1 * time.Second)
			c.logger.Info(fmt.Sprintf("worker %v completed task '%s'", c.consumerId, task.Title))
			c.logger.Info(fmt.Sprintf("worker %v hibernated", c.consumerId))
			time.Sleep(1 * time.Second)

		case <-c.stopChannel:
			c.logger.Info(fmt.Sprintf("worker %v stopped", c.consumerId))
			return
		}
	}
}

func StartWorkerPool(db *gorm.DB, queue lib.Queue[model.Task], logger logger.AppLogger, stopChan chan struct{}) {
	for i := 0; i < constants.NumberOfConsumer; i++ {
		consumerId := i + 1
		consumerObj := NewConsumer(consumerId, db, queue, logger, stopChan)
		go consumerObj.PerformTask()
	}
}

func NewConsumer(consumerId int, db *gorm.DB, queue lib.Queue[model.Task], logger logger.AppLogger, stopChan chan struct{}) Consumer {
	return Consumer{
		consumerId:  consumerId,
		db:          db,
		queue:       queue,
		logger:      logger,
		stopChannel: stopChan,
	}
}
