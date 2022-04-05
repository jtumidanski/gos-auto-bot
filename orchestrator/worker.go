package orchestrator

import (
	"context"
	"github.com/sirupsen/logrus"
	"gos-auto-bot/bot"
	"gos-auto-bot/bot/database"
	"gos-auto-bot/macro"
	"gos-auto-bot/step"
	"gos-auto-bot/task"
	time2 "gos-auto-bot/time"
	"os"
	"sync"
	"time"
)

type Instruction struct {
	Name string
	Task string
}

func NewOrchestrator(fl logrus.FieldLogger) func(ctx context.Context, wg *sync.WaitGroup, instructions <-chan Instruction) {
	l := fl.WithField("thread", "orchestrator")
	return func(ctx context.Context, wg *sync.WaitGroup, instructions <-chan Instruction) {
		wg.Add(1)

		ms, err := macro.ReadSettings(l)
		if err != nil {
			return
		}

		names, err := bot.ReadNames(l)

		botTaskChan := make(map[string]chan task.Task, 0)
		botOutChan := make(map[string]chan string, 0)
		botTasks := make(map[string]map[string]bool, 0)
		for _, name := range names {
			in := make(chan task.Task, 9)
			botTaskChan[name] = in
			out := make(chan string, 9)
			botOutChan[name] = out

			botTasks[name] = make(map[string]bool, 0)
			for _, t := range []string{task.AcademyTask, task.AdTask, task.CoalitionTask, task.CouncilTask, task.DivinationTask, task.HaremTask, task.LevyTask, task.RankingTask, task.UnionTask, task.DailyCheckInTask} {
				botTasks[name][t] = false
			}

			bs, err := bot.Read(l)(name)
			if err != nil {
				l.WithError(err).Fatalf("Unable to load bot settings.")
			}

			go bot.NewWorker(l)(bs, in, out)
		}

		running := true
		for running {
			select {
			case <-ctx.Done():
				running = false
			case instruction := <-instructions:
				addInstruction(l)(ms, instruction.Name, instruction.Task, botTasks, botTaskChan)
			case <-time.After(15 * time.Second):
				orchestratorStep(l)(ms, names, botTasks, botTaskChan, botOutChan)
			}
		}

		wg.Done()
		l.Infof("Shutting down orchestrator")
	}
}

func addInstruction(l logrus.FieldLogger) func(ms *macro.Settings, name string, taskName string, botTasks map[string]map[string]bool, botTaskChan map[string]chan task.Task) {
	return func(ms *macro.Settings, name string, taskName string, botTasks map[string]map[string]bool, botTaskChan map[string]chan task.Task) {
		bs, err := bot.Read(l)(name)
		if err != nil {
			l.WithError(err).Fatalf("Unable to load bot settings.")
		}

		var t task.Task
		taskCreator := bot.CreateTask(bs.Name, bs.NoxId, ms.Startup.Duration)

		switch taskName {
		case task.AcademyTask:
			t = taskCreator(task.AcademyTask, step.CreateAcademySteps(l)(bs.AcademySeats))
		case task.AdTask:
			t = taskCreator(task.AdTask, step.CreateAdSteps(l))
		case task.CoalitionTask:
			t = taskCreator(task.CoalitionTask, step.CreateCoalitionSteps(l))
		case task.CouncilTask:
			t = taskCreator(task.CouncilTask, step.CreateCouncilSteps(l)(bs.VIP, bs.Envoys))
		case task.DivinationTask:
			t = taskCreator(task.DivinationTask, step.CreateDivinationSteps(l))
		case task.HaremTask:
			t = taskCreator(task.HaremTask, step.CreateHaremSteps(l)(bs.Scripts.Harem.Total))
		case task.LevyTask:
			t = taskCreator(task.LevyTask, step.CreateLevySteps(l)(bs.VIP, bs.Scripts.Levy.Total))
		case task.RankingTask:
			t = taskCreator(task.RankingTask, step.CreateRankingsSteps(l))
		case task.UnionTask:
			t = taskCreator(task.UnionTask, step.CreateUnionSteps(l))
		case task.DailyCheckInTask:
			t = taskCreator(task.DailyCheckInTask, step.CreateDailyCheckInSteps(l))
		case task.ImperialCarnivalHelpTask:
			t = taskCreator(task.ImperialCarnivalHelpTask, step.CreateImperialCarnivalHelpSteps(l))
		}

		aq := botTasks[name][t.Task]
		if !aq {
			botTasks[name][t.Task] = true
			botTaskChan[name] <- t
			l.Infof("Scheduling task %s for %s.", t.Task, bs.Name)
		}

	}
}

