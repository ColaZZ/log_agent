package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logagent/log_agent/tailf"
)

func main(){
	filename := "./conf/logagent.conf"
	err := loadConf("ini", filename)
	if err != nil {
		fmt.Printf("load conf failed %s\n", err)
		panic("load conf failed")
		return
	}

	err = initLogger()
	if err != nil {
		fmt.Printf("load logger failed %s\n", err)
		panic("load logger failed")
		return
	}

	logs.Debug("load conf succ, config:%v", appConfig)

	err = tailf.InitTail(appConfig.collectConf, appConfig.chanSize)
	if err != nil {
		logs.Error("init tail failed, err:%v", err)
		return
	}

	err = ServerRun()
	if err != nil {
		logs.Error("serverRUn failed, err:%v", err)
		return
	}

}
