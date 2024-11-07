package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ChuhC/crontab/worker/logger"
	"github.com/ChuhC/crontab/worker/worker"
)

func main() {

	flagSet := flag.NewFlagSet("crontab worker", flag.ExitOnError)
	flagSet.String("config", "", "path to config file")

	flagSet.Parse(os.Args[1:])

	// load configuration
	cfg := worker.ConfigData{}
	configFile := flagSet.Lookup("config").Value.String()
	if configFile != "" {
		_, err := toml.DecodeFile(configFile, &cfg)
		if err != nil {
			panic(fmt.Sprintf("error: failed to load config file %s - %s", configFile, err.Error()))
		}
	}

	// init logger
	logger.InitLogger(cfg.Log.FileName, cfg.Log.Daily, cfg.Log.MaxDays, cfg.Log.Console, cfg.Log.Level)

	// init scheduler
	worker.InitScheduler()

	// init job mgr
	worker.InitJobMgr(cfg.Etcd.EndPoints, cfg.Etcd.DialTimeout)

	// catch ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	for {
		select {
		case <-c:
			return
		default:
			time.Sleep(time.Second)
		}
	}

}
