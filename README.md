# Task Manager

## Prerequisites

- Go 1.23

## Installation

1. Unzip the project:
2. Enter command "make". This will create a binary executable "task_manager"

To install dependencies

enter command "go mod tidy"

This application can run in standalone way in both methods API or CLI.

**In CLI mode**

**Commands**

1. task_manager add "task title"
   It adds tasks to sqlite in pending state
2. task_manager list
   It lists all the tasks created
3. task_manager process
   It start consumer worker pool of 4(in my case) and fetches pending tasks from database and adds them in queue, therefore workers start consuming tasks from queue and marks them completed

**In API mode**

GET /api/task - List all the tasks in normal view or paginated view

GET /api/task/:id - Get detail task for particular task

POST /api/task - Create task

    body :

{

"title":"task created 1"

}

PUT /api/task/:id - Update task title

    body:

{

"title":"cli task"

}

DELETE /api/task/:id - Delete a task

DELETE /api/task/all - Delete all tasks

GET /api/worker-start - Start Consumer/Worker pool

GET /api/worker-stop - Stop Consumer/Worker pool

GET /api/task/:id/run - process a single task from api level by giving its id in param
