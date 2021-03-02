package odkcmdapp

import (
	"errors"
	"strings"
)

// 2020/5/27 by oldkingnearby@gmail.com

// 解析命令行参数 并处理特定任务
type OdkCmd struct {
	Text     string //命令基本文本
	Reply    string //回复文本
	UserName string
	sep      string   //分割符 默认为空格
	method   string   //方法
	params   []string //参数
	status   int      //当前状态 handler_in handler_out
}

func (oc *OdkCmd) Method() string {
	return oc.method
}
func (oc *OdkCmd) Params() []string {
	return oc.params
}
func (oc *OdkCmd) Status() int {
	return oc.status
}
func (oc *OdkCmd) Sep() string {
	return oc.sep
}
func (oc *OdkCmd) QuitStatus() {
	oc.status = HANDLER_OUT
}

// 初始化命令
func InitOdkCmd(username, text, sep string) (ret OdkCmd, err error) {
	ret.UserName = username
	ret.Text = text
	ret.sep = sep
	strArr := strings.Split(strings.TrimSpace(text), sep)
	if len(strArr) < 1 {
		err = errors.New("空文本")
		return
	}
	ret.method = strings.ToLower(strArr[0])
	if len(strArr) > 1 {
		ret.params = strArr[1:]
	}
	return
}
func InitOdkCmdSpace(username, text string) (ret OdkCmd, err error) {
	ret, err = InitOdkCmd(username, text, " ")
	return
}

const (
	HANDLER_NEXT  = 1
	HANDLER_ABORT = 2
	HANDLER_IN    = 3 //进入Handler状态
	HANDLER_OUT   = 4 //退出此Handler状态
)

type OdkCmdHandlerFun func(*OdkCmd) int

type OdkCmdApp struct {
	handlers     []OdkCmdHandlerFun
	userStatus   map[string]int //用户状态 指向第几个handler
	DefaultReply string
}

// 初始化处理函数
func (oca *OdkCmdApp) InitHandlers(handlers ...OdkCmdHandlerFun) {
	oca.handlers = handlers
	oca.userStatus = make(map[string]int)

}

// 添加处理函数
func (oca *OdkCmdApp) AddHandlers(handlers ...OdkCmdHandlerFun) {
	oca.handlers = append(oca.handlers, handlers...)
}

// 处理一条命令
func (oca *OdkCmdApp) ParseOneCmd(cmd *OdkCmd) {
	handlerIndex, ok := oca.userStatus[cmd.UserName]
	if ok {
		cmd.status = HANDLER_IN
		status := oca.handlers[handlerIndex](cmd)
		switch status {
		case HANDLER_ABORT:
			return
		case HANDLER_IN:
			oca.userStatus[cmd.UserName] = handlerIndex
		case HANDLER_OUT:
			delete(oca.userStatus, cmd.UserName)
		}
		return
	}
	for i, handler := range oca.handlers {
		status := handler(cmd)
		switch status {
		case HANDLER_ABORT:
			return
		case HANDLER_NEXT:
			continue
		case HANDLER_IN:
			oca.userStatus[cmd.UserName] = i
		case HANDLER_OUT:
			delete(oca.userStatus, cmd.UserName)
		}
	}
	if cmd.Reply == "" {
		cmd.Reply = oca.DefaultReply
	}
}

func Help(cmd *OdkCmd) (status int) {
	if cmd.method == "help" || cmd.method == "/help" {
		cmd.Reply = "你正在访问帮助文档"
		status = HANDLER_ABORT
		return
	}
	status = HANDLER_NEXT
	return
}
