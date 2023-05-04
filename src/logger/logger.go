package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var log = logrus.New()

func init() {
	log.Out = os.Stdout
	nowTime := time.Now()
	year, month, day := nowTime.Date()
	hour, min, sec := nowTime.Clock()
	log.WithFields(logrus.Fields{
		"year":  year,
		"month": month,
		"day":   day,
		"hour":  hour,
		"min":   min,
		"sec":   sec,
	}).Info("日志系统启动")
}

func GetLogger() *logrus.Logger {
	return log
}
