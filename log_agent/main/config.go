package main

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
	"logagent/log_agent/tailf"
)

var (
	appConfig *Config
)

type Config struct {
	logLevel string
	logPath  string

	chanSize    int
	kafkaAddr   string
	collectConf []tailf.CollectConf

	etcdAddr string
	etcdKey  string
}

func loadConf(configType, filename string) (err error) {
	//New config
	conf, err := config.NewConfig(configType, filename)
	if err != nil {
		fmt.Println("new config failed ", err)
	}

	// conf.String("loglevel")
	appConfig := &Config{}
	appConfig.logLevel = conf.String("logs::log_level")
	if len(appConfig.logLevel) == 0 {
		appConfig.logLevel = "debug"
	}

	//conf.String("log_path")
	appConfig.logPath = conf.String("logs::log_path")
	if len(appConfig.logPath) == 0 {
		appConfig.logPath = "./logs"
	}

	appConfig.chanSize, err = conf.Int("logs::chan_size")
	if err != nil {
		appConfig.chanSize = 100
	}

	appConfig.kafkaAddr = conf.String("kafka::server_addr")
	if len(appConfig.kafkaAddr) == 0 {
		err = fmt.Errorf("invalid kafka addr")
		return
	}

	appConfig.etcdAddr = conf.String("etcd::addr")
	if len(appConfig.etcdAddr) == 0 {
		err = fmt.Errorf("invalid etcd addr")
		return
	}

	appConfig.etcdKey = conf.String("etcd::configKey")
	if len(appConfig.etcdKey) == 0 {
		err = fmt.Errorf("invalid etcd key")
	}

	err = loadCollectConf(conf)
	if err != nil {
		fmt.Println("collect conf failed ", err)
		return
	}
	return
}

func loadCollectConf(conf config.Configer) (err error) {
	var cc tailf.CollectConf

	cc.LogPath = conf.String("logs::log_path")
	if len(cc.LogPath) == 0 {
		err = errors.New("invalid collect::log_path")
		return
	}

	cc.Topic = conf.String("logs::topic")
	if len(cc.Topic) == 0 {
		err = errors.New("invalid collect::topic")
		return
	}

	appConfig.collectConf = append(appConfig.collectConf, cc)
	return
}
