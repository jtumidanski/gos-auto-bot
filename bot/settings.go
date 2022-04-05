package bot

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"os"
)

type Settings struct {
	NoxId        string `json:"noxId"`
	Name         string `json:"name"`
	DeviceId     string `json:"deviceId"`
	VIP          int    `json:"vip"`
	AcademySeats int    `json:"academy_seats"`
	Envoys       int    `json:"envoys"`
	Scripts      struct {
		Levy         intervalSetting  `json:"levy"`
		Divination   singleUseSetting `json:"divination"`
		Council      intervalSetting  `json:"council"`
		Academy      intervalSetting  `json:"academy"`
		Rankings     singleUseSetting `json:"rankings"`
		Harem        intervalSetting  `json:"harem"`
		Coalition    singleUseSetting `json:"coalition"`
		Ads          singleUseSetting `json:"ads"`
		Union        singleUseSetting `json:"union"`
		DailyCheckIn singleUseSetting `json:"dailyCheckIn"`
	} `json:"scripts"`
}

type singleUseSetting struct {
	Enabled bool `json:"enabled"`
}

type intervalSetting struct {
	Enabled  bool `json:"enabled"`
	Total    int  `json:"total"`
	Interval int  `json:"interval"`
}

func Read(l logrus.FieldLogger) func(name string) (*Settings, error) {
	return func(name string) (*Settings, error) {
		botPath, ok := os.LookupEnv("BOT_PATH")
		if !ok {
			l.Errorf("Unable to lookup BOT_PATH")
			return nil, errors.New("env error")
		}

		botDir := botPath + "\\" + name + "\\settings.json"
		f, err := os.Open(botDir)
		if err != nil {
			l.WithError(err).Errorf("Unable to read bot settings.")
			return nil, err
		}
		var bs = &Settings{}
		d := json.NewDecoder(f)
		err = d.Decode(bs)
		if err != nil {
			l.WithError(err).Fatalf("Unable to parse bot settings.")
			return nil, err
		}
		err = f.Close()
		if err != nil {
			l.WithError(err).Fatalf("Unable to close bot settings file.")
			return nil, err
		}
		return bs, nil
	}
}
