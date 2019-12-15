package com

import (
	"github.com/hanjingo/util"
)

const START_TIME_STR string = "2019-01-01 00:00:00"

var START_TIME = util.TimeStampToTime(START_TIME_STR)

//未知编码id
const UNKNOWN_CODECID uint8 = 0

//掩码
const (
	MASK_ACTION uint32 = 0xff000000 //动作掩码
)

//动作
const (
	//0~7 基础功能
	ACTION_ROUTE         uint32 = 0 //路由
	ACTION_PING          uint32 = 1 //ping
	ACTION_PONG          uint32 = 2 //pong
	ACTION_NEW_AGENT     uint32 = 3 //建立端点请求
	ACTION_NEW_AGENT_RSP uint32 = 4 //建立端点返回

	//8～15 客户端
	ACTION_REQ   uint32 = 8  //请求
	ACTION_SUB   uint32 = 9  //订阅
	ACTION_UNSUB uint32 = 10 //取消订阅

	//32～63 服务器
	ACTION_RSP          uint32 = 32 //服务器返回消息
	ACTION_REG_SERVER   uint32 = 33 //注册服务
	ACTION_UNREG_SERVER uint32 = 34 //取消注册服务
	ACTION_MULTCAST     uint32 = 35 //多播
	ACTION_BROADCAST    uint32 = 36 //广播
	ACTION_PUB          uint32 = 37 //发布

	//128~254 控制码
	ACTION_CLOSE_AGENT              uint32 = 128 //关闭节点
	ACTION_SET_AGENT_FILTED         uint32 = 129 //过滤
	ACTION_SET_AGENT_UNFILTED       uint32 = 130 //解除过滤
	ACTION_SET_AGENT_LIMIT_STREAM   uint32 = 131 //设置节点控流限制
	ACTION_SET_AGENT_UNLIMIT_STREAM uint32 = 132 //设置节点不限流
	ACTION_CHANGE_PERM              uint32 = 133 //变更权限

	//255 系统码
	ACTION_SYS uint32 = 255
)

//权限
const (
	PERM_NORMAL  uint32 = 7   //普通权限
	PERM_CLI     uint32 = 15  //客户端权限
	PERM_SERVER  uint32 = 63  //服务器权限
	PERM_CONTROL uint32 = 254 //控制权限
	PERM_SYS     uint32 = 255 //系统权限
)
