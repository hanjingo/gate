package com

import (
	pc4 "github.com/hanjingo/protocol/v4"
)

//通用消息
type Msg pc4.Messager

//未知编码id
const UNKNOWN_CODECID uint8 = 0

//掩码
const (
	MASK_ACTION uint32 = 0xff000000 //动作掩码
)

//动作
const (
	//0~7 基础功能
	ACTION_ROUTE uint32 = 0 //路由
	ACTION_PING  uint32 = 1 //ping
	ACTION_PONG  uint32 = 2 //pong

	//8～15 客户端
	ACTION_REQ uint32 = 8 //请求
	ACTION_SUB uint32 = 9 //订阅

	//32～63 服务器
	ACTION_REG_SERVER uint32 = 32 //注册服务
	ACTION_RSP        uint32 = 33 //服务器返回消息
	ACTION_MULTCAST   uint32 = 34 //多播
	ACTION_BROADCAST  uint32 = 35 //广播
	ACTION_PUB        uint32 = 36 //发布

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
