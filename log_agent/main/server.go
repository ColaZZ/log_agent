package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logagent/log_agent/kafka"
	"logagent/log_agent/tailf"
	"time"
)

func ServerRun() (err error) {
	for {
		msg := tailf.GetOneLine()
		err = sendToKafka(msg)
		if err != nil {
			logs.Error("send kafka failed, err ：%v", err)
			time.Sleep(time.Second)
			continue
		}
	}
}

func sendToKafka(msg *tailf.TextMsg) (err error) {
	fmt.Println("read msg:%s, topic:%s", msg.Msg, msg.Topic)
	err = kafka.SendToKafka(msg.Msg, msg.Topic)
	return
}
