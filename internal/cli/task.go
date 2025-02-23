package cli

import (
	"fmt"
	"sync"

	"github.com/niteshswarnakar/task_manager/internal/constants.go"
	"github.com/niteshswarnakar/task_manager/internal/infrastructure/consumer"
	"github.com/niteshswarnakar/task_manager/internal/infrastructure/database"
	"github.com/niteshswarnakar/task_manager/internal/lib"
	"github.com/niteshswarnakar/task_manager/internal/logger"
	"github.com/niteshswarnakar/task_manager/internal/model"
	"gorm.io/gorm"
)

type AppCli struct {
	db       *gorm.DB
	logger   logger.AppLogger
	queue    lib.Queue[model.Task]
	stopChan chan struct{}
}

func (a AppCli) Add(title string) error {
	task := model.Task{
		ID:     model.NewID(),
		Title:  title,
		Status: model.StatusPending,
	}

	err := database.Create[model.Task](a.db, task)
	if err != nil {
		a.logger.Error("Cli: Add: Create task", err)
		return err
	}
	a.logger.Info("Task created : " + task.Title)
	return nil
}

func (a AppCli) List() ([]model.Task, error) {
	tasks, err := database.FindAll[model.Task](a.db.Order("created_at desc"))
	if err != nil {
		a.logger.Error("Cli: List: Find all tasks", err)
		return nil, err
	}
	return tasks, nil
}

func (a AppCli) Process() error {
	tasks, err := database.FindPendingTasks[model.Task](a.db)
	if err != nil {
		a.logger.Error("Cli: Process: Find pending tasks", err)
		return err
	}

	if len(tasks) == 0 {
		a.logger.Info("No task to process...")
		return nil
	}

	wg := sync.WaitGroup{}

	totalWorker := constants.NumberOfConsumer
	if len(tasks) < constants.NumberOfConsumer {
		totalWorker = len(tasks)
	}

	// i implemented worker pool pattern for concurrency
	for i := 0; i < totalWorker; i++ {
		wg.Add(1)
		consumerId := i + 1
		newConsumer := consumer.NewConsumer(consumerId, a.db, a.queue, a.logger, a.stopChan)
		go func() {
			defer wg.Done()
			newConsumer.PerformTask()
		}()
	}

	a.logger.Info(fmt.Sprintf("Processing %v tasks...v", len(tasks)))

	for _, task := range tasks {
		a.logger.Info("CLI-> Task added to queue : " + task.Title)
		a.queue.Put(task)
	}

	for i := 0; i < totalWorker; i++ {
		a.stopChan <- struct{}{}
	}

	wg.Wait()
	return nil
}

func NewAppCli(db *gorm.DB, logger logger.AppLogger, queue lib.Queue[model.Task], stopChan chan struct{}) AppCli {
	return AppCli{
		db:       db,
		logger:   logger,
		queue:    queue,
		stopChan: stopChan,
	}
}
