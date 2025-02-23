package constants

const NumberOfConsumer int = 4

const DBName string = "task.db"

type CommandType string

const (
	CommandType_Add     CommandType = "add"
	CommandType_List    CommandType = "list"
	CommandType_Process CommandType = "process"
	CommandType_Api     CommandType = "api"
)
