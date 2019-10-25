package tailf

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"sync"
	"time"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

type CollectConf struct {
	LogPath string
	Topic   string
}

type TailObj struct {
	tail     *tail.Tail
	conf     CollectConf
	status   int
	exitChan chan int
}

type TextMsg struct {
	Msg   string
	Topic string
}

type TailObjMgr struct {
	tailObjs []*TailObj
	msgChan  chan *TextMsg
	lock     sync.Mutex
}

var (
	tailObjMgr *TailObjMgr
)

func GetOneLine() (msg *TextMsg) {
	msg = <-tailObjMgr.msgChan
	return msg
}

func UpdateConfig(confs []CollectConf) (err error) {
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()
	for _, oneConf := range confs {
		var isRunning = false
		for _, obj := range tailObjMgr.tailObjs {
			if oneConf.LogPath == obj.conf.LogPath {
				isRunning = true
				break
			}
		}
		if isRunning {
			continue
		}
		createNewTask(oneConf)
	}

	var tailObjs []*TailObj
	for _, obj := range tailObjMgr.tailObjs {
		obj.status = StatusDelete
		for _, oneConf := range confs {
			if oneConf.LogPath == obj.conf.LogPath {
				obj.status = StatusNormal
				break
			}
		}
		if obj.status == StatusDelete {
			obj.exitChan <- 1
			continue
		}
		tailObjs = append(tailObjs, obj)
	}
	tailObjMgr.tailObjs = tailObjs
	return
}

func InitTail(conf []CollectConf, chanSize int) (err error) {
	if len(conf) == 0 {
		err = fmt.Errorf("invalid config for log collect:%v", conf)
		return
	}

	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *TextMsg, chanSize),
	}
	for _, v := range conf {
		createNewTask(v)
	}

	return
}

func createNewTask(oneconf CollectConf) {
	obj := &TailObj{
		conf: oneconf,
	}

	tails, errTail := tail.TailFile(oneconf.LogPath, tail.Config{
		ReOpen:    true,
		MustExist: false,
		Poll:      true,
		Follow:    false,
	})

	if errTail != nil {
		logs.Error("tail file err:", errTail)
		return
	}

	obj.tail = tails
	tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)

	go readFromTail(obj)
}

func readFromTail(tailObj *TailObj) {
	for true {
		line, ok := <-tailObj.tail.Lines
		if !ok {
			logs.Warn("tail file close reopne, filenam:%s\n", tailObj)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		textMsg := &TextMsg{
			Msg:   line.Text,
			Topic: tailObj.conf.Topic,
		}

		tailObjMgr.msgChan <- textMsg
	}
}
