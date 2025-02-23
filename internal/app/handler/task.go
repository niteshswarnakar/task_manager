package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/niteshswarnakar/task_manager/internal/app"
	"github.com/niteshswarnakar/task_manager/internal/app/handler/view"
	"github.com/niteshswarnakar/task_manager/internal/constants.go"
	"github.com/niteshswarnakar/task_manager/internal/infrastructure/consumer"
	"github.com/niteshswarnakar/task_manager/internal/infrastructure/database"
	"github.com/niteshswarnakar/task_manager/internal/lib"
	"github.com/niteshswarnakar/task_manager/internal/logger"
	"github.com/niteshswarnakar/task_manager/internal/model"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db       *gorm.DB
	logger   logger.AppLogger
	queue    lib.Queue[model.Task]
	stopChan chan struct{}
}

func (t TaskHandler) Create(c *gin.Context) {

	request, err := view.Bind[view.Task](c)
	if err != nil {
		t.logger.Error("TaskHandler: Create: Bind request", err)
		app.InvalidRequest(c)
		return
	}

	task := model.Task{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		Title:     request.Title,
		Status:    model.StatusPending,
	}

	if err := t.db.Create(&task).Error; err != nil {
		t.logger.Error("TaskHandler: Create: Create task", err)
		app.InternalServerError(c)
		return
	}

	app.SuccessResponse(c, task)
}

func (t TaskHandler) List(c *gin.Context) {

	pagination, err := view.Bind[view.PaginationQuery](c)
	var isPagingated bool = true
	if err != nil {
		t.logger.Info("Pagination not found")
		isPagingated = false
	}

	count, err := database.Count[model.Task](t.db)
	if err != nil {
		t.logger.Error("TaskHandler: List: Count tasks", err)
		app.InternalServerError(c)
		return
	}

	ordering := "desc"
	if pagination.Order == "asc" {
		ordering = "asc"
	}

	tx := t.db.Order(fmt.Sprintf("created_at %s", ordering))

	if isPagingated {
		tx = tx.Limit(pagination.GetLimit()).Offset(pagination.GetOffset())
	}

	tasks, err := database.FindAll[model.Task](tx)
	if err != nil {
		t.logger.Error("TaskHandler: List: Find all tasks", err)
		app.InternalServerError(c)
		return
	}

	app.PaginatedResponse(c, tasks, count, pagination.GetPage(), pagination.GetLimit())
}

func (t TaskHandler) Detail(c *gin.Context) {
	id := c.Param("id")

	task, err := database.First[model.Task](t.db, id)
	if database.IsNotFound(err) {
		app.NotFoundResponse(c, "task not found")
		return
	} else if err != nil {
		t.logger.Error("TaskHandler: Detail: Find task", err)
		app.InternalServerError(c)
		return
	}
	app.SuccessResponse(c, task)
}

func (t TaskHandler) Update(c *gin.Context) {
	id := c.Param("id")
	request, err := view.Bind[view.Title](c)
	if err != nil {
		t.logger.Error("TaskHandler: Update: Bind request", err)
		app.InvalidRequest(c)
		return
	}

	_, err = database.First[model.Task](t.db, id)
	if database.IsNotFound(err) {
		app.NotFoundResponse(c, "task not found to update")
		return
	} else if err != nil {
		t.logger.Error("TaskHandler: Update: Find task", err)
		app.InternalServerError(c)
		return
	}

	err = database.UpdateTask[model.Task](t.db.Where("id = ?", id), map[string]interface{}{"title": request.Title})
	if err != nil {
		// check unique constraint error
		if database.IsUniqueConstraintError(err) {
			app.BadRequest(c, "task with this title already exists")
			return
		}
		t.logger.Error("TaskHandler: Update: Update task", err)
		app.InternalServerError(c)
		return
	}

	app.SuccessResponse(c, gin.H{"message": "task updated"})
}

func (t TaskHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	_, err := database.First[model.Task](t.db, id)
	if database.IsNotFound(err) {
		app.NotFoundResponse(c, "task not found to delete")
		return
	}

	err = database.Delete[model.Task](t.db, id)
	if err != nil {
		t.logger.Error("TaskHandler: Delete: Delete task", err)
		app.InternalServerError(c)
		return
	}
	app.SuccessResponse(c, gin.H{"message": "task deleted"})
}

func (t TaskHandler) DeleteAll(c *gin.Context) {
	err := database.DeleteAll[model.Task](t.db)
	if err != nil {
		t.logger.Error("TaskHandler: DeleteAll: Delete all tasks", err)
		app.InternalServerError(c)
		return
	}

	app.SuccessResponse(c, gin.H{"message": "all tasks deleted"})
}

func (t TaskHandler) RunTask(c *gin.Context) {
	id := c.Param("id")

	t.logger.Info("task id: " + id)

	task, err := database.First[model.Task](t.db, id)
	if database.IsNotFound(err) {
		app.NotFoundResponse(c, "task not found")
		return
	} else if err != nil {
		t.logger.Error("TaskHandler: RunTask: Find task", err)
		app.InternalServerError(c)
		return
	}

	if task.Status == model.StatusCompleted {
		t.logger.Info("TaskHandler: RunTask: task already completed")
		app.BadRequest(c, "task already completed")
		return
	}

	fmt.Println("Task adding into queue : " + task.Title)

	queue := t.queue.GetChannel()

	select {
	case queue <- task:
		t.logger.Info("Task added to queue")
		app.SuccessResponse(c, gin.H{"message": "task added to queue"})
	default:
		t.logger.Info("no worker is available to process the task")
		app.SuccessResponse(c, gin.H{"message": "no worker available to process the task"})
	}

}

func (t TaskHandler) StopWorker(c *gin.Context) {
	for i := 0; i < constants.NumberOfConsumer; i++ {
		select {
		case t.stopChan <- struct{}{}:

		default:
			t.logger.Info("No worker to stop")
		}
	}
	app.SuccessResponse(c, gin.H{"message": "worker stopped"})
}

func (t TaskHandler) StartWorker(c *gin.Context) {
	go consumer.StartWorkerPool(t.db, t.queue, t.logger, t.stopChan)
	app.SuccessResponse(c, gin.H{"message": "worker started"})

}

func NewTaskHandler(db *gorm.DB, logger logger.AppLogger, queue lib.Queue[model.Task], stopChan chan struct{}) TaskHandler {
	return TaskHandler{
		db:       db,
		logger:   logger,
		queue:    queue,
		stopChan: stopChan,
	}
}
