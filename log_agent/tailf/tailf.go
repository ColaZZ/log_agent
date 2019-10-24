package tailf

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"time"
)

type CollectConf struct {
	LogPath string
	Topic   string
}

type TailObj struct {
	tail *tail.Tail
	conf CollectConf
}

type TestMsg struct {
	Msg   string
	Topic string
}

type TailObjMgr struct {
	tailObjs []*TailObj
	msgChan  chan *TestMsg
}

var (
	tailObjMgr *TailObjMgr
)

func InitTail(conf []CollectConf, chanSize int) (err error) {
	if len(conf) == 0 {
		err = fmt.Errorf("invalid config for log collect:%v", conf)
		return
	}

	tailObjMgr = &TailObjMgr{
		msgChan : make(chan *TestMsg, chanSize),
	}
	for _, v := range conf {
		obj := &TailObj{
			conf: v,
		}

		tails, errTail := tail.TailFile(v.LogPath, tail.Config{
			ReOpen:    true,
			MustExist: false,
			Poll:      true,
			Follow:    false,
		})

		if errTail != nil {
			err = errTail
			logs.Error("tail file err:", err)
			return
		}

		obj.tail = tails
		tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)

		go readFromTail(obj)
	}

	return
}

func readFromTail(tailObj *TailObj) {
	for true {
		line, ok := <- tailObj.tail.Lines
		if !ok {
			logs.Warn("tail file close reopne, filenam:%s\n", tailObj)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		textMsg := &TestMsg{
			Msg:   line.Text,
			Topic: tailObj.conf.Topic,
		}

		tailObjMgr.msgChan <- textMsg
	}
}
