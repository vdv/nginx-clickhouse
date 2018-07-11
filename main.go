package main

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/hpcloud/tail"
	"github.com/vdv/nginx-clickhouse/clickhouse"
	configParser "github.com/vdv/nginx-clickhouse/config"
	"github.com/vdv/nginx-clickhouse/nginx"
)

var (
	locker sync.Mutex
	logs   []string
)

func main() {

	// Read config & incoming flags
	config := configParser.Read()

	// Update config with environment variables if exist
	config.SetEnvVariables()

	nginxParser, err := nginx.GetParser(config)

	if err != nil {
		logrus.Fatal("Can`t parse nginx log format: ", err)
	}

	logs = []string{}
	t, err := tail.TailFile(config.Settings.LogPath, tail.Config{Follow: true, ReOpen: true, MustExist: true})

	if err != nil {
		logrus.Fatal("Can`t tail logfile: ", err)
	}

	logrus.Info("Opening logfile: " + config.Settings.LogPath)

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(config.Settings.Interval))

			if len(logs) > 0 {

				logrus.Info("Preparing to save ", len(logs), " new log entries.")
				locker.Lock()
				err := clickhouse.Save(config, nginx.ParseLogs(nginxParser, logs))

				if err != nil {
					logrus.Error("Can`t save logs: ", err)
				} else {
					logrus.Info("Saved ", len(logs), " new logs.")
					logs = []string{}
				}

				locker.Unlock()
			}
		}
	}()

	// Push new log entries to array
	for line := range t.Lines {
		locker.Lock()
		logs = append(logs, line.Text)
		locker.Unlock()
	}
}
