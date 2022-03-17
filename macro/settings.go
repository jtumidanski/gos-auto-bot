package macro

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
)

type Settings struct {
	Startup struct {
		Duration int `json:"duration"`
	} `json:"startup"`
	AdsOffset int `json:"ads_offset"`
	Scripts   struct {
		Levy         Setting `json:"levy"`
		Divination   Setting `json:"divination"`
		Council      Setting `json:"council"`
		Academy      Setting `json:"academy"`
		Rankings     Setting `json:"rankings"`
		Harem        Setting `json:"harem"`
		Coalition    Setting `json:"coalition"`
		Ads          Setting `json:"ads"`
		Union        Setting `json:"union"`
		DailyCheckIn Setting `json:"dailyCheckIn"`
	} `json:"scripts"`
}

type Setting struct {
	Times   int  `json:"times"`
	Enabled bool `json:"enabled"`
}

func ReadSettings(l logrus.FieldLogger) (*Settings, error) {
	macroSettings, ok := os.LookupEnv("MACRO_SETTINGS")
	if !ok {
		l.Fatalf("Unable to lookup MACRO_SETTINGS")
	}

	f, err := os.Open(macroSettings)
	if err != nil {
		l.WithError(err).Errorf("Unable to read macro settings.")
		return nil, err
	}

	var ms = &Settings{}
	d := json.NewDecoder(f)
	err = d.Decode(ms)
	if err != nil {
		l.WithError(err).Errorf("Unable to parse macro settings.")
		return nil, err
	}
	err = f.Close()
	if err != nil {
		l.WithError(err).Errorf("Unable to close macro settings file.")
		return nil, err
	}
	return ms, nil
}

func WriteAdOffset(l logrus.FieldLogger) func(value int) error {
	macroSettings, ok := os.LookupEnv("MACRO_SETTINGS")
	if !ok {
		l.Fatalf("Unable to lookup MACRO_SETTINGS")
	}

	return func(value int) error {
		s, err := ReadSettings(l)
		if err != nil {
			l.WithError(err).Errorf("Error opening the settings file for updating.")
			return err
		}
		s.AdsOffset = value

		f, err := os.Create(macroSettings)
		if err != nil {
			l.WithError(err).Errorf("Unable to open macro settings.")
			return err
		}
		defer f.Close()

		e := json.NewEncoder(f)
		err = e.Encode(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to write macro settings.")
			return err
		}
		return nil
	}
}
