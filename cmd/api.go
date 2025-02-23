package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/niteshswarnakar/task_manager/internal/app/route"
	"github.com/niteshswarnakar/task_manager/internal/lib"
	"github.com/niteshswarnakar/task_manager/internal/logger"
	"github.com/niteshswarnakar/task_manager/internal/model"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func ApiServerCmd(db *gorm.DB, logger logger.AppLogger, queue lib.Queue[model.Task], stopChan chan struct{}) *cobra.Command {

	return &cobra.Command{
		Use:   "api",
		Short: "starts the api server to perform task operations",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			server := gin.Default()

			routeGroup := server.Group("/api")
			route.RegisterTaskRoutes(routeGroup, db, logger, queue, stopChan)

			err := server.Run(":5000")
			if err != nil {
				logger.Panic(err)
			}
		},
	}
}
