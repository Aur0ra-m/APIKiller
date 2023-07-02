package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type RotateFileConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Level      logrus.Level
	LocalTime  bool
	Formatter  logrus.Formatter
}

type RotateFileHook struct {
	Config         RotateFileConfig
	nextRotateTime time.Time
	logWriter      *lumberjack.Logger
}

func NewRotateFileHook(config RotateFileConfig) (logrus.Hook, error) {
	hook := RotateFileHook{
		Config: config,
	}

	// load rotate log system
	err := hook.rotateLogFile()
	if err != nil {
		return nil, err
	}

	return &hook, nil
}

func (hook *RotateFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.Config.Level+1]
}

func (hook *RotateFileHook) Fire(entry *logrus.Entry) (err error) {
	b, err := hook.Config.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.logWriter.Write(b)
	return
}

func (hook *RotateFileHook) rotateLogFile() error {
	now := time.Now()
	if now.After(hook.nextRotateTime) {
		// close current log file
		if hook.logWriter != nil {
			err := hook.logWriter.Close()
			if err != nil {
				return err
			}
		}

		// calculate next rotate time
		hook.nextRotateTime = now.Truncate(24 * time.Hour).Add(24 * time.Hour)

		// rename log filename according to date
		date := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())

		_, err1 := os.Stat(hook.Config.Filename)
		if err1 == nil {
			err := os.Rename(hook.Config.Filename, "./log/"+date+".log")
			if err != nil {
				return err
			}
		}

		// create new log file
		hook.logWriter = &lumberjack.Logger{
			Filename:   hook.Config.Filename,
			MaxSize:    hook.Config.MaxSize,
			MaxBackups: hook.Config.MaxBackups,
			MaxAge:     hook.Config.MaxAge,
			LocalTime:  hook.Config.LocalTime,
		}
	}
	return nil
}
