package database

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	time2 "gos-auto-bot/time"
	"os"
	"time"
)

func GetOrCreate(l logrus.FieldLogger) func(path string) (*Model, error) {
	return func(path string) (*Model, error) {
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			loc, _ := time.LoadLocation("America/New_York")
			defaultTime, err := time.ParseInLocation(time2.TimeFormat, "2021-01-29 10:51:09.518996", loc)
			if err != nil {
				return nil, err
			}
			defaultTimeStr := defaultTime.Format(time2.TimeFormat)

			nd := &Model{
				Levy: Task{
					Execution: defaultTimeStr,
				},
				Divination: Task{
					Execution: defaultTimeStr,
				},
				Council: Task{
					Execution: defaultTimeStr,
				},
				Academy: Task{
					Execution: defaultTimeStr,
				},
				Rankings: Task{
					Execution: defaultTimeStr,
				},
				Harem: Task{
					Execution: defaultTimeStr,
				},
				Coalition: Task{
					Execution: defaultTimeStr,
				},
				Ads: Task{
					Execution: defaultTimeStr,
				},
				Union: Task{
					Execution: defaultTimeStr,
				},
			}
			f, err := os.Create(path)
			if err != nil {
				l.WithError(err).Errorf("Unable to create database file.")
				return nil, err
			}
			e := json.NewEncoder(f)
			err = e.Encode(nd)
			if err != nil {
				l.WithError(err).Errorf("Unable to write new database file.")
				return nil, err
			}
			err = f.Close()
			if err != nil {
				l.WithError(err).Errorf("Unable to close new database file.")
				return nil, err
			}
		}
		f, err := os.Open(path)
		if err != nil {
			l.WithError(err).Errorf("Unable to open database file.")
			return nil, err
		}
		db := &Model{}
		d := json.NewDecoder(f)
		err = d.Decode(db)
		if err != nil {
			l.WithError(err).Errorf("Unable to parse bot settings.")
			return nil, err
		}
		err = f.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close bot settings file.")
			return nil, err
		}
		return db, nil
	}
}
