package tailf

import (
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"sync"
	"time"
)

//状态
const (
	StatusNormal = 1
	StatusDelete = 2
)

//方便其他包调用
type CollectConf struct {
	LogPath string `json:"log_path"`
	Topic   string `json:"topic"`
}

//存入 Collect
type TailObj struct {
	tail 		*tail.Tail
	conf 		CollectConf
	status 		int
	exitChan 	chan int
}

//定义message信息
type TextMsg struct {
	Msg		string
	Topic	string
}

//管理系统所有的tail对象
type TailObjMgr struct {
	TailObjs	[]*TailObj
	msgChan		chan *TextMsg
	lock		sync.Locker
}

//定义全局变量
var (
	tailObjMgr *TailObjMgr
)

//初始化
func InitTail(conf []CollectConf, chanSize int) (err error) {

	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *TextMsg, chanSize), //创建chan
	}
	//加载配置
	if len(conf) == 0 {
		logs.Error("无效的collect配置：%v, err：%v",conf, err)
	}

	//循环导入
	for _, v := range conf {
		createNewTask(v)
	}

	return
}

//
func createNewTask(conf CollectConf) {

	//初始化tail实例
	tails, err := tail.TailFile(conf.LogPath, tail.Config{
		ReOpen: 	true,
		Follow: 	true,
		Location: 	&tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: 	false,
		Poll: 		true,
	})

	if err != nil {
		logs.Error("收集日志文件[%s]错误, err:", conf.LogPath, err)
		return
	}

	//导入配置
	obj := &TailObj{
		conf:		conf,
		exitChan: 	make(chan int, 1),
	}

	obj.tail = tails
	//放到里面管理
	tailObjMgr.TailObjs = append(tailObjMgr.TailObjs, obj)

	go readFromTail(obj)
}

func readFromTail(tailObj *TailObj) {
	for true {
		select {
		//取出日志
		case msg, ok := <-tailObj.tail.Lines:
			if !ok {
				logs.Warn("文件关闭从新打开， filename:%s\n", tailObj.tail.Filename)
				time.Sleep(100 * time.Millisecond) // 0.1s
				continue
			}
			//消息
			TextMsg := &TextMsg{
				Msg:	msg.Text,
				Topic: 	tailObj.conf.Topic,
			}
			//放到chan
			tailObjMgr.msgChan <- TextMsg
		// 如果exitChan==1,删除对应配置
		case <-tailObj.exitChan:
			logs.Warn("tail obj 退出，配置为conf:%v", tailObj.conf)
			return
		}

	}

}

// 新增etcd配置项
func UpdateConfig(confs []CollectConf) (err error) {
	//// 加入锁, 防止多个goroutine工作
	//tailObjMgr.lock.Lock()
	//defer tailObjMgr.lock.Unlock()

	// 创建新的tailtask
	for _, oneConf := range confs {
		// 对于已经运行的所有实例, 路径是否一样
		var isRuning = false
		for _, obj := range tailObjMgr.TailObjs {
			// 路径一样则证明是同一实例
			if oneConf.LogPath == obj.conf.LogPath {
				isRuning = true
				obj.status = StatusNormal
				break
			}
		}

		// 检查是否已经存在
		if isRuning {
			continue
		}

		// 如果不存在该配置项 新建一个tailtask任务
		createNewTask(oneConf)
	}

	// 遍历所有查看是否存在删除操作
	var tailObjs []*TailObj
	for _, obj := range tailObjMgr.TailObjs {
		obj.status = StatusDelete
		for _, oneConf := range confs {
			if oneConf.LogPath == obj.conf.LogPath {
				obj.status = StatusNormal
				break
			}
		}
		// 如果status为删除, 则将exitChan置为1
		if obj.status == StatusDelete {
			obj.exitChan <- 1
		}
		// 将obj存入临时的数组中
		tailObjs = append(tailObjs, obj)
	}
	// 将临时数组传入tailsObjs中
	tailObjMgr.TailObjs = tailObjs
	return
}

//获取一条
func GetOneLine() (msg *TextMsg) {

	msg = <-tailObjMgr.msgChan
	return
}

