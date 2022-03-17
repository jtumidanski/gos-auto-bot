package bot

import (
	"gos-auto-bot/step"
	"gos-auto-bot/task"
)

func CreateTask(botName string, botId string, startup int) func(name string, producer step.Producer) task.Task {
	return func(name string, producer step.Producer) task.Task {
		return task.Create(name, botName, botId, startup, producer)
	}
}