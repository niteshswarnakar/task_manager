/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/aquasecurity/table"
	"github.com/niteshswarnakar/task_manager/internal/cli"
	"github.com/niteshswarnakar/task_manager/internal/constants.go"
	"github.com/niteshswarnakar/task_manager/internal/infrastructure/database"
	"github.com/niteshswarnakar/task_manager/internal/lib"
	"github.com/niteshswarnakar/task_manager/internal/logger"
	"github.com/niteshswarnakar/task_manager/internal/model"
	"github.com/spf13/cobra"
)

func init() {
	appLogger := logger.NewAppLogger()
	appLogger.Info("Task Manager started")

	db, err := database.InitDB(constants.DBName)
	if err != nil {
		appLogger.Panic(err)
	}

	err = database.AutoMigrate(db)
	if err != nil {
		appLogger.Panic(err)
	}

	queue := lib.NewQueue[model.Task]()
	stopChan := make(chan struct{})

	appCli := cli.NewAppCli(db, appLogger, queue, stopChan)
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "adds task into database(eg. task_manager add 'task title')",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			title := args[0]
			err := appCli.Add(title)
			if err != nil {
				appLogger.Error("CLI->", err)
				return
			}
			appLogger.Info("CLI-> Task created : " + title)
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "list all the tasks in the database",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := appCli.List()
			if err != nil {
				appLogger.Error("Cli: List Command: Find all tasks", err)
			}

			display := table.New(os.Stdout)
			display.SetRowLines(false)
			display.SetHeaders("ID", "Title", "Status", "Created At")

			// pending emoji
			for _, task := range tasks {
				status := "⏳"
				if task.Status == model.StatusCompleted {
					status = "✅"
				}
				display.AddRow(task.ID, task.Title, string(fmt.Sprintf("%s %s", status, string(task.Status))), task.CreatedAt.Format("2006-01-02 15:04:05"))
				// appLogger.Info(fmt.Sprintf("%s - %s", task.Title, task.Status))
			}

			display.Render()
		},
	}

	var processCmd = &cobra.Command{
		Use:   "process",
		Short: "perform task by the consumer",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := appCli.Process()
			if err != nil {
				appLogger.Error("Cli: Process Command: Perform task", err)
			}
		},
	}

	apiServerCommand := ApiServerCmd(db, appLogger, queue, stopChan)

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(processCmd)
	rootCmd.AddCommand(apiServerCommand)
}
