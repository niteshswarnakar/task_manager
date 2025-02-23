package route

import (
	"github.com/gin-gonic/gin"
	"github.com/niteshswarnakar/task_manager/internal/app/handler"
	"github.com/niteshswarnakar/task_manager/internal/lib"
	"github.com/niteshswarnakar/task_manager/internal/logger"
	"github.com/niteshswarnakar/task_manager/internal/model"
	"gorm.io/gorm"
)

func RegisterTaskRoutes(group *gin.RouterGroup, db *gorm.DB, logger logger.AppLogger, queue lib.Queue[model.Task], stopChan chan struct{}) {
	taskHandler := handler.NewTaskHandler(db, logger, queue, stopChan)
	group.POST("/task", taskHandler.Create)
	group.GET("/task", taskHandler.List)
	group.GET("/task/:id", taskHandler.Detail)
	group.GET("/task/:id/run", taskHandler.RunTask)
	group.PUT("/task/:id", taskHandler.Update)
	group.DELETE("/task/:id", taskHandler.Delete)
	group.DELETE("/task/all", taskHandler.DeleteAll)
	group.GET("/worker-stop", taskHandler.StopWorker)
	group.GET("/worker-start", taskHandler.StartWorker)
}
