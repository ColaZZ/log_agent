package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
)

func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = appConfig.logPath
	config["level"] = convertLog(appConfig.logLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("initLogger failed, marshal err:", err)
		return
	}

	_ = logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}

func convertLog(level string) int {
	switch level {
	case "Debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}
	return logs.LevelTrace
}