func orchestratorStep(l logrus.FieldLogger) func(ms *macro.Settings, names []string, botTasks map[string]map[string]bool, botTaskChan map[string]chan task.Task, botOutChan map[string]chan string) {
	botPath, ok := os.LookupEnv("BOT_PATH")
	if !ok {
		l.Fatalf("Unable to lookup BOT_PATH")
	}

	return func(ms *macro.Settings, names []string, botTasks map[string]map[string]bool, botTaskChan map[string]chan task.Task, botOutChan map[string]chan string) {
		for _, name := range names {
			bs, err := bot.Read(l)(name)
			if err != nil {
				l.WithError(err).Fatalf("Unable to load bot settings.")
			}
			l.Infof("Checking in on %s.", bs.Name)

			db, err := database.GetOrCreate(l)(botPath + "\\" + name + "\\db")
			if err != nil {
				l.WithError(err).Fatalf("Unable to get or create bot database.")
			}

			tasks := make([]task.Task, 0)
			taskCreator := bot.CreateTask(bs.Name, bs.NoxId, ms.Startup.Duration)

			if canExecute(task.AcademyTask, ms.Scripts.Academy.Enabled, ms.Scripts.Academy.Times, db.Academy, bs.Scripts.Academy.Enabled, bs.Scripts.Academy.Total, bs.Scripts.Academy.Interval) {
				tasks = append(tasks, taskCreator(task.AcademyTask, step.CreateAcademySteps(l)(bs.AcademySeats)))
			}
			if canExecute(task.AdTask, ms.Scripts.Ads.Enabled, ms.Scripts.Ads.Times, db.Ads, bs.Scripts.Ads.Enabled, 0, 0) {
				tasks = append(tasks, taskCreator(task.AdTask, step.CreateAdSteps(l)))
			}
			if canExecute(task.CoalitionTask, ms.Scripts.Coalition.Enabled, ms.Scripts.Coalition.Times, db.Coalition, bs.Scripts.Coalition.Enabled, 0, 0) {
				tasks = append(tasks, taskCreator(task.CoalitionTask, step.CreateCoalitionSteps(l)))
			}
			if canExecute(task.CouncilTask, ms.Scripts.Council.Enabled, ms.Scripts.Council.Times, db.Council, bs.Scripts.Council.Enabled, bs.Scripts.Council.Total, bs.Scripts.Council.Interval) {
				tasks = append(tasks, taskCreator(task.CouncilTask, step.CreateCouncilSteps(l)(bs.VIP, bs.Envoys)))
			}
			if canExecute(task.DivinationTask, ms.Scripts.Divination.Enabled, ms.Scripts.Divination.Times, db.Divination, bs.Scripts.Divination.Enabled, 0, 0) {
				tasks = append(tasks, taskCreator(task.DivinationTask, step.CreateDivinationSteps(l)))
			}
			if canExecute(task.HaremTask, ms.Scripts.Harem.Enabled, ms.Scripts.Harem.Times, db.Harem, bs.Scripts.Harem.Enabled, bs.Scripts.Harem.Total, bs.Scripts.Harem.Interval) {
				tasks = append(tasks, taskCreator(task.HaremTask, step.CreateHaremSteps(l)(bs.Scripts.Harem.Total)))
			}
			if canExecute(task.LevyTask, ms.Scripts.Levy.Enabled, ms.Scripts.Levy.Times, db.Levy, bs.Scripts.Levy.Enabled, bs.Scripts.Levy.Total, bs.Scripts.Levy.Interval) {
				tasks = append(tasks, taskCreator(task.LevyTask, step.CreateLevySteps(l)(bs.VIP, bs.Scripts.Levy.Total)))
			}
			if canExecute(task.RankingTask, ms.Scripts.Rankings.Enabled, ms.Scripts.Rankings.Times, db.Rankings, bs.Scripts.Rankings.Enabled, 0, 0) {
				tasks = append(tasks, taskCreator(task.RankingTask, step.CreateRankingsSteps(l)))
			}
			if canExecute(task.UnionTask, ms.Scripts.Union.Enabled, ms.Scripts.Union.Times, db.Union, bs.Scripts.Union.Enabled, 0, 0) {
				tasks = append(tasks, taskCreator(task.UnionTask, step.CreateUnionSteps(l)))
			}
			if canExecute(task.DailyCheckInTask, ms.Scripts.DailyCheckIn.Enabled, ms.Scripts.DailyCheckIn.Times, db.DailyCheckIn, bs.Scripts.DailyCheckIn.Enabled, 0, 0) {
				tasks = append(tasks, taskCreator(task.DailyCheckInTask, step.CreateDailyCheckInSteps(l)))
			}

			if len(tasks) > 0 {
				for _, t := range tasks {
					aq := botTasks[name][t.Task]
					if !aq {
						botTasks[name][t.Task] = true
						botTaskChan[name] <- t
						l.Infof("Scheduling task %s for %s.", t.Task, bs.Name)
					}
				}
			}
		}
		for _, name := range names {
			out := botOutChan[name]
			reading := true
			for reading {
				select {
				case t := <-out:
					botTasks[name][t] = false
					l.Infof("Task %s for %s removed from schedule.", t, name)
				case <-time.After(1 * time.Second):
					reading = false
				}
			}
		}
	}
}

func canExecute(name string, macroEnabled bool, maxCount int, db database.Task, botEnabled bool, total int, interval int) bool {
	if !macroEnabled {
		return false
	}
	if !botEnabled {
		return false
	}
	loc, _ := time.LoadLocation("America/New_York")
	lastExecution, err := time.ParseInLocation(time2.TimeFormat, db.Execution, loc)
	if err != nil {
		return false
	}
	if lastExecution.YearDay() != time.Now().YearDay() {
		return true
	}

	if name == task.LevyTask || name == task.HaremTask || name == task.AcademyTask || name == task.CouncilTask {
		offset := total * interval
		nextExecution := lastExecution.Add(time.Duration(offset) * time.Second)
		if nextExecution.After(time.Now().UTC()) {
			return false
		}
		if name == task.CouncilTask && db.Count >= maxCount {
			return false
		}
	} else if db.Count >= maxCount {
		return false
	}
	return true
}
