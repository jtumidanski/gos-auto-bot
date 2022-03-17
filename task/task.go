package task

import (
	"gos-auto-bot/step"
)

const (
	AcademyTask              = "academy"
	AdTask                   = "ads"
	CoalitionTask            = "coalition"
	CouncilTask              = "council"
	DivinationTask           = "divination"
	HaremTask                = "harem"
	LevyTask                 = "levy"
	RankingTask              = "rankings"
	UnionTask                = "union"
	DailyCheckInTask         = "dailyCheckIn"
	ImperialCarnivalHelpTask = "imperialCarnivalHelpTask"
)

type Task struct {
	Name    string
	NoxId   string
	Task    string
	Startup int
	Steps   []step.Step
}

func Create(name string, botName string, botId string, startup int, stepProducer step.Producer) Task {
	return Task{
		Name:    botName,
		NoxId:   botId,
		Task:    name,
		Startup: startup,
		Steps:   stepProducer(),
	}
}
