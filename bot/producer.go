package bot

import (
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func ReadNames(l logrus.FieldLogger) ([]string, error) {
	botPath, ok := os.LookupEnv("BOT_PATH")
	if !ok {
		l.Errorf("Unable to lookup BOT_PATH")
		return nil, errors.New("env error")
	}

	f, err := os.Open(botPath)
	if err != nil {
		l.WithError(err).Errorf("Unable to bot path.")
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, errors.New("bot path provided is a file not a directory")
	}

	fs, err := ioutil.ReadDir(botPath)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	for _, cf := range fs {
		l.Infof("Found %s for parsing.", cf.Name())
		names = append(names, cf.Name())
	}
	return names, nil
}
