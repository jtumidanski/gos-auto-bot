package bot

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gos-auto-bot/bot/database"
	"gos-auto-bot/identity"
	"gos-auto-bot/task"
	time2 "gos-auto-bot/time"
	"os"
	"os/exec"
	"time"
)

func NewWorker(fl logrus.FieldLogger) func(settings *Settings, tasks chan task.Task, out chan string) {
	return func(settings *Settings, tasks chan task.Task, out chan string) {
		l := fl.WithField("thread", "bot-"+settings.Name)
		l.Infof("Bot %s initializing.", settings.Name)
		for true {
			select {
			case t := <-tasks:
				botTask(l)(settings, t)
				out <- t.Task
			}
		}
	}
}

func botTask(l logrus.FieldLogger) func(settings *Settings, t task.Task) {
	botPath, ok := os.LookupEnv("BOT_PATH")
	if !ok {
		l.Fatalf("Unable to lookup BOT_PATH")
	}
	noxExe, ok := os.LookupEnv("NOX_EXE")
	if !ok {
		l.Fatalf("Unable to lookup NOX_EXE")
	}

	return func(settings *Settings, t task.Task) {
		l.Infof("Bot %s received task %s.", settings.Name, t.Task)

		launchCmd := exec.Command(noxExe, "-clone:"+t.NoxId, "-package:com.arlosoft.macrodroid")
		err := launchCmd.Start()
		if err != nil {
			l.WithError(err).Errorf("Unable to launch NOX instance for %s.", t.Name)
		}
		l.Infof("Waiting %ds for %s to startup. pid=%d", t.Startup, t.Name, launchCmd.Process.Pid)
		time.Sleep(time.Duration(t.Startup) * time.Second)
		l.Infof("Starting %s's task %s.", t.Name, t.Task)
		i := identity.BotIdentity(settings.NoxId, settings.Name, settings.DeviceId)
		for _, s := range t.Steps {
			s.StepFunc(i)
			time.Sleep(time.Duration(s.Delay) * time.Millisecond)
		}
		l.Infof("%s's task %s completed.", t.Name, t.Task)

		db, err := database.GetOrCreate(l)(botPath + "\\" + settings.Name + "\\db")
		if err != nil {
			l.WithError(err).Fatalf("Unable to read %s's database.", t.Name)
			return
		}
		var dbt database.Task
		switch t.Task {
		case task.AcademyTask:
			dbt = db.Academy
		case task.AdTask:
			dbt = db.Ads
		case task.CoalitionTask:
			dbt = db.Coalition
		case task.CouncilTask:
			dbt = db.Council
		case task.DivinationTask:
			dbt = db.Divination
		case task.HaremTask:
			dbt = db.Harem
		case task.LevyTask:
			dbt = db.Levy
		case task.RankingTask:
			dbt = db.Rankings
		case task.UnionTask:
			dbt = db.Union
		case task.DailyCheckInTask:
			dbt = db.DailyCheckIn
		}

		count := dbt.Count
		loc, _ := time.LoadLocation("America/New_York")
		lastExecution, err := time.ParseInLocation(time2.TimeFormat, dbt.Execution, loc)
		if err != nil {
			return
		}
		if lastExecution.YearDay() != time.Now().YearDay() {
			count = 0
		}
		dbt.Execution = time.Now().In(loc).Format(time2.TimeFormat)
		dbt.Count = count + 1

		switch t.Task {
		case task.AcademyTask:
			db.Academy = dbt
		case task.AdTask:
			db.Ads = dbt
		case task.CoalitionTask:
			db.Coalition = dbt
		case task.CouncilTask:
			db.Council = dbt
		case task.DivinationTask:
			db.Divination = dbt
		case task.HaremTask:
			db.Harem = dbt
		case task.LevyTask:
			db.Levy = dbt
		case task.RankingTask:
			db.Rankings = dbt
		case task.UnionTask:
			db.Union = dbt
		case task.DailyCheckInTask:
			db.DailyCheckIn = dbt
		}

		dbPath := botPath + "\\" + settings.Name + "\\db"
		f, err := os.Create(dbPath)
		if err != nil {
			l.WithError(err).Errorf("Unable to create database file.")
			return
		}
		e := json.NewEncoder(f)
		err = e.Encode(db)
		if err != nil {
			l.WithError(err).Errorf("Unable to write new database file.")
			return
		}

		cmd := exec.Command(noxExe, "-clone:"+t.NoxId, "-quit")
		err = cmd.Start()
		if err != nil {
			l.WithError(err).Errorf("Unable to exit NOX instance for %s.", t.Name)
		}
		time.Sleep(time.Duration(5) * time.Second)
		_ = launchCmd.Process.Kill()
		time.Sleep(time.Duration(1) * time.Second)
		return
	}
}
